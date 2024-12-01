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
	return s.Srv_Repository.CreateNewService(ctx, service)
}

func (s *Srv_Service) GetServices(ctx context.Context, limit int, offset int) (*[]Service, error) {
	return s.Srv_Repository.GetAllServices(ctx, limit, offset)
}

func (s *Srv_Service) UpdateService(ctx context.Context, service *Service, id string) (*Service, error) {
	return s.Srv_Repository.UpdateServiceByID(ctx, service, id)
}
