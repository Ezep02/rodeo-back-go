package repository

import (
	"context"
	"encoding/json"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormReviewRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormReviewRepo(db *gorm.DB, redis *redis.Client) domain.ReviewRepository {
	return &GormReviewRepository{db, redis}
}

func (r *GormReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *GormReviewRepository) Update(ctx context.Context, review *domain.Review) error {
	return nil
}
func (r *GormReviewRepository) Delete(ctx context.Context, id uint) error {

	return nil
}

func (r *GormReviewRepository) List(ctx context.Context) ([]domain.Appointment, error) {
	var (
		appointments []domain.Appointment
		revCacheKey  string = "review"
	)

	// 1. Recuperar productos del cache
	servicesInCache, err := r.redis.Get(ctx, revCacheKey).Result()

	if err == nil {
		json.Unmarshal([]byte(servicesInCache), &appointments)
		return appointments, nil
	}

	if err := r.db.WithContext(ctx).
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Slot", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "date")
		}).
		Preload("Review").
		Where("appointments.id IN (?)", r.db.Table("reviews").Select("appointment_id")).
		Limit(6).
		Find(&appointments).Error; err != nil {
		return nil, err
	}

	return appointments, nil
}
