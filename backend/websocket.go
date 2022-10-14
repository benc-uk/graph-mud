package main

import (
	"log"
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

var clientConnections = make(map[string]*websocket.Conn)

type Message struct {
	Source string `json:"source"`
	Text   string `json:"text"`
}

type ConnectRequest struct {
	Username string `json:"username"`
}

func sender() {

	// send message to client every 5 seconds
	for {
		time.Sleep(2 * time.Second)
		for u, wsConn := range clientConnections {
			if wsConn == nil {
				continue
			}

			msg := Message{
				Source: "server",
				Text:   "PING " + u + " you " + time.Now().String(),
			}
			if err := wsConn.WriteJSON(msg); err != nil {
				log.Println(err)
			}
		}
	}

}
func wsConnect(resp http.ResponseWriter, req *http.Request) {
	log.Println("### wsConnect")
	wsConn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("### wsConnect - upgraded")

	var connectRequest ConnectRequest
	err = wsConn.ReadJSON(&connectRequest)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("### New connection from:", connectRequest.Username)

	// hold this connection for this user
	clientConnections[connectRequest.Username] = wsConn

	msg := Message{
		Source: "server",
		Text:   "Welcome!",
	}

	if err := wsConn.WriteJSON(msg); err != nil {
		log.Println(err)
		return
	}
}
