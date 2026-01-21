package repository

import (
	"base-project/internal/domain"
	"context"

	"github.com/google/uuid"
)

type User interface {
	GetUsers(ctx context.Context) ([]domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	CreateUser(ctx context.Context, user domain.UserData) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
