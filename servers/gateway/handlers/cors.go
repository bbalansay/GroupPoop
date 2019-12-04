package handlers

import (
    "net/http"
)

/* TODO: implement a CORS middleware handler, as described
in https://drstearns.github.io/tutorials/cors/ that responds
with the following headers to all requests:

    Access-Control-Allow-Origin: *
    Access-Control-Allow-Methods: GET, PUT, POST, PATCH, DELETE
    Access-Control-Allow-Headers: Content-Type, Authorization
    Access-Control-Expose-Headers: Authorization
    Access-Control-Max-Age: 600
*/

// EnsureCORS is a middleware handler that enables requests to be callable cross-origin
type EnsureCORS struct {
    handler http.Handler
}

// ServeHTTP passes request to real handler after setting response headers.
// Despite not being listed in 'Access-Control-Allow-Methods' OPTIONS is also supported,
// otherwise preflight requests will always be denied. 
func (ec *EnsureCORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    w.Header().Set("Access-Control-Expose-Headers", "Authorization")
    w.Header().Set("Access-Control-Max-Age", "600")

    if (r.Method == "OPTIONS") {
        w.WriteHeader(http.StatusOK)
        return
    }

    ec.handler.ServeHTTP(w, r)
}

// NewEnsureCORS constructs a new EnsureCORS middleware handler
func NewEnsureCORS(handlerToWrap http.Handler) *EnsureCORS {
    ec := &EnsureCORS{handlerToWrap}
    return ec
}

