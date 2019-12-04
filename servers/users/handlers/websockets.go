package handlers

import (
	"assignments-zhouyifan0904/servers/gateway/models/users"
	"assignments-zhouyifan0904/servers/gateway/sessions"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
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

// rabit mq event format
type Message struct {
	UpdateType string `json:"type"`
	Message    struct {
		ID        int    `json:"id"`
		ChannelID int    `json:"channelID"`
		Body      string `json:"body"`
		CreatedAt string `json:"createdAt"`
		Creator   struct {
			ID        int    `json:"id"`
			UserName  string `json:"user_name"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			PhotoURL  string `json:"photo_url"`
		} `json:"creator"`
		EditedAt string `json:"editedAt"`
	} `json:"message"`
	MessageID string `json:"messageID"`
	Channel   struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Private     string `json:"private"`
		CreatedAt   string `json:"createdAt"`
		Creator     string `json:"creator"`
		EditedAt    string `json:"editedAt"`
	} `json:"channel"`
	ChannelID string  `json:"channelID"`
	UserIDs   []int64 `json:"userIDs"`
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
			messageType, _, err := conn.ReadMessage()
			if messageType == TextMessage || messageType == BinaryMessage {
				writeError := conn.WriteMessage(TextMessage, []byte("Hello from server!"))
				if writeError != nil {
					log.Printf("Error sending through socket")
				}
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

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket
func (ctx *HandlerContext) Broadcast(ch *amqp.Channel, q amqp.Queue) {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
	}

	go func() {
		for d := range msgs {
			d.Ack(false)
			message := &Message{}
			json.Unmarshal(d.Body, message)
			// extract appropriate data field to
			dataToSend := []byte("Unrecognized message type")
			switch message.UpdateType {
			case "channel-new":
				dataToSend, err = json.Marshal(message.Channel)
			case "channel-update":
				dataToSend, err = json.Marshal(message.Channel)
			case "channel-delete":
				dataToSend = []byte(message.ChannelID)
			case "message-new":
				dataToSend, err = json.Marshal(message.Message)
			case "message-update":
				dataToSend, err = json.Marshal(message.Message)
			case "message-delete":
				dataToSend = []byte(message.MessageID)
			}
			userList := message.UserIDs
			if userList == nil || len(userList) == 0 {
				// public channel, send to all connections
				for _, socketConn := range ctx.SocketStore.Connections {
					writeError := socketConn.WriteMessage(TextMessage, dataToSend)
					if writeError != nil {
					}
				}
			} else {
				// private channel, send to all members
				for _, id := range userList {
					socketConn := ctx.SocketStore.Connections[id]
					writeError := socketConn.WriteMessage(TextMessage, dataToSend)
					if writeError != nil {
						log.Printf("Error writing message: %v\n", writeError)
					}
				}
			}
			log.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
		}
	}()
}
