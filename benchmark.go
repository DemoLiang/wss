package main

//import (
//	"log"
//
//	"fmt"
//	"github.com/rgamba/evtwebsocket"
//	"golang.org/x/net/websocket"
//)
//
//func main() {
//	for i := 0; i < 1000; i++ {
//		go func() {
//			c := evtwebsocket.Conn{
//
//				// When connection is established
//				OnConnected: func(w *websocket.Conn) {
//					fmt.Println("Connected")
//				},
//
//				// When a message arrives
//				OnMessage: func(msg []byte) {
//					log.Printf("Received uncatched message: %s\n", msg)
//				},
//
//				// When the client disconnects for any reason
//				OnError: func(err error) {
//					fmt.Printf("** ERROR **\n%s\n", err.Error())
//				},
//
//				// This is used to match the request and response messages
//				MatchMsg: func(req, resp []byte) bool {
//					return string(req) == string(rest)
//				},
//			}
//
//			// Connect
//			if err := c.Dial("ws://<server_ip:server_port>/ws"); err != nil {
//				log.Fatal(err)
//			}
//
//			// Create the message with a callback
//			msg := evtwebsocket.Msg{
//				Body: nil,
//				Callback: func(resp []byte) {
//					fmt.Printf("Got back: %s\n", resp)
//				},
//			}
//
//			log.Printf("%s\n", msg.Body)
//		}()
//	}
//	select {}
//}
