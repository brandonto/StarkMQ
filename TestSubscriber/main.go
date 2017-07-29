package main

import (
    "bufio"
    //"fmt"
    "os"
    "strings"

    "github.com/brandonto/StarkMQ/go-client"
)

func isQuitMessage(msg string) bool {
    return strings.ToLower(strings.TrimSpace(msg)) == "quit"
}

func main() {
    starkmq.Connect()
    starkmq.Subscribe()

    reader := bufio.NewReader(os.Stdin)
    for {
        text, _ := reader.ReadString('\n')
        if isQuitMessage(text) {
            starkmq.Close()
            break
        }
        starkmq.Publish(text)
    }

    starkmq.Close()
}
