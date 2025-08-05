package repository

import (
	"context"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormCategoryRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormCategoryRepo(db *gorm.DB, redis *redis.Client) domain.CategoryRepository {
	return &GormCategoryRepository{db, redis}
}

func (r *GormCategoryRepository) Create(ctx context.Context, category *domain.Category) error {

	return r.db.WithContext(ctx).Create(category).Error
}

func (r *GormCategoryRepository) Update(ctx context.Context, category *domain.Category) error {

	return r.db.WithContext(ctx).Save(category).Error
}

func (r *GormCategoryRepository) Delete(ctx context.Context, id uint) error {

	// Optionally, remove the category from Redis cache
	return r.db.WithContext(ctx).Delete(&domain.Category{}, id).Error
}

func (r *GormCategoryRepository) List(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
