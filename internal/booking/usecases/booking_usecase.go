package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ezep02/rodeo/internal/booking/domain/booking"
)

type BookingService struct {
	bookingRepo booking.BookingRepository
}

func NewBookingService(bookingRepo booking.BookingRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo}
}

func (s *BookingService) CreateBooking(ctx context.Context, b *booking.Booking) error {
	if b == nil {
		return errors.New("booking es nil")
	}
	return s.bookingRepo.Create(ctx, b)
}

func (s *BookingService) UpdateBooking(ctx context.Context, b *booking.Booking) error {
	if b == nil {
		return errors.New("booking es nil")
	}
	return s.bookingRepo.Update(ctx, b)
}

func (s *BookingService) UpdateBookingStatus(ctx context.Context, bookingID uint, status string) error {
	if status == "" {
		return errors.New("status no puede ser vacío")
	}
	return s.bookingRepo.UpdateStatus(ctx, bookingID, status)
}

func (s *BookingService) GetBookingByID(ctx context.Context, bookingID uint) (*booking.Booking, error) {
	return s.bookingRepo.GetByID(ctx, bookingID)
}

func (s *BookingService) MarkAsPaid(ctx context.Context, bookingID uint) error {

	if bookingID == 0 {
		return errors.New("el id de la reserva no puede ser nulo")
	}

	return s.bookingRepo.MarkAsPaid(ctx, bookingID)
}

func (s *BookingService) MarkAsRejected(ctx context.Context, bookingID uint) error {

	if bookingID == 0 {
		return errors.New("el id de la reserva no puede ser nulo")
	}

	return s.bookingRepo.MarkAsRejected(ctx, bookingID)
}

// PARA BARBEROS
func (s *BookingService) Upcoming(ctx context.Context, barberID uint, date time.Time, status string) ([]booking.Booking, error) {

	if barberID <= 0 {
		return nil, errors.New("ingrese un id valido")
	}

	return s.bookingRepo.Upcoming(ctx, barberID, date, status)
}

func (s *BookingService) StatsByBarberID(ctx context.Context, barberID uint) (*booking.BookingStats, error) {
	if barberID == 0 {
		return nil, errors.New("el id del usuario no puede ser nulo")
	}

	return s.bookingRepo.StatsByBarberID(ctx, barberID)
}

func (s *BookingService) AllPendingPayment(ctx context.Context) ([]booking.Booking, error) {
	return s.bookingRepo.AllPendingPayment(ctx)
}
