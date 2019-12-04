package handlers

import (
    "net/http"
    "net/http/httptest"
	"testing"
	"io"
)

func headerAllowed(header string) bool {
	allowedHeaders := [7]string{"Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers",
								"Access-Control-Expose-Headers", "Access-Control-Max-Age", "Origin", "Host"}
	for _, b := range allowedHeaders {
		if header == b {
			return true
		}
	}
	return false
}

func TestCORSMiddleware(t *testing.T) {
	cases := []struct {
		method         string
		expectedStatus int
	}{
		{
			"OPTIONS",
			200,
		},
		{
			"GET",
			201,
		},
		{
			"POST",
			201,
		},
	}

	for _, c := range cases {
		fn := func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
					io.WriteString(w, "<html><body>Hello World!</body></html>")
		}

		handler := NewEnsureCORS(http.HandlerFunc(fn))

		// Create a request to pass to our handler. 
		req, err := http.NewRequest(c.method, "/mock", nil)
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != c.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, c.expectedStatus)
		}

		result := rr.Result()


		for key, _ := range result.Header {
			if allowed := headerAllowed(key); allowed == false {
				t.Errorf("Header not expected: %v", key)
			}
		}

		if ctype := rr.Header().Get("Access-Control-Allow-Origin"); ctype != "*" {
			t.Errorf("content type header does not match: got %v want %v",
				ctype, "*")
		}
		
		if ctype := rr.Header().Get("Access-Control-Allow-Methods"); ctype != "GET, PUT, POST, PATCH, DELETE" {
			t.Errorf("content type header does not match: got %v want %v",
				ctype, "GET, PUT, POST, PATCH, DELETE")
		}

		if ctype := rr.Header().Get("Access-Control-Allow-Headers"); ctype != "Content-Type, Authorization" {
			t.Errorf("content type header does not match: got %v want %v",
				ctype, "Content-Type, Authorization")
		}

		if ctype := rr.Header().Get("Access-Control-Expose-Headers"); ctype != "Authorization" {
			t.Errorf("content type header does not match: got %v want %v",
				ctype, "Authorization")
		}

		if ctype := rr.Header().Get("Access-Control-Max-Age"); ctype != "600" {
			t.Errorf("content type header does not match: got %v want %v",
				ctype, "600")
		}
	}
}
