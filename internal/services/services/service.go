package services

import (
	"context"

	"github.com/ezep02/rodeo/internal/services/models"
	"github.com/ezep02/rodeo/internal/services/repository"
)

type Srv_Service struct {
	Srv_Repository *repository.ServiceRepository
}

func NewSrvRepository(srv_r *repository.ServiceRepository) *Srv_Service {
	return &Srv_Service{
		Srv_Repository: srv_r,
	}
}

func (s *Srv_Service) CreateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	return s.Srv_Repository.CreateNewService(ctx, service)
}

func (s *Srv_Service) GetServices(ctx context.Context, limit int, offset int) (*[]models.Service, error) {
	return s.Srv_Repository.GetServices(ctx, limit, offset)
}

func (s *Srv_Service) GetBarberServices(ctx context.Context, limit int, offset int, barberID int) (*[]models.Service, error) {
	return s.Srv_Repository.GetBarberServices(ctx, limit, offset, barberID)
}

func (s *Srv_Service) UpdateService(ctx context.Context, service *models.Service, id string) (*models.Service, error) {
	return s.Srv_Repository.UpdateServiceByID(ctx, service, id)
}

func (s *Srv_Service) DeleteServiceByID(ctx context.Context, serviceID int) error {
	return s.Srv_Repository.DeleteServiceByID(ctx, serviceID)
}

// GetPopularServices obtiene los servicios populares
func (s *Srv_Service) GetPopularServices(ctx context.Context) ([]models.PopularServices, error) {
	return s.Srv_Repository.GetPopularServices(ctx)
}
