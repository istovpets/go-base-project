package app

import (
	"base-project/internal/config"
	"base-project/internal/delivery/rest"
	"base-project/internal/infrastructure/repository/postgres"
	"base-project/internal/pkg/logger"
	"fmt"
	"log/slog"
)

type App struct {
	config     *config.Config
	log        *slog.Logger
	repository *postgres.Database
	rest       *rest.Rest
}

// Getters

func (a *App) Config() *config.Config {
	if a.config == nil {
		var err error
		a.config, err = config.New()
		if err != nil {
			panicError(fmt.Errorf("failed to load config: %w", err))
		}
	}

	return a.config
}

func (a *App) Log() *slog.Logger {
	if a.log == nil {
		var err error
		a.log, err = logger.New(a.Config().LogLevel.Level())
		if err != nil {
			panicError(fmt.Errorf("failed to initialize logger: %w", err))
		}
	}

	return a.log
}

func (a *App) Repository() *postgres.Database {
	if a.repository == nil {
		a.repository = postgres.New(a.Log())
	}

	return a.repository
}

func (a *App) Rest() *rest.Rest {
	if a.rest == nil {
		a.rest = rest.New()
	}

	return a.rest
}

// Start/Stop

func (a *App) Stop() {
	a.Log().Debug("application stopped")
}

func (a *App) Run() {
	a.Log().Debug("application started")
	a.Rest().Run()
}

func New() *App {
	return &App{}
}

type ErrorPanic struct {
	Err error
}

func (e ErrorPanic) Error() string {
	return e.Err.Error()
}

func panicError(err error) {
	panic(ErrorPanic{Err: err})
}

func (a *App) Recover() {
	if r := recover(); r != nil {
		if err, ok := r.(ErrorPanic); ok {
			if a.log == nil {
				fmt.Printf("ERROR: Failed to initialize application. err: %v\n", err)
			} else {
				a.log.Error("failed to initialize application", slog.String("err", err.Error()))
			}

			return
		}

		panic(r)
	}
}
