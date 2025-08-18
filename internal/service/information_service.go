package service

import (
	"context"

	"github.com/ezep02/rodeo/internal/domain"
)

type InformationService struct {
	informationRepo domain.InformationRepository
}

func NewInfoRepository(infoRepo domain.InformationRepository) *InformationService {
	return &InformationService{infoRepo}
}

func (s *InformationService) Information(ctx context.Context) (*domain.BarberInformation, error) {
	return s.informationRepo.BarberInformation(ctx)
}
