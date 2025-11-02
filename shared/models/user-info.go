package models

type UserInfo struct {
	ID              int    `json:"id"`
	UserName        string `json:"user_name"`
	OrganizationId  *int   `json:"organization_id"`
	TotalKilometers int    `json:"total_kilometers"`
	NumberOfRides   int    `json:"number_of_rides"`
	Rating          string `json:"rating"`
}
