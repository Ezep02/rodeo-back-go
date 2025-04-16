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

func (r *ServiceRepository) CreateNewService(ctx context.Context, service *models.Service) (*models.Service, error) {

	result := r.Connection.WithContext(ctx).Create(service)

	if result.Error != nil {
		return nil, result.Error
	}

	return &models.Service{
		Model:            service.Model,
		Title:            service.Title,
		Created_by_id:    service.Created_by_id,
		Description:      service.Description,
		Service_Duration: service.Service_Duration,
		Price:            service.Price,
	}, nil
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

// devuelve una lista de servicios activos y creados por el barbero
func (r *ServiceRepository) GetBarberServices(ctx context.Context, limit int, offset int, barberID int) (*[]models.Service, error) {
	var services *[]models.Service

	result := r.Connection.WithContext(ctx).Where("created_by_id = ?", barberID).Order("created_at desc").Limit(limit).Offset(offset).Find(&services)

	if result.Error != nil {
		return nil, result.Error
	}

	return services, nil
}

func (r *ServiceRepository) UpdateServiceByID(ctx context.Context, service *models.Service, id string) (*models.Service, error) {
	// Iniciar transacción
	tx := r.Connection.WithContext(ctx).Begin()

	if tx.Error != nil {
		log.Println("[UPDATE SHIFT] Error al iniciar la transacción")
		return nil, tx.Error
	}

	// Construir los datos a actualizar
	updatedService := &models.Service{
		Model:            service.Model,
		Created_by_id:    service.Created_by_id,
		Title:            service.Title,
		Price:            service.Price,
		Description:      service.Description,
		Service_Duration: service.Service_Duration,
	}

	// Ejecutar la actualización
	result := tx.Model(&models.Service{}).Where("id = ?", id).Updates(updatedService)

	if result.Error != nil {
		log.Println("[UPDATE SHIFT] Error al actualizar registro, realizando rollback:", result.Error)
		tx.Rollback()
		return nil, result.Error
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		log.Println("[UPDATE SHIFT] Error al confirmar la transacción:", err)
		return nil, err
	}

	return updatedService, nil
}

func (r *ServiceRepository) DeleteServiceByID(ctx context.Context, serviceID int) error {

	result := r.Connection.WithContext(ctx).Where("id = ?", serviceID).Delete(&models.Service{})

	if result.Error != nil {
		return result.Error
	}

	return nil
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

	err := r.Connection.WithContext(ctx).Raw(`
		SELECT 
			title, 
			COUNT(*) AS total_orders
		FROM orders 
		WHERE mp_status = ?
		AND EXTRACT(MONTH FROM schedule_day_date) = EXTRACT(MONTH FROM CURRENT_DATE)
		GROUP BY title
		ORDER BY total_orders DESC
		LIMIT 4
	`, statusApproved).Scan(&monthlyPopularServices).Error

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
