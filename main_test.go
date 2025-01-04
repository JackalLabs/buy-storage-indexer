package main

import (
	"testing"

	"github.com/gorilla/websocket"
)

var (
	blank = `{"jsonrpc":"2.0","id":0,"result":{}}`
	url   = "ws://127.0.0.1:26657/websocket"
)

func TestSubscribe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ws, _, _ := websocket.DefaultDialer.Dial(url, nil)
		send(ws, "subscribe")
		msg := receive(ws)
		if string(msg) != blank {
			t.Fatalf("Unexpected Output %s\n", msg)
		}
	})
}
