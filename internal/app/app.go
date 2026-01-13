package app

import (
	"base-project/internal/config"
	"base-project/internal/delivery/rest"
	"base-project/internal/infrastructure/repository"
	"base-project/internal/usecase"
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type App struct {
	config     *config.Config
	repository repository.Repository
	usecase    *usecase.Usecase
	rest       *rest.Rest
}

func New() *App {
	a := &App{}
	initLog(&a.Config().LogLevel)

	return a
}

func initLog(level slog.Leveler) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: level,
		}))
	slog.SetDefault(logger)
}

// Start/Stop

func (a *App) Start(cancel context.CancelCauseFunc) error {
	if err := a.Repository().Connect(a.Config().DBConnStr, 5*time.Second); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	a.Rest().Start(cancel)

	return nil
}

func (a *App) Stop() error {
	defer a.Repository().Close()

	return a.Rest().Stop()
}

func (a *App) CheckHealth(ctx context.Context) error {
	return a.Rest().Ping(ctx)
}

func (a *App) Wait(ctx context.Context) error {
	<-ctx.Done()

	if err := context.Cause(ctx); err != nil && err != context.Canceled {
		return err
	}

	return nil
}
