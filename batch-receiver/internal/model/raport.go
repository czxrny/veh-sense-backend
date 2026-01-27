package model

type RideEventType string

type RideEvent struct {
	Timestamp int64         `json:"timestamp"`
	Type      RideEventType `json:"type"`
	Value     float64       `json:"value,omitempty"` // np. m/s^2, rpm, load, km/h
}
