package rest

import (
	"base-project/internal/domain"
	"errors"

	"github.com/go-fuego/fuego"
)

func errorHandler(err error) error {
	if errors.Is(err, domain.ErrNotFound) {
		return fuego.NotFoundError{Title: "Not found", Detail: err.Error(), Err: err}
	}

	return fuego.ErrorHandler(err)
}
