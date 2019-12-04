package sessions

import (
	"crypto/rand"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID.
//This is a base64 URL encoded string created from a byte slice
//where the first `idLength` bytes are crytographically random
//bytes representing the unique session ID, and the remaining bytes
//are an HMAC hash of those ID bytes (i.e., a digital signature).
//The byte slice layout is like so:
//+-----------------------------------------------------+
//|...32 crypto random bytes...|HMAC hash of those bytes|
//+-----------------------------------------------------+
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
//Bradley: See in-line comments in function
func NewSessionID(signingKey string) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("signingKey must not be empty")
	}

	// Create a slice of 32 random bytes
	idLength := make([]byte, 32)
	_, err := rand.Read(idLength)
	if err != nil {
		return InvalidSessionID, ErrInvalidID
	}

	// Create a new HMAC hasher using the signing key 
	h := hmac.New(sha256.New, []byte(signingKey))

	// Write idLength into it
	h.Write(idLength)

	// Append the signature to the random bytes
	byteArray := append(idLength, h.Sum(nil)...)

	// Encode byte slice using base64 URL Encoding and return as SessionID
	return SessionID(base64.URLEncoding.EncodeToString(byteArray)), nil
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
//Bradley: See below again
func ValidateID(id string, signingKey string) (SessionID, error) {

	//TODO: validate the `id` parameter using the provided `signingKey`.
	//base64 decode the `id` parameter, HMAC hash the
	//ID portion of the byte slice, and compare that to the
	//HMAC hash stored in the remaining bytes. If they match,
	//return the entire `id` parameter as a SessionID type.
	//If not, return InvalidSessionID and ErrInvalidID.

	//base64 decode the given id
	decodedID, _ := base64.URLEncoding.DecodeString(id)
	if len(decodedID) != 64 {
		return InvalidSessionID, ErrInvalidID
	}

	//Take the first 32 bytes to get the randomly assigned id
	idLength := decodedID[0:32]

	//Recreate HMAC signature
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write(idLength)

	//If the signature is equal to the rest of the decoded ID
	//return SessionID as validation
	if hmac.Equal(h.Sum(nil), decodedID[32:]) {
		return SessionID(id), nil
	}

	return InvalidSessionID, ErrInvalidID
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}