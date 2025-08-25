package models

type Organization struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	Name          string `json:"name" validate:"required"`
	Address       string `json:"address" validate:"required"`
	City          string `json:"city" validate:"required"`
	Country       string `json:"country" validate:"required"`
	ZipCode       string `json:"zip_code" validate:"required,numeric"`
	CountryCode   string `json:"country_code" validate:"required,len=2,alpha"`
	ContactNumber string `json:"contact_number" validate:"required,numeric"`
	Email         string `json:"email" validate:"required,email"`
}

type OrganizationUpdate struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	City          string `json:"city"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
	CountryCode   string `json:"country_code"`
	ContactNumber string `json:"contact_number"`
	Email         string `json:"email"`
}

type OrganizationFilter struct {
	City    string
	Country string
}
