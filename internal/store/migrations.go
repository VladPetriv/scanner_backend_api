package store

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	"github.com/VladPetriv/scanner_backend_api/pkg/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var ErrNoDBURL = errors.New("no db url provided")

func runMigrations(cfg *config.Config) error {
	if cfg.DBURL == "" {
		return ErrNoDBURL
	}

	m, err := migrate.New(
		cfg.DBMigrationsPath,
		cfg.DBURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrations: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
