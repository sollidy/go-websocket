package ws

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type AppWs struct {
	log  *slog.Logger
	Conn *websocket.Conn
}

func New(log *slog.Logger) *AppWs {
	var conn *websocket.Conn
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := handleWebSocket(w, r)
		if err != nil {
			return
		}
		defer conn.Close()
	})
	return &AppWs{
		log:  log,
		Conn: conn,
	}
}

func (a *AppWs) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *AppWs) Run() error {
	const op = "ws.Run"
	a.log.With((slog.String("op", op))).Info("WebSocket server is running on :5050")
	return http.ListenAndServe(":5050", nil)
}

func (a *AppWs) Close() error {
	const op = "ws.Close"
	if a.Conn != nil {
		err := a.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			return fmt.Errorf("error sending close message: %v", err)
		}
		err = a.Conn.Close()
		if err != nil {
			return fmt.Errorf("error closing connection: %v", err)
		}
		a.Conn = nil
	}
	a.log.With(slog.String("op", op)).Info("DISCONNECTED from websocket server")
	return nil
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
