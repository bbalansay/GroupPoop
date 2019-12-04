package handlers

import (
	"GroupPoop/servers/auth/models/users"
	"GroupPoop/servers/auth/sessions"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

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
