package nocoin

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Commands
// generate address
// transfer <amount> <addr>
// deploy <asm_path>
// call <args> <fn:addr> <limit>

func StartConsole(node *Node) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s ", ">")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		mnemonics := strings.Split(text, " ")

		if text == "exit" {
			fmt.Println("closing input channel")
			os.Exit(0)
		} else if text[0] == 'B' {
			node.Broadcast(text)
		} else if text == "my address" {
			myAddress()
		} else if mnemonics[0] == "transfer" {
			txStr := transfer(mnemonics)
			node.Broadcast(fmt.Sprintf("TRANSFER %s", txStr))
		}
	}
}

func myAddress() {
	addr := &Addr{}
	addr.LoadFromFile()
	str := addr.PubKeyToHexStr()
	address := addr.Get()

	fmt.Printf("Public Key :: %s \n", str)
	fmt.Printf("Address :: %s \n", address)
}

// Generate a transfer TX from a string input
// <amount> <recv:addr>
// eg: 20 D80C9BF910F144738EF983724BC04BD6BD3F17C5C83ED57BEDEE1B1B9278E811
func transfer(mnemonics []string) string {
	if len(mnemonics) > 3 || len(mnemonics) < 3 {
		log.Print("wrong number of inputs expect format <amount> <addr>")
		return ""
	}
	amount, err := strconv.Atoi(mnemonics[1])
	if err != nil {
		log.Print("first input must be a number")
		return ""
	}
	addr := mnemonics[2]
	tx := NewTxTransfer(amount, addr)
	return tx.String()
}
