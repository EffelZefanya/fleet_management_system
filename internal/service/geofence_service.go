package service

import (
	"armada_management_system/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type GeofenceService struct {
	channel *amqp.Channel
}

func NewGeofenceService(channel *amqp.Channel) (*GeofenceService, error) {
	err := channel.ExchangeDeclare("fleet.events", "fanout", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}
	return &GeofenceService{channel: channel}, nil
}

func (s *GeofenceService) CheckAndPublish(loc models.VehicleLocation) error {
	targetLat := -6.2088
	targetLon := 106.8456

	distance := s.haversine(loc.Latitude, loc.Longitude, targetLat, targetLon)

	if distance <= 50 {
		event := models.GeofenceEvent{
			VehicleID: loc.VehicleID,
			Event:     "geofence_entry",
			Location: models.LocationPayload{
				Latitude:  loc.Latitude,
				Longitude: loc.Longitude,
			},
			Timestamp: time.Now().Unix(),
		}

		if err := event.Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
		if err := event.Location.Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		body, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		err = s.channel.PublishWithContext(context.Background(), "fleet.events", "", false, false,
			amqp.Publishing{ContentType: "application/json", Body: body})
		if err != nil {
			return fmt.Errorf("failed to publish geofence event: %w", err)
		}
	}
	return nil
}

func (s *GeofenceService) haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)
	lat1Rad := lat1 * (math.Pi / 180.0)
	lat2Rad := lat2 * (math.Pi / 180.0)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}