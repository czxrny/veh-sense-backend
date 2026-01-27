package repository

import (
	"github.com/czxrny/veh-sense-backend/shared/models"
	c "github.com/czxrny/veh-sense-backend/shared/repository"
	"gorm.io/gorm"
)

type ReportDataRepository struct {
	db *gorm.DB
	*c.CommonRepository[models.RawRideRecord]
}

func NewReportDataRepository(db *gorm.DB) *ReportDataRepository {
	return &ReportDataRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[models.RawRideRecord](db),
	}
}
