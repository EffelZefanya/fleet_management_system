CREATE TABLE IF NOT EXISTS vehicle_locations (
    id SERIAL PRIMARY KEY,
    vehicle_id VARCHAR(50) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_vehicle_id_timestamp ON vehicle_locations (vehicle_id, timestamp);