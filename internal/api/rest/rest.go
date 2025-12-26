package rest

import (
	"github.com/go-fuego/fuego"
)

type Rest struct {
	srv *fuego.Server
}

func New() *Rest {
	r := &Rest{}
	r.srv = fuego.NewServer(
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(
				fuego.OpenAPIConfig{
					JSONFilePath:     "doc/openapi.json",
					PrettyFormatJSON: true,
				},
			),
		),
	)

	// Register handlers
	r.addHandlers()

	return r
}

// Start/Stop

func (r *Rest) Start() {

}

func (r *Rest) Stop() {

}

func (r *Rest) Run() {
	r.srv.Run()
}

// Endpoint /ping
// @Summary      Ping the server
// @Description  Simple health-check endpoint
// @Tags         health
// @Success      200  {object}  PingResponse
// @Router       /ping [get]
func pingHandler(c fuego.ContextNoBody) (PingResponse, error) {
	return PingResponse{Message: "pong"}, nil
}

type PingResponse struct {
	Message string `json:"message" example:"pong"`
}
