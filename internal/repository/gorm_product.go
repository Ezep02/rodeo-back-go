package repository

import (
	"context"
	"encoding/json"
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

func (r *GormProductRepository) List(ctx context.Context) ([]domain.Product, error) {
	var (
		appt         []domain.Product
		prodCacheKey string = "products"
	)

	// 1. Recuperar productos del cache
	servicesInCache, err := r.redis.Get(ctx, prodCacheKey).Result()

	if err == nil {
		json.Unmarshal([]byte(servicesInCache), &appt)
		return appt, nil
	}

	// 2. Si no estaba en el cache, realizar consulta sql
	if err := r.db.WithContext(ctx).Find(&appt).Error; err != nil {
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
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Omit("created_at").Save(Product).Error
}

func (r *GormProductRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

func (r *GormProductRepository) Popular(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product

	err := r.db.WithContext(ctx).
		Raw(`
			SELECT * FROM products
			ORDER BY
				CASE WHEN number_of_reviews > 0 THEN rating_sum / number_of_reviews ELSE 0 END DESC
			LIMIT 3
		`).
		Scan(&products).Error

	if err != nil {
		return nil, err
	}
	return products, nil
}
