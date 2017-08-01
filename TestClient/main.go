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

func msgRxCb(msg string, topic string) int {
    fmt.Printf("Message received: %v from topic %v\n", msg, topic)
    return 0
}

func main() {
    starkmq.Init()

    starkmq.Connect()

    starkmq.RegisterRxCallback(msgRxCb)
    starkmq.Subscribe("default")

    reader := bufio.NewReader(os.Stdin)
    for {
        text, err := reader.ReadString('\n')
        if err != nil || isQuitMessage(text) {
            starkmq.Close()
            break
        } else if isSubscribeMessage(text) {
            starkmq.Subscribe("default")
        } else if isUnsubscribeMessage(text) {
            starkmq.Unsubscribe("default")
        } else {
            starkmq.Publish(text, "default")
        }
    }

    starkmq.Close()
}
