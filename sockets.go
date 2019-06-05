package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func readSocket(conn *websocket.Conn) error {
	for {
		msgtype, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		fmt.Printf("Socket message: %s\n", msg)

		if err = conn.WriteMessage(msgtype, msg); err != nil {
			return err
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleSocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true } // allow all requests

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error:", err)
		}
		readSocket(ws)
	}
}
