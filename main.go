package main

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

var (
	rpc      = `{"jsonrpc": "2.0", "method": "%s", "id": 0, "params": {"query": "%s"}}`
	endpoint = url.URL{Scheme: "ws", Host: "127.0.0.1:26657", Path: "/websocket"}
	query    = "tm.event = 'Tx' AND message.action = '/canine_chain.storage.MsgBuyStorage'"
	ws, _, _ = websocket.DefaultDialer.Dial(endpoint.String(), nil)
)

func send(ws *websocket.Conn, method string) {
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(rpc, method, query)))
}

func receive(ws *websocket.Conn) (msg []byte, err error) {
	_, msg, err = ws.ReadMessage()
	return msg, err
}

func main() {
	defer ws.Close()
	go func() { // receive messages
		for {
			msg, err := receive(ws)
			if err != nil { // handle exit error
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
