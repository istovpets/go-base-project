package postgres

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func New(log *slog.Logger) *Database {
	return &Database{
		log: log,
	}
}
