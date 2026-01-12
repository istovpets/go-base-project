package app

import (
	"base-project/internal/config"
	"base-project/internal/delivery/rest"
	"base-project/internal/infrastructure/repository"
	"base-project/internal/infrastructure/repository/postgres"
	"base-project/internal/infrastructure/repository/postgres/queries"
	"base-project/internal/usecase"
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/istovpets/pgxhelper/sqlsetpgxhelper"
	"github.com/istovpets/sqlset"
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

func (a *App) Repository() repository.Repository {
	if a.repository == nil {
		sqlSet, err := sqlset.New(queries.QueriesFS)
		if err != nil {
			panicError(fmt.Errorf("failed to load queries: %v", err))
		}

		a.repository = postgres.New(sqlsetpgxhelper.New(sqlSet))
	}

	return a.repository
}

func (a *App) Usecase() *usecase.Usecase {
	if a.usecase == nil {
		a.usecase = usecase.New(a.Repository())
	}

	return a.usecase
}

func (a *App) Rest() *rest.Rest {
	if a.rest == nil {
		a.rest = rest.New(a.Usecase())
	}

	return a.rest
}

// Start/Stop

func (a *App) Stop() error {
	return a.Rest().Stop()
}

func (a *App) Start(cancel context.CancelCauseFunc) error {
	a.Rest().Start(cancel)

	return nil

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

// Errors

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
			slog.Error("failed to initialize application", slog.String("err", err.Error()))

			return
		}

		panic(r)
	}
}
