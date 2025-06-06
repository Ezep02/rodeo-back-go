package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/services/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	Connection      *gorm.DB
	RedisConnection *redis.Client
}

func NewServiceRepository(DATABASE *gorm.DB, REDIS *redis.Client) *ServiceRepository {
	return &ServiceRepository{
		Connection:      DATABASE,
		RedisConnection: REDIS,
	}
}

// recupera los servicios activos para el usuario
func (r *ServiceRepository) GetServices(ctx context.Context, limit int, offset int) (*[]models.CustomerServices, error) {
	var services *[]models.CustomerServices

	servicesCacheKey := fmt.Sprintf("services:limit-%d:offset-%d", limit, offset)

	// recuperar servicios desde cache
	servicesInCache, cacheErr := r.RedisConnection.Get(ctx, servicesCacheKey).Result()

	// Datos en cache, recuperar
	if cacheErr == nil {
		json.Unmarshal([]byte(servicesInCache), &services)
		return services, nil
	}

	err := r.Connection.WithContext(ctx).Raw(`
	SELECT 
		s.id,
		s.title,
		s.description,
		s.price,
		s.service_duration,
		s.category,
		s.preview_url,
		s.created_by_id
		IFNULL(AVG(r.rating), 0) AS rating,
		COUNT(r.rating) AS reviews
	FROM services s
		LEFT JOIN orders o ON o.service_id = s.id AND o.deleted_at IS NULL
		LEFT JOIN reviews r ON r.order_id = o.id AND r.schedule_id = o.shift_id AND r.user_id = o.user_id
	WHERE s.deleted_at IS NULL
	GROUP BY s.id
	ORDER BY s.created_at DESC
	LIMIT ? OFFSET ?`, limit, offset).Scan(&services).Error

	if err != nil {
		log.Println("Algo no fue bien recuperando los servicios")
		return nil, err
	}

	// Hacer caching de datos
	data, _ := json.Marshal(services)
	r.RedisConnection.Set(ctx, servicesCacheKey, data, 3*time.Minute)

	return services, nil
}

// devuelve una lista de servicios populares
func (r *ServiceRepository) GetPopularServices(ctx context.Context) ([]models.PopularServices, error) {

	var (
		popularServices []models.PopularServices
		redisCacheKey   string = "customer:popular-services"
		statusApproved  string = "approved"
	)

	if cachedPopularServices, err := r.RedisConnection.Get(ctx, redisCacheKey).Result(); err == nil {
		json.Unmarshal([]byte(cachedPopularServices), &popularServices)
		return popularServices, nil
	}

	// veces elegido + cortes totales en el mes / 2
	err := r.Connection.WithContext(ctx).Raw(`
		SELECT 
			s.title,
			s.description,
			s.service_duration,
			s.price,
			s.preview_url,
			COUNT(DISTINCT o.user_id) * 100.0 / total.total_users AS total_avg,
			AVG(r.rating) AS rating
		FROM services s
		JOIN orders o ON o.service_id = s.id AND o.mp_status = ?
		LEFT JOIN reviews r ON r.order_id = o.id AND r.schedule_id = o.shift_id AND r.user_id = o.user_id
		CROSS JOIN (
			SELECT COUNT(DISTINCT user_id) AS total_users
			FROM orders
			WHERE mp_status = ?
		) AS total
		GROUP BY s.id, s.title, s.description, s.service_duration, s.price, s.preview_url, total.total_users
		ORDER BY total_avg DESC
		LIMIT 3`, statusApproved, statusApproved).Scan(&popularServices).Error

	if err != nil {
		return nil, err
	}

	//cachear la informacion
	if popularServicesBytes, _ := json.Marshal(popularServices); popularServicesBytes != nil {
		r.RedisConnection.Set(ctx, redisCacheKey, popularServicesBytes, 5*time.Minute)
		return popularServices, nil
	}

	return popularServices, nil
}
