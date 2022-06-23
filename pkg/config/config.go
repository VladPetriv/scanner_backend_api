package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL            string
	DBMigrationsPath string
	LogLevel         string
	Port             string
}

func Get() (*Config, error) {
	if err := godotenv.Load("configs/.config.env"); err != nil {
		return nil, fmt.Errorf("failed to get .env file: %w", err)
	}

	return &Config{
		DBURL:            os.Getenv("DATABASE_URL"),
		DBMigrationsPath: os.Getenv("DB_MIGRATION_PATH"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
		Port:             os.Getenv("PORT"),
	}, nil
}
