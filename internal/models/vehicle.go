package models

type Vehicle struct {
	ID             int     `json:"id"`
	OwnerID        *int    `json:"owner_id,omitempty"`
	Private        bool    `json:"private,omitempty"`
	Brand          string  `json:"brand" validate:"required"`
	Model          string  `json:"model" validate:"required"`
	Year           int     `json:"year" validate:"required"`
	EngineCapacity int     `json:"engine_capacity" validate:"required,gte=0"`
	EnginePower    int     `json:"engine_power" validate:"required,gte=0"`
	Plates         string  `json:"plates,omitempty"`
	ExpectedFuel   float64 `json:"expected_fuel" validate:"required,gte=0"`
}
