package repository

import (
	"context"

	c "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/common/repository"
	"github.com/czxrny/veh-sense-backend/shared/models"
	"gorm.io/gorm"
)

type UserAuthRepository struct {
	db *gorm.DB
	*c.CommonRepository[models.UserAuth]
}

func NewUserAuthRepository(db *gorm.DB) *UserAuthRepository {
	return &UserAuthRepository{
		db:               db,
		CommonRepository: c.NewCommonRepository[models.UserAuth](db),
	}
}

func (r *UserAuthRepository) FindByEmail(ctx context.Context, email string) (*models.UserAuth, error) {
	var found models.UserAuth
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&found).Error; err != nil {
		return nil, err
	}
	return &found, nil
}
