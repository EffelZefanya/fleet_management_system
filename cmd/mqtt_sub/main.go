package main

import (
	"armada_management_system/internal/config"
	"armada_management_system/internal/handler"
	"armada_management_system/internal/repository"
	"armada_management_system/internal/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.Load()

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[MQTT Sub] Unable to create connection pool: %v", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("[MQTT Sub] Unable to ping database: %v", err)
	}
	log.Println("[MQTT Sub] Connected to PostgreSQL successfully.")

	repo := repository.NewVehicleRepository(dbpool)

	rabbitConn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("[MQTT Sub] Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("[MQTT Sub] Failed to open RabbitMQ channel: %v", err)
	}
	defer rabbitCh.Close()

	geofenceService, err := service.NewGeofenceService(rabbitCh)
	if err != nil {
		log.Fatalf("[MQTT Sub] Failed to initialize geofence service: %v", err)
	}

	locationService := service.NewLocationService(repo, geofenceService)
	mqttHandler := handler.NewMqttHandler(locationService)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTTBrokerURL)
	opts.SetClientID("armada_management_subscriber")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[MQTT Sub] MQTT connect error: %v", token.Error())
	}
	defer client.Disconnect(250)

	topic := "/fleet/vehicle/+/location"
	if token := client.Subscribe(topic, 1, mqttHandler.MessageHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("[MQTT Sub] MQTT subscription failed: %v", token.Error())
	}

	log.Printf("Subscribed to %s. Waiting for data...", topic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}