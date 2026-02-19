package app

import (
	"base-project/internal/config"
	"base-project/internal/delivery/rest"
	"base-project/internal/infrastructure/repository"
	"base-project/internal/infrastructure/repository/postgres"
	"base-project/internal/infrastructure/repository/postgres/queries"
	"base-project/internal/usecase"
	"fmt"

	"github.com/istovpets/pgxhelper/sqlsetpgxhelper"
	"github.com/istovpets/sqlset"
)

func (a *App) Config() *config.Config {
	a.configOnce.Do(func() {
		var err error
		a.config, err = config.New()
		if err != nil {
			panicError(fmt.Errorf("failed to load config: %w", err))
		}
	})

	return a.config
}

func (a *App) Repository() repository.Repository {
	a.repositoryOnce.Do(func() {
		sqlSet, err := sqlset.New(queries.QueriesFS)
		if err != nil {
			panicError(fmt.Errorf("failed to load queries: %v", err))
		}

		a.repository = postgres.New(sqlsetpgxhelper.New(sqlSet))
	})

	return a.repository
}

func (a *App) Usecase() *usecase.Usecase {
	a.usecaseOnce.Do(func() {
		a.usecase = usecase.New(a.Repository())
	})

	return a.usecase
}

func (a *App) Rest() *rest.Rest {
	a.restOnce.Do(func() {
		a.rest = rest.New(a.Config().Port, a.Usecase())
	})

	return a.rest
}
