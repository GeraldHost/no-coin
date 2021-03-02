package nocoin

import (
	"flag"
	"fmt"
)

//////////////////////////////////////////
// Global state
//////////////////////////////////////////

// total number of credits in the network
var marketCap int = 9223372036854775807

// Unspent Output Pool map of publicKeyStr => hex string (vout)
var utxoPool map[string][]*Utxo = make(map[string][]*Utxo, 0)

// Transaction Pool map of TxHash => Tx
var txPool map[string]*Tx = make(map[string]*Tx)

// Hard coded seeds for the node to make initial connection with
var seeds []string = []string{"localhost:3001", "localhost:3002"}

// tracks the latest block height as we sync
var latestBlockHeight int = 0

// My address
var myAddr *Addr

//////////////////////////////////////////
// App
//////////////////////////////////////////

type config struct {
	port string
}

func parseFlags() *config {
	conf := &config{}
	flag.StringVar(&conf.port, "p", "5000", "Port eg. 80")
	flag.Parse()
	return conf
}

func SetupAddr() {
	fmt.Println("setting up address")
	addr := &Addr{}
	ok := addr.LoadFromFile()
	if !ok {
		fmt.Println("no address found. Generating new address")
		addr.generate()
	}
	myAddr = addr
}

// When you first run nocoin we push the vendor uxto into the uxto pool
// this is the vendor that owns all the market cap of credits in order to
// transfer value between miners and users
func SetupVendor() {
	vendorAddr := "4d265138333dfdfa3b22454fd654e581052688b8a3592dd4306e1426f4bbc6ed"
	utxo := &Utxo { addr: vendorAddr, amount: marketCap }
	AddToUtxoPool(utxo)
}

func Start() {
	config := parseFlags()
	node := NewNode("JAKE", config.port)
	node.DiscoverAndConnect()
	SetupAddr()
	SetupVendor()
	go StartConsole(&node)
	node.Serve()
}
