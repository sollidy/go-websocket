package ws

import (
	"log/slog"
	"net/http"
)

type AppWs struct {
	log *slog.Logger
}

func New(log *slog.Logger) *AppWs {
	return &AppWs{
		log: log,
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

// func (a *AppWs) Close() error {
// 	const op = "ws.Close"
// 	if a.Conn != nil {
// 		err := a.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
// 		if err != nil {
// 			return fmt.Errorf("error sending close message: %s: %w", op, err)
// 		}
// 		err = a.Conn.Close()
// 		if err != nil {
// 			return fmt.Errorf("error closing connection: %s: %w", op, err)
// 		}
// 		a.Conn = nil
// 	}
// 	a.log.With(slog.String("op", op)).Info("DISCONNECTED from websocket server")
// 	return nil
// }
