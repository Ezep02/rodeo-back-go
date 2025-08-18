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

type GormInfoRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormInfoRepo(db *gorm.DB, redis *redis.Client) domain.InformationRepository {
	return &GormInfoRepository{db, redis}
}

func (r *GormInfoRepository) BarberInformation(ctx context.Context) (*domain.BarberInformation, error) {

	var (
		info             *domain.BarberInformation
		bussinessInfoKey string = "elrodeoinfokey"
	)

	// 1. Recuperar productos del cache
	infoInCache, err := r.redis.Get(ctx, bussinessInfoKey).Result()

	if err == nil {
		json.Unmarshal([]byte(infoInCache), &info)
		return info, nil
	}

	if err := r.db.WithContext(ctx).Raw(`
		SELECT 
			(SELECT COUNT(*) FROM users) AS member,
			(SELECT COUNT(*) FROM appointments) AS total_appointment,
			COALESCE((SELECT AVG(rating) FROM reviews), 0) AS promedy
		`).Scan(&info).Error; err != nil {
		return nil, err
	}

	// 3. Cachear los datos recuperados
	data, err := json.Marshal(info)
	if err != nil {
		log.Println("Error realizando cache de los productos")
	}
	r.redis.Set(ctx, bussinessInfoKey, data, 20*time.Minute)

	return info, nil
}
