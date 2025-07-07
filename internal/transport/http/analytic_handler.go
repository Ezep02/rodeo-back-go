package http

import (
	"net/http"

	"github.com/ezep02/rodeo/internal/service"
	"github.com/gin-gonic/gin"
)

type AnalyticHandler struct {
	svc *service.AnalyticService
}

func NewAnalyticHandler(analyticSvc *service.AnalyticService) *AnalyticHandler {
	return &AnalyticHandler{svc: analyticSvc}
}

func (h *AnalyticHandler) BookingOcupationRate(c *gin.Context) {

	// PASS

	// 1. Consultar analiticas de la tasa de ocupacion de los slots por mes
	bookingRate, err := h.svc.BookingOcupationRate(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fue posible recuperar la informacion",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "solicitud exitosa",
		"ocupation_rate": bookingRate,
	})
}

func (h *AnalyticHandler) MonthBookingCount(c *gin.Context) {

	//

	// 1. Analiticas de numero de citas por mes
	MonthBookingCount, err := h.svc.MonthBookingCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fue posible recuperar la informacion",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "solicitud exitosa",
		"month_booking_count": MonthBookingCount,
	})
}

func (h *AnalyticHandler) MonthlyRevenue(c *gin.Context) {

	// 1. Analiticas del total de ingresos por mes
	MonthlyRevenue, err := h.svc.MonthlyRevenue(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fue posible recuperar la informacion",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "solicitud exitosa",
		"monthly_revenue": MonthlyRevenue,
	})
}

func (h *AnalyticHandler) NewClientRate(c *gin.Context) {

	// 1.  Analiticas de nuevos clientes por mes
	ClientRate, err := h.svc.NewClientRate(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fue posible recuperar la informacion",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "solicitud exitosa",
		"month_client_rate": ClientRate,
	})
}

func (h *AnalyticHandler) PopularTimeSlot(c *gin.Context) {
	// 1.  Analiticas de la franja horaria mas popular
	popularTimeSlot, err := h.svc.PopularTimeSlot(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fue posible recuperar la informacion",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "solicitud exitosa",
		"popular_time_slot": popularTimeSlot,
	})
}

func (h *AnalyticHandler) WeeklyBookingRate(c *gin.Context) {
	// 1.  Analiticas del promedio de citas por semana
	weeklyBookingRate, err := h.svc.WeeklyBookingRate(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fue posible recuperar la informacion",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "solicitud exitosa",
		"weekly_booking_rate": weeklyBookingRate,
	})
}
