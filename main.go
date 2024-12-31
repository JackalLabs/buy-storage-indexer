package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	rpc   = `{"jsonrpc": "2.0", "method": "%s", "id": 0, "params": {"query": "%s"}}`
	query = "tm.event = 'Tx' AND message.action = '/canine_chain.storage.MsgBuyStorage'"
)

var (
	endpoint = url.URL{Scheme: "ws", Host: "127.0.0.1:26657", Path: "/websocket"}
	ws, _, _ = websocket.DefaultDialer.Dial(endpoint.String(), nil)
)

func send(ws *websocket.Conn, method string) {
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(rpc, method, query)))
}

func receive(ws *websocket.Conn) (msg []byte) {
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	return msg
}

func main() {
	defer ws.Close()
	go func() { // receive messages
		for {
			msg := receive(ws)
			if msg == nil {
				break
			}
			fmt.Printf("%s\n", msg)
		}
	}()

	var cmd int
loop:
	for {
		fmt.Println("1: Subscribe 2: Unsubscribe 3: Exit")
		switch fmt.Scan(&cmd); cmd {
		case 1:
			send(ws, "subscribe")
		case 2:
			send(ws, "unsubscribe")
		case 3:
			break loop
		}
	}
}
