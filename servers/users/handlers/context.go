package handlers

import (
	"assignments-zhouyifan0904/servers/gateway/models/users"
	"assignments-zhouyifan0904/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

//HandlerContext used to share variables and values with handlers
type HandlerContext struct {
	SigningKey   string
	SessionStore sessions.Store
	UserStore    users.Store
	SocketStore  SocketStore
}

//NewHandlerContext constructs a new HandlerContext,
//ensuring that the dependencies are valid values
func NewHandlerContext(SigningKey string, SessionStore sessions.Store, UserStore users.Store, SocketStore SocketStore) *HandlerContext {
	if SigningKey == "" {
		panic("empty SigningKey")
	}

	if SessionStore == nil {
		panic("nil SessionStore!")
	}

	if UserStore == nil {
		panic("nil UserStore!")
	}

	if &SocketStore == nil {
		panic("nil SocketStore!")
	}

	return &HandlerContext{SigningKey, SessionStore, UserStore, SocketStore}
}
