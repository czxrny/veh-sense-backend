package models

import "time"

type Raport struct {
	ID                  int       `json:"id"`
	UserID              int       `json:"user_id"`
	OrganizationID      *int      `json:"organization_id"`
	VehicleID           int       `json:"vehicle_id"`
	FrameTime           time.Time `json:"frame_time"`
	StopTime            time.Time `json:"stop_time"`
	AvgFuel             float64   `json:"avg_fuel"`
	MaxFuel             float64   `json:"max_fuel"`
	ConsumedFuel        float64   `json:"consumed_fuel"`
	KilometersTravelled float64   `json:"kilometers_travelled"`
	AccelerationStyle   string    `json:"acceleration_style"`
	BrakingStyle        string    `json:"braking_style"`
	AverageSpeed        float64   `json:"average_speed"`
}

var Styles = []string{
	"Gentle",
	"Moderate",
	"Aggressive",
}
