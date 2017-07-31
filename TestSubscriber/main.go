package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/brandonto/StarkMQ/go-client"
)

func isQuitMessage(msg string) bool {
    return strings.ToLower(strings.TrimSpace(msg)) == "quit"
}

func isSubscribeMessage(msg string) bool {
    return strings.ToLower(strings.TrimSpace(msg)) == "subscribe"
}

func isUnsubscribeMessage(msg string) bool {
    return strings.ToLower(strings.TrimSpace(msg)) == "unsubscribe"
}

func msgRxCb(msg string) int {
    fmt.Printf("Message received: %v\n", msg)
    return 0
}

func main() {
    starkmq.Init()

    starkmq.Connect()

    starkmq.RegisterRxCallback(msgRxCb)
    starkmq.Subscribe()

    reader := bufio.NewReader(os.Stdin)
    for {
        text, err := reader.ReadString('\n')
        if err != nil || isQuitMessage(text) {
            starkmq.Close()
            break
        } else if isSubscribeMessage(text) {
            starkmq.Subscribe()
        } else if isUnsubscribeMessage(text) {
            starkmq.Unsubscribe()
        } else {
            starkmq.Publish(text)
        }
    }

    starkmq.Close()
}
