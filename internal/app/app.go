package app

import (
	"context"
	"go-ws/internal/app/ws-server"
	"go-ws/internal/lib/logger/sl"
	"go-ws/internal/storage"
	"go-ws/internal/ws/handler"
	"log/slog"
)

type App struct {
	Storage    *storage.Postgres
	Ws         *ws.Server
	ws_handler *handler.MessageHandler
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
	app.ws_handler = handler.New(log)
	app.Ws = ws.StartServer(app.ws_handler.Handle, log)
	return app
}
