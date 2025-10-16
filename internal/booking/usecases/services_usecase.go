package usecases

import (
	"context"

	"github.com/ezep02/rodeo/internal/booking/domain/services"
)

type ServicesService struct {
	svcRepo services.ServicesRepository
}

// Constructor
func NewServicesService(svcRepo services.ServicesRepository) *ServicesService {
	return &ServicesService{svcRepo}
}

func (s *ServicesService) GetByID(ctx context.Context, id uint) (*services.Service, error) {
	return s.svcRepo.GetByID(ctx, id)
}

func (s *ServicesService) GetTotalPriceByIDs(ctx context.Context, serviceIDs []uint) (float64, error) {
	return s.svcRepo.GetTotalPriceByIDs(ctx, serviceIDs)
}
