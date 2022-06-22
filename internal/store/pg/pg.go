package pg

import (
	"github.com/VladPetriv/scanner_backend_api/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var ErrNoDBURL = errors.New("no db url provided")

type DB struct {
	*sqlx.DB
}

func Dial() (*DB, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config")
	}

	if cfg.DBURL == "" {
		return nil, ErrNoDBURL
	}

	db, err := sqlx.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}

	_, err = db.Exec("SELECT 1;")
	if err != nil {
		return nil, errors.Wrap(err, "db is not accessible")
	}

	return &DB{db}, nil
}
