package main

import (
    //"bufio"
    "encoding/json"
    "errors"
    "fmt"
    "net"
    "time"

    "github.com/brandonto/StarkMQ/common"
)

type Connection struct {
    id int
    conn net.Conn
    topicsSubscribedTo []string
    encoder *json.Encoder
}

func NewConnection(id int, conn net.Conn) *Connection {
    connection := Connection {
        id: id,
        conn: conn,
        topicsSubscribedTo: nil,
        encoder: json.NewEncoder(conn),
    }
    return &connection
}

func (conn *Connection) isSubscribedToTopic(topic string) bool {
    for _, t := range conn.topicsSubscribedTo {
        if t == topic {
            return true
        }
    }
    return false
}

func (conn *Connection) addTopicToList(topic string) {
    conn.topicsSubscribedTo = append(conn.topicsSubscribedTo, topic)
}

func (conn *Connection) removeTopicFromList(topic string) {
    index, err := conn.getIndexOfTopicInList(topic)
    if err != nil {
        return
    }

    conn.topicsSubscribedTo = append(conn.topicsSubscribedTo[:index], conn.topicsSubscribedTo[index+1:]...)
}

func (conn *Connection) getIndexOfTopicInList(topic string) (int, error) {
    for index, t := range conn.topicsSubscribedTo {
        if t == topic {
            return index, nil
        }
    }

    return -1, errors.New("index not found")
}

func (conn *Connection) send(msg common.StarkMQMsg) {
    err := conn.encoder.Encode(msg)
    if err != nil {
        fmt.Println("unable to send message")
        return
    }
}

type ConnectionLoop struct {
    connections map[int]*Connection
    topicSubscriptions map[string][]int
    newConnectionChan chan net.Conn
    nextConnectionID int
}

func (cl *ConnectionLoop) addConnection(conn net.Conn) *Connection {
    id := cl.nextConnectionID
    cl.connections[id] = NewConnection(id, conn)
    cl.nextConnectionID++
    fmt.Printf("Connection with id %v has been added\n", id)
    return cl.connections[id]
}

func (cl *ConnectionLoop) removeConnection(id int) {
    // Remove connection from all topics it is subscribed to
    for _, topic := range cl.connections[id].topicsSubscribedTo {
        cl.removeSubscriptionFromTopic(id, topic)
    }

    // Remove connection from connection map
    delete(cl.connections, id)
    fmt.Printf("Connection with id %v has been removed\n", id)
}

func (cl *ConnectionLoop) addSubscriptionToTopic(id int, topic string) bool {
    // Topic exists in connection topic list, do nothing
    if cl.connections[id].isSubscribedToTopic(topic) {
        fmt.Printf("Subscribe failed: Connection with id %v is already subscribed to topic %v\n", id, topic)
        return false
    }

    _, err := cl.getIndexOfSubscriptionInTopic(id, topic)

    // Already exist in topic, do nothing
    if err == nil {
        fmt.Printf("Synchronization issue: Connection with id %v exists in topic %v but doesn't know it is\n", id, topic)
        return false
    }

    cl.connections[id].addTopicToList(topic)
    cl.topicSubscriptions[topic] = append(cl.topicSubscriptions[topic], id)
    fmt.Printf("Connection with id %v has been subscribed to %v\n", id, topic)
    return true
}

func (cl *ConnectionLoop) removeSubscriptionFromTopic(id int, topic string) bool {
    // Topic doesn't exists in connection topic list, do nothing
    if !cl.connections[id].isSubscribedToTopic(topic) {
        fmt.Printf("Unsubscribe failed: connection with id %v is not subscribed to topic %v\n", id, topic)
        return false
    }

    index, err := cl.getIndexOfSubscriptionInTopic(id, topic)

    // Doesn't exist in topic, do nothing
    if err != nil {
        return false
    }

    cl.connections[id].removeTopicFromList(topic)
    cl.topicSubscriptions[topic] = append(cl.topicSubscriptions[topic][:index], cl.topicSubscriptions[topic][index+1:]...)
    fmt.Printf("Connection with id %v has been removed from topic %v\n", id, topic)
    return true
}

func (cl *ConnectionLoop) getIndexOfSubscriptionInTopic(id int, topic string) (int, error) {
    for index, connId := range cl.topicSubscriptions[topic] {
        if id == connId {
            return index, nil
        }
    }

    return -1, errors.New("index not found")
}

func (cl *ConnectionLoop) publishMessageToTopic(msg common.StarkMQMsg, topic string) {
    fmt.Printf("%v is being published to all subscribers of %v\n", msg.Payload, topic)
    for _, id := range cl.topicSubscriptions[topic] {
        fmt.Printf("Publishing to %v\n", id)
        cl.connections[id].send(msg)
    }
}

func (cl *ConnectionLoop) handleMsg(connMsg ConnectionMsg) {
    // Try to gracefully handle disconnection
    if connMsg.disconnected {
        cl.removeConnection(connMsg.connection.id)
        return
    }

    switch connMsg.msg.MsgType {
    case common.SUBSCRIBE:
        cl.addSubscriptionToTopic(connMsg.connection.id, "default")
    case common.UNSUBSCRIBE:
        cl.removeSubscriptionFromTopic(connMsg.connection.id, "default")
    case common.PUBLISH:
        cl.publishMessageToTopic(connMsg.msg, "default")
    case common.QUIT:
        cl.removeConnection(connMsg.connection.id)
    default:
        fmt.Println("unsupported message type")
    }
}

func (cl *ConnectionLoop) Exec() {
    connectionChan := make(chan ConnectionMsg, 128)
    for {
        select {
        case conn := <-cl.newConnectionChan:
            connection := cl.addConnection(conn)
            go connectionHandler(connection, connectionChan)
        case connMsg := <-connectionChan:
            cl.handleMsg(connMsg)
        default:
            //for i := range cl.connections {
            //    fmt.Println(i)
            //}
            time.Sleep(100 * time.Millisecond)
        }
    }
}

type ConnectionMsg struct {
    connection *Connection
    msg common.StarkMQMsg
    disconnected bool
}

func connectionHandler(connection *Connection, connectionChan chan ConnectionMsg) {
    decoder := json.NewDecoder(connection.conn)
    for {
        var msg common.StarkMQMsg
        err := decoder.Decode(&msg)

        // Try to gracefully handle disconnection
        if err != nil {
            connectionChan <- ConnectionMsg{connection: connection, disconnected: true}
            break
        }

        fmt.Printf("%v\n", msg.String())
        connectionChan <- ConnectionMsg{connection: connection, msg: msg, disconnected: false}
        if msg.MsgType == common.QUIT {
            break
        }
    }
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
    cl := ConnectionLoop {
        connections: make(map[int]*Connection),
        topicSubscriptions: make(map[string][]int),
        newConnectionChan: newConnectionChan,
        nextConnectionID: 0,
    }
    go cl.Exec()

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
            break
        }

        newConnectionChan <- conn
    }

    ln.Close()
}
