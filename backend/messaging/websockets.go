package messaging

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

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

// GameMessage is a message sent to a user
type GameMessage struct {
	Source string `json:"source"`
	Text   string `json:"text"`
	Type   string `json:"type"`
}

type ConnectRequest struct {
	Username string `json:"username"`
}

func SenderTestLoop() {
	// send message to client every 5 seconds
	for {
		time.Sleep(4 * time.Second)
		for u, wsConn := range clientConnections {
			if wsConn == nil {
				continue
			}

			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			SendToUser(u, fmt.Sprintf("Hello from server %d", r1.Int()), "server", "ping")
		}
	}
}

func UserConnect(resp http.ResponseWriter, req *http.Request) {
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

	SendToUser(connectRequest.Username, "OK", "server", "connection")
}

// Send a message to a specific user
func SendToUser(username string, message string, source string, typeStr string) {
	conn := clientConnections[username]
	if conn == nil {
		return
	}

	err := conn.WriteJSON(GameMessage{
		Source: source,
		Type:   typeStr,
		Text:   message,
	})

	if err != nil {
		log.Println(err)
	}
}
