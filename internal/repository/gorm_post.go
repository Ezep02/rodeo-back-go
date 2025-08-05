package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/domain"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormPostRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormPostRepo(db *gorm.DB, redis *redis.Client) domain.PostRepository {
	return &GormPostRepository{db, redis}
}

func (r *GormPostRepository) List(ctx context.Context, offset int) ([]domain.Post, error) {

	var (
		posts    []domain.Post
		cacheKey = fmt.Sprintf("posts:%d", offset)
	)

	// 1. Verificar si los posts están en cache
	cachedPosts, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Si están en caché, deserializar y retornar
		if err := json.Unmarshal([]byte(cachedPosts), &posts); err == nil {
			return posts, nil
		}
		log.Println("Error deserializing cached posts:", err)
	}

	// 2. Si no están en cache, consultar la base de datos
	if err := r.db.WithContext(ctx).Offset(offset).Limit(10).Find(&posts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
	}

	// 3. Si hay mas resultados, guardar en cache
	if len(posts) > 0 {
		data, err := json.Marshal(posts)
		if err != nil {
			log.Println("Error serializing posts for cache:", err)
		}

		if err := r.redis.Set(ctx, cacheKey, data, 2*time.Minute).Err(); err != nil {
			log.Println("Error setting posts in cache:", err)
		}
	}

	return posts, nil
}

func (r *GormPostRepository) Create(ctx context.Context, post *domain.Post) error {

	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		log.Println("Error creating post:", err)
	}

	// Invalidate cache after creating a new post
	cacheKey := fmt.Sprintf("posts:%d", 0)
	if err := r.redis.Del(ctx, cacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after post creation:", err)
	}

	return nil
}
func (r *GormPostRepository) Update(ctx context.Context, post *domain.Post, post_id uint) error {
	var (
		cacheKey = fmt.Sprintf("posts:%d", 0)
	)

	if err := r.db.WithContext(ctx).Model(&domain.Post{}).Where("id = ?", post_id).Updates(post).Error; err != nil {
		log.Println("Error updating post:", err)
		return err
	}

	// Invlidar el cache después de actualizar un post
	if err := r.redis.Del(ctx, cacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after post update:", err)
		return err
	}

	return nil
}

func (r *GormPostRepository) Delete(ctx context.Context, post_id uint) error {
	var (
		cacheKey = fmt.Sprintf("posts:%d", 0)
	)

	r.redis.Del(ctx, cacheKey)

	return r.db.WithContext(ctx).Delete(&domain.Post{}, post_id).Error
}

func (r *GormPostRepository) GetByID(ctx context.Context, id uint) (*domain.Post, error) {
	var post domain.Post
	if err := r.db.WithContext(ctx).First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (r *GormPostRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Post{}).Count(&count).Error; err != nil {
		log.Println("Error counting posts:", err)
		return 0, err
	}
	return count, nil
}
