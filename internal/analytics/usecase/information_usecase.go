package usecase

import (
	"context"

	"github.com/ezep02/rodeo/internal/analytics/domain/information"
)

type InformationService struct {
	informationRepo information.InformationRepository
}

func NewInfoService(infoRepo information.InformationRepository) *InformationService {
	return &InformationService{infoRepo}
}

func (s *InformationService) Information(ctx context.Context) (*information.BarberInformation, error) {
	return s.informationRepo.BarberInformation(ctx)
}
