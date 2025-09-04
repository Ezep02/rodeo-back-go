package appointment

import (
	"context"
	"time"
)

type AppointmentRepository interface {
	Create(ctx context.Context, appt *Appointment) error
	Update(ctx context.Context, appt *Appointment, slot_id uint) error
	Delete(ctx context.Context, id uint) error
	ListByDateRange(ctx context.Context, start, end time.Time) ([]Appointment, error)
	GetByID(ctx context.Context, id uint) (*Appointment, error)
	GetByUserID(ctx context.Context, id uint, offset int) ([]Appointment, error)
}
