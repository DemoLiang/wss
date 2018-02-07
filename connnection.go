package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var online int

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The web socket connection
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// reader pumps messages from the websocket connection to the hub.
func (c *connection) reader() {
	for {
		fmt.Printf("start\n")
		_, message, err := c.ws.ReadMessage()
		fmt.Printf("message:%s\n",message)
		if err != nil {
			fmt.Println(err)
			break
		}
		h.broadcast <- message
	}
	c.ws.Close()
	online -= 1
	fmt.Println(online)
}

// write writes a message with the given message type and payload.
func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	c.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool {
	return true
}}

// wsHandler handles websocket requests from the peer.
func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	online += 1
	fmt.Println(online)
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}
