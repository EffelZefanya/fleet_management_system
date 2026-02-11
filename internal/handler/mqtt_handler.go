package handler

import (
	"armada_management_system/internal/models"
	"armada_management_system/internal/service"
	"context"
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttHandler struct {
	locationService *service.LocationService
}

func NewMqttHandler(locationService *service.LocationService) *MqttHandler {
	return &MqttHandler{locationService: locationService}
}

func (h *MqttHandler) MessageHandler(client mqtt.Client, msg mqtt.Message) {
	var location models.VehicleLocation
	if err := json.Unmarshal(msg.Payload(), &location); err != nil {
		log.Printf("ERROR: Received invalid JSON on topic %s: %v", msg.Topic(), err)
		return
	}

	if err := h.locationService.ProcessLocation(context.Background(), location); err != nil {
		log.Printf("ERROR: Failed to process location for vehicle %s: %v", location.VehicleID, err)
		return
	}

	log.Printf("[%s] Location processed successfully: lat %.4f, Long %.4f", location.VehicleID, location.Latitude, location.Longitude)
}