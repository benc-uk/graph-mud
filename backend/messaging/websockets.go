package messaging

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Map of all connected users, by username
var clientConnections = make(map[string]*websocket.Conn)
var Version = "0.0.1"

// GameMessage is a message sent to a user
type GameMessage struct {
	Source    string    `json:"source"`
	Text      string    `json:"text"`
	Type      string    `json:"type"`
	Value     string    `json:"value"`
	TimeStamp time.Time `json:"timestamp"`
}

type ConnectRequest struct {
	Username string `json:"username"`
}

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/connect", userConnect)
}

func userConnect(resp http.ResponseWriter, req *http.Request) {
	log.Println("### Player connecting...")
	wsConn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	var connectRequest ConnectRequest
	err = wsConn.ReadJSON(&connectRequest)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("### Player '%s' connected OK", connectRequest.Username)

	// Store connection for this user
	clientConnections[connectRequest.Username] = wsConn

	hostname, _ := os.Hostname()
	msg := fmt.Sprintf("⚔️ Welcome to Nano Realms v%s - you are connected to server %s", Version, hostname)
	SendToUser(connectRequest.Username, "⚔️ "+connectRequest.Username+" connected", "server", "connection")
	SendToUser(connectRequest.Username, msg, "server", "connection")
}

// Send a message to a specific user
func SendToUser(username string, message string, source string, typeStr string) {
	conn := clientConnections[username]
	if conn == nil {
		return
	}

	log.Printf("### Sending message to user '%s': %s", username, message)

	err := conn.WriteJSON(GameMessage{
		Source:    source,
		Type:      typeStr,
		Text:      message,
		TimeStamp: time.Now(),
	})

	if err != nil {
		log.Println(err)
	}
}
