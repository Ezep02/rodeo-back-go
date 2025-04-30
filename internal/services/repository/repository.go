package repository

import (
	"context"
	"encoding/json"
	"fmt"
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
func (r *ServiceRepository) GetServices(ctx context.Context, limit int, offset int) (*[]models.Service, error) {
	var services *[]models.Service

	servicesCacheKey := fmt.Sprintf("services:limit-%d:offset-%d", limit, offset)

	// recuperar servicios desde cache
	servicesInCache, cacheErr := r.RedisConnection.Get(ctx, servicesCacheKey).Result()

	// Datos en cache, recuperar
	if cacheErr == nil {
		json.Unmarshal([]byte(servicesInCache), &services)
		return services, nil
	}

	result := r.Connection.WithContext(ctx).Order("created_at desc").Offset(offset).Limit(limit).Find(&services)

	if result.Error != nil {
		return nil, result.Error
	}

	// Hacer caching de datos
	data, _ := json.Marshal(services)
	r.RedisConnection.Set(ctx, servicesCacheKey, data, 20*time.Minute)

	return services, nil
}

// devuelve una lista de servicios populares
func (r *ServiceRepository) GetPopularServices(ctx context.Context) ([]models.PopularServices, error) {

	var (
		monthlyPopularServices []models.PopularServices
		redisCacheKey          string = "customer:popular-services"
		statusApproved         string = "approved"
	)

	if cachedPopularServices, err := r.RedisConnection.Get(ctx, redisCacheKey).Result(); err == nil {
		json.Unmarshal([]byte(cachedPopularServices), &monthlyPopularServices)
		return monthlyPopularServices, nil
	}

	// veces elegido + cortes totales en el mes / 2
	err := r.Connection.WithContext(ctx).Raw(`
		SELECT 
			o.title,
			(COUNT(DISTINCT o.user_id) * 100.0 / total.total_users) AS total_avg
		FROM orders o
		CROSS JOIN (
			SELECT COUNT(DISTINCT user_id) AS total_users
			FROM orders
			WHERE mp_status = ?
			AND EXTRACT(MONTH FROM schedule_day_date) = EXTRACT(MONTH FROM CURRENT_DATE)
		) AS total
		WHERE o.mp_status = ?
		AND EXTRACT(MONTH FROM o.schedule_day_date) = EXTRACT(MONTH FROM CURRENT_DATE)
		GROUP BY o.title, total.total_users
		ORDER BY total_avg DESC
		LIMIT 3
	`, statusApproved, statusApproved).Scan(&monthlyPopularServices).Error

	if err != nil {
		return nil, err
	}

	//cachear la informacion
	if popularServicesBytes, _ := json.Marshal(monthlyPopularServices); popularServicesBytes != nil {
		r.RedisConnection.Set(ctx, redisCacheKey, popularServicesBytes, 30*time.Minute)
		return monthlyPopularServices, nil
	}

	return monthlyPopularServices, nil
}
