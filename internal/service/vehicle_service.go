package service

import (
	"armada_management_system/internal/models"
	"armada_management_system/internal/repository"
	"context"
)

type VehicleService struct {
	repo *repository.VehicleRepository
}

func NewVehicleService(repo *repository.VehicleRepository) *VehicleService {
	return &VehicleService{repo: repo}
}

func (s *VehicleService) GetLastLocation(ctx context.Context, vehicleID string) (*models.VehicleLocation, error) {
	return s.repo.GetLastLocation(ctx, vehicleID)
}

func (s *VehicleService) GetLocationHistory(ctx context.Context, vehicleID string, start, end int64) ([]models.VehicleLocation, error) {
	return s.repo.GetLocationHistory(ctx, vehicleID, start, end)
}