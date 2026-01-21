package integration

import (
	"base-project/internal/config"
	"base-project/internal/delivery/rest"
	"base-project/internal/infrastructure/repository/postgres"
	"base-project/internal/infrastructure/repository/postgres/queries"
	"base-project/internal/usecase"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/istovpets/pgxhelper/sqlsetpgxhelper"
	"github.com/istovpets/sqlset"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	dbConnStr := "postgres://postgres:postgres@localhost:5432/dadata_v2?sslmode=disable"
	// // Get DB connection string from environment
	// dbConnStr := os.Getenv("DB_CONN_STR")
	// if dbConnStr == "" {
	// 	t.Fatal("DB_CONN_STR environment variable not set")
	// }

	cfg := &config.Config{
		LogLevel:  slog.LevelVar{},
		Port:      9999,
		DBConnStr: dbConnStr,
	}

	// Initialize repository
	sqlSet, err := sqlset.New(queries.QueriesFS)
	require.NoError(t, err, "failed to load queries")
	repo := postgres.New(sqlsetpgxhelper.New(sqlSet))
	err = repo.Connect(cfg.DBConnStr, 5*time.Second)
	require.NoError(t, err, "failed to connect to database")

	// Initialize usecase
	uc := usecase.New(repo)

	// Initialize rest
	r := rest.New(cfg.Port, uc)

	// Create a new test server
	srv := httptest.NewServer(r.Mux())

	// Cleanup function
	cleanup := func() {
		srv.Close()
		repo.Close()
	}

	return srv, cleanup
}
