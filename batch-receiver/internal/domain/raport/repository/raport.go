package repository

import (
	"github.com/czxrny/veh-sense-backend/shared/models"
	c "github.com/czxrny/veh-sense-backend/shared/repository"
	"gorm.io/gorm"
)

type RaportRepository struct {
	db *gorm.DB
	*c.CommonRepository[models.Raport]
}

func NewRaportRepository(db *gorm.DB) *RaportRepository {
	return &RaportRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[models.Raport](db),
	}
}
