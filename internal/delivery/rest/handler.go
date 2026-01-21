package rest

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func (r *Rest) addHandlers() {
	fuego.Get(r.srv, "/ping", r.pingHandler,
		option.Summary("Ping"),
		option.Description("Check connection to the server"),
		option.Tags("generic"),
	)

	fuego.Get(r.srv, "/user", r.getUsers,
		option.Summary("Get users"),
		option.Description("Get all users"),
		option.Tags("user"),
	)

	fuego.Get(r.srv, "/user/{id}", r.getUser,
		option.Summary("Get user"),
		option.Description("Get user by ID"),
		option.Tags("user"),
		option.Path("id", "ID of the user", param.Required(), param.Example("id", "9046a593-8a27-4c4b-971a-22661274aa60")),
	)

	fuego.Post(r.srv, "/user", r.createUser,
		option.Summary("Create user"),
		option.Description("Create new user"),
		option.Tags("user"),
	)

	fuego.Put(r.srv, "/user/{id}", r.updateUser,
		option.Summary("Update user"),
		option.Description("Update existing user"),
		option.Tags("user"),
		option.Path("id", "ID of the user", param.Required(), param.Example("id", "9046a593-8a27-4c4b-971a-22661274aa60")),
	)

	fuego.Delete(r.srv, "/user/{id}", r.deleteUser,
		option.Summary("Delete user"),
		option.Description("Delete existing user"),
		option.Tags("user"),
		option.Path("id", "ID of the user", param.Required(), param.Example("id", "9046a593-8a27-4c4b-971a-22661274aa60")),
	)
}

type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

func (r *Rest) pingHandler(c fuego.ContextNoBody) (PingResponse, error) {
	return PingResponse{Message: "pong"}, nil
}
