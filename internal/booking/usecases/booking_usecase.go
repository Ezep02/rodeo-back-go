package usecases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ezep02/rodeo/internal/booking/domain/booking"
	"github.com/ezep02/rodeo/internal/booking/domain/payments"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type BookingService struct {
	bookingRepo booking.BookingRepository
	paymentRepo payments.PaymentRepository
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

// PARA CLIENTES
func (s *BookingService) GetByUserID(ctx context.Context, userID uint, offset int64) ([]booking.Booking, error) {

	if userID == 0 {
		return nil, errors.New("el id del usuario no puede ser nulo")
	}

	return s.bookingRepo.GetByUserID(ctx, userID, offset)
}

func (s *BookingService) Reschedule(ctx context.Context, bookingID, slotID uint) (*booking.RescheduleResponse, error) {

	if bookingID == 0 {
		return nil, errors.New("el id de la reserva es necesario")
	}

	if slotID == 0 {
		return nil, errors.New("el id del turno es necesario")
	}

	// 1. Recuperar booking
	existing, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, errors.New("no fue posible recuperar la cita")
	}

	// 2. Validar si ya ocurrió
	now := time.Now().UTC()
	if existing.Slot.Start.UTC().Before(now) {
		return nil, errors.New("la cita ya ocurrió")
	}

	// 3. ¿Está dentro de las 24hs? → helper
	isWithin := IsWithin24Hours(existing.Slot.Start)

	// --- CASE A — Dentro de 24h → requiere pago ----
	if isWithin {

		payment, err := s.paymentRepo.GetByBookingID(ctx, existing.ID)
		if err != nil {
			return nil, errors.New("no fue posible recuperar el pago asociado")
		}

		surcharge := GetSurcharge(payment.Status, int64(payment.Amount))

		percentage := 0
		switch payment.Status {
		case "parcial":
			percentage = 50
		case "total":
			percentage = 25
		}

		initPoint, err := CreateReschedulePref(*existing, *payment, slotID)
		if err != nil {
			return nil, errors.New("no fue posible crear el link de pago")
		}

		return &booking.RescheduleResponse{
			RequiresPayment: true,
			Amount:          float64(surcharge),
			Percentage:      percentage,
			InitPoint:       initPoint,
			Free:            false,
			Reprogrammed:    false,
			Message:         fmt.Sprintf("La reprogramación es dentro de las 24 horas. Se aplicará un recargo del %d%% (monto: $%d).", percentage, surcharge),
		}, nil
	}

	// --- CASE B — Reprogramación gratuita ----
	// err = s.bookingRepo.UpdateSlot(ctx, bookingID, slotID)
	// if err != nil {
	//     return nil, errors.New("no fue posible reprogramar la cita")
	// }

	return &booking.RescheduleResponse{
		RequiresPayment: false,
		Free:            true,
		Reprogrammed:    true,
		Message:         "La reprogramación fue realizada exitosamente sin costo.",
	}, nil
}

func CreateReschedulePref(booking booking.Booking, payment payments.Payment, slotID uint) (string, error) {

	var (
		MP_ACCESS_TOKEN  = os.Getenv("MP_ACCESS_TOKEN")
		notification_url = ""
	)

	// 1. Analizar instancia
	totalAmount := GetSurcharge(payment.Status, int64(payment.Amount))

	// 4. Configurar Mercado Pago
	cfg, err := config.New(MP_ACCESS_TOKEN)
	if err != nil {
		return "", errors.New("error al configurar Mercado Pago")
	}

	client := preference.NewClient(cfg)

	mpRequest := preference.Request{
		Items: []preference.ItemRequest{
			{
				Title:     "Reprogramacion del turno",
				UnitPrice: float64(totalAmount),
				Quantity:  1,
			},
		},
		Payer: &preference.PayerRequest{
			Name:    booking.Client.Name,
			Surname: booking.Client.Surname,
		},
		NotificationURL: fmt.Sprintf("%s/api/v1/mercado_pago/notification", notification_url),
		Metadata: map[string]any{
			"booking_id": booking.ID,
			"payment_id": payment.ID,
			"slot_id":    slotID,
		},
		BackURLs: &preference.BackURLsRequest{
			Success: "http://localhost:5173",
		},
	}

	preferenceRes, err := client.Create(context.Background(), mpRequest)
	if err != nil {
		log.Println("[DEBUG]", err.Error())

		return "", errors.New("error al crear preferencia en Mercado Pago")
	}

	return preferenceRes.InitPoint, nil
}

func IsWithin24Hours(slotStart time.Time) bool {
	now := time.Now().UTC()
	return slotStart.UTC().Sub(now) <= 24*time.Hour
}

func GetSurcharge(paymentStatus string, totalPaid int64) int64 {
	switch paymentStatus {
	case "parcial":
		return totalPaid * 50 / 100
	case "total":
		return totalPaid * 25 / 100
	default:
		return 0
	}
}
