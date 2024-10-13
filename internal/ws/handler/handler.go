package handler

import (
	"encoding/json"
	"log"
	"log/slog"

	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	log *slog.Logger
}

type messageType struct {
	Event string `json:"event"`
}

func New(log *slog.Logger) *MessageHandler {
	return &MessageHandler{
		log: log,
	}
}

func (m *MessageHandler) Handle(message []byte, conn *websocket.Conn) {
	var msg messageType
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Println(err)
	}

	switch msg.Event {
	case "event1":
		m.getEventFirst(conn)
	default:
		message := []byte("Invalid event")
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err)
			return
		}
	}
}

func (m *MessageHandler) getEventFirst(conn *websocket.Conn) {
	message := []byte("Handling event 1")
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println(err)
	}
}
