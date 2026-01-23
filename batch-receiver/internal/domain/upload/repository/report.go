package repository

import (
	"github.com/czxrny/veh-sense-backend/shared/models"
	c "github.com/czxrny/veh-sense-backend/shared/repository"
	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
	*c.CommonRepository[models.Report]
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[models.Report](db),
	}
}
