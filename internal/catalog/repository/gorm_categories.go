package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/catalog/domain/categorie"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormCategorieRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormCategorieRepo(db *gorm.DB, redis *redis.Client) categorie.CategorieRepository {
	return &GormCategorieRepository{db, redis}
}

func (r *GormCategorieRepository) Create(ctx context.Context, categorie *categorie.Categorie) error {

	var (
		promoCacheKey string = fmt.Sprintf("categorie-page:%d", 0)
	)

	// Invalidate cache
	if err := r.redis.Del(ctx, promoCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after Service creation:", err)
	}

	return r.db.WithContext(ctx).Create(categorie).Error
}

func (r *GormCategorieRepository) Update(ctx context.Context, id uint, data *categorie.Categorie) error {

	var (
		promoCacheKey string = fmt.Sprintf("categorie-page:%d", 0)
	)

	// Invalidate cache
	if err := r.redis.Del(ctx, promoCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after Service creation:", err)
	}

	updates := map[string]any{
		"name":        data.Name,
		"preview_url": data.PreviewURL,
		"updated_at":  time.Now(),
	}

	if err := r.db.WithContext(ctx).Model(&categorie.Categorie{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Println("Error updating post:", err)
		return err
	}

	return nil
}

func (r *GormCategorieRepository) Delete(ctx context.Context, id uint) error {
	var (
		promoCacheKey string = fmt.Sprintf("categorie-page:%d", 0)
	)

	// Invalidate cache
	if err := r.redis.Del(ctx, promoCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after Service creation:", err)
	}

	return r.db.WithContext(ctx).Delete(&categorie.Categorie{}, id).Error
}

func (r *GormCategorieRepository) List(ctx context.Context) ([]categorie.Categorie, error) {
	var (
		cacheKey   = "categorie-page"
		categories []categorie.Categorie
	)
	// Cache
	if cachedData, err := r.redis.Get(ctx, cacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cachedData), &categories); err == nil {
			return categories, nil
		}
		log.Println("Error deserializando cache:", err)
	}

	// limit := 10
	// offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).
		Find(&categories).Error; err != nil {
		return nil, err
	}

	// Guardar cache
	if data, err := json.Marshal(categories); err == nil {
		r.redis.Set(ctx, cacheKey, data, 3*time.Minute)
	} else {
		log.Println("Error cacheando Services:", err)
	}

	return categories, nil
}

func (r *GormCategorieRepository) GetByID(ctx context.Context, id uint) (*categorie.Categorie, error) {
	var categorie categorie.Categorie
	if err := r.db.WithContext(ctx).First(&categorie, id).Error; err != nil {
		return nil, err
	}

	return &categorie, nil
}
