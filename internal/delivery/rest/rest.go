package rest

import (
	"base-project/internal/usecase"

	"github.com/go-fuego/fuego"
)

type Rest struct {
	srv     *fuego.Server
	usecase *usecase.Usecase
}

func New(usecase *usecase.Usecase) *Rest {
	r := &Rest{
		usecase: usecase,
	}
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

func pingHandler(c fuego.ContextNoBody) (PingResponse, error) {
	return PingResponse{Message: "pong"}, nil
}

type PingResponse struct {
	Message string `json:"message" example:"pong"`
}
