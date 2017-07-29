package starkmq

import (
    "fmt"
    "net"
)

type starkMQClient struct {
    conn net.Conn
}

type StarkMQMsg struct {
    MsgType StarkMQMsgType
    Payload string
}

type StarkMQMsgType int

const (
    SUBSCRIBE StarkMQMsgType = iota
    UNSUBSCRIBE
    PUBLISH
    QUIT
)

var client starkMQClient

func Connect() {
    var err error
    client.conn, err = net.Dial("tcp", ":3005")
    if err != nil {
        // handle error
        fmt.Println(err)
        return
    }
    fmt.Println("Connected to server on port 3005.")
}

func Subscribe() {
    msg := StarkMQMsg{MsgType: SUBSCRIBE, Payload: "subscribe\n"}
    send(msg)
}

func Unsubscribe() {
    msg := StarkMQMsg{MsgType: UNSUBSCRIBE, Payload: "unsubscribe\n"}
    send(msg)
}

func Publish(text string) {
    msg := StarkMQMsg{MsgType: PUBLISH, Payload: text}
    send(msg)
}

func send(msg StarkMQMsg) {
    switch msg.MsgType {
    case SUBSCRIBE:
        fallthrough
    case UNSUBSCRIBE:
        fallthrough
    case PUBLISH:
        fallthrough
    case QUIT:
        fmt.Fprintf(client.conn, msg.Payload)
    default:
        // handle error
        fmt.Println("message type unsupported")
        return
    }
}

func Close() {
    msg := StarkMQMsg{MsgType: QUIT, Payload: "quit\n"}
    send(msg)
    client.conn.Close()
}
