package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/DemoLiang/wss/golib"
)

// reader pumps messages from the websocket connection to the hub.
func (c *Connection) ReaderHandler() {
	for {
		_, message, err := c.Ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		golib.Log("len:%v\n",len(h.Connections))

		c.HandlerMessage(message)
		//解析消息，此链接已经在游戏大厅，
		// TODO 如果消息是创建，则新建房间
		// TODO 如果是加入房间，则处理加入房间消息
		// TODO 如果是游戏消息，则处理游戏时间，交互游戏厅消息

		// TODO 全球大厅广播消息
		h.Broadcast <- message
	}
	c.Close()
}

// write writes a message with the given message type and payload.
func (c *Connection) WriterHandler() {
	for message := range c.Send {
		err := c.Ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			golib.Log(err.Error())
			break
		}
	}
	c.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 2048, WriteBufferSize: 2048, CheckOrigin: func(r *http.Request) bool {
	return true
}}

func ConnHandler(c *Connection) {
	//先启动写回程序
	go c.WriterHandler()
	//启动读取程序
	go c.ReaderHandler()
}


// wsHandler handles websocket requests from the peer.
func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := &Connection{Send: make(chan []byte, 1024), Ws: ws}
	h.Register <- c
	ConnHandler(c)
}

func (c *Connection)Close()  {
	c.Ws.Close()
	h.Unregister <- c
}

