package main

import (
	"armada_management_system/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("armada_management_subscriber")

	messageHandler := func(client mqtt.Client, msg mqtt.Message){
		var location models.VehicleLocation

		err := json.Unmarshal(msg.Payload(), &location)
		if err != nil{
			fmt.Printf("error: received invalid JSON on topic %s: %v", msg.Topic(), err)
			return
		}

		if err := location.Validate(); err != nil {
			fmt.Printf("error: validation failed: %v\n", err)
			return
		}

		fmt.Printf("[%s] received location: lat %.4f, Long %.4f at %d\n", location.VehicleID, location.Latitude, location.Longitude, location.Timestamp)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil{
		panic(token.Error())
	}
	defer client.Disconnect(250)

	topic := "/fleet/vehicle/+/location"
	if token := client.Subscribe(topic, 1, messageHandler); token.Wait() && token.Error() != nil{
		fmt.Printf("error: subscription failed %v\n", token.Error())
		return
	}

	fmt.Printf("subscribed to %s. Waiting for data...\n", topic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}