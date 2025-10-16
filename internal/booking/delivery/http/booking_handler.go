package http

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/booking/usecases"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingSvc *usecases.BookingService
	paymentSvc *usecases.PaymentService
}

func NewBookingHandler(
	bookingSvc *usecases.BookingService,
	paymentSvc *usecases.PaymentService) *BookingHandler {
	return &BookingHandler{bookingSvc, paymentSvc}
}

func (b *BookingHandler) Upcoming(c *gin.Context) {

	var (
		dateStr     = c.Param("date")
		barberIDStr = c.Param("barber")

		auth_token = os.Getenv("AUTH_TOKEN")
		status     = c.Query("status") // "" si no viene
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
