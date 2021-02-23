package nocoin

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartConsole(node *Node) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s ", ">")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Println("Closing input channel")
			os.Exit(0)
		} else if text[0] == 'B' {
			node.Broadcast(text)
		}

		fmt.Printf("Input :: %s \n", text)
	}
}
