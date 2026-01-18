package repository

import (
	"context"

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

func (r *RaportRepository) FindAll(ctx context.Context, filter models.RaportFilter) ([]models.Raport, error) {
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
		db = db.Where("ogranization_id = ?", filter.OrganizationID)
	}

	var raports []models.Raport
	if err := db.Find(&raports).Error; err != nil {
		return nil, err
	}
	return raports, nil
}
