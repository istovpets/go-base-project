package postgres

import (
	"base-project/internal/domain"
	"base-project/internal/infrastructure/repository/postgres/queries"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (db *Database) GetUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User

	err := db.Select(ctx, &users, queries.User.GetAll)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db *Database) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}

	err := db.Get(ctx, user, queries.User.Get, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errUserNotFound(id)
		}

		return nil, err
	}

	return user, nil
}

func (db *Database) CreateUser(ctx context.Context, user domain.UserData) (*domain.User, error) {
	createdUser := &domain.User{}

	err := db.Get(ctx, createdUser, queries.User.Create, user.Name)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (db *Database) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	n, err := db.Exec(ctx, queries.User.Update, user.ID, user.Name)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, errUserNotFound(user.ID)
	}

	return &user, nil
}

func (db *Database) DeleteUser(ctx context.Context, id uuid.UUID) error {
	n, err := db.Exec(ctx, queries.User.Delete, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return errUserNotFound(id)
	}

	return nil
}

// common errors

func errUserNotFound(id uuid.UUID) error {
	return fmt.Errorf("user %v not found: %w", id, domain.ErrNotFound)
}
