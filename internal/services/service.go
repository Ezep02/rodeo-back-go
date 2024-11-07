package services

import "context"

type Srv_Service struct {
	Srv_Repository *ServiceRepository
}

func NewSrvRepository(srv_r *ServiceRepository) *Srv_Service {
	return &Srv_Service{
		Srv_Repository: srv_r,
	}
}

func (s *Srv_Service) CreateService(ctx context.Context, service *Service) (*Service, error) {
	return s.Srv_Repository.CreateService(ctx, service)
}

func (s *Srv_Service) GetServices(ctx context.Context) ([]*Service, error) {
	return s.Srv_Repository.GetAllServices(ctx)
}
