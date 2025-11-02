package repository

import (
	"context"

	"github.com/czxrny/veh-sense-backend/shared/models"
	c "github.com/czxrny/veh-sense-backend/shared/repository"
	"gorm.io/gorm"
)

type UserInfoRepository struct {
	db *gorm.DB
	*c.CommonRepository[models.UserInfo]
}

func NewUserInfoRepository(db *gorm.DB) *UserInfoRepository {
	return &UserInfoRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[models.UserInfo](db),
	}
}

func (r *UserInfoRepository) FindByOrganizationId(ctx context.Context, organizationId int) ([]models.UserInfo, error) {
	var users []models.UserInfo
	if err := r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
