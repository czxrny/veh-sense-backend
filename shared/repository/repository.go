package common

import (
	"context"

	"gorm.io/gorm"
)

type CommonRepository[T any] struct {
	db *gorm.DB
}

func NewCommonRepository[T any](db *gorm.DB) *CommonRepository[T] {
	return &CommonRepository[T]{db: db}
}

func (r *CommonRepository[T]) Add(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *CommonRepository[T]) GetByID(ctx context.Context, id int) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *CommonRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *CommonRepository[T]) DeleteById(ctx context.Context, id int) error {
	var zero T
	return r.db.WithContext(ctx).Delete(&zero, id).Error
}

func (r *CommonRepository[T]) List(ctx context.Context) ([]T, error) {
	var entities []T
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}
