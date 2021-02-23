package nocoin

import "github.com/gorilla/websocket"

// Represents a TCP/IP connection
type Connection struct {
  conn *websocket.Conn
}

func (conn *Connection) Send() {}

func (conn *Connection) Stop() {}

