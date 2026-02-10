package handler

import (
	"armada_management_system/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type VehicleHandler struct {
	repo *repository.VehicleRepository
}

func NewVehicleHandler(repo *repository.VehicleRepository) *VehicleHandler {
	return &VehicleHandler{repo: repo}
}

func (h *VehicleHandler) GetLastLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")

	loc, err := h.repo.GetLastLocation(c.Request.Context(), vehicleID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve location"})
		return
	}

	c.JSON(http.StatusOK, loc)
}

func (h *VehicleHandler) GetHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end query parameters are required"})
		return
	}

	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start timestamp"})
		return
	}

	end, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end timestamp"})
		return
	}

	history, err := h.repo.GetLocationHistory(c.Request.Context(), vehicleID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
