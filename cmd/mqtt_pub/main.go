package main

import (
	"armada_management_system/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func simulateVehicle(vehicleID string, baseLat, baseLong float64, client mqtt.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Starting simulation for vehicle: %s", vehicleID)

	for {
		data := models.VehicleLocation{
			VehicleID: vehicleID,
			Latitude:  baseLat + (rand.Float64()-0.5)*0.002,
			Longitude: baseLong + (rand.Float64()-0.5)*0.002,
			Timestamp: time.Now().Unix(),
		}

		payload, err := json.Marshal(data)
		if err != nil {
			log.Printf("[Publisher %s] ERROR: Error marshalling JSON: %v", vehicleID, err)
			continue
		}

		topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)
		token := client.Publish(topic, 1, false, payload)
		token.Wait()

		log.Printf("Published from %s to %s: %s", vehicleID, topic, string(payload))

		time.Sleep(2 * time.Second)
	}
}

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("multi_vehicle_simulator")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT: %v", token.Error())
	}
	defer client.Disconnect(250)

	var wg sync.WaitGroup
	wg.Add(2)

	go simulateVehicle("B1234XYZ", -6.2088, 106.8456, client, &wg)
	go simulateVehicle("B5678ABC", -6.2188, 106.8556, client, &wg)

	wg.Wait()
}