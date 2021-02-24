package nocoin

import (
	"flag"
)

// Global state
var txPool map[string]*Tx
var seeds []string = []string{"localhost:3001", "localhost:3002"}

type config struct {
	port string
}

func parseFlags() *config {
	conf := &config{}
	flag.StringVar(&conf.port, "p", "5000", "Port eg. 80")
	flag.Parse()
	return conf
}

func Start() {
	config := parseFlags()
	node := Node{Id: "JAKE", port: config.port, outbound_connections: make(map[string]*Connection)}
	node.DiscoverAndConnect()
	go StartConsole(&node)
	node.Serve()
}
