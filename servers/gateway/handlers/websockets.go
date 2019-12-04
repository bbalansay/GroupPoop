package handlers

import (
	"GroupPoop/servers/gateway/models/users"
	"GroupPoop/servers/gateway/sessions"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutines.

// A simple store to store all the connections
type SocketStore struct {
	// maps user id to a connection
	Connections map[int64]*websocket.Conn
	lock        sync.Mutex
}

// Control messages for websocket
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

// Thread-safe method for putting a userId : connection mapping
func (s *SocketStore) PutConnection(userId int64, conn *websocket.Conn) {
	s.lock.Lock()
	s.Connections[userId] = conn
	s.lock.Unlock()
}

// Thread-safe method for removing a user's web socket connection
func (s *SocketStore) RemoveConnection(userId int64) {
	s.lock.Lock()
	// insert socket connection
	delete(s.Connections, userId)
	s.lock.Unlock()
}

func (s *SocketStore) Broadcast(message []byte) {
	for userID := range s.Connections {
		writeError := s.Connections[userID].WriteMessage(TextMessage, message)
		if writeError != nil {
			log.Printf("Error sending through socket")
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// This function's purpose is to reject websocket upgrade requests if the
		// origin of the websockete handshake request is coming from unknown domains.
		// This prevents some random domain from opening up a socket with your server.
		// TODO: make sure you modify this for your HW to check if r.Origin is your host
		return true
	},
}

// WebsocketConnectionHandler handles requests to /v1/ws. Upgrades the client
// to a websocket and start the goroutine to send messages to the client
func (ctx *HandlerContext) WebsocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// check the user is authorized before upgrading
	_, err := sessions.GetSessionID(r, ctx.SigningKey)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// get the user info
	user := &users.User{}
	_, err = sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, user)
	if err != nil {
		http.Error(w, "Could not retrieve profile for authenticated user", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", http.StatusUnauthorized)
		return
	}

	// Insert our connection onto our datastructure for ongoing usage
	ctx.SocketStore.PutConnection(user.ID, conn)

	// Invoke a goroutine for handling control messages from this connection
	go (func(conn *websocket.Conn, userId int64) {
		defer conn.Close()
		defer ctx.SocketStore.RemoveConnection(userId)

		for {
			// read every messages
			messageType, message, err := conn.ReadMessage()
			if messageType == TextMessage {
				ctx.SocketStore.Broadcast(message)
			} else if messageType == CloseMessage {
				break
			} else if err != nil {
				log.Printf("Error reading from socket")
				// break
			}
			// ignore ping and pong messages
		}

	})(conn, user.ID)

}
