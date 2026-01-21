package rest

import (
	"base-project/internal/domain"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

func (r *Rest) getUsers(c fuego.ContextNoBody) ([]domain.User, error) {
	return r.usecase.GetUsers(c.Context())
}

func (r *Rest) getUser(c fuego.ContextNoBody) (*domain.User, error) {
	s := c.PathParam("id")
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid uuid", Detail: s, Err: err}
	}

	return r.usecase.GetUser(c.Context(), id)
}

func (r *Rest) createUser(c fuego.ContextWithBody[domain.UserData]) (*domain.User, error) {
	body, err := c.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Err: err}
	}

	return r.usecase.CreateUser(c.Context(), body)
}

func (r *Rest) updateUser(c fuego.ContextWithBody[domain.UserData]) (*domain.User, error) {
	s := c.PathParam("id")
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid uuid", Detail: s, Err: err}
	}

	body, err := c.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Err: err}
	}

	return r.usecase.UpdateUser(c.Context(), domain.User{ID: id, UserData: body})
}

func (r *Rest) deleteUser(c fuego.ContextNoBody) (any, error) {
	s := c.PathParam("id")
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid uuid", Detail: s, Err: err}
	}

	return nil, r.usecase.DeleteUser(c.Context(), id)
}
