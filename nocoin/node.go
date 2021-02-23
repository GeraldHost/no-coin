package nocoin

import (
  "log"
  "fmt"
  "sync"
  "net/url"
  "net/http"

  "github.com/gorilla/websocket"
)

var seeds []string = []string{"localhost:3001", "localhost:3002"}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type Node struct {
  sync.Mutex
  ID string
  port string
  connections []*Connection
}

// func (node *Node) Broadcast() {}
// remove closed connections
// func (node *Node) Prune() {}

func (node *Node) ConnectToNode(host string) {
  u := url.URL{Scheme: "ws", Host: host, Path: "/"}
  conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
  if err != nil {
    log.Printf("Error connecting to node :: %s", err)
  }
  node.Lock()
  defer node.Unlock()
  connection := &Connection{ conn }
  node.connections = append(node.connections, connection)
}

// Based on the seeds create new connections and
// push them into the connection pool
func (node *Node) DiscoverAndConnect() {
  for _, seed := range seeds {
    node.ConnectToNode(seed)
  }
}

func (node *Node) WelcomeMessage(conn *websocket.Conn) {
  msg := fmt.Sprintf("Welcome :: %s", node.ID)
  conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (node *Node) HandleConn(res http.ResponseWriter, req *http.Request) {
  conn, err := upgrader.Upgrade(res, req, nil)
  if err != nil {
      log.Println(err)
      return
  }

  node.WelcomeMessage(conn)

  for {
    _, msg, err := conn.ReadMessage()
    log.Printf("Recieved: %s", msg)
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
