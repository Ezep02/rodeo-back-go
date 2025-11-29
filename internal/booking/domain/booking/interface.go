package booking

import (
	"context"
	"time"
)

type BookingRepository interface {
	Create(ctx context.Context, b *Booking) error
	UpdateStatus(ctx context.Context, bookingID uint, status string) error
	UpdateSlot(ctx context.Context, bookingID, slotID uint) error
	Cancel(ctx context.Context, bookingID uint) error
	GetByID(ctx context.Context, bookingID uint) (*Booking, error)
	StartBookingCleanupJob(interval time.Duration)
	MarkAsPaid(ctx context.Context, bookingID uint) error
	MarkAsRejected(ctx context.Context, bookingID uint) error
	MarkAsRescheduled(ctx context.Context, bookingID uint) error
	Upcoming(ctx context.Context, barberID uint, date time.Time, status string) ([]Booking, error)
	GetByUserID(ctx context.Context, userID uint, offset int64) ([]Booking, error)
	StatsByBarberID(ctx context.Context, barberID uint) (*BookingStats, error)
	AllPendingPayment(ctx context.Context) ([]Booking, error)
}
