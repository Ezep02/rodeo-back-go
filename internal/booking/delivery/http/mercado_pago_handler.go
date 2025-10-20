package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/booking/domain/booking"
	"github.com/ezep02/rodeo/internal/booking/domain/payments"
	"github.com/ezep02/rodeo/internal/booking/usecases"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type MepaHandler struct {
	bookingSvc  *usecases.BookingService
	paymentSvc  *usecases.PaymentService
	couponSvc   *usecases.CouponService
	servicesSvc *usecases.ServicesService
}

var (
	notification_url string = "https://bc1dd7937fd2.ngrok-free.app" // URL de notificación
)

func NewMepaHandler(
	bookingSvc *usecases.BookingService,
	paymentSvc *usecases.PaymentService,
	couponSvc *usecases.CouponService,
	servicesSvc *usecases.ServicesService) *MepaHandler {
	return &MepaHandler{bookingSvc, paymentSvc, couponSvc, servicesSvc}
}

type CreatePreferenceRequest struct {
	SlotID            uint   `json:"slot_id"`
	ServicesID        []uint `json:"services_id"`
	PaymentPercentage int64  `json:"payment_percentage"` // 50 para seña, 100 para total
	CouponCode        string `json:"coupon_code"`
}

func (h *MepaHandler) CreatePreference(c *gin.Context) {
	var (
		req             CreatePreferenceRequest
		MP_ACCESS_TOKEN = os.Getenv("MP_ACCESS_TOKEN")
		AUTH_TOKEN      = os.Getenv("AUTH_TOKEN")
	)

	if MP_ACCESS_TOKEN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "falta mp acces token"})
		return
	}

	if AUTH_TOKEN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "falta auth token"})
		return
	}

	// Parsear request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validar sesión del usuario
	authenticatedUser, err := jwt.VerifyUserSession(c, AUTH_TOKEN)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 2. Crear booking
	totalAmount, err := h.servicesSvc.GetTotalPriceByIDs(c.Request.Context(), req.ServicesID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no fue posible recuperar los servicios"})
		return
	}

	// 3. Crear booking
	booking := &booking.Booking{
		SlotID:      req.SlotID,
		ClientID:    authenticatedUser.ID,
		Status:      "pendiente_pago",
		TotalAmount: totalAmount,
		ExpiresAt: func() *time.Time {
			t := time.Now().Add(1 * time.Hour)
			return &t
		}(),
	}

	if err := h.bookingSvc.CreateBooking(c, booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear reserva"})
		return
	}

	// 3. Crear payment (seña o total)
	paymentAmount := totalAmount
	paymentType := "total"
	if req.PaymentPercentage < 100 {
		paymentAmount = paymentAmount * float64(req.PaymentPercentage) / 100
		paymentType = "seña"
	}

	payment := &payments.Payment{
		BookingID: booking.ID,
		Amount:    paymentAmount,
		Type:      paymentType,
		Method:    "mercadopago",
		Status:    "pendiente",
	}

	if err := h.paymentSvc.CreatePayment(c, payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear pago"})
		return
	}

	// 4. Configurar Mercado Pago
	cfg, err := config.New(MP_ACCESS_TOKEN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al configurar Mercado Pago"})
		return
	}
	client := preference.NewClient(cfg)

	mpRequest := preference.Request{
		Items: []preference.ItemRequest{
			{
				Title:     "Tus servicios",
				UnitPrice: totalAmount,
				Quantity:  1,
			},
		},
		Payer: &preference.PayerRequest{
			Name:    authenticatedUser.Name,
			Surname: authenticatedUser.Surname,
		},
		NotificationURL: fmt.Sprintf("%s/api/v1/mercado_pago/notification", notification_url),
		Metadata: map[string]any{
			"booking_id":         booking.ID,
			"payment_id":         payment.ID,
			"slot_id":            req.SlotID,
			"user_id":            authenticatedUser.ID,
			"payment_percentage": req.PaymentPercentage,
		},
		BackURLs: &preference.BackURLsRequest{
			Success: "http://localhost:5173",
		},
	}

	preferenceRes, err := client.Create(c.Request.Context(), mpRequest)
	if err != nil {
		log.Println("[DEBUG]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear preferencia en Mercado Pago"})
		return
	}

	// 5. Guardar URL de preferencia en payment
	payment.PaymentURL = &preferenceRes.InitPoint
	if err := h.paymentSvc.UpdatePayment(c, payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando pago"})
		return
	}

	c.JSON(http.StatusOK, preferenceRes.InitPoint)
}

func (h *MepaHandler) HandleNotification(c *gin.Context) {
	var (
		payload         map[string]any
		MP_ACCESS_TOKEN = os.Getenv("MP_ACCESS_TOKEN")
	)

	// 2. Decodificar payload enviado por mp
	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// 3. Recuperar del payload el campo id almacenado dentro de data
	data, ok := payload["data"].(map[string]any)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'data' inválido"})
		return
	}

	paymentStr := fmt.Sprintf("%v", data["id"])
	mpPaymentID, err := strconv.ParseInt(paymentStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de pago inválido"})
		return
	}

	// 4. Inicializar el cliente de Mercado Pago
	cfg, err := config.New(MP_ACCESS_TOKEN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo configurar el cliente"})
		return
	}

	// 5. Consultar pago utilizanodo el ID
	paymentClient := payment.NewClient(cfg)

	paymentInfo, err := paymentClient.Get(context.Background(), int(mpPaymentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pago no encontrado"})
		return
	}

	// 	// Leer metadata
	bookingID := uint(paymentInfo.Metadata["booking_id"].(float64))
	paymentID := uint(paymentInfo.Metadata["payment_id"].(float64))

	// Actualizar en base
	if paymentInfo.Status == "approved" {

		go func(bookingID, paymentID uint) {
			ctx := context.Background()
			// 1. Actualizar payment
			if err := h.paymentSvc.MarkAsPaid(ctx, paymentID, paymentInfo.Order.ID); err != nil {
				log.Println("Err", err.Error())
				log.Println("payment fallo actualizando status a pagado")
			}

			if err := h.bookingSvc.MarkAsPaid(ctx, bookingID); err != nil {
				log.Println("Err", err.Error())
				log.Println("booking fallo actualizando status a confirmado")
			}

		}(bookingID, paymentID)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
