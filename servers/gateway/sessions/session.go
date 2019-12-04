package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

/*
sessions maintains states with specific users by assigning an ID to each of them in the HTML header and keeping track of recent users.
*/

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	//Create a new SessionID
	sid, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	//Set the session state
	err = store.Save(sid, sessionState)
	if err != nil {
		return InvalidSessionID, err
	}

	//Add authentication header in the form "Bearer <sid>"
	w.Header().Set(headerAuthorization, schemeBearer + sid.String())

	return sid, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//TODO: get the value of the Authorization header,
	//or the "auth" query string parameter if no Authorization header is present,
	//and validate it. If it's valid, return the SessionID. If not
	//return the validation error.

	//Get auth header if present
	authHeader := r.Header.Get(headerAuthorization)
	authToken := ""

	//If not present, attempt to get "auth" query string parameter
	if authHeader == "" {
		authHeader = r.FormValue("auth")
		if authHeader == "" {
			return InvalidSessionID, errors.New("invalid authorization parameter")
		}
		authToken = authHeader
	} else {
		//Split header into Bearer and sid
		auth := strings.SplitN(authHeader, " ", 2)

		if len(auth) != 2 || auth[0] != "Bearer" {
			return InvalidSessionID, errors.New("invalid authorization scheme")
		}

		authToken = auth[1]
	}

	//Validate and return sid
	sid, err := ValidateID(authToken, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	
	return sid, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	sid, err := GetSessionID(r, signingKey)
	
	if err != nil {
		return InvalidSessionID, err
	}

	//Get data associated with sid
	err = store.Get(sid, sessionState) 
	if err != nil {
		return InvalidSessionID, err
	}

	return sid, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	//Delete data associated with sid
	err = store.Delete(sid)
	if err != nil {
		return InvalidSessionID, err
	}
	
	return sid, nil
}
