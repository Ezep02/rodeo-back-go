package services

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type ServiceRepository struct {
	Connection *gorm.DB
}

func NewServiceRepository(DATABASE *gorm.DB) *ServiceRepository {
	return &ServiceRepository{
		Connection: DATABASE,
	}
}

func (r *ServiceRepository) CreateNewService(ctx context.Context, service *Service) (*Service, error) {

	result := r.Connection.WithContext(ctx).Create(service)

	if result.Error != nil {
		return nil, result.Error
	}

	return &Service{
		Model:            service.Model,
		Title:            service.Title,
		Created_by_id:    service.Created_by_id,
		Description:      service.Description,
		Service_Duration: service.Service_Duration,
		Price:            service.Price,
		Preview_url:      service.Preview_url,
	}, nil
}

// para usuarios
func (r *ServiceRepository) GetServices(ctx context.Context, limit int, offset int) (*[]Service, error) {
	var services *[]Service

	result := r.Connection.WithContext(ctx).Order("created_at desc").Offset(offset).Limit(limit).Find(&services)

	if result.Error != nil {
		return nil, result.Error
	}

	return services, nil
}

func (r *ServiceRepository) GetBarberServices(ctx context.Context, limit int, offset int, barberID int) (*[]Service, error) {
	var services *[]Service

	result := r.Connection.WithContext(ctx).Where("created_by_id = ?", barberID).Order("created_at desc").Limit(limit).Offset(offset).Find(&services)

	if result.Error != nil {
		return nil, result.Error
	}

	return services, nil
}

func (r *ServiceRepository) UpdateServiceByID(ctx context.Context, service *Service, id string) (*Service, error) {
	// Iniciar transacción
	tx := r.Connection.WithContext(ctx).Begin()

	if tx.Error != nil {
		log.Println("[UPDATE SHIFT] Error al iniciar la transacción")
		return nil, tx.Error
	}

	log.Println("service:", service)

	// Construir los datos a actualizar
	updatedService := &Service{
		Model:            service.Model,
		Created_by_id:    service.Created_by_id,
		Title:            service.Title,
		Price:            service.Price,
		Description:      service.Description,
		Service_Duration: service.Service_Duration,
	}

	// Ejecutar la actualización
	result := tx.Model(&Service{}).Where("id = ?", id).Updates(updatedService)

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

	result := r.Connection.WithContext(ctx).Where("id = ?", serviceID).Delete(&Service{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
