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

func (r *VehicleRepository) GetLastLocation(ctx context.Context, vehicleID string) (*models.VehicleLocation, error) {
	query := `
		SELECT vehicle_id, latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`
	var loc models.VehicleLocation
	var ts time.Time

	err := r.db.QueryRow(ctx, query, vehicleID).Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &ts)
	if err != nil {
		return nil, err
	}
	loc.Timestamp = ts.Unix()
	return &loc, nil
}

func (r *VehicleRepository) GetLocationHistory(ctx context.Context, vehicleID string, start, end int64) ([]models.VehicleLocation, error) {
	query := `
		SELECT vehicle_id, latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1 AND timestamp >= $2 AND timestamp <= $3
		ORDER BY timestamp ASC
	`
	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	rows, err := r.db.Query(ctx, query, vehicleID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query location history: %w", err)
	}
	defer rows.Close()

	locations := make([]models.VehicleLocation, 0)
	for rows.Next() {
		var loc models.VehicleLocation
		var ts time.Time
		if err := rows.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &ts); err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		loc.Timestamp = ts.Unix()
		locations = append(locations, loc)
	}
	return locations, nil
}