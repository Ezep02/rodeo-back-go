package repository

import (
	"context"
	"log"

	"github.com/ezep02/rodeo/internal/catalog/domain/media"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormMediaRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormMediaRepo(db *gorm.DB, redis *redis.Client) media.MediaRepository {
	return &GormMediaRepository{db, redis}
}

func (r *GormMediaRepository) Create(ctx context.Context, data *media.Medias) error {

	return r.db.WithContext(ctx).Create(data).Error
}

func (r *GormMediaRepository) Delete(ctx context.Context, id uint) error {

	return r.db.WithContext(ctx).Delete(&media.Medias{}, id).Error

}

func (r *GormMediaRepository) Update(ctx context.Context, id uint, data *media.Medias) error {
	log.Printf("Updating media ID=%d, URL=%s\n", id, data.URL)

	result := r.db.WithContext(ctx).
		Model(&media.Medias{}).
		Where("id = ?", id).
		Update("url", data.URL)

	if result.Error != nil {
		log.Println("Error updating media:", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Println("No media record found for ID:", id)
	}

	return nil
}

func (r *GormMediaRepository) ListByServiceId(ctx context.Context, id uint) ([]media.Medias, error) {

	var (
		listById []media.Medias
	)

	// 2. Recuperar de sql
	if err := r.db.WithContext(ctx).
		Where("service_id = ?", id).
		Find(&listById).Error; err != nil {
		return nil, err
	}

	return listById, nil
}
