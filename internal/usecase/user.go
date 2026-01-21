package usecase

import (
	"base-project/internal/domain"
	"context"

	"github.com/google/uuid"
)

func (u *Usecase) GetUsers(ctx context.Context) ([]domain.User, error) {
	return u.repo.GetUsers(ctx)
}

func (u *Usecase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return u.repo.GetUser(ctx, id)
}

func (u *Usecase) CreateUser(ctx context.Context, user domain.UserData) (*domain.User, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *Usecase) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	return u.repo.UpdateUser(ctx, user)
}

func (u *Usecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteUser(ctx, id)
}
