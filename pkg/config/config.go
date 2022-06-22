package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	DBURL            string
	DBMigrationsPath string
	LogLevel         string
}

func Get() (*Config, error) {
	if err := godotenv.Load("configs/.config.env"); err != nil {
		return nil, errors.Wrap(err, "failed to get .env file")
	}

	return &Config{
		DBURL:            os.Getenv("DATABASE_URL"),
		DBMigrationsPath: os.Getenv("DB_MIGRATION_PATH"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
	}, nil
}
