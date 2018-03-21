package main

var h = Hub{
	Broadcast:    make(chan []byte),
	Register:     make(chan *Connection),
	Unregister:   make(chan *Connection),
	RegisterRoom: make(chan *GameRoom),
	//JoinGameRoom :make(chan *Connection),
	GameRooms:   make(map[string]*GameRoom),
	Connections: make(map[*Connection]bool),
}

//游戏大厅
func (h *Hub) run() {
	for {
		select {
		case c := <-h.Register:
			h.Connections[c] = true
		case c := <-h.Unregister:
			if _, ok := h.Connections[c]; ok {
				delete(h.Connections, c)
				close(c.Send)
			}
		case m := <-h.Broadcast:
			for c := range h.Connections {
				select {
				case c.Send <- m:
				default:
					delete(h.Connections, c)
					close(c.Send)
				}
			}
		case room := <-h.RegisterRoom:
			h.GameRooms[room.Id] = room
		}
	}
}
