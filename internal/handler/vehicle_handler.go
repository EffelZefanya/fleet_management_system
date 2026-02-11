package handler

import (
	"armada_management_system/internal/models"
	"armada_management_system/internal/service"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	service *service.VehicleService
}

func NewVehicleHandler(service *service.VehicleService) *VehicleHandler {
	return &VehicleHandler{service: service}
}

func (h *VehicleHandler) GetLastLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")

	loc, err := h.service.GetLastLocation(c.Request.Context(), vehicleID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle location not found"})
			return
		}
		log.Printf("ERROR: Failed to retrieve location for vehicle %s: %v", vehicleID, err)
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

	history, err := h.service.GetLocationHistory(c.Request.Context(), vehicleID, start, end)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve history for vehicle %s: %v", vehicleID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
