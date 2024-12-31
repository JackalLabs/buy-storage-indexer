package main

import (
	"testing"
)

var blank = `{"jsonrpc":"2.0","id":0,"result":{}}`

func TestSubscribe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		send(ws, "subscribe")
		msg, _ := receive(ws)
		if string(msg) != blank {
			t.Fatalf("Unexpected Output %s\n", msg)
		}
	})
}
