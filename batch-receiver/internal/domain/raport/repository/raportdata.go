package repository

import (
	"github.com/czxrny/veh-sense-backend/shared/models"
	c "github.com/czxrny/veh-sense-backend/shared/repository"
	"gorm.io/gorm"
)

type RaportDataRepository struct {
	db *gorm.DB
	*c.CommonRepository[models.Raport]
}

func NewRaportDataRepository(db *gorm.DB) *RaportDataRepository {
	return &RaportDataRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[models.Raport](db),
	}
}
