package usecases

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/booking/domain/booking"
	"github.com/ezep02/rodeo/internal/booking/domain/payments"
	"github.com/ezep02/rodeo/internal/booking/domain/services"
)

type MepService struct {
	bookingRepo booking.BookingRepository
	paymentRepo payments.PaymentRepository
	//couponRepo  coupon.CouponRepository
	svcRepo services.ServicesRepository
}

func NewMepService(
	bookingRepo booking.BookingRepository,
	paymentRepo payments.PaymentRepository,
	svcRepo services.ServicesRepository,
) *MepService {
	return &MepService{bookingRepo, paymentRepo, svcRepo}
}

type MepaPreference struct {
	SlotID            uint   `json:"slot_id"`
	ServicesID        []uint `json:"services_id"`
	PaymentPercentage int64  `json:"payment_percentage"` // 50 para seña, 100 para total
	CouponCode        string `json:"coupon_code"`
}

func (s *MepService) CreateMpPreference(ctx context.Context, pref MepaPreference, clientID uint) (*booking.Booking, *payments.Payment, float64, error) {

	// 1. Calcular el precio final
	totalAmount, err := s.svcRepo.GetTotalPriceByIDs(ctx, pref.ServicesID)
	if err != nil {
		return nil, nil, 0, errors.New("no fue posible recuperar los servicios")
	}

	// 3. Crear booking
	booking := &booking.Booking{
		SlotID:      pref.SlotID,
		ClientID:    clientID,
		Status:      "pendiente_pago",
		TotalAmount: totalAmount,
		ExpiresAt: func() *time.Time {
			t := time.Now().Add(5 * time.Minute)
			return &t
		}(),
	}

	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, nil, 0, errors.New("no fue posible creando reserva")
	}

	// 3. Crear payment (seña o total)
	paymentAmount := totalAmount
	paymentType := "total"
	if pref.PaymentPercentage < 100 {
		paymentAmount = paymentAmount * float64(pref.PaymentPercentage) / 100
		paymentType = "parcial"
	}

	payment := &payments.Payment{
		BookingID: booking.ID,
		Amount:    paymentAmount,
		Type:      paymentType,
		Method:    "mercadopago",
		Status:    "pendiente",
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, nil, 0, errors.New("no fue posible instanciar la preferencia de pago")
	}

	// Instanciar servicios en segundo plano
	go func(svcIds []uint) {

		log.Println("[Iniciando proceso de almacenamiento de servicios]")
		var selectedSvc []services.BookingServices = make([]services.BookingServices, 0)

		ctx := context.Background()

		for _, id := range svcIds {

			if existing, _ := s.svcRepo.GetByID(ctx, id); existing != nil {
				selectedSvc = append(selectedSvc, services.BookingServices{
					BookingID: booking.ID,
					ServiceID: uint(existing.ID),
				})
			}
		}

		if err := s.svcRepo.SetBookingServices(ctx, selectedSvc); err != nil {
			log.Println("Error alamacenando los servicios seleccionados")
			return
		}

	}(pref.ServicesID)

	return booking, payment, totalAmount, nil
}
