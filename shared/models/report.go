package models

type ReportFrame struct {
	Rpm          int `json:"rpm"`
	VehicleSpeed int `json:"vehicle_speed"`
	EngineLoad   int `json:"engine_load"`
}

type Report struct {
	ID                  int     `json:"id"`
	UserID              int     `json:"user_id"`
	OrganizationID      *int    `json:"organization_id"`
	VehicleID           int     `json:"vehicle_id"`
	StartTime           int64   `json:"start_time"`
	StopTime            int64   `json:"stop_time"`
	AccelerationStyle   string  `json:"acceleration_style"`
	BrakingStyle        string  `json:"braking_style"`
	AverageSpeed        float64 `json:"average_speed"`
	MaxSpeed            float64 `json:"max_speed"`
	KilometersTravelled float64 `json:"kilometers_travelled"`
}

type ReportFilter struct {
	CreatedAfter   string
	CreatedBefore  string
	UserID         int
	OrganizationID *int
	Role           string
}

// Replacing user id with user name
type AdminReport struct {
	ID                  int     `json:"id"`
	Username            int     `json:"user_name"`
	UserID              int     `json:"user_id"`
	VehicleID           int     `json:"vehicle_id"`
	StartTime           int64   `json:"start_time"`
	StopTime            int64   `json:"stop_time"`
	AccelerationStyle   string  `json:"acceleration_style"`
	BrakingStyle        string  `json:"braking_style"`
	AverageSpeed        float64 `json:"average_speed"`
	MaxSpeed            float64 `json:"max_speed"`
	KilometersTravelled float64 `json:"kilometers_travelled"`
}
