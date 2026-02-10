package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	RabbitMQURL string
}

func Load() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/armada_db?sslmode=disable"
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	return &Config{
		DatabaseURL: dbURL,
		RabbitMQURL: rabbitURL,
	}
}