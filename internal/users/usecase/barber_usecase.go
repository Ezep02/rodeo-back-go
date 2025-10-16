package usecase

import (
	"context"

	"github.com/ezep02/rodeo/internal/users/domain/barber"
)

type BarberService struct {
	barberRepo barber.BarberRepository
}

func NewBarberService(barberRepo barber.BarberRepository) *BarberService {
	return &BarberService{barberRepo}
}

func (s *BarberService) GetByID(ctx context.Context, id uint) (*barber.Barber, error) {
	return s.barberRepo.GetByID(ctx, id)
}

func (s *BarberService) List(ctx context.Context) ([]barber.BarberWithUser, error) {
	return s.barberRepo.List(ctx)
}
