package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/catalog/domain/promotions"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormPromoRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormPromoRepo(db *gorm.DB, redis *redis.Client) promotions.PromoRepository {
	return &GormPromoRepository{db, redis}
}

func (r *GormPromoRepository) Create(ctx context.Context, data *promotions.Promotion) error {

	var (
		promoCacheKey string = fmt.Sprintf("promotion-page:%d", 0)
	)

	// Invalidate cache after creating a new Service
	if err := r.redis.Del(ctx, promoCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after Service creation:", err)
	}

	return r.db.WithContext(ctx).Create(data).Error
}

func (r *GormPromoRepository) ListByServiceId(ctx context.Context, id uint, page int) ([]promotions.Promotion, error) {

	var (
		promoCacheKey string = fmt.Sprintf("promotion-svc:%d-page:%d", id, page)
		listById      []promotions.Promotion
	)

	limit := 5
	offset := (page - 1) * limit

	// 1. Verificar del cache
	if cachedData, err := r.redis.Get(ctx, promoCacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cachedData), &listById); err == nil {
			return listById, nil
		}
		log.Println("Error deserializando cache:", err)
	}

	// 2. Recuperar de sql
	if err := r.db.WithContext(ctx).
		Where("service_id = ?", id).
		Offset(offset).
		Limit(limit).
		Find(&listById).Error; err != nil {
		return nil, err
	}

	return listById, nil
}

func (r *GormPromoRepository) Update(ctx context.Context, id uint, data *promotions.Promotion) error {

	var (
		promoCacheKey string = fmt.Sprintf("promotion-svc:%d-page:%d", 0, 0)
	)

	// Invalidate cache after updating a Service
	if err := r.redis.Del(ctx, promoCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after promotion update:", err)
	}

	updates := map[string]any{
		"discount":   data.Discount,
		"start_date": data.StartDate,
		"end_date":   data.EndDate,
		"type":       data.Type,
		"updated_at": time.Now(),
	}

	if err := r.db.WithContext(ctx).Model(&promotions.Promotion{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Println("Error updating promotion:", err)
		return err
	}

	return nil
}

func (r *GormPromoRepository) Delete(ctx context.Context, id uint) error {
	var (
		promoCacheKey string = fmt.Sprintf("promotion-svc:%d-page:%d", 0, 0)
	)

	// 1. Liberar datos de la cache
	if err := r.redis.Del(ctx, promoCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after promotion update:", err)
	}

	return r.db.WithContext(ctx).Delete(&promotions.Promotion{}, id).Error
}
