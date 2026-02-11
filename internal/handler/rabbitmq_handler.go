package handler

import (
	"armada_management_system/internal/models"
	"armada_management_system/internal/service"
	"encoding/json"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type GeofenceAlertHandler struct {
	alertService *service.GeofenceAlertService
}

func NewGeofenceAlertHandler(s *service.GeofenceAlertService) *GeofenceAlertHandler {
	return &GeofenceAlertHandler{alertService: s}
}

func (h *GeofenceAlertHandler) Handle(msgs <-chan amqp.Delivery) {
	const workerCount = 2
	var wg sync.WaitGroup

	log.Printf("Starting %d workers for geofence alerts...", workerCount)

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go h.worker(i, msgs, &wg)
	}

	wg.Wait()
	log.Println("All workers have finished. Handler stopping.")
}

func (h *GeofenceAlertHandler) worker(id int, msgs <-chan amqp.Delivery, wg *sync.WaitGroup) {
    defer wg.Done()
    log.Printf("Worker %d started", id)

    for d := range msgs {
        var event models.GeofenceEvent
        
        if err := json.Unmarshal(d.Body, &event); err != nil {
            log.Printf("Worker %d ERROR: JSON decode failed: %s", id, err)
            d.Nack(false, false) 
            continue
        }

        log.Printf("Worker %d processing vehicle: %s", id, event.VehicleID)
        
        // 2. Handle Processing Errors
        if err := h.alertService.ProcessAlert(event); err != nil {
            log.Printf("Worker %d ERROR: Service failed: %v", id, err)
            d.Nack(false, true)
            continue
        }

        if err := d.Ack(false); err != nil {
            log.Printf("Worker %d ERROR: Ack failed: %v", id, err)
        }
    }
}