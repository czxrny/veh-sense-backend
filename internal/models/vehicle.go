package models

type Vehicle struct {
	ID             int64   `json:"id"`
	OwnerID        *int64  `json:"owner_id,omitempty"`
	Private        bool    `json:"private,omitempty"`
	Brand          string  `json:"brand" validate:"required"`
	Model          string  `json:"model" validate:"required"`
	Year           int64   `json:"year" validate:"required"`
	EngineCapacity int64   `json:"engine_capacity" validate:"required,gte=0"`
	EnginePower    int64   `json:"engine_power" validate:"required,gte=0"`
	Plates         string  `json:"plates,omitempty"`
	ExpectedFuel   float64 `json:"expected_fuel" validate:"required,gte=0"`
}
