package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "time"
)

type MessageQueue struct {
    subscribers []*Subscriber
}

func (mq *MessageQueue) AddSubscriber(sub *Subscriber) {
    fmt.Printf("Subscriber added!\n")
    mq.subscribers = append(mq.subscribers, sub)
}

func (mq *MessageQueue) Publish(msg string) {
    for _, sub := range mq.subscribers {
        sub.OnReceive(msg)
    }
}

type Subscriber struct {
    id int
}

func (sub *Subscriber) String() string {
    return fmt.Sprintf("Subscriber{%v}", sub.id)
}

func (sub *Subscriber) OnReceive(msg string) {
    fmt.Printf("%v received: %s\n", sub, msg)
}

func connectionLoop(newConnectionChan chan net.Conn) {
    var connections []net.Conn
    connectionChan := make(chan string, 128)
    for {
        select {
        case conn := <-newConnectionChan:
            connections = append(connections, conn)
            go connectionHandler(conn, connectionChan)
            fmt.Println("omg a connection")
        case msg := <-connectionChan:
            fmt.Println(msg)
        default:
            for i := range connections {
                fmt.Println(i)
            }
            time.Sleep(1 * time.Second)
        }
    }
}

func connectionHandler(conn net.Conn, connectionChan chan string) {
    reader:= bufio.NewReader(conn)
    for {
        text, _ := reader.ReadString('\n')
        connectionChan <- text
        //fmt.Println(text)
        if strings.ToLower(strings.TrimSpace(text)) == "quit" {
            break
        }
    }
}

func isQuitMessage() {
    // TODO
}

func main() {
    ln, err := net.Listen("tcp", ":3005")
    if err != nil {
        // handle error
        fmt.Println(err)
        return
    }
    fmt.Println("Started server listening on port 3005.")

    newConnectionChan := make(chan net.Conn, 5)
    go connectionLoop(newConnectionChan)

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
            return
        }

        newConnectionChan <- conn
    }
}
