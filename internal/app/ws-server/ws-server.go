package ws

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

type Server struct {
	clients       map[*websocket.Conn]bool
	handleMessage func(message []byte, connection *websocket.Conn) // хандлер новых сообщений
	log           *slog.Logger
}

func StartServer(handleMessage func(message []byte, connection *websocket.Conn), log *slog.Logger) *Server {
	const op = "ws.StartServer"
	server := Server{
		make(map[*websocket.Conn]bool),
		handleMessage,
		log,
	}

	http.HandleFunc("/", server.echo)
	go func() {
		err := http.ListenAndServe(":5050", nil)
		if err != nil {
			log.With(slog.String("op", op)).Error("Failed to start WebSocket server", "error", err)
			return
		}
	}()
	log.With(slog.String("op", op)).Info("WebSocket server is running on :5050")

	return &server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	server.clients[connection] = true        // Сохраняем соединение, используя его как ключ
	defer delete(server.clients, connection) // Удаляем соединение

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
		}

		go server.handleMessage(message, connection)
	}
}

func (server *Server) WriteMessage(message []byte) {
	for conn := range server.clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			server.log.Error("Failed to write message", "error", err)
		}
	}
}
