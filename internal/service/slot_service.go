package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
)

type SlotService struct {
	slotRepo domain.SlotRepository
}

func NewSlotService(slotRepo domain.SlotRepository) *SlotService {
	return &SlotService{slotRepo}
}

func (s *SlotService) Create(ctx context.Context, slot *[]domain.Slot) error {

	// 1. Validar que la cita nos sea en el pasado
	// for _, s := range *slot {
	// 	if s.Date.Before(time.Now()) {
	// 		return fmt.Errorf("no podes crear un turno en el pasado %s", s.Date)
	// 	}
	// }

	return s.slotRepo.Create(ctx, slot)
}

func (s *SlotService) Update(ctx context.Context, updatedSlot *[]domain.Slot) error {

	// 1. Validar que haya horarios
	if len(*updatedSlot) == 0 {
		return errors.New("no hay horarios para actualizar")
	}

	existing, err := s.slotRepo.ListByDate(ctx, (*updatedSlot)[0].Date)
	if err != nil {
		return domain.ErrNotFound
	}

	// 2. Validar que la cita nos sea en el pasado
	for _, slot := range *updatedSlot {

		if slot.Date.Before(time.Now()) {
			return fmt.Errorf("no podes instanciar un turno en el pasado %s", slot.Date)
		}
	}

	for _, existingSlot := range existing {
		// Si el horario es igual y el ID es diferente, significa que ya existe otro
		if s.slotOverlap(existingSlot, *updatedSlot) {
			return errors.New("ya existe una cita en esa fecha y hora")
		}
	}

	return s.slotRepo.Update(ctx, updatedSlot)
}

func (s *SlotService) Delete(ctx context.Context, slots []domain.Slot) error {
	// 1. Verificar que el producto exista
	if len(slots) == 0 {
		return errors.New("no hay turnos para eliminar")
	}

	// 2. Eliminar el producto
	return s.slotRepo.Delete(ctx, &slots)
}

func (s *SlotService) ListByDate(ctx context.Context, date time.Time) ([]domain.Slot, error) {
	return s.slotRepo.ListByDate(ctx, date)
}

func (s *SlotService) ListByDateRange(ctx context.Context, start time.Time, end time.Time) ([]domain.Slot, error) {
	return s.slotRepo.ListByDateRange(ctx, start, end)
}

func (s *SlotService) GetByID(ctx context.Context, id uint) (*domain.Slot, error) {
	return s.slotRepo.GetByID(ctx, id)
}

func (s *SlotService) slotOverlap(existing domain.Slot, updatedSlot []domain.Slot) bool {
	// Compara si las citas se superponen
	for _, slot := range updatedSlot {
		if slot.Time == existing.Time {
			return true
		}
	}
	return false
}
