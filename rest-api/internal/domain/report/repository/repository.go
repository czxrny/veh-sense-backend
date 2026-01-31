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

	db = r.queryFilter(filter, db)

	var Reports []models.Report
	if err := db.Find(&Reports).Error; err != nil {
		return nil, err
	}
	return Reports, nil
}

func (r *ReportRepository) FindAllAdmin(ctx context.Context, filter models.ReportFilter) ([]models.AdminReport, error) {
	db := r.db.WithContext(ctx)

	db = db.Where("reports.organization_id = ?", filter.OrganizationID)

	var reports []models.AdminReport
	err := db.Table("reports").
		Select(`
		reports.id,
		reports.user_id,
		user_infos.user_name,
		reports.vehicle_id,
		reports.start_time,
		reports.stop_time,
		reports.acceleration_style,
		reports.braking_style,
		reports.average_speed,
		reports.max_speed,
		reports.kilometers_travelled
	`).
		Joins("JOIN user_infos ON user_infos.id = reports.user_id").
		Find(&reports).Error
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *ReportRepository) queryFilter(filter models.ReportFilter, db *gorm.DB) *gorm.DB {
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

	return db
}
