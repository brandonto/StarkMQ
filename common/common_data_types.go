package common

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

func NewStarkMQMsg(msgType StarkMQMsgType, payload string) StarkMQMsg {
    return StarkMQMsg{MsgType: msgType, Payload: payload}
}

func Serialize(msg StarkMQMsg) string {
    return msg.Payload
}

func Deserialize(rawMsg string) StarkMQMsg {
    var msg StarkMQMsg
    switch rawMsg {
    case "subscribe":
        msg.MsgType = SUBSCRIBE
    case "unsubscribe":
        msg.MsgType = UNSUBSCRIBE
    case "publish":
        msg.MsgType = PUBLISH
    case "quit":
        msg.MsgType = QUIT
    }
    msg.Payload = rawMsg
    return msg
}
