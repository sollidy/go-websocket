package app

import (
	"context"
	"go-ws/internal/app/ws-server"
	"go-ws/internal/lib/logger/sl"
	"go-ws/internal/storage"
	"log/slog"

	"github.com/gorilla/websocket"
)

type App struct {
	Storage *storage.Postgres
	Ws      *ws.Server
}

func New(
	ctx context.Context,
	log *slog.Logger,
	pgConnString string,
) *App {

	app := &App{}

	storage, err := storage.NewPG(ctx, pgConnString, log)
	if err != nil {
		panic(err)
	}
	// ping database
	err = storage.Ping(ctx)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
	}

	app.Storage = storage
	app.Ws = ws.StartServer(messageHandler, log)
	return app
}

func messageHandler(message []byte, connection *websocket.Conn) {
	connection.WriteMessage(websocket.TextMessage, message)
}
