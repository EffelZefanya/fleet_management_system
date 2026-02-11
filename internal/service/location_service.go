package service

import (
	"armada_management_system/internal/models"
	"armada_management_system/internal/repository"
	"context"
	"fmt"
)

type LocationService struct {
	repo     *repository.VehicleRepository
	geofence *GeofenceService
}

func NewLocationService(repo *repository.VehicleRepository, geofence *GeofenceService) *LocationService {
	return &LocationService{
		repo:     repo,
		geofence: geofence,
	}
}

func (s *LocationService) ProcessLocation(ctx context.Context, location models.VehicleLocation) error {
	if err := location.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.repo.SaveLocation(ctx, location); err != nil {
		return fmt.Errorf("failed to save location to DB: %w", err)
	}

	// We can decide if a geofence check failure should be a fatal error or just logged.
	return s.geofence.CheckAndPublish(location)
}