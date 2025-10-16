package repository

import (
	"context"

	"github.com/ezep02/rodeo/internal/booking/domain/services"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormServiceRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormServiceRepo(db *gorm.DB, redis *redis.Client) services.ServicesRepository {
	return &GormServiceRepository{db, redis}
}

func (r *GormServiceRepository) GetByID(ctx context.Context, id uint) (*services.Service, error) {
	var appt services.Service

	if err := r.db.WithContext(ctx).
		Preload("Medias").
		Preload("Categories").
		Preload("Promotions").
		First(&appt, id).Error; err != nil {
		return nil, err
	}

	return &appt, nil
}

// GetTotalPriceByIDs devuelve la suma de los precios de los servicios
func (r *GormServiceRepository) GetTotalPriceByIDs(ctx context.Context, serviceIDs []uint) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&services.Service{}).
		Select("SUM(price)").
		Where("id IN ?", serviceIDs).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
