package service

import "context"

// TODO REEMPLAZAR POR SERVICES
type ServiceRepository interface {
	Create(ctx context.Context, data *Service) error
	Update(ctx context.Context, id uint, data *Service) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page int) ([]Service, error)
	GetByID(ctx context.Context, id uint) (*Service, error)
	GetServiceStats(ctx context.Context) (*ServiceStats, error)
	Popular(ctx context.Context) ([]Service, error)
	AddCategories(ctx context.Context, id uint, categories_ids []uint) error
	RemoveCategories(ctx context.Context, id uint, categories_ids []uint) error
}
