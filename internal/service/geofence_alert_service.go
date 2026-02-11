package service

import (
	"armada_management_system/internal/models"
	"log"
)

type GeofenceAlertService struct {

}

func NewGeofenceAlertService() *GeofenceAlertService {
	return &GeofenceAlertService{}
}

func (s *GeofenceAlertService) ProcessAlert(event models.GeofenceEvent) error {
	log.Printf(" [x] ALERT: Vehicle %s entered geofence at %d. Location: [%.4f, %.4f]",
		event.VehicleID, event.Timestamp, event.Location.Latitude, event.Location.Longitude)
	return nil
}