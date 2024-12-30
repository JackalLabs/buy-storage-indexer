package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var rpc_msg = `{"jsonrpc": "2.0", "method": "%s", "id": 0, "params": {"query": "tm.event = '%s'"}}`

func main() {
	ws, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:26657/websocket", nil)
	defer ws.Close()
	go func() {
		for {
			_, msg, _ := ws.ReadMessage()
			fmt.Println(string(msg))
		}
	}()

	var cmd int
	for {
		fmt.Println("(1) Subscribe (2) Unsubscribe (3) Exit")
		fmt.Scan(&cmd)
		if cmd == 1 {
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(rpc_msg, "subscribe", "NewBlock"))) // use NewBlock to test
		} else if cmd == 2 {
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(rpc_msg, "unsubscribe", "NewBlock")))
		} else if cmd == 3 {
			break
		}
	}
}
