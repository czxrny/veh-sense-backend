package repository

import (
	"context"

	"github.com/czxrny/veh-sense-backend/shared/models"
	c "github.com/czxrny/veh-sense-backend/shared/repository"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

type OrganizationRepository struct {
	db *gorm.DB
	*c.CommonRepository[models.Organization]
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[models.Organization](db),
	}
}

func (r *OrganizationRepository) FindAll(ctx context.Context, filter models.OrganizationFilter) ([]models.Organization, error) {
	db := r.db.WithContext(ctx)

	if filter.City != "" {
		db = db.Where("city = ?", filter.City)
	}
	if filter.City != "" {
		db = db.Where("country = ?", filter.Country)
	}

	var organizations []models.Organization
	if err := db.Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

func (r *OrganizationRepository) UpdatePartial(ctx context.Context, id int, update *models.OrganizationUpdate) error {
	return r.db.WithContext(ctx).
		Model(&models.Organization{}).
		Where("id = ?", id).
		Updates(update).Error
}
