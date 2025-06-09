package postgres

import (
	"context"

	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *Storage

func init() {
	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		},
	).With().Timestamp().Logger().With().Caller().Logger()
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var files []string
	err := filepath.Walk(filepath.Join("../../../SQL/current"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(info.Name()) == ".sql" {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		log.Fatal().Msgf("failed to walk through files: %v", err)
	}

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:15.3-alpine",
		postgres.WithInitScripts(files...),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
		testcontainers.WithLogger(&log.Logger),
	)
	if err != nil {
		log.Fatal().Msgf("failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal().Msgf("failed to get postgres connection string: %v", err)
	}

	pgConn, err := New(connStr)
	if err != nil {
		log.Fatal().Msgf("failed to create postgres connection: %v", err)
	}

	testDB = pgConn

	m.Run()

	err = pgContainer.Terminate(ctx)
	if err != nil {
		log.Fatal().Msgf("failed to terminate postgres container: %v", err)
	}
}
