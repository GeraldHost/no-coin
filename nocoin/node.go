package nocoin

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Node struct {
	sync.Mutex
	Id                   string
	port                 string
	outbound_connections map[string]*Connection
	inbound_connections  map[string]*Connection
}

func (node *Node) Connections() map[string]*Connection {
	connections := node.inbound_connections
	for addr, connection := range node.outbound_connections {
		connections[addr] = connection
	}
	return connections
}

func (node *Node) Broadcast(msg string) {
	for _, connection := range node.Connections() {
		connection.Send(msg)
	}
}

// remove closed connections
// func (node *Node) Prune() {}

func (node *Node) ConnectToNode(host string) {
	u := url.URL{Scheme: "ws", Host: host, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("Error connecting to node :: %s", err)
		return
	}
	node.Lock()
	defer node.Unlock()
	connection := &Connection{conn}
	addr := conn.RemoteAddr().String()
	node.outbound_connections[addr] = connection
}

// Based on the seeds create new connections and
// push them into the connection pool
func (node *Node) DiscoverAndConnect() {
	for _, seed := range seeds {
		node.ConnectToNode(seed)
	}
}

func (node *Node) WelcomeMessage(conn *websocket.Conn) {
	msg := fmt.Sprintf("Welcome :: %s", node.Id)
	conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (node *Node) HandleConn(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	node.WelcomeMessage(conn)

	node.Lock()
	addr := conn.RemoteAddr().String()
	node.inbound_connections[addr] = &Connection{conn}
	node.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		log.Printf("%sRecieved: %s", TERMINAL_CLEAR_LINE, msg)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (node *Node) Serve() {
	http.HandleFunc("/", node.HandleConn)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", node.port), nil))
}
