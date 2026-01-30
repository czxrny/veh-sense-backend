package repository

import (
	"context"

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

func (r *ReportRepository) FindAll(ctx context.Context, filter models.ReportFilter) ([]models.Report, error) {
	db := r.db.WithContext(ctx)

	if filter.CreatedAfter != "" {
		db = db.Where("start_time >= ?", filter.CreatedAfter)
	}
	if filter.CreatedBefore != "" {
		db = db.Where("start_time <= ?", filter.CreatedBefore)
	}

	switch filter.Role {
	case "user":
		db = db.Where("user_id = ?", filter.UserID)
	case "admin":
		db = db.Where("organization_id = ?", filter.OrganizationID)
	}

	var Reports []models.Report
	if err := db.Find(&Reports).Error; err != nil {
		return nil, err
	}
	return Reports, nil
}
