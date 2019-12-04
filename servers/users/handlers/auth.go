package handlers

import (
	"assignments-zhouyifan0904/servers/gateway/models/users"
	"assignments-zhouyifan0904/servers/gateway/sessions"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

// UsersHandler handles requests for the "users" resource
// Accepts POST requests to create new user accounts
func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost: // If request method is POST
		// Content-Type must be JSON
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			http.Error(w, "Request body must be in JSON.", http.StatusUnsupportedMediaType)
			return
		}

		// Decode JSON into NewUser struct
		newUser := &users.NewUser{}
		err := json.NewDecoder(r.Body).Decode(newUser)
		if err != nil {
			http.Error(w, "Could not decode JSON", http.StatusBadRequest)
			return
		}

		// Validate NewUser and create User
		user, err := newUser.ToUser()
		if err != nil {
			http.Error(w, "Could not validate user", http.StatusBadRequest)
			return
		}

		// Insert user into user store
		user, err = ctx.UserStore.Insert(user)
		if err != nil {
			http.Error(w, "Could not enter user into database", http.StatusBadRequest)
			return
		}

		// Begin session for new user in session store
		_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, user, w)
		if err != nil {
			http.Error(w, "Could not validate user", http.StatusBadRequest)
			return
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificUserHandler handles requests for a specific user.
// The resource path will be /v1/users/{UserID}
// Implement functionality so that going to /v1/users/me
// will refer to the UserID of the currently-authenticated user.
// User must be authenticated to call this handler regardless of HTTP method
func (ctx *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: // If request method is GET
		// Parse ID from resource path
		idString := path.Base(r.URL.Path)
		
		user := &users.User{}
		_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, user)
		if err != nil {
			http.Error(w, "Could not retrieve profile for authenticated user", http.StatusUnauthorized)
			return
		}
		id := user.ID

		if idString != "me" {
			idInt, err := strconv.Atoi(idString)
			id = int64(idInt)
			if err != nil {
				http.Error(w, "Did not provide {UserId} as a number in /v1/users/{UserId}, please provide the correct ID", http.StatusBadRequest)
				return
			}
		}

		// Retrieve user from user store
		user, _ = ctx.UserStore.GetByID(id)
		if user == nil {
			http.Error(w, "Could not find user with that ID", http.StatusNotFound)
			return
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	case http.MethodPatch: // If request method is PATCH
		// Retrieve user from store
		user := &users.User{}
		_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, user)
		if err != nil {
			http.Error(w, "Could not retrieve profile for authenticated user", http.StatusUnauthorized)
			return
		}

		// Parse ID from resource path
		path := r.URL.Path
		segments := strings.Split(path, "/")
		id := segments[len(segments)-1]
		if id != "me" {
			id, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				http.Error(w, "Did not provide {UserId} as a number in /v1/users/{UserId}, please provide the correct ID", http.StatusBadRequest)
				return
			}
			if id != user.ID {
				http.Error(w, "{UserID} does not match profile in /v1/users/{UserId}, use resource path /v1/users/me", http.StatusForbidden)
				return
			}
		}

		// Content-Type must be JSON
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			http.Error(w, "Request body must be in JSON.", http.StatusUnsupportedMediaType)
			return
		}

		// Decode JSON into Updates struct
		updates := &users.Updates{}
		err = json.NewDecoder(r.Body).Decode(updates)
		if err != nil {
			http.Error(w, "Could not decode JSON", http.StatusBadRequest)
			return
		}

		// Update user in user store
		user, err = ctx.UserStore.Update(user.ID, updates)

		if err != nil {
			http.Error(w, "Could not update user into database", http.StatusBadRequest)
			return
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
}

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
