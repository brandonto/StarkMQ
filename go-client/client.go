package starkmq

import (
    "fmt"
    "net"

    "github.com/brandonto/StarkMQ/common"
)

type starkMQClient struct {
    conn net.Conn
}

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
    msg := common.NewStarkMQMsg(common.SUBSCRIBE, "subscribe\n")
    send(msg)
}

func Unsubscribe() {
    msg := common.NewStarkMQMsg(common.UNSUBSCRIBE, "unsubscribe\n")
    send(msg)
}

func Publish(text string) {
    msg := common.NewStarkMQMsg(common.PUBLISH, text)
    send(msg)
}

func send(msg common.StarkMQMsg) {
    switch msg.MsgType {
    case common.SUBSCRIBE:
        fallthrough
    case common.UNSUBSCRIBE:
        fallthrough
    case common.PUBLISH:
        fallthrough
    case common.QUIT:
        fmt.Fprintf(client.conn, common.Serialize(msg))
    default:
        // handle error
        fmt.Println("message type unsupported")
        return
    }
}

func Close() {
    msg := common.NewStarkMQMsg(common.QUIT, "quit\n")
    send(msg)
    client.conn.Close()
}
