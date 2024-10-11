package app

import (
	"context"
	"go-ws/internal/app/ws-app"
	"go-ws/internal/lib/logger/sl"
	"go-ws/internal/storage"
	handler "go-ws/internal/ws"
	"log/slog"
)

type App struct {
	Storage *storage.Postgres
	Ws      *ws.AppWs
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
	app.Ws = ws.New(log)
	handler.Init()
	return app
}
