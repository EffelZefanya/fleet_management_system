package main

import (
	"armada_management_system/internal/config"
	"armada_management_system/internal/handler"
	"armada_management_system/internal/repository"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	repo := repository.NewVehicleRepository(dbpool)
	h := handler.NewVehicleHandler(repo)

	r := gin.Default()

	r.GET("/vehicles/:vehicle_id/location", h.GetLastLocation)
	r.GET("/vehicles/:vehicle_id/history", h.GetHistory)

	//TODO: Refactor so the handler will call service, instead of repository directly.
	fmt.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
