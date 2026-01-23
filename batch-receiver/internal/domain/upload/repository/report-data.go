package repository

import (
	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/model"
	c "github.com/czxrny/veh-sense-backend/shared/repository"
	"gorm.io/gorm"
)

type ReportDataRepository struct {
	db *gorm.DB
	*c.CommonRepository[model.RawRideRecord]
}

func NewReportDataRepository(db *gorm.DB) *ReportDataRepository {
	return &ReportDataRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[model.RawRideRecord](db),
	}
}
