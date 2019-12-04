package handlers

import (
	"GroupPoop/servers/users/models/users"
	"assignments-zhouyifan0904/servers/gateway/sessions"
	"encoding/json"
	"net/http"
	"strings"
)

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
