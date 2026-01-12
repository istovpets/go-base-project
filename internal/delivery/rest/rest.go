package rest

import (
	"base-project/internal/usecase"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-fuego/fuego"
)

type Rest struct {
	srv     *fuego.Server
	stopCh  chan struct{}
	usecase *usecase.Usecase
	once    sync.Once
}

func New(usecase *usecase.Usecase) *Rest {
	r := &Rest{
		stopCh:  make(chan struct{}),
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

func (r *Rest) Start(cancel context.CancelCauseFunc) {
	ready := make(chan struct{})

	go func() {
		slog.Info("Starting HTTP server...")
		close(ready)

		err := r.srv.Run()
		if err != nil && err != http.ErrServerClosed {
			cancel(err)
		}

		close(r.stopCh)
	}()

	<-ready // waiting for the goroutine to start
}

func (r *Rest) Stop() error {
	if r.srv == nil {
		return nil
	}

	slog.Info("Stopping HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	r.once.Do(func() {
		err = r.srv.Shutdown(ctx)
	})
	<-r.stopCh

	return err

}

func pingHandler(c fuego.ContextNoBody) (PingResponse, error) {
	return PingResponse{Message: "pong"}, nil
}

type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

func (r *Rest) Ping(ctx context.Context) error {
	const (
		pingTimeout    = 5 * time.Second
		requestTimeout = 200 * time.Millisecond
		retryTimeout   = 100 * time.Millisecond
	)

	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	client := http.Client{
		Timeout: requestTimeout,
	}

	url := "http://" + r.srv.Addr + "/ping"

	for {
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				return nil
			}

			err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		select {
		case <-ctx.Done():
			if err != nil {
				return err
			}

			return ctx.Err()
		case <-time.After(retryTimeout):
		}
	}
}
