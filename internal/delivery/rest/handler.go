package rest

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func (r *Rest) addHandlers() {
	fuego.Get(r.srv, "/ping", pingHandler,
		option.Summary("Ping"),
		option.Description("Check connection to the server"),
		option.Tags("generic"),
	)
}

type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

func pingHandler(c fuego.ContextNoBody) (PingResponse, error) {
	return PingResponse{Message: "pong"}, nil
}
