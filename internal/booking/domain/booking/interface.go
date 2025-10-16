package booking

import (
	"context"
	"time"
)

type BookingRepository interface {
	Create(ctx context.Context, b *Booking) error
	UpdateStatus(ctx context.Context, bookingID uint, status string) error
	Update(ctx context.Context, b *Booking) error
	GetByID(ctx context.Context, bookingID uint) (*Booking, error)
	StartBookingCleanupJob(interval time.Duration)
	MarkAsPaid(ctx context.Context, bookingID uint) error
	Upcoming(ctx context.Context, barberID uint, date time.Time, status string) ([]Booking, error)

	StatsByBarberID(ctx context.Context, barberID uint) (*BookingStats, error)
}
