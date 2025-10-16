package domain

import (
	"context"
	"time"
)

type SlotRepository interface {
	CreateInBatches(ctx context.Context, slot *[]Slot) error
	Update(ctx context.Context, slot *Slot, slot_id uint) error
	// Delete(ctx context.Context, id uint) error
	ListByDateRange(ctx context.Context, barber_id uint, start, end time.Time) ([]SlotWithStatus, error)
	// GetByID(ctx context.Context, id uint) (*Slot, error)
	// GetByUserID(ctx context.Context, id uint, offset int) ([]Slot, error)
}
