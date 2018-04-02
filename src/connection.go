package main

import (
	"encoding/json"
	"fmt"
	"github.com/DemoLiang/wss/golib"
	"github.com/gorilla/websocket"
	"net/http"
)

//读取消息进行处理，如果没有登录服务器，用微信code交换服务器code，则关闭链接
func (c *Connection) ReaderHandler() {
	for {
		_, message, err := c.Ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		golib.Log("message:%s len(connections):%v\n", message, len(h.Connections))

		err = c.HandlerMessage(message)
		if err != nil {
			golib.Log("发送信息错误，关闭连接")
			break
		}
	}
	c.Close()
}

//写出消息到客户端，统一以文本形式写出
func (c *Connection) WriterHandler() {
	for message := range c.Send {
		err := c.Ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			golib.Log("write message to client error:%v", err.Error())
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

//HTTP请求升级为websocket请求
func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := &Connection{
		Send:           make(chan []byte, 1024),
		Ws:             ws,
		ConfirDataChan: make(chan bool, 1),
	}
	h.Register <- c
	ConnHandler(c)
}

func (c *Connection) Close() {
	c.Ws.Close()
	h.Unregister <- c
}

//客户端发送消息
func (c *Connection) SendMessage(data interface{}) (err error) {
	dataByte, _ := json.Marshal(&data)
	c.Send <- dataByte
	return nil
}
