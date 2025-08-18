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

type GormProductRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormProductRepo(db *gorm.DB, redis *redis.Client) domain.ProductRepository {
	return &GormProductRepository{db, redis}
}

func (r *GormProductRepository) Create(ctx context.Context, Product *domain.Product) error {
	return r.db.WithContext(ctx).Create(Product).Error
}

func (r *GormProductRepository) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	var appt domain.Product

	if err := r.db.WithContext(ctx).First(&appt, id).Error; err != nil {

		// Check if the error is a record not found error
		// If so, return a custom error indicating that the Product was not found
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &appt, nil
}

func (r *GormProductRepository) List(ctx context.Context, offset int) ([]domain.Product, error) {
	var (
		appt         []domain.Product
		prodCacheKey string = fmt.Sprintf("products-page:%d", offset)
	)

	// 1. Recuperar productos del cache
	servicesInCache, err := r.redis.Get(ctx, prodCacheKey).Result()

	if err == nil {
		json.Unmarshal([]byte(servicesInCache), &appt)
		return appt, nil
	}

	// 2. Si no estaba en el cache, realizar consulta sql
	if err := r.db.WithContext(ctx).Preload("Category").Where("promotion_end_date IS NULL OR promotion_end_date > NOW()").Offset(offset).Find(&appt).Error; err != nil {
		return nil, err
	}

	// 3. Cachear los datos recuperados
	data, err := json.Marshal(appt)
	if err != nil {
		log.Println("Error realizando cache de los productos")
	}
	r.redis.Set(ctx, prodCacheKey, data, 3*time.Minute)

	return appt, nil
}

func (r *GormProductRepository) Update(ctx context.Context, Product *domain.Product) error {

	var (
		prodCacheKey string = "products"
	)

	// Invalidate cache after updating a product
	if err := r.redis.Del(ctx, prodCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after product update:", err)
	}

	updates := map[string]any{
		"name":               Product.Name,
		"price":              Product.Price,
		"description":        Product.Description,
		"category_id":        Product.CategoryID,
		"preview_url":        Product.PreviewUrl,
		"promotion_discount": Product.PromotionDiscount,
		"promotion_end_date": Product.PromotionEndDate,
		"has_promotion":      Product.HasPromotion,
		"updated_at":         Product.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Model(&domain.Product{}).Where("id = ?", Product.ID).Updates(updates).Error; err != nil {
		log.Println("Error updating post:", err)
		return err
	}

	return nil
}

func (r *GormProductRepository) Delete(ctx context.Context, id uint) error {
	var (
		prodCacheKey string = "products"
	)

	// Invalidate cache after updating a product
	if err := r.redis.Del(ctx, prodCacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after product update:", err)
	}

	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

func (r *GormProductRepository) Popular(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product

	err := r.db.WithContext(ctx).
		Raw(`
		SELECT p.* 
		FROM products p
		JOIN appointment_products ap ON ap.product_id = p.id
		JOIN appointments a ON a.id = ap.appointment_id
		WHERE a.status <> 'cancelled'
		GROUP BY p.id
		ORDER BY COUNT(*) DESC
		LIMIT 3
	`).
		Scan(&products).Error

	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *GormProductRepository) Promotion(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product

	err := r.db.WithContext(ctx).
		Raw(`
		SELECT * FROM products 
		WHERE has_promotion = ? 
		AND promotion_end_date > NOW()
	`, true).Scan(&products).Error

	if err != nil {
		return nil, err
	}
	return products, nil
}
