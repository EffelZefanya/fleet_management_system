package models

import "github.com/go-playground/validator/v10"

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"gte=-90,lte=90"`
	Longitude float64 `json:"longitude" validate:"gte=-180,lte=180"`
	Timestamp int64   `json:"timestamp" validate:"gt=0"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (v *VehicleLocation) Validate() error {
	return validate.Struct(v)
}