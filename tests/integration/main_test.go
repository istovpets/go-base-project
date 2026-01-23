package integration

import (
	"base-project/internal/delivery/rest"
	"base-project/internal/infrastructure/repository"
	"base-project/internal/infrastructure/repository/postgres"
	"base-project/internal/infrastructure/repository/postgres/queries"
	"base-project/internal/pkg/utils"
	"base-project/internal/usecase"
	"context"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/istovpets/pgxhelper/sqlsetpgxhelper"
	"github.com/istovpets/sqlset"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Test settings and timeouts.
const (
	testTimeout    = 30 * time.Second
	connTimeout    = 5 * time.Second
	serverPort     = 9999
	startupTimeout = 10 * time.Second
)

// Postgres container settings.
const (
	pgImage    = "postgres:16-alpine"
	pgUser     = "postgres"
	pgPassword = "postgres"
	pgDatabase = "testdb"
)

var (
	srv  *httptest.Server
	repo repository.Repository
)

// TestMain sets up the test environment before running tests and tears it down afterward.
func TestMain(m *testing.M) {
	os.Exit(run(m))
}

// run orchestrates the setup and teardown of the test environment.
func run(m *testing.M) int {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	pg, err := initPostgresContainer(ctx)
	if err != nil {
		log.Printf("failed to start postgres container: %v", err)

		return 1
	}
	defer pg.Terminate(context.Background()) //nolint:errcheck

	connStr, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Printf("failed to get connection string: %v", err)

		return 1
	}

	repo, err = initRepository(connStr)
	if err != nil {
		log.Printf("failed to initialize repository: %v", err)

		return 1
	}
	defer repo.Close()

	srv, err = initApp(repo)
	if err != nil {
		log.Printf("failed to initialize application: %v", err)

		return 1
	}
	defer srv.Close()

	return m.Run()
}

// initPostgresContainer starts a Postgres container for testing.
func initPostgresContainer(ctx context.Context) (*pgcontainer.PostgresContainer, error) {
	pg, err := pgcontainer.Run(ctx,
		pgImage,
		pgcontainer.WithDatabase(pgDatabase),
		pgcontainer.WithUsername(pgUser),
		pgcontainer.WithPassword(pgPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(startupTimeout),
		),
	)
	if err != nil {
		return nil, err
	}

	return pg, nil
}

// initRepository initializes the database connection and applies migrations.
func initRepository(connStr string) (*postgres.Database, error) {
	sqlSet, err := sqlset.New(queries.QueriesFS)
	if err != nil {
		return nil, err
	}

	db := postgres.New(sqlsetpgxhelper.New(sqlSet))
	if err = db.Connect(connStr, connTimeout); err != nil {
		return nil, err
	}

	if err = applyMigrations(db.Pool()); err != nil {
		db.Close() // Ensure connection is closed on migration failure.

		return nil, err
	}

	return db, nil
}

// applyMigrations runs database migrations.
func applyMigrations(pool *pgxpool.Pool) error {
	goose.SetDialect("postgres") //nolint:errcheck
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close() //nolint:errcheck

	dir := filepath.Join(utils.ProjectRoot(), "migrations", "postgres", "main")
	if err := goose.Up(db, dir); err != nil {
		return err
	}

	return nil
}

// initApp sets up the application's HTTP server for testing.
func initApp(repo repository.Repository) (*httptest.Server, error) {
	usacase := usecase.New(repo)
	rest := rest.New(serverPort, usacase)
	srv := httptest.NewServer(rest.Mux())

	return srv, nil
}
