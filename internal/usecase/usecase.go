package usecase

import "base-project/internal/infrastructure/repository"

type Usecase struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}
