# Armada Management System

A real-time vehicle tracking and fleet management system built with Go. This system handles high-frequency vehicle location updates via MQTT, processes geofencing logic, and provides an API for historical and real-time data.

## Architecture

The system consists of four main components:

1.  **MQTT Publisher (Vehicle Simulator)**: Simulates vehicles moving and publishing location data.
2.  **MQTT Subscriber (Ingestor)**: Subscribes to MQTT topics, persists location data to PostgreSQL, and checks geofences.
3.  **Geofence Worker**: Consumes geofence alert events from RabbitMQ.
4.  **API Server**: Exposes REST endpoints for frontend or external consumers.

## Prerequisites

- **Go** (1.20+) (Optional, for local run)
- **Docker & Docker Compose** (Recommended)
- **PostgreSQL** (with PostGIS recommended for future spatial queries)
- **RabbitMQ**
- **MQTT Broker** (e.g., Mosquitto, EMQX)

## Setup

### 1. Database Schema

Ensure your PostgreSQL database has the required table. The system relies on a unique constraint to prevent duplicate data processing.

```sql
CREATE TABLE IF NOT EXISTS vehicle_locations (
    vehicle_id VARCHAR(50) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(vehicle_id, timestamp)
);
```

### 2. Configuration

The application loads configuration from environment variables (or a config file, depending on your `internal/config` implementation). Ensure the following are reachable:

- Database URL (PostgreSQL)
- RabbitMQ URL
- MQTT Broker URL

### 2. Configuration

The application loads configuration from environment variables (or a config file, depending on your `internal/config` implementation). Ensure the following are reachable:
The application loads configuration from the `.env` file. A default `.env` is provided for local development (pointing to `localhost`).

- Database URL (PostgreSQL)
- RabbitMQ URL
- MQTT Broker URL
When running via Docker Compose, the `docker-compose.yaml` overrides specific variables (like `DATABASE_URL`) to communicate between containers.

## How to Run

### Migrate the PSQL

```
# we'll run the PSQL first to populate it

docker compose -f 'docker-compose.yaml' up -d --build 'postgres'

cat migrations/000001_create_vehicle_locations.up.sql | sudo docker exec -i armada_postgres psql -U postgres -d armada_db
```

### Run Docker Compose

```bash
docker compose up --build
```

This starts Postgres, RabbitMQ, Mosquitto, and all Go services (API, Subscriber, Worker, Simulator).

## API Endpoints

- **Get Last Location**: `GET /vehicles/:vehicle_id/location`
- **Get History**: `GET /vehicles/:vehicle_id/history?start=<unix_timestamp>&end=<unix_timestamp>`

## Notes

### Explanation

- **2 pub to simulate 2 vehicle at once**: The simulator runs concurrent goroutines to mimic multiple vehicles.
- **QoS 1 of MQTT**: Used to ensure data is sent at least once but remains lightweight.
- **Concurrency in RabbitMQ Worker**: The worker handles event consumption asynchronously.
- **Gatekeeping**: Ensures no duplicate data using SQL unique index (`ON CONFLICT DO NOTHING`).

### Limitations

- **Multiple target stations**: Currently limited. Future improvements could use Spatial Index with R-Tree and "Cheap Checking" (e.g., skip if `abs(lat_diff) > 0.001`).
  - Alternatively, use a PostGIS query: `SELECT id, name FROM bus_stops WHERE ST_DWithin(geom, ST_MakePoint(lon, lat)::geography, 50);`
- **Shared Subscription**: Not yet implemented (useful for load balancing MQTT consumers).
