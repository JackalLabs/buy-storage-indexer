package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var rpc_msg = `{"jsonrpc": "2.0", "method": "%s", "id": 0, "params": {"query": "tm.event = 'Tx' AND message.action = '/canine_chain.storage.MsgBuyStorage'"}}`

func main() {
	ws, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:26657/websocket", nil)
	defer ws.Close()
	go func() { // display new messages
		for {
			_, msg, _ := ws.ReadMessage()
			println(string(msg))
		}
	}()

	var cmd int
	for {
		fmt.Println("(1) Subscribe (2) Unsubscribe (3) Exit")
		fmt.Scan(&cmd)
		if cmd == 1 {
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(rpc_msg, "subscribe")))
		} else if cmd == 2 {
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(rpc_msg, "unsubscribe")))
		} else if cmd == 3 {
			break
		}
	}
}
