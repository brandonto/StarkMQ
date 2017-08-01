package common

import (
    "encoding/json"
    "fmt"
)

type StarkMQMsgType int

const (
    SUBSCRIBE StarkMQMsgType = iota
    UNSUBSCRIBE
    PUBLISH
    QUIT
)

type StarkMQPayload []byte

type StarkMQMsg struct {
    MsgType StarkMQMsgType `json:"msgtype"`
    Payload StarkMQPayload `json:"payload"`
}

func NewStarkMQMsg(msgType StarkMQMsgType, payload StarkMQPayload) StarkMQMsg {
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

type StarkMQSubscribePayload struct {
    Topic string `json:"topic"`
}

func NewStarkMQSubscribePayload(topic string) StarkMQSubscribePayload {
    return StarkMQSubscribePayload{Topic: topic}
}

type StarkMQUnsubscribePayload struct {
    Topic string `json:"topic"`
}

func NewStarkMQUnsubscribePayload(topic string) StarkMQUnsubscribePayload {
    return StarkMQUnsubscribePayload{Topic: topic}
}

type StarkMQPublishPayload struct {
    Topic string `json:"topic"`
    Text string `json:"text"`
}

func NewStarkMQPublishPayload(topic string, text string) StarkMQPublishPayload {
    return StarkMQPublishPayload{Topic: topic, Text: text}
}

