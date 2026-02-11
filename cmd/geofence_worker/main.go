package main

import (
	"armada_management_system/internal/config"
	"armada_management_system/internal/handler"
	"armada_management_system/internal/service"
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
		log.Fatalf("[Geofence Worker] Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[Geofence Worker] Failed to open a channel: %v", err)
	}
	defer ch.Close()
	err = ch.Qos(10, 0, false)
	if err != nil{
		log.Fatalf("[Geofence Worker] Failed to set QoS: %v", err)
	}

	err = ch.ExchangeDeclare("fleet.events", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[Geofence Worker] Failed to declare exchange: %v", err)
	}

	q, err := ch.QueueDeclare("geofence_alerts", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[Geofence Worker] Failed to declare queue: %v", err)
	}

	err = ch.QueueBind(q.Name, "", "fleet.events", false, nil)
	if err != nil {
		log.Fatalf("[Geofence Worker] Failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("[Geofence Worker] Failed to register a consumer: %v", err)
	}

	alertService := service.NewGeofenceAlertService()
	alertHandler := handler.NewGeofenceAlertHandler(alertService)

	log.Printf(" [*] Waiting for geofence alerts. To exit press CTRL+C")

	done := make(chan struct{})
	go func() {
		defer close(done)
		alertHandler.Handle(msgs)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down geofence worker... closing RabbitMQ channel.")
	ch.Close()
	<-done
	log.Println("Geofence worker shut down gracefully.")
}
