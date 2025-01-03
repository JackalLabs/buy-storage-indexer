// buy-storage-indexer monitors the Jackal Protocol for buy storage events through WebSocket
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	rpc   = `{"jsonrpc": "2.0", "method": "%s", "id": 0, "params": {"query": "%s"}}`
	query = "tm.event = 'Tx' AND message.action = '/canine_chain.storage.MsgBuyStorage'"
)

var ws *websocket.Conn

type Response struct {
	Result struct {
		Events struct {
			Buyer    []string `json:"buy_storage.buyer"`
			Receiver []string `json:"buy_storage.receiver"`
			Bytes    []string `json:"buy_storage.bytes_bought"`
			Hours    []string `json:"buy_storage.hours_bought"`
			Tx       []string `json:"tx.hash"`
		} `json:"events"`
	} `json:"result"`
}

func send(ws *websocket.Conn, c string) { // send one command
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(rpc, c, query)))
}

func receive(ws *websocket.Conn) (m []byte) { // receive one message
	_, m, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	return m
}

func init() {
	if len(os.Args) < 2 {
		log.Fatal("./buy-storage-indexer [rpc ip:port]")
	}

	f, err := os.OpenFile("indexer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	url := "ws://" + os.Args[1] + "/websocket"
	ws, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer ws.Close()

	go func() { // receive loop
		for {
			m := receive(ws)
			if m == nil {
				break
			}
			log.Println(string(m))

			var response Response
			json.Unmarshal(m, &response)
			if len(response.Result.Events.Bytes) > 0 {
				b, _ := strconv.ParseFloat(response.Result.Events.Bytes[0], 64)
				h, _ := strconv.ParseFloat(response.Result.Events.Hours[0], 64)
				fmt.Printf(
					"%s | %.2f gb %.2f days  %s buyer %s receiver %s tx\n",
					time.Now().Format("2006-01-02 15:04:05"),
					b/(1<<30), h/24,
					response.Result.Events.Buyer[0],
					response.Result.Events.Receiver[0],
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
