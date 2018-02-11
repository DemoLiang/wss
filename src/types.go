package main

import "github.com/gorilla/websocket"

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The web socket connection
	Ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}

// hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	Connections map[*Connection]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Connection

	// Unregister requests from clients.
	Unregister chan *Connection
}


type GameRoom struct {
	// Registered game user clients.
	connections map[*Connection]bool

	// Inbound messages from the other clients. and broadcast to other's clients
	broadcast chan []byte

	// Register clients from the pool
	register chan *Connection

	// Unregister clinets from game room
	unregister chan *Connection
}
