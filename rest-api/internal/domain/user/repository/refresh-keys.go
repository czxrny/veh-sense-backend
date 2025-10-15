package repository

import (
	"context"

	"github.com/czxrny/veh-sense-backend/shared/models"
	c "github.com/czxrny/veh-sense-backend/shared/repository"
	"gorm.io/gorm"
)

type RefreshKeyRepository struct {
	db *gorm.DB
	*c.CommonRepository[models.RefreshInfo]
}

func NewRefreshKeyRepository(db *gorm.DB) *RefreshKeyRepository {
	return &RefreshKeyRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[models.RefreshInfo](db),
	}
}

func (r *RefreshKeyRepository) FindByMatchingKey(ctx context.Context, key string) (*models.RefreshInfo, error) {
	var result models.RefreshInfo
	if err := r.db.WithContext(ctx).First(&result).Where("refresh_key = ?", key).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *RefreshKeyRepository) DeleteByUserID(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.RefreshInfo{}).Error
}
