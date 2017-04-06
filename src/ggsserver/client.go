package main

import (
	"encoding/binary"
	"log"
	"os"
	"sync"
	"time"

	"ggs/conf"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/kr/pretty"

	"ggsserver/msg"
)

var (
	conn *websocket.Conn
)

func init() {
	_conn, _, err := websocket.DefaultDialer.Dial("ws://"+conf.Env.WSAddr, nil)
	conn = _conn
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//massConnect(10)
	connect()
}

func connect() {
	getReceive := readMsg()

	//sendHelloRequest()
	sendLoginRequest()

	<-getReceive
	return
}

func massConnect(n int) {
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go connectServer(wg)
	}
	wg.Wait()
}

func connectServer(wg *sync.WaitGroup) {
	getReceive := readMsg()

	<-getReceive
	wg.Done()
}

func msgMarshal(id int, data []byte) (msg []byte) {
	msg = make([]byte, 2+len(data))
	binary.BigEndian.PutUint16(msg, uint16(id))

	copy(msg[2:], data)
	return
}

func readMsg() (getReceive chan bool) {
	getReceive = make(chan bool)

	go func() {
		for {
			_, msgStr, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
			}

			msgData := getMsg(msgStr)
			pretty.Println("receive msg:", msgData, "\n")
		}
		getReceive <- true
	}()

	return
}

func getMsg(data []byte) interface{} {
	m, err := msg.Processor.Unmarshal(data)
	if err != nil {
		log.Fatal()
	}
	return m
}

func sendMsg(pb proto.Message, id int) {
	data, err := proto.Marshal(pb)
	if err != nil {
		log.Fatal(err)
	}

	msg := msgMarshal(id, data)
	err = conn.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func sendHelloRequest() {
	sendMsg(&msg.Hello{
		Name: proto.String("hello"),
	}, 0)
}

func sendLoginRequest() {
	token := os.Args[3]
	sendMsg(&msg.LoginRequest{
		Token:    proto.String(token),
		DeviceId: proto.Uint32(uint32(time.Now().Unix())),
		Version:  proto.Uint32(1),
	}, 24)
}
