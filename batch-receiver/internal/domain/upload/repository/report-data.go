package repository

import (
	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/model"
	c "github.com/czxrny/veh-sense-backend/shared/repository"
	"gorm.io/gorm"
)

type RaportDataRepository struct {
	db *gorm.DB
	*c.CommonRepository[model.RawRideRecord]
}

func NewRaportDataRepository(db *gorm.DB) *RaportDataRepository {
	return &RaportDataRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[model.RawRideRecord](db),
	}
}
