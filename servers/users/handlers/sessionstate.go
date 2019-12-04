package handlers

import (
	"assignments-zhouyifan0904/servers/gateway/models/users"
	"time"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

// SessionState struct to keep track of when the session was started and the user
// associated with the session
type SessionState struct {
	StartTime time.Time `json:"time"`
	user      *users.User
}
