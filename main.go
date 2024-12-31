package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	rpc      = `{"jsonrpc": "2.0", "method": "%s", "id": 0, "params": {"query": "%s"}}`
	query    = "tm.event = 'Tx' AND message.action = '/canine_chain.storage.MsgBuyStorage'" // "tm.event = 'NewBlock'"
	endpoint = "ws://127.0.0.1:26657/websocket"
)

var ws, _, _ = websocket.DefaultDialer.Dial(endpoint, nil)

type Response struct {
	Result struct {
		Events struct {
			Buy_storage_bytes_bought []string `json:"buy_storage.bytes_bought"`
		} `json:"events"`
	} `json:"result"`
}

func send(ws *websocket.Conn, method string) { // send one msg
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(rpc, method, query)))
}

func receive(ws *websocket.Conn) (msg []byte) { // receive one msg
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	return msg
}

func main() {
	defer ws.Close()
	go func() { // receive loop
		for {
			msg := receive(ws)
			if msg == nil {
				break
			}

			var response Response
			json.Unmarshal(msg, &response)
			if len(response.Result.Events.Buy_storage_bytes_bought) > 0 {
				fmt.Println(response.Result.Events.Buy_storage_bytes_bought[0])
			}
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
