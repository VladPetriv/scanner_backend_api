package store

import (
	"github.com/VladPetriv/scanner_backend_api/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var ErrNoDBURL = errors.New("no db url provided")

func runMigrations() error {
	cfg, err := config.Get()
	if err != nil {
		return errors.WithStack(err)
	}

	if cfg.DBURL == "" {
		return ErrNoDBURL
	}

	m, err := migrate.New(
		cfg.DBMigrationsPath,
		cfg.DBURL,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create migrations")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to run migrations")
	}

	return nil
}
