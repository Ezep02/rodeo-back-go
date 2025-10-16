package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ezep02/rodeo/internal/booking/domain/payments"
)

type PaymentService struct {
	paymentRepo payments.PaymentRepository
}

// Constructor
func NewPaymentService(paymentRepo payments.PaymentRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo}
}

func (s *PaymentService) CreatePayment(ctx context.Context, p *payments.Payment) error {
	if p == nil {
		return errors.New("payment es nil")
	}
	return s.paymentRepo.Create(ctx, p)
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, paymentID uint) (*payments.Payment, error) {
	return s.paymentRepo.GetByID(ctx, paymentID)
}

func (s *PaymentService) GetPaymentsByBookingID(ctx context.Context, bookingID uint) ([]payments.Payment, error) {
	return s.paymentRepo.GetByBookingID(ctx, bookingID)
}

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, paymentID uint, status string, paidAt *time.Time) error {
	if status == "" {
		return errors.New("status no puede ser vacío")
	}
	return s.paymentRepo.UpdateStatus(ctx, paymentID, status, paidAt)
}

func (s *PaymentService) UpdatePayment(ctx context.Context, p *payments.Payment) error {
	if p == nil {
		return errors.New("payment es nil")
	}
	return s.paymentRepo.Update(ctx, p)
}

func (s *PaymentService) MarkAsPaid(ctx context.Context, paymentID uint, mpPaymentID string) error {
	if paymentID == 0 {
		return errors.New("el id del pago no puede ser nulo")
	}

	return s.paymentRepo.MarkAsPaid(ctx, paymentID, mpPaymentID)
}
