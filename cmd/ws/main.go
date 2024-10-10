package main

import (
	"context"
	"go-ws/internal/app"
	"go-ws/internal/lib/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Env          string
	PgConnString string
}

func main() {
	cfg := Config{
		Env:          "local",
		PgConnString: "connection string",
	}

	ctx, cancel := context.WithCancel(context.Background())
	log := logger.SetupLogger(cfg.Env)
	log.Debug("starting application", slog.Any("cfg", cfg))

	application := app.New(ctx, log, cfg.PgConnString)

	shutdown(cancel, application, log) // graceful shutdown
}

func shutdown(cancel context.CancelFunc, app *app.App, log *slog.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	cancel()
	app.Storage.Close()

	log.Warn("STOPED application", slog.String("signal", sign.String()))
}
