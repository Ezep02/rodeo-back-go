package usecase

import (
	"context"
	"time"

	"github.com/ezep02/rodeo/internal/slots/domain"
)

type SlotUsecase struct {
	slotRepo domain.SlotRepository
}

func NewSlotUsecase(slotRepo domain.SlotRepository) *SlotUsecase {
	return &SlotUsecase{slotRepo}
}

func (s *SlotUsecase) CreateInBatches(ctx context.Context, slot *[]domain.Slot) error {
	return s.slotRepo.CreateInBatches(ctx, slot)
}

func (s *SlotUsecase) Update(ctx context.Context, slot *domain.Slot, id uint) error {
	return s.slotRepo.Update(ctx, slot, id)
}

func (s *SlotUsecase) GetByDateRange(ctx context.Context, barber_id uint, start, end time.Time) ([]domain.SlotWithStatus, error) {

	// 1. validar que exista un barber id
	return s.slotRepo.ListByDateRange(ctx, barber_id, start, end)
}
