package storage

import (
	"context"
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
	pgOnce.Do(func() {
		db, dbErr := pgxpool.New(ctx, connString)
		if dbErr != nil {
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
