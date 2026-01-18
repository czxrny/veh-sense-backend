package model

import "time"

type RideEventType string

type RideEvent struct {
	Timestamp time.Time     `json:"timestamp"`
	Type      RideEventType `json:"type"`
	Value     float64       `json:"value,omitempty"` // np. m/s^2, rpm, load, km/h
}

type RideRecord struct {
	RaportID  int    `json:"raport_id"`
	Data      string `json:"data"`
	EventData string `json:"event_data"`
}
