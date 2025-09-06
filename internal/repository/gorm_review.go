package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/domain/review"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormReviewRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormReviewRepo(db *gorm.DB, redis *redis.Client) review.ReviewRepository {
	return &GormReviewRepository{db, redis}
}

func (r *GormReviewRepository) Create(ctx context.Context, review *review.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *GormReviewRepository) Update(ctx context.Context, review *review.Review) error {
	return nil
}
func (r *GormReviewRepository) Delete(ctx context.Context, id, user_id uint) error {
	var (
		revCacheKey string = fmt.Sprintf("review:user:%d-offset:%d", user_id, 0)
	)

	// 1. Eliminar cache
	if err := r.redis.Del(ctx, revCacheKey).Err(); err != nil {
		log.Println("Error deleting user from cache:", err)
	}

	if err := r.db.WithContext(ctx).Delete(&review.Review{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormReviewRepository) List(ctx context.Context, offset int) ([]review.ReviewDetail, error) {
	var (
		reviews     []review.ReviewDetail
		revCacheKey string = fmt.Sprintf("review:all-offset:%d", offset)
	)

	// 1. Recuperar productos del cache
	servicesInCache, err := r.redis.Get(ctx, revCacheKey).Result()

	if err == nil {
		json.Unmarshal([]byte(servicesInCache), &reviews)
		return reviews, nil
	}

	if err := r.db.WithContext(ctx).
		Table("reviews as r").
		Select(`
			r.id as review_id, 
			r.rating, 
			r.comment, 
			r.created_at,
            a.id as appointment_id, 
			a.client_name, 
			a.client_surname, 
			a.status as appointment_status,
            u.id as user_id, 
			u.name as user_name, 
			u.email, 
			u.avatar`).
		Joins("JOIN appointments a ON r.appointment_id = a.id").
		Joins("JOIN users u ON a.user_id = u.id").
		Order("r.rating DESC").
		Offset(offset).
		Limit(10).
		Scan(&reviews).Error; err != nil {
		return nil, err
	}

	// 3. Cachear resultados opcionalmente
	data, _ := json.Marshal(reviews)
	r.redis.Set(ctx, revCacheKey, data, time.Hour)

	return reviews, nil
}

func (r *GormReviewRepository) ListByProductID(ctx context.Context, productID uint) ([]review.Review, error) {
	var reviews []review.Review

	if err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Find(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *GormReviewRepository) ListByUserID(ctx context.Context, userID uint, offset int) ([]review.ReviewDetail, error) {
	var (
		reviews     []review.ReviewDetail
		revCacheKey string = fmt.Sprintf("review:user:%d-offset:%d", userID, offset) // cache por userID
	)

	log.Println("[DEBUG] revCacheKey:", revCacheKey)
	// 1. Recuperar citas del cache
	servicesInCache, err := r.redis.Get(ctx, revCacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(servicesInCache), &reviews)
		return reviews, nil
	}

	// 2. Consulta a la base de datos
	if err := r.db.WithContext(ctx).
		Table("reviews as r").
		Select(`
			r.id as review_id, 
			r.rating, 
			r.comment, 
			r.created_at,
            a.id as appointment_id, 
			a.client_name, 
			a.client_surname, 
			a.status as appointment_status,
            u.id as user_id, 
			u.username,
			u.email, 
			u.avatar`).
		Joins("JOIN appointments a ON r.appointment_id = a.id").
		Joins("JOIN users u ON a.user_id = u.id").
		Where("a.user_id = ?", userID). // Aquí se filtra por userID
		Offset(offset).
		Limit(10).
		Find(&reviews).Error; err != nil {
		return nil, err
	}

	// 3. Cachear resultados opcionalmente
	data, _ := json.Marshal(reviews)
	r.redis.Set(ctx, revCacheKey, data, 5*time.Minute)

	return reviews, nil
}

func (r *GormReviewRepository) ReviewRatingStats(ctx context.Context) (*review.ReviewRatingStats, error) {
	var (
		stats       review.ReviewRatingStats
		ratingCount map[int]int = make(map[int]int)
	)

	// 1. Consulta a la base de datos
	type RatingCountResult struct {
		Rating int
		Count  int
	}

	var results []RatingCountResult
	if err := r.db.WithContext(ctx).
		Model(&review.Review{}).
		Select("rating, COUNT(*) as count").
		Group("rating").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	totalReviews := 0
	sumRatings := 0
	for _, res := range results {
		ratingCount[res.Rating] = res.Count
		totalReviews += res.Count
		sumRatings += res.Rating * res.Count
	}

	averageRating := 0.0
	if totalReviews > 0 {
		averageRating = float64(sumRatings) / float64(totalReviews)
	}

	stats = review.ReviewRatingStats{
		TotalReviews:  totalReviews,
		AverageRating: averageRating,
		RatingCount:   ratingCount,
	}
	return &stats, nil
}

func (r *GormReviewRepository) GetByID(ctx context.Context, id uint) (*review.Review, error) {
	var rev review.Review
	if err := r.db.WithContext(ctx).First(&rev, id).Error; err != nil {
		return nil, err
	}
	return &rev, nil
}
