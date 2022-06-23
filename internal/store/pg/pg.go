package pg

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/VladPetriv/scanner_backend_api/pkg/config"

	_ "github.com/lib/pq"
)

var ErrNoDBURL = errors.New("no db url provided")

type DB struct {
	*sqlx.DB
}

func Dial(cfg *config.Config) (*DB, error) {
	if cfg.DBURL == "" {
		return nil, ErrNoDBURL
	}

	db, err := sqlx.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	_, err = db.Exec("SELECT 1;")
	if err != nil {
		return nil, fmt.Errorf("db is not accessible: %w", err)
	}

	return &DB{db}, nil
}
