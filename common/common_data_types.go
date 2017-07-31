package common

import (
    "encoding/json"
    "fmt"
)

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

func (msg *StarkMQMsg) String() string {
    data, err := json.Marshal(msg)
    if err != nil {
        fmt.Println("unable to convert message to string")
        return ""
    }

    return string(data)
}
