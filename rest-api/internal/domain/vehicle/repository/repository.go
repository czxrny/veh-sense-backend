package repository

import (
	"context"

	c "github.com/czxrny/veh-sense-backend/rest-api/internal/repositories/common"
	"github.com/czxrny/veh-sense-backend/shared/models"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

type VehicleRepository struct {
	db *gorm.DB
	*c.CommonRepository[models.Vehicle]
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[models.Vehicle](db),
	}
}

func (r *VehicleRepository) FindAll(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error) {
	db := r.db.WithContext(ctx)

	if filter.Brand != "" {
		db = db.Where("brand = ?", filter.Brand)
	}
	if filter.MinCapacity != "" {
		db = db.Where("engine_capacity >= ?", filter.MinCapacity)
	}
	if filter.MaxCapacity != "" {
		db = db.Where("engine_capacity <= ?", filter.MaxCapacity)
	}
	if filter.MinEnginePower != "" {
		db = db.Where("engine_power >= ?", filter.MinEnginePower)
	}
	if filter.MaxEnginePower != "" {
		db = db.Where("engine_power <= ?", filter.MaxEnginePower)
	}

	switch filter.Role {
	case "user":
		if filter.OrganizationID != nil {
			db = db.Where("(owner_id = ? OR (organization_id = ? AND owner_id IS NULL))",
				filter.UserID,
				filter.OrganizationID,
			)
		} else {
			db = db.Where("owner_id = ?", filter.UserID)
		}
	case "admin":
		db = db.Where("organization_id = ?", filter.OrganizationID)
	}

	var vehicles []models.Vehicle
	if err := db.Find(&vehicles).Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (r *VehicleRepository) UpdatePartial(ctx context.Context, id int, update *models.VehicleUpdate) error {
	return r.db.WithContext(ctx).
		Model(&models.Vehicle{}).
		Where("id = ?", id).
		Updates(update).Error
}
