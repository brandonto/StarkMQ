package starkmq

import (
    "bufio"
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
}

var client starkMQClient

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

    go listen()

    fmt.Println("Connected to server on port 3005.")
}

func listen() {
    reader:= bufio.NewReader(client.conn)
    for {
        text, err := reader.ReadString('\n')

        // Try to gracefully handle disconnection
        if err != nil {
            break
        }

        msg := common.Deserialize(text)

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
        fmt.Fprintf(client.conn, common.Serialize(msg))
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
