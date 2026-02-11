package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	RabbitMQURL string
	MQTTBrokerURL string
}

func Load() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("WARN: psql db url isn't found, reverting to dev configuration")
		dbURL = "postgres://postgres:postgres@localhost:5432/armada_db?sslmode=disable"
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Println("WARN: rabbitmq url isn't found, reverting to dev configuration")
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	mqttURL := os.Getenv("MQTT_BROKER_URL")
	if mqttURL == "" {
		log.Println("WARN: mqtt broker url isn't found, reverting to dev configuration")
		mqttURL = "tcp://localhost:1883"
	}

	return &Config{
		DatabaseURL: dbURL,
		RabbitMQURL: rabbitURL,
		MQTTBrokerURL: mqttURL,
	}
}