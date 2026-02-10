package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
}

func Load() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/armada_db?sslmode=disable"
	}
	return &Config{
		DatabaseURL: dbURL,
	}
}