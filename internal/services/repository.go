package services

import (
	"context"

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

func (r *ServiceRepository) CreateService(ctx context.Context, service *Service) (*Service, error) {

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
	}, nil
}

func (r *ServiceRepository) GetAllServices(ctx context.Context) ([]*Service, error) {
	var services []Service
	result := r.Connection.WithContext(ctx).Find(&services)

	if result.Error != nil {
		return nil, result.Error
	}

	servicePtrs := make([]*Service, len(services))

	// se recorre la respuesta y se cargan a servicesPtrs como punteros
	for i := range services {
		servicePtrs[i] = &services[i]
	}

	return servicePtrs, nil
}

func (r *ServiceRepository) UpdateService() {

}

func (r *ServiceRepository) DeleteServiceByID() {

}
