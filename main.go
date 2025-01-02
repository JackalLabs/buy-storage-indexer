package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

const (
	rpc      = `{"jsonrpc": "2.0", "method": "%s", "id": 0, "params": {"query": "%s"}}`
	query    = "tm.event = 'Tx' AND message.action = '/canine_chain.storage.MsgBuyStorage'"
	endpoint = "ws://127.0.0.1:26657/websocket"
)

var ws, _, _ = websocket.DefaultDialer.Dial(endpoint, nil)

type Response struct {
	Result struct {
		Events struct {
			Bytes []string `json:"buy_storage.bytes_bought"`
			Hours []string `json:"buy_storage.hours_bought"`
			Buyer []string `json:"buy_storage.buyer"`
			Tx    []string `json:"tx.hash"`
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
			if len(response.Result.Events.Bytes) > 0 {
				bytes, _ := strconv.ParseFloat(response.Result.Events.Bytes[0], 64)
				hours, _ := strconv.ParseFloat(response.Result.Events.Hours[0], 64)
				fmt.Printf(
					"%.2f gb %.2f days - buyer %s in %s\n",
					bytes/(1<<30),
					hours/24,
					response.Result.Events.Buyer[0],
					response.Result.Events.Tx[0],
				)
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
