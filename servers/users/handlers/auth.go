package handlers

import (
	"GroupPoop/servers/users/models/users"
	"GroupPoop/servers/users/sessions"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// UsersHandler handles requests for the "users" resource
// Accepts POST requests to create new user accounts
func (ctx *HandlerContext) UserHandler(w http.ResponseWriter, r *http.Request) {
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

		//retrieve reviews associated with user
		reviews, _ := ctx.UserStore.GetReviews(id)

		result := UserAndReview{
			User:    user,
			Reviews: reviews,
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
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
	case http.MethodDelete:
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

		// delete user in user store
		err = ctx.UserStore.Delete(user.ID)
		if err != nil {
			http.Error(w, "Could not delete user from database", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
}
