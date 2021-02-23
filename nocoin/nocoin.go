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
  node := Node { ID: "JAKE", port: config.port }
  node.DiscoverAndConnect()
  go StartConsole()
  node.Serve()
}
