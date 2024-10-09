package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Event string `json:"event"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var msg Message
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Println(err)
			continue
		}

		switch msg.Event {
		case "event1":
			handleEvent1(conn)
		case "event2":
			handleEvent2(conn)
		default:
			message := []byte("Invalid event")
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func handleEvent1(conn *websocket.Conn) {
	time.Sleep(500 * time.Millisecond)
	message := []byte("Handling event 1")
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println(err)
	}
}
func handleEvent2(conn *websocket.Conn) {
	message := []byte("Handling event 2")
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/", handleWebSocket)
	fmt.Println("WebSocket server is running on :5050")
	log.Fatal(http.ListenAndServe(":5050", nil))
}
