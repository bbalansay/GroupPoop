package handlers

import (
	"assignments-zhouyifan0904/servers/gateway/models/users"
	"assignments-zhouyifan0904/servers/gateway/sessions"
	"encoding/json"
	"net/http"
)

const headerAuthorization = "Authorization"

// EnsureAuth is a middleware handler that authenticates specific HTTP requests
type EnsureAuth struct {
	handler      http.Handler
	signingKey   string
	sessionStore sessions.Store
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get(headerAuthorization)
	delete(r.Header, "X-User")

	if authHeader != "" {
		user := &users.User{}
		_, err := sessions.GetState(r, ea.signingKey, ea.sessionStore, user)
		if err != nil {
			http.Error(w, "Could not authenticate user", http.StatusUnauthorized)
			return
		}

		userJSON, err := json.Marshal(user)

		if err != nil {
			http.Error(w, "Could not authenticate user", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User", string(userJSON))
	}

	ea.handler.ServeHTTP(w, r)
}

func NewEnsureAuth(handlerToWrap http.Handler, signingKey string, sessionStore sessions.Store) *EnsureAuth {
	ea := &EnsureAuth{handlerToWrap, signingKey, sessionStore}
	return ea
}
