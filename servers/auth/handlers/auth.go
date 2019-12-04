package handlers

import (
	"GroupPoop/servers/auth/models/users"
	"GroupPoop/servers/auth/sessions"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// SessionsHandler starts a new session using credentials provided through a
// POST request
func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if ctx.UserStore == nil {
		http.Error(w, "invalid context", http.StatusUnauthorized)
		return
	}
	if ctx.SessionStore == nil {
		http.Error(w, "invalid context", http.StatusUnauthorized)
		return
	}
	if len(ctx.SigningKey) == 0 {
		http.Error(w, "invalid context", http.StatusUnauthorized)
		return
	}
	// check request type
	if r.Method != http.MethodPost {
		http.Error(w, "only support POST request", http.StatusMethodNotAllowed)
		return
	}

	// check header
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
		return
	}

	// decode into a users.Credentials struct
	cred := &users.Credentials{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(cred)
	if err != nil {
		http.Error(w, "error decoding JSON", http.StatusUnsupportedMediaType)
		return
	}

	// authenticate
	usr, err := ctx.UserStore.GetByEmail(cred.Email)
	if err != nil {
		// wait for some time and notify user
		time.Sleep(1 * time.Second)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	err = usr.Authenticate(cred.Password)
	if err != nil {
		// wait for some time and notify userftime
		time.Sleep(1 * time.Second)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// begin new session
	_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, usr, w)
	if err != nil {
		// error starting session
		http.Error(w, "failed to start session", http.StatusUnauthorized)
		return
	}

	// log session sign-in
	err = ctx.UserStore.Log(time.Now(), r.RemoteAddr)
	if err != nil {
		// error logging sign in
		http.Error(w, "Error logging sign in", http.StatusUnauthorized)
		return
	}

	// success, respond to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usr)
}

// SpecificSessionHandler ends a session
func (ctx *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		// only support DELETE for sign out

		// check last path segment
		path := r.URL.Path

		segments := strings.Split(path, "/")
		last := segments[len(segments)-1]
		if last != "mine" {
			http.Error(w, "trying to delete a session that's not \"mine\"", http.StatusForbidden)
			return
		}

		// end session
		_, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore)
		if err != nil {
			http.Error(w, "failed to end session", http.StatusForbidden)
			return
		}
		w.Write([]byte("signed out\n"))
	} else {
		http.Error(w, "unsupported http method", http.StatusMethodNotAllowed)
	}
}
