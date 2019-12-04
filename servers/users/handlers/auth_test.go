package handlers

import (
	"GroupPoop/servers/users/models/users"
	"GroupPoop/servers/users/sessions"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUsersHandler(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	ms := users.NewMySQLStore(db)

	dur, err := time.ParseDuration("3m")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating memstore", err)
	}
	memStore := sessions.NewMemStore(dur, dur)

	ctx := NewHandlerContext("mocksigningkey", memStore, ms)

	// TEST NORMAL
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(users.NewUser{
		Email:        "MyEmailAddress@example.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "username",
		FirstName:    "first",
		LastName:     "last",
	})

	mock.ExpectExec("insert into users").
		WillReturnResult(sqlmock.NewResult(1, 1))

	r, _ := http.NewRequest(http.MethodPost, "/v1/users", body)
	r.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status := rr.Code
	if status != http.StatusCreated {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusCreated)
	}
	// check response body (a users.User object)
	expected := `{"id":1,"userName":"username","firstName":"first","lastName":"last","photoURL":"https://www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346"}` + "\n"

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if rr.Body.String() != expected {
		t.Errorf("UsersHandler returned unexpected body: got \n%vinstead of \n%v\n", rr.Body.String(), expected)
	}

	// CONTENT-TYPE NOT JSON USERS
	r, _ = http.NewRequest(http.MethodPost, "/v1/users", body)
	r.Header.Add("Content-Type", "application/text")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusUnsupportedMediaType {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusUnsupportedMediaType)
	}

	// METHOD NOT ALLOWED
	r, _ = http.NewRequest(http.MethodDelete, "/v1/users", body)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusMethodNotAllowed)
	}

	r, _ = http.NewRequest(http.MethodPost, "/v1/users", body)
	r.Header.Add("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusBadRequest {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusBadRequest)
	}

	// CAN'T VALIDATE USER
	body = new(bytes.Buffer)
	updates := users.Updates{
		FirstName: "NewFirstName",
		LastName:  "NewLastName",
	}
	json.NewEncoder(body).Encode(updates)

	r, _ = http.NewRequest(http.MethodPost, "/v1/users", body)
	r.Header.Add("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.UsersHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusBadRequest {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusBadRequest)
	}
}

func TestSpecificUserHandler(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	ms := users.NewMySQLStore(db)

	dur, err := time.ParseDuration("3m")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating memstore", err)
	}
	memStore := sessions.NewMemStore(dur, dur)
	ctx := NewHandlerContext("mocksigningkey", memStore, ms)

	body := new(bytes.Buffer)
	user := users.User{
		ID:        1,
		Email:     "test@test.com",
		PassHash:  []byte("password"),
		UserName:  "username",
		FirstName: "first",
		LastName:  "last",
		PhotoURL:  "photo.com",
	}
	json.NewEncoder(body).Encode(user)

	// TEST NOT AUTHORIZED GET
	r, _ := http.NewRequest(http.MethodGet, "/v1/users/me", nil)
	r.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status := rr.Code
	if status != http.StatusUnauthorized {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusUnauthorized)
	}

	// TEST NOT AUTHORIZED PATCH
	r, _ = http.NewRequest(http.MethodPatch, "/v1/users/me", nil)
	r.Header.Add("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusUnauthorized {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusUnauthorized)
	}

	// SAVE SESSION
	sid, err := sessions.NewSessionID("mocksigningkey")
	ctx.SessionStore.Save(sid, user)

	// TEST BASIC
	columns := []string{"id", "email", "passhash", "user_name", "first_name", "last_name", "photo_url"}
	mock.ExpectQuery("select (.+) from users where id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL))

	r, _ = http.NewRequest(http.MethodGet, "/v1/users/1", nil)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer " + sid.String())

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusOK {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusOK)
	}
	// check response body (a users.User object)
	expected := `{"id":1,"userName":"username","firstName":"first","lastName":"last","photoURL":"photo.com"}` + "\n"

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if rr.Body.String() != expected {
		t.Errorf("SpecificUserHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// TEST GET FOR /me ENDPOINT
	mock.ExpectQuery("select (.+) from users").
		WithArgs(&user.ID).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL))

	r, _ = http.NewRequest(http.MethodGet, "/v1/users/me", body)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+sid.String())

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusOK {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	expected = `{"id":1,"userName":"username","firstName":"first","lastName":"last","photoURL":"photo.com"}` + "\n"

	if rr.Body.String() != expected {
		t.Errorf("SpecificUserHandler returned unexpected body: got \n%vinstead of \n%v\n", rr.Body.String(), expected)
	}

	// TEST PATCH FOR /me ENDPOINT
	body = new(bytes.Buffer)
	updates := users.Updates{
		FirstName: "NewFirstName",
		LastName:  "NewLastName",
	}
	json.NewEncoder(body).Encode(updates)

	mock.ExpectExec("update users").
		WithArgs(updates.FirstName, updates.LastName, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("select (.+) from users").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(&user.ID, &user.Email, &user.PassHash, &user.UserName, "NewFirstName", "NewLastName", &user.PhotoURL))

	r, _ = http.NewRequest(http.MethodPatch, "/v1/users/1", body)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+sid.String())

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusOK {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusOK)
	}

	expected = `{"id":1,"userName":"username","firstName":"NewFirstName","lastName":"NewLastName","photoURL":"photo.com"}` + "\n"

	if rr.Body.String() != expected {
		t.Errorf("SpecificUserHandler returned unexpected body: got \n%vinstead of \n%v\n", rr.Body.String(), expected)
	}


	// TEST USER DOES NOT EXIST
	mock.ExpectQuery("select (.+) from users where id = ?").
		WithArgs(999).
		WillReturnError(errors.New("user not found"))

	r, _ = http.NewRequest(http.MethodGet, "/v1/users/999", nil)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer " + sid.String())

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusNotFound {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusNotFound)
	}


	// TEST METHOD NOT ALLOWED
	r, _ = http.NewRequest(http.MethodDelete, "/v1/users/1", nil)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusMethodNotAllowed)
	}

	// TEST IMPROPER RESOURCE PATH
	r, _ = http.NewRequest(http.MethodGet, "/v1/users/1x", nil)
	r.Header.Add("Authorization", "Bearer "+sid.String())

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusBadRequest {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusBadRequest)
	}

	// IMPROPER RESOURCE PATH PATCH
	r, _ = http.NewRequest(http.MethodPatch, "/v1/users/1x", body)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+sid.String())

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusBadRequest {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusBadRequest)
	}

	r, _ = http.NewRequest(http.MethodPatch, "/v1/users/19999", body)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+sid.String())

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusForbidden {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusForbidden)
	}

	// CONTENT-TYPE NOT JSON
	r, _ = http.NewRequest(http.MethodPatch, "/v1/users/me", body)
	r.Header.Add("Authorization", "Bearer "+sid.String())

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificUserHandler)
	handler.ServeHTTP(rr, r)

	// check status code
	status = rr.Code
	if status != http.StatusUnsupportedMediaType {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusUnsupportedMediaType)
	}
}