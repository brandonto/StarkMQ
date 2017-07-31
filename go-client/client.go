package starkmq

import (
    "encoding/json"
    "fmt"
    "net"

    "github.com/brandonto/StarkMQ/common"
)

type RxCallbackFunc func(msg string) int

func defaultRxCallback(msg string) int {
    fmt.Println("Received message but no callback was registered")
    return 0
}

type starkMQClient struct {
    conn net.Conn
    cb RxCallbackFunc
    encoder *json.Encoder
}

var client starkMQClient

func (cl *starkMQClient) send(msg common.StarkMQMsg) {
    err := cl.encoder.Encode(msg)
    if err != nil {
        fmt.Println("unable to send message")
        return
    }
}

func Init() {
    client.cb = defaultRxCallback
}

func RegisterRxCallback(cb RxCallbackFunc) {
    client.cb = cb
}

func Connect() {
    var err error
    client.conn, err = net.Dial("tcp", ":3005")
    if err != nil {
        // handle error
        fmt.Println(err)
        return
    }

    client.encoder = json.NewEncoder(client.conn)

    go listen()

    fmt.Println("Connected to server on port 3005.")
}

func listen() {
    decoder := json.NewDecoder(client.conn)
    for {
        var msg common.StarkMQMsg
        err := decoder.Decode(&msg)

        // Try to gracefully handle disconnection
        if err != nil {
            break
        }

        //fmt.Println(msg.String())
        switch msg.MsgType {
        case common.PUBLISH:
            client.cb(msg.Payload)
        default:
            fmt.Println("Error unsupported message type")
            return
        }
    }
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
        client.send(msg)
        //fmt.Printf("Message sent: %v\n", msg)
    default:
        // handle error
        fmt.Println("unable to serialize message: message type unsupported")
        return
    }
}

func Close() {
    msg := common.NewStarkMQMsg(common.QUIT, "quit\n")
    send(msg)
    client.conn.Close()
}
