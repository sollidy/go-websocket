package app

import (
	"context"
	"go-ws/internal/storage"
	"log/slog"
)

type App struct {
	Storage *storage.Postgres
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
	app.Storage = storage

	return app
}
