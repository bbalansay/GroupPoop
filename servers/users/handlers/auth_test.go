package handlers

import (
	"assignments-zhouyifan0904/servers/gateway/models/users"
	"assignments-zhouyifan0904/servers/gateway/sessions"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"
	"time"

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

func TestSessionsHandler(t *testing.T) {

	// test successful sign in

	// make request body (a users.Credentials)
	testBody := new(bytes.Buffer)
	json.NewEncoder(testBody).Encode(users.Credentials{
		Email:    "email@abc.com",
		Password: "password",
	})
	// construct request
	req, err := http.NewRequest("POST", "/sessions", testBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	// set up mock db so that GetByEmail, Authenticate and BeginSession succeeds
	// db, mock, err := sqlmock.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	ms := users.NewMySQLStore(db)
	// set up memstore with duration and purge interval of 3 minutes
	dur, err := time.ParseDuration("3m")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating memstore", err)
	}
	memStore := sessions.NewMemStore(dur, dur)
	// set up: construct context
	ctx := HandlerContext{"mocksigningkey", memStore, ms}
	// expect getting user with email
	columns := []string{"id", "email", "passhash", "user_name", "first_name", "last_name", "photo_url"}
	mock.ExpectQuery("select (.+) from users where email = ?").
		WithArgs("email@abc.com").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(123456, "email@abc.com", "$2a$11$MXOOO1JYngri2arcL6Cic.KuBujhqgz.B2ri6szqN2/cfsdiQa7se", "abc", "a", "b", "abc.com/123"))
	// expect insert log
	mock.ExpectExec("insert into logs").
		WillReturnResult(sqlmock.NewResult(0, 1))
	// serve and record response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	// check status code
	status := rr.Code
	if status != http.StatusCreated {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusCreated)
	}
	// check response body (a users.User object) - newline at the end
	expected := []byte(`{"id":123456,"userName":"abc","firstName":"a","lastName":"b","photoURL":"abc.com/123"}` + "\n")
	if !bytes.Equal(rr.Body.Bytes(), expected) {
		t.Errorf("SessionsHandler returned unexpected body: got \n%v \ninstead of \n%v", string(rr.Body.Bytes()), string(expected))
	}

	// test wrong http method

	req, err = http.NewRequest("GET", "/sessions", testBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	// serve and record response
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	// expect http.StatusUnsupportedMediaType
	status = rr.Code
	if status != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusMethodNotAllowed)
	}

	// test request not json

	req, err = http.NewRequest("POST", "/sessions", testBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/html")
	// serve and record response
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	// expect http.StatusUnsupportedMediaType
	status = rr.Code
	if status != http.StatusUnsupportedMediaType {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusUnsupportedMediaType)
	}

	// test user not found

	testBody = new(bytes.Buffer)
	json.NewEncoder(testBody).Encode(users.Credentials{
		Email:    "email@abc.com",
		Password: "password",
	})
	req, err = http.NewRequest("POST", "/sessions", testBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	// expect getting user with email, but return 0 row
	columns = []string{"id", "email", "passhash", "user_name", "first_name", "last_name", "photo_url"}
	mock.ExpectQuery("select (.+) from users where email = ?").
		WithArgs("email@abc.com").
		WillReturnRows(sqlmock.NewRows(columns))
	// serve and record response
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	// check status code - expect StatusUnauthorized
	status = rr.Code
	if status != http.StatusUnauthorized {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusUnauthorized)
	}

	// test user db error

	testBody = new(bytes.Buffer)
	json.NewEncoder(testBody).Encode(users.Credentials{
		Email:    "email@abc.com",
		Password: "password",
	})
	req, err = http.NewRequest("POST", "/sessions", testBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	// serve and record response
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	// check status code - expect StatusUnauthorized
	status = rr.Code
	if status != http.StatusUnauthorized {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusUnauthorized)
	}

	// test wrong password

	testBody = new(bytes.Buffer)
	json.NewEncoder(testBody).Encode(users.Credentials{
		Email:    "email@abc.com",
		Password: "wordpass",
	})
	req, err = http.NewRequest("POST", "/sessions", testBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	// expect getting user with email, but return 0 row
	columns = []string{"id", "email", "passhash", "user_name", "first_name", "last_name", "photo_url"}
	mock.ExpectQuery("select (.+) from users where email = ?").
		WithArgs("email@abc.com").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(123456, "email@abc.com", "$2a$11$MXOOO1JYngri2arcL6Cic.KuBujhqgz.B2ri6szqN2/cfsdiQa7se", "abc", "a", "b", "abc.com/123"))
	// serve and record response
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	// check status code - expect StatusUnauthorized
	status = rr.Code
	if status != http.StatusUnauthorized {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusUnauthorized)
	}

	// test invalid json

	testBody = bytes.NewBuffer([]byte{0, 1, 2, 3})

	req, err = http.NewRequest("POST", "/sessions", testBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	expectedString := "error decoding JSON\n"
	if rr.Body.String() != expectedString {
		t.Errorf("wrong error returned: got %v instead of %v\n", rr.Body.String(), expectedString)
	}

	// test nil user store

	testBody = new(bytes.Buffer)
	json.NewEncoder(testBody).Encode(users.Credentials{
		Email:    "email@abc.com",
		Password: "wordpass",
	})
	req, err = http.NewRequest("POST", "/sessions", testBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	ctx.UserStore = nil
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	if rr.Body.String() != "invalid context\n" {
		t.Errorf("wrong error returned: got %v instead of %v\n", rr.Body.String(), "invalid context\n")
	}

	// test nil session store

	testBody = new(bytes.Buffer)
	json.NewEncoder(testBody).Encode(users.Credentials{
		Email:    "email@abc.com",
		Password: "wordpass",
	})
	req, err = http.NewRequest("POST", "/sessions", testBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	ctx = HandlerContext{"mocksigningkey", nil, ms}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	if rr.Body.String() != "invalid context\n" {
		t.Errorf("wrong error returned: got %v instead of %v\n", rr.Body.String(), "invalid context\n")
	}

	// test nil session store

	testBody = new(bytes.Buffer)
	json.NewEncoder(testBody).Encode(users.Credentials{
		Email:    "email@abc.com",
		Password: "wordpass",
	})
	req, err = http.NewRequest("POST", "/sessions", testBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	ctx = HandlerContext{"", memStore, ms}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SessionsHandler)
	handler.ServeHTTP(rr, req)
	if rr.Body.String() != "invalid context\n" {
		t.Errorf("wrong error returned: got %v instead of %v\n", rr.Body.String(), "invalid context\n")
	}
}

func TestSpecificSessionHandler(t *testing.T) {

	// test successful deletion

	dur, err := time.ParseDuration("3m")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating memstore", err)
	}
	memStore := sessions.NewMemStore(dur, dur)
	ctx := HandlerContext{"mocksigningkey", memStore, nil}
	// begin a mock session
	rr := httptest.NewRecorder()
	sessID, err := sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, ctx.UserStore, rr)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a stub session", err)
	}
	// construct delete request
	req, err := http.NewRequest("DELETE", "/v1/sessions/mine", nil)
	if err != nil {
		t.Fatal(err)
	}
	sessIDHeader := "Bearer " + sessID
	req.Header.Set("Authorization", sessIDHeader.String())
	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.SpecificSessionHandler)
	handler.ServeHTTP(rr, req)
	res := rr.Body.String()
	expected := "signed out\n"
	if res != expected {
		t.Errorf("wrong result returned: got %v instead of %v\n", res, expected)
	}

	// test wrong http method

	req, err = http.NewRequest("POST", "/v1/sessions/mine", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificSessionHandler)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	if status != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusMethodNotAllowed)
	}

	// test wrong url

	sessID, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, ctx.UserStore, rr)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a stub session", err)
	}
	req, err = http.NewRequest("DELETE", "/v1/sessions/min", nil)
	if err != nil {
		t.Fatal(err)
	}
	sessIDHeader = "Bearer " + sessID
	req.Header.Set("Authorization", sessIDHeader.String())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificSessionHandler)
	handler.ServeHTTP(rr, req)
	status = rr.Code
	if status != http.StatusForbidden {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusForbidden)
	}

	// test empty memstore
	sessID, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, ctx.UserStore, rr)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a stub session", err)
	}
	req, err = http.NewRequest("DELETE", "/v1/sessions/mine", nil)
	if err != nil {
		t.Fatal(err)
	}
	sessIDHeader = "Bearer " + "invalidsessionid"
	req.Header.Set("Authorization", sessIDHeader.String())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ctx.SpecificSessionHandler)
	handler.ServeHTTP(rr, req)
	status = rr.Code
	if status != http.StatusForbidden {
		t.Errorf("wrong status code returned: got %v instead of %v\n", status, http.StatusForbidden)
	}
	expected = "failed to end session\n"
	if rr.Body.String() != expected {
		t.Errorf("wrong error returned: got %v instead of %v\n", rr.Body.String(), expected)
	}
}
