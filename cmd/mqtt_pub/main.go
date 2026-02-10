package main

import (
	"armada_management_system/internal/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main(){
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("vehicle_simulator_B1234XYZ")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil{
		panic(token.Error())
	}
	defer client.Disconnect(250)

	vehicleID := "B1234XYZ"
	baseLat := -6.2088
	baseLong := 106.8456

	fmt.Printf("Starting simulation for vehicle: %s\n", vehicleID)

	for {
		data := models.VehicleLocation{
			VehicleID: vehicleID,
			Latitude: baseLat + (rand.Float64()-0.5)*0.002,
			Longitude: baseLong + (rand.Float64()-0.5)*0.002,
			Timestamp: time.Now().Unix(),
		}

		payload, err := json.Marshal(data)
		if err != nil{
			fmt.Printf("Error marshalling JSON: %v\n", err)
			continue
		}

		topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)
		token := client.Publish(topic, 1, false, payload)
		token.Wait()

		fmt.Printf("Published to %s: %s\n", topic, string(payload))

		time.Sleep(2 * time.Second)
	}
}