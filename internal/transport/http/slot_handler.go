package http

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/service"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type SlotHandler struct {
	svc *service.SlotService
}

func NewSlotHandler(apptService *service.SlotService) *SlotHandler {
	return &SlotHandler{apptService}
}

type CreateSlotRequest struct {
	Date time.Time `json:"date"`
	Time string    `json:"time"`
}

type UpdateSlotRequest struct {
	ID        uint      `json:"id"`
	Date      time.Time `json:"date"`
	Time      string    `json:"time"`
	Is_booked bool      `json:"is_booked"`
}

type DeleteSlotRequest struct {
	ID        uint      `json:"id"`
	Date      time.Time `json:"date"`
	Time      string    `json:"time"`
	Is_booked bool      `json:"is_booked"`
}

func (h *SlotHandler) Create(c *gin.Context) {

	var (
		req        []CreateSlotRequest
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	cookie, err := c.Cookie(auth_token)
	if err != nil {
		log.Println("error", err)
	}

	// 4. Validar la cookie
	user, err := jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	if !user.Is_barber {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	// 1. Crear arreglo de slot
	slot := make([]domain.Slot, len(req))

	for i, s := range req {
		slot[i] = domain.Slot{
			Date:     s.Date,
			Time:     s.Time,
			IsBooked: false,
			BarberID: user.ID,
		}
	}

	if err := h.svc.Create(c.Request.Context(), &slot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Producto creado exitosamente",
		"slot":    slot,
	})
}

func (h *SlotHandler) Update(c *gin.Context) {

	var (
		req        []UpdateSlotRequest
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	cookie, err := c.Cookie(auth_token)
	if err != nil {
		log.Println("error", err)
	}

	// 4. Validar la cookie
	user, err := jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	if !user.Is_barber {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	// 1. Crear arreglo de slot
	slot := make([]domain.Slot, len(req))

	for i, s := range req {
		slot[i] = domain.Slot{
			ID:       s.ID,
			Date:     s.Date,
			Time:     s.Time,
			IsBooked: s.Is_booked,
			BarberID: user.ID,
		}
	}

	if err := h.svc.Update(c.Request.Context(), &slot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Slot actualizado exitosamente",
		"slot":    slot,
	})

}

func (h *SlotHandler) Delete(c *gin.Context) {

	var (
		req []DeleteSlotRequest
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Convertir de []DeleteSlotRequest a []domain.Slot
	slots := make([]domain.Slot, len(req))
	for i, s := range req {
		slots[i] = domain.Slot{
			Date: s.Date,
			Time: s.Time,
		}
	}

	if err := h.svc.Delete(c.Request.Context(), slots); err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "slot no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "eliminado exitosamente"})
}

func (h *SlotHandler) ListByDate(c *gin.Context) {

	idStr := c.Param("id")

	date, err := time.Parse(time.RFC3339, idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	slots, err := h.svc.ListByDate(c.Request.Context(), date)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "slots no encontrados"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"slots": slots,
	})
}

func (h *SlotHandler) ListByDateRange(c *gin.Context) {

	// 1. Obtener params de la consulta
	startStr := c.Param("start")
	endStr := c.Param("end")

	// 2. Parsear los horarios
	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start date invalid"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end date invalid"})
		return
	}

	// 3. Verificar que no sean fechas en el pasado
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end date must be after start date"})
		return
	}

	slots, err := h.svc.ListByDateRange(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error recuperando slots"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"slots": slots,
		"total": len(slots),
	})
}

func (h *SlotHandler) GetByID(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	slot, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cita no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching appointment"})
		return
	}

	c.JSON(http.StatusOK, slot)
}
