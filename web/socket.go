package web

import (
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ReadSocket is function that will not stick around so don't use it.
func ReadSocket(conn *websocket.Conn) (chan []byte, chan error) {
	ch := make(chan []byte)
	errch := make(chan error)
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				errch <- err
			}
			ch <- msg
		}
	}()
	return ch, errch
}
