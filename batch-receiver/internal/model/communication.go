package model

type ObdFrame struct {
	Rpm          int `json:"rpm" validate:"required"`
	EngineLoad   int `json:"engine_load" validate:"required"`
	VehicleSpeed int `json:"vehicle_speed" validate:"required"`
}

type UploadRideRequest struct {
	VehicleID int64  `json:"vehicle_id" validate:"required"`
	Data      string `json:"data" validate:"required"`
}
