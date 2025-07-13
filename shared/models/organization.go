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
