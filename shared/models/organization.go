package models

type Organization struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	Name          string `json:"name" validate:"required"`
	Address       string `json:"address" validate:"required"`
	City          string `json:"city" validate:"required"`
	Country       string `json:"country" validate:"required"`
	ZipCode       string `json:"zip_code" validate:"required"`
	CountryCode   string `json:"country_code" validate:"required"`
	ContactNumber string `json:"contact_number" validate:"required"`
}

type OrganizationUpdate struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	City          string `json:"city"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
	CountryCode   string `json:"country_code"`
	ContactNumber string `json:"contact_number"`
}

type OrganizationFilter struct {
	City    string
	Country string
}
