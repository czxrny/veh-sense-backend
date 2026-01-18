package models

import "time"

type RaportFrame struct {
	Rpm          int `json:"rpm"`
	VehicleSpeed int `json:"vehicle_speed"`
	EngineLoad   int `json:"engine_load"`
}

type Raport struct {
	ID                  int       `json:"id"`
	UserID              int       `json:"user_id"`
	OrganizationID      *int      `json:"organization_id"`
	VehicleID           int       `json:"vehicle_id"`
	StartTime           time.Time `json:"start_time"`
	StopTime            time.Time `json:"stop_time"`
	AccelerationStyle   string    `json:"acceleration_style"`
	BrakingStyle        string    `json:"braking_style"`
	AverageSpeed        float64   `json:"average_speed"`
	MaxSpeed            float64   `json:"max_speed"`
	KilometersTravelled float64   `json:"kilometers_travelled"`
}

type RaportFilter struct {
	CreatedAfter   string
	CreatedBefore  string
	UserID         int
	OrganizationID *int
	Role           string
}

var Styles = []string{
	"Gentle",
	"Moderate",
	"Aggressive",
}
