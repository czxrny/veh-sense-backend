package models

type Organization struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	Name          string `json:"name" validate:"required"`
	Address       string `json:"address" validate:"required"`
	City          string `json:"city" validate:"required"`
	ZipCode       string `json:"zip_code" validate:"required"`
	ContactNumber int    `json:"contact_number" validate:"required"`
}

type OrganizationUpdate struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	City          string `json:"city" validate:"required"`
	ZipCode       string `json:"zip_code" validate:"required"`
	ContactNumber int    `json:"contact_number" validate:"required"`
}
