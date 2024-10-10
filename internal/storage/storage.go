package storage

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	log *slog.Logger
	Db  *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
	dbErr      error
)

func NewPG(ctx context.Context, connString string, log *slog.Logger) (*Postgres, error) {
	const op = "storage.NewPG"
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, "")
		if err != nil {
			dbErr = fmt.Errorf("%s: %w", op, err)
			return
		}
		pgInstance = &Postgres{Db: db, log: log}
	})
	return pgInstance, dbErr
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.Db.Ping(ctx)
}

func (pg *Postgres) Close() {
	const op = "storage.Close"
	pg.Db.Close()
	pg.log.With(slog.String("op", op)).Info("DISCONNECTED from database")
}
