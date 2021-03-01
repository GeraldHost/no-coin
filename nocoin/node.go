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

func NewNode(id string, port string) Node {
	node := Node{Id: id, port: port}
	node.outbound_connections = make(map[string]*Connection)
	node.inbound_connections = make(map[string]*Connection)
	return node
}

// Get outbound and inbound connections to nodes
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

// Connect over websockets to a single node
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

// When we create a connection we send a welcome message
// to the connecting node with our ID
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
	// add connection to our inbound connection pool
	addr := conn.RemoteAddr().String()
	node.inbound_connections[addr] = &Connection{conn}
	node.Unlock()

	for {
		// Recieve a message from a node.
		// TODO: implement a way of effeciently parsing these recieved messages
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
