package repository

import (
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
