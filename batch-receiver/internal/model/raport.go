package model

import "time"

type IntermediateRaport struct {
	UserID             int
	OrganizationID     *int
	VehicleID          int
	FrameTime          time.Time
	LastFrameTime      time.Time
	Rpm                int
	EngineLoad         int
	VehicleSpeed       int
	SpeedSum           float64
	AggressiveBrakings int
	EventCount         int
}

type UpperDeviation struct {
	Timestamp time.Time
	ZValue    float64
	Load      float64
}
