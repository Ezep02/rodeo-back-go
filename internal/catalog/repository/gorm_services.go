package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/catalog/domain/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormServiceRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormServiceRepo(db *gorm.DB, redis *redis.Client) service.ServiceRepository {
	return &GormServiceRepository{db, redis}
}

func (r *GormServiceRepository) Create(ctx context.Context, data *service.Service) error {

	var (
		prodCacheKey string = fmt.Sprintf("services-page:%d", 0)
	)

	// Invalidate cache after creating a new Service
	if err := r.redis.Del(ctx, prodCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after Service creation:", err)
	}

	return r.db.WithContext(ctx).Create(data).Error
}

func (r *GormServiceRepository) GetByID(ctx context.Context, id uint) (*service.Service, error) {
	var appt service.Service

	if err := r.db.WithContext(ctx).
		Preload("Medias").
		Preload("Categories").
		Preload("Promotions").
		First(&appt, id).Error; err != nil {
		return nil, err
	}

	return &appt, nil
}

func (r *GormServiceRepository) List(ctx context.Context, page int) ([]service.Service, error) {
	var svcList []service.Service
	cacheKey := fmt.Sprintf("services-page:%d", page)

	// Cache
	if cachedData, err := r.redis.Get(ctx, cacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cachedData), &svcList); err == nil {
			return svcList, nil
		}
		log.Println("Error deserializando cache:", err)
	}

	limit := 10
	offset := (page - 1) * limit

	// DB con Preload
	if err := r.db.WithContext(ctx).
		Preload("Medias").
		Preload("Categories").
		Preload("Promotions").
		Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&svcList).Error; err != nil {
		return nil, err
	}

	// Guardar cache
	if data, err := json.Marshal(svcList); err == nil {
		r.redis.Set(ctx, cacheKey, data, 3*time.Minute)
	} else {
		log.Println("Error cacheando Services:", err)
	}

	return svcList, nil
}

func (r *GormServiceRepository) Update(ctx context.Context, id uint, data *service.Service) error {

	var (
		prodCacheKey string = fmt.Sprintf("services-page:%d", 0)
	)

	// Invalidate cache after updating a Service
	if err := r.redis.Del(ctx, prodCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after Service update:", err)
	}

	updates := map[string]any{
		"name":        data.Name,
		"price":       data.Price,
		"description": data.Description,
		"preview_url": data.PreviewURL,
		"is_active":   data.IsActive,
		"updated_at":  time.Now(),
	}

	if err := r.db.WithContext(ctx).Model(&service.Service{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Println("Error updating post:", err)
		return err
	}

	return nil
}

func (r *GormServiceRepository) Delete(ctx context.Context, id uint) error {
	var (
		prodCacheKey string = fmt.Sprintf("services-page:%d", 0)
	)

	// Invalidate cache after updating a Service
	if err := r.redis.Del(ctx, prodCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after Service update:", err)
	}

	return r.db.WithContext(ctx).Delete(&service.Service{}, id).Error
}

func (r *GormServiceRepository) GetServiceStats(ctx context.Context) (*service.ServiceStats, error) {

	var (
		journeysKey string = "journeys-stats"
		stats       service.ServiceStats
	)

	//1. Intentar recuperar del cache
	if infoInCache, err := r.redis.Get(ctx, journeysKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(infoInCache), &stats); err == nil {
			return &stats, nil
		}
	}

	// Total journeys
	if err := r.db.WithContext(ctx).Model(&service.Service{}).Count(&stats.TotalServices).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *GormServiceRepository) AddCategories(ctx context.Context, serviceID uint, categories_ids []uint) error {
	// Invalidate cache
	cacheKey := fmt.Sprintf("services-page:%d", 0)
	if err := r.redis.Del(ctx, cacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after Service creation:", err)
	}

	// Generar los pares correctos (service_id = serviceID)
	type Relation struct {
		ServiceID  uint `gorm:"column:service_id"`
		CategoryID uint `gorm:"column:category_id"`
	}

	relations := make([]Relation, 0, len(categories_ids))
	for _, catID := range categories_ids {
		log.Println("[ID TO ADD]", catID)
		relations = append(relations, Relation{ServiceID: serviceID, CategoryID: catID})
	}

	// Batch insert con GORM
	err := r.db.WithContext(ctx).
		Table("service_categories").
		Clauses(clause.Insert{Modifier: "IGNORE"}).
		Create(&relations).Error

	if err != nil {
		log.Printf("Error batch insert service-categories: %v", err)
		return err
	}

	return nil
}

func (r *GormServiceRepository) RemoveCategories(ctx context.Context, id uint, categoryIDs []uint) error {

	// Invalidate cache
	cacheKey := fmt.Sprintf("services-page:%d", 0)
	if err := r.redis.Del(ctx, cacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after removing categories:", err)
	}

	err := r.db.WithContext(ctx).
		Exec(`DELETE FROM service_categories WHERE service_id = ? AND category_id IN ?`, id, categoryIDs).
		Error

	if err != nil {
		log.Printf("Error removing categories %v from service %d: %v", categoryIDs, id, err)
		return err
	}

	return nil
}

// TODO: verificar utilidad
func (r *GormServiceRepository) Popular(ctx context.Context) ([]service.Service, error) {
	var Services []service.Service

	err := r.db.WithContext(ctx).
		Raw(`
		SELECT p.* 
		FROM Services p
		JOIN appointment_Services ap ON ap.Service_id = p.id
		JOIN appointments a ON a.id = ap.appointment_id
		WHERE a.status <> 'cancelled'
		GROUP BY p.id
		ORDER BY COUNT(*) DESC
		LIMIT 3
	`).
		Scan(&Services).Error

	if err != nil {
		return nil, err
	}
	return Services, nil
}
