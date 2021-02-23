package nocoin

import (
  "fmt"
  "bufio"
  "os"
  "strings"
)

func StartConsole() {
  reader := bufio.NewReader(os.Stdin);
  for {
    fmt.Printf("%s ", ">")
    text, _ := reader.ReadString('\n')
    if strings.TrimSpace(text) == "exit" {
      fmt.Println("Closing input channel")
      return
    }
    fmt.Printf("Input :: %s", text)
  }
}
