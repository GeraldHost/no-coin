package nocoin

import "github.com/gorilla/websocket"

// Represents a TCP/IP connection
type Connection struct {
	conn *websocket.Conn
}

func (c *Connection) Send(msg string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (conn *Connection) Stop() {}
