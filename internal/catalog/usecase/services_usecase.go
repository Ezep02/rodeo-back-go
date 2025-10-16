package usecase

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/catalog/domain/service"
)

type ServicesService struct {
	svcRepo service.ServiceRepository
}

func NewServiceService(svcRepo service.ServiceRepository) *ServicesService {
	return &ServicesService{svcRepo}
}

func (s *ServicesService) Create(ctx context.Context, service *service.Service) error {

	// 1. Validar que el producto tenga un nombre
	if service.Name == "" {
		return errors.New("el producto debe tener un nombre")
	}

	// 2. Validar que el producto tenga un precio positivo
	if service.Price <= 0 {
		return errors.New("el producto debe tener un precio mayor o igual a cero")
	}

	return s.svcRepo.Create(ctx, service)
}

func (s *ServicesService) GetByID(ctx context.Context, id uint) (*service.Service, error) {
	return s.svcRepo.GetByID(ctx, id)
}

func (s *ServicesService) ListAll(ctx context.Context, offset int) ([]service.Service, error) {
	return s.svcRepo.List(ctx, offset)
}

func (s *ServicesService) Update(ctx context.Context, id uint, service *service.Service) error {
	// 1. Validar que el servicio tenga un nombre
	if service.Name == "" {
		return errors.New("el servicio debe tener un nombre")
	}

	// 2. Validar que el servicio tenga un precio positivo
	if service.Price <= 0 {
		return errors.New("el servicio debe tener un precio mayor o igual a cero")
	}

	return s.svcRepo.Update(ctx, id, service)
}

func (s *ServicesService) Delete(ctx context.Context, id uint) error {

	if id == 0 {
		return errors.New("error recupeando servicio")
	}

	// 2. Eliminar el servicio
	return s.svcRepo.Delete(ctx, id)
}

func (s *ServicesService) Stats(ctx context.Context) (*service.ServiceStats, error) {
	return s.svcRepo.GetServiceStats(ctx)
}

func (s *ServicesService) AddCategories(ctx context.Context, id uint, categories_ids []uint) error {
	if id == 0 {
		return errors.New("no fue posible recuperar el servicio")
	}

	if len(categories_ids) == 0 {
		return errors.New("no se encontraron elementos para agregar")
	}

	return s.svcRepo.AddCategories(ctx, id, categories_ids)
}

func (s *ServicesService) RemoveCategories(ctx context.Context, id uint, categories_ids []uint) error {

	if id == 0 {
		return errors.New("no fue posible recuperar el servicio")
	}

	if len(categories_ids) == 0 {
		return errors.New("no se encontraron elementos para remover")
	}

	return s.svcRepo.RemoveCategories(ctx, id, categories_ids)
}

func (s *ServicesService) Popular(ctx context.Context) ([]service.Service, error) {
	return s.svcRepo.Popular(ctx)
}
