package repository

import (
	"armada_management_system/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VehicleRepository struct {
	db *pgxpool.Pool
}

func NewVehicleRepository(db *pgxpool.Pool) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) SaveLocation(ctx context.Context, loc models.VehicleLocation) error {
	query := `
		INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp)
		VALUES ($1, $2, $3, $4)
	`
	ts := time.Unix(loc.Timestamp, 0)

	_, err := r.db.Exec(ctx, query, loc.VehicleID, loc.Latitude, loc.Longitude, ts)
	if err != nil {
		return fmt.Errorf("failed to insert vehicle location: %w", err)
	}
	return nil
}