package payments

import (
	"context"
	"time"
)

type PaymentRepository interface {
	// Crear un nuevo pago (puede ser seña o total)
	Create(ctx context.Context, payment *Payment) error

	// Obtener todos los pagos de un booking
	GetByBookingID(ctx context.Context, bookingID uint) (*Payment, error)

	// Actualizar el status del pago (pendiente, aprobado, rechazado, reembolsado)
	UpdateStatus(ctx context.Context, paymentID uint, status string, paidAt *time.Time) error

	// Actualizar datos del pago, como amount, method, url o MercadoPagoID
	Update(ctx context.Context, payment *Payment) error

	MarkAsPaid(ctx context.Context, paymentID uint, mpPaymentID string) error
}
