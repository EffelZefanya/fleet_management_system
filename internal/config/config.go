package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	RabbitMQURL string
}

func Load() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("ERROR: psql db url isn't found, reverting to dev configuration")
		dbURL = "postgres://postgres:postgres@localhost:5432/armada_db?sslmode=disable"
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Println("ERROR: rabbitmq url isn't found, reverting to dev configuration")
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	return &Config{
		DatabaseURL: dbURL,
		RabbitMQURL: rabbitURL,
	}
}