package main

import (
	"armada_management_system/internal/config"
	"armada_management_system/internal/models"
	"armada_management_system/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("Unable to create connection pool: %v\n", err))
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("Unable to ping database: %v\n", err))
	}
	fmt.Println("Connected to PostgreSQL successfully.")

	repo := repository.NewVehicleRepository(dbpool)

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("armada_management_subscriber")

	//TODO: Create Handler, Service, and Repository layer
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		// TODO: Delete when refactoring or consider using proper log
		fmt.Printf("DEBUG: Received message on topic: %s\n", msg.Topic())

		var location models.VehicleLocation

		err := json.Unmarshal(msg.Payload(), &location)
		if err != nil {
			fmt.Printf("error: received invalid JSON on topic %s: %v", msg.Topic(), err)
			return
		}

		if err := location.Validate(); err != nil {
			fmt.Printf("error: validation failed: %v\n", err)
			return
		}

		if err := repo.SaveLocation(context.Background(), location); err != nil {
			fmt.Printf("error: failed to save location to DB: %v\n", err)
			return
		}

		fmt.Printf("[%s] Data saved to DB: lat %.4f, Long %.4f\n", location.VehicleID, location.Latitude, location.Longitude)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	topic := "/fleet/vehicle/+/location"
	if token := client.Subscribe(topic, 1, messageHandler); token.Wait() && token.Error() != nil {
		fmt.Printf("error: subscription failed %v\n", token.Error())
		return
	}

	fmt.Printf("subscribed to %s. Waiting for data...\n", topic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}