package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
    conn, err := net.Dial("tcp", ":3005")
    if err != nil {
        // handle error
        fmt.Println(err)
        return
    }
    fmt.Println("Connected to server on port 3005.")

    reader := bufio.NewReader(os.Stdin)
    for {
        text, _ := reader.ReadString('\n')
        fmt.Fprintf(conn, text)
        if strings.ToLower(strings.TrimSpace(text)) == "quit" {
            break
        }
    }

    conn.Close()
}
