package models

type UserInfo struct {
	ID              int    `json:"id"`
	UserName        string `json:"username"`
	OrganizationId  *int   `json:"organization_id"`
	TotalKilometers int    `json:"total_kilometers"`
}
