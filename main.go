package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "time"
)

//type MessageQueue struct {
//    subscribers []*Subscriber
//}
//
//func (mq *MessageQueue) AddSubscriber(sub *Subscriber) {
//    fmt.Printf("Subscriber added!\n")
//    mq.subscribers = append(mq.subscribers, sub)
//}
//
//func (mq *MessageQueue) Publish(msg string) {
//    for _, sub := range mq.subscribers {
//        sub.OnReceive(msg)
//    }
//}
//
//type Subscriber struct {
//    id int
//}
//
//func (sub *Subscriber) String() string {
//    return fmt.Sprintf("Subscriber{%v}", sub.id)
//}
//
//func (sub *Subscriber) OnReceive(msg string) {
//    fmt.Printf("%v received: %s\n", sub, msg)
//}


type Connection struct {
    id int
    conn net.Conn
}

type ConnectionLoop struct {
    connections map[int]Connection
    newConnectionChan chan net.Conn
    nextConnectionID int
}

func (cl *ConnectionLoop) AddConnection(conn net.Conn) Connection {
    id := cl.nextConnectionID
    cl.connections[id] = Connection{id, conn}
    cl.nextConnectionID++
    return cl.connections[id]
}

func (cl *ConnectionLoop) DeleteConnection(connection Connection) {
    delete(cl.connections, connection.id)
}

func (cl *ConnectionLoop) handleMsg(msg ) {
    fmt.Println(connMsg.msg)
    if isQuitMessage(connMsg.msg) {
        cl.DeleteConnection(connMsg.connection)
    }
}

func (cl *ConnectionLoop) Exec() {
    connectionChan := make(chan ConnectionMsg, 128)
    for {
        select {
        case conn := <-cl.newConnectionChan:
            connection := cl.AddConnection(conn)
            go connectionHandler(connection, connectionChan)
            fmt.Println("omg a connection")
        case connMsg := <-connectionChan:
            cl.handleMsg(connMsg)
        default:
            for i := range cl.connections {
                fmt.Println(i)
            }
            time.Sleep(1 * time.Second)
        }
    }
}

type ConnectionMsg struct {
    connection Connection
    msg string
}

func connectionHandler(connection Connection, connectionChan chan ConnectionMsg) {
    reader:= bufio.NewReader(connection.conn)
    for {
        text, _ := reader.ReadString('\n')
        connectionChan <- ConnectionMsg{connection, text}
        //fmt.Println(text)
        if isQuitMessage(text) {
            break
        }
    }
}

func isQuitMessage(msg string) bool {
    return strings.ToLower(strings.TrimSpace(msg)) == "quit"
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
        connections: make(map[int]Connection),
        newConnectionChan: newConnectionChan,
        nextConnectionID: 0,
    }
    go cl.Exec()

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
            return
        }

        newConnectionChan <- conn
    }
}
