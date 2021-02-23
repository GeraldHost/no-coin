package nocoin

import (
	"flag"
)

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
