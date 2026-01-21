package rest

import (
	"base-project/internal/usecase"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/go-fuego/fuego"
)

type Rest struct {
	srv     *fuego.Server
	port    int16
	stopCh  chan struct{}
	usecase *usecase.Usecase
	once    sync.Once
}

func New(port int16, usecase *usecase.Usecase) *Rest {
	r := &Rest{
		port:    port,
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
			fuego.WithErrorHandler(errorHandler),
		),
	)

	fuego.Use(r.srv, recoverMiddleware())

	// Register handlers
	r.addHandlers()

	return r
}

// recover middleware
func recoverMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					slog.Error("panic recovered in handler",
						"panic", rec,
						"stack", string(debug.Stack()),
						"method", r.Method,
						"path", r.URL.Path,
					)

					http.Error(w, "Internal Server Error (panic recovered)", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func (r *Rest) Mux() http.Handler {
	return r.srv.Mux
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

// Ping

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
			resp.Body.Close() //nolint:errcheck

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
