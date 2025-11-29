package http

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/booking/domain/booking"
	"github.com/ezep02/rodeo/internal/booking/domain/payments"
	"github.com/ezep02/rodeo/internal/booking/usecases"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingSvc  *usecases.BookingService
	paymentSvc  *usecases.PaymentService
	couponSvc   *usecases.CouponService
	servicesSvc *usecases.ServicesService
}

func NewBookingHandler(
	bookingSvc *usecases.BookingService,
	paymentSvc *usecases.PaymentService,
	couponSvc *usecases.CouponService,
	servicesSvc *usecases.ServicesService,
) *BookingHandler {
	return &BookingHandler{bookingSvc, paymentSvc, couponSvc, servicesSvc}
}

func (b *BookingHandler) Upcoming(c *gin.Context) {

	var (
		dateStr     = c.Param("date")
		barberIDStr = c.Param("barber")

		auth_token = os.Getenv("AUTH_TOKEN")
		status     = c.Query("status")
	)

	// 1. Verificar la sesion del usuario
	authenticated, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !authenticated.IsBarber {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usted no tiene permiso suficiente"})
		return
	}

	// 3. Parsing de fechas
	startDateParsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parseando fecha de inicio en parametros de la consulta"})
		return
	}

	barberID, err := strconv.Atoi(barberIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de barbero inválido"})
		return
	}

	// // Llamar al repositorio con los filtros
	bookings, err := b.bookingSvc.Upcoming(c.Request.Context(), uint(barberID), startDateParsed, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (b *BookingHandler) StatsByBarberID(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		idStr      = c.Param("id")
	)

	if auth_token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "algo fue mal"})
		return
	}

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "algo fue mal recuperando el id"})
		return
	}

	// 1. Verificar sesion del usuario
	authenticated, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !authenticated.IsBarber {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usted no tiene permiso suficiente"})
		return
	}

	// 2. Parsear el id
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no fue posible parsear el id"})
		return
	}

	// 3. Consulta
	barberStats, err := b.bookingSvc.StatsByBarberID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no fue posible recuperar las estadisticas"})
		return
	}

	c.JSON(http.StatusOK, barberStats)
}

func (b *BookingHandler) AllPendingPayment(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	// 1. Verificar sesion del usuario
	authenticated, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !authenticated.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usted no tiene permiso suficiente"})
		return
	}

	// 2. Consulta
	bookings, err := b.bookingSvc.AllPendingPayment(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no fue posible recuperar las reservas pendientes de pago"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

type CreateBookingRequest struct {
	SlotID            uint   `json:"slot_id"`
	ServicesID        []uint `json:"services_id"`
	PaymentPercentage int64  `json:"payment_percentage"` // 50 para seña, 100 para total
	CouponCode        string `json:"coupon_code"`
}

func (b *BookingHandler) Create(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		req        CreateBookingRequest
	)

	// 1. Verificar sesion del usuario
	authenticated, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Parsear request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 2. Crear booking
	totalAmount, err := b.servicesSvc.GetTotalPriceByIDs(c.Request.Context(), req.ServicesID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no fue posible recuperar los servicios"})
		return
	}

	// 3. Crear booking
	booking := &booking.Booking{
		SlotID:      req.SlotID,
		ClientID:    authenticated.ID,
		Status:      "pendiente_pago",
		TotalAmount: totalAmount,
	}

	if err := b.bookingSvc.CreateBooking(c, booking); err != nil {
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
		Method:    "transferencia",
		Status:    "pendiente",
	}

	if err := b.paymentSvc.CreatePayment(c, payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear pago"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (b *BookingHandler) MarkAsPaid(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		idStr      = c.Param("id")
	)

	// 1. Verificar sesion del usuario
	authenticated, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !authenticated.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usted no tiene permiso suficiente"})
		return
	}

	// 2. Parsear el id
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no fue posible parsear el id"})
		return
	}

	// 3. Marcar como pagado
	if err := b.bookingSvc.MarkAsPaid(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no fue posible marcar la reserva como pagada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reserva aceptada exitosamente"})
}

func (b *BookingHandler) MarkAsRejected(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		idStr      = c.Param("id")
	)

	// 1. Verificar sesion del usuario
	authenticated, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !authenticated.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usted no tiene permiso suficiente"})
		return
	}

	// 2. Parsear el id
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no fue posible parsear el id"})
		return
	}

	// 3. Marcar como rechazado
	if err := b.bookingSvc.MarkAsRejected(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no fue posible marcar la reserva como rechazada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reserva rechazada exitosamente"})
}

func (b *BookingHandler) BookingPayment(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		idStr      = c.Param("id")
	)

	// 1. Verificar sesion del usuario
	if _, err := jwt.VerifyUserSession(c, auth_token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Parsear el id
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no fue posible parsear el id"})
		return
	}

	paymentInfo, err := b.paymentSvc.GetByBookingID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fue posible parsear el id"})
		return
	}

	c.JSON(http.StatusOK, paymentInfo)
}
