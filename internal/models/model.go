package models

import "github.com/go-playground/validator/v10"

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude float64 `json:"longitude" validate:"gte=-180,lte=180"`
	Timestamp int64   `json:"timestamp" validate:"gt=0"`
}

type LocationPayload struct {
	Latitude  float64 `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude float64 `json:"longitude" validate:"gte=-180,lte=180"`
}

type GeofenceEvent struct {
	VehicleID string          `json:"vehicle_id" validate:"required"`
	Event     string          `json:"event" validate:"required"`
	Location  LocationPayload `json:"location" validate:"required"`
	Timestamp int64           `json:"timestamp" validate:"required"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (v *VehicleLocation) Validate() error {
	return validate.Struct(v)
}

func (lp *LocationPayload) Validate() error {
	return validate.Struct(lp)
}

func (ge *GeofenceEvent) Validate() error {
	return validate.Struct(ge)
}