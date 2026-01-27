package models

type RideRecord struct {
	ReportID  int    `json:"report_id"`
	Data      string `json:"data"`
	EventData string `json:"event_data"`
}

type RawRideRecord struct {
	ReportID  int
	Data      []byte
	EventData []byte
}
