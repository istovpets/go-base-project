package repository

import (
	"context"
	"time"
)

type Repository interface {
	Ping(ctx context.Context) error
	Connect(connStr string, timeout time.Duration) error
	Close()
	User
}
