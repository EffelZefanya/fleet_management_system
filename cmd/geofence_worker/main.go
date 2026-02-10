package main

import (
	"armada_management_system/internal/config"
	"armada_management_system/internal/models"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.Load()

	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("fleet.events", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	q, err := ch.QueueDeclare("geofence_alerts", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	err = ch.QueueBind(q.Name, "", "fleet.events", false, nil)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf(" [*] Waiting for geofence alerts. To exit press CTRL+C")

	go func() {
		for d := range msgs {
			var event models.GeofenceEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}
			log.Printf(" [x] ALERT: Vehicle %s entered geofence at %d. Location: [%.4f, %.4f]",
				event.VehicleID, event.Timestamp, event.Location.Latitude, event.Location.Longitude)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down geofence worker...")
}
