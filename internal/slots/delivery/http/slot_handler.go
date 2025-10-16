package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/slots/domain"
	"github.com/ezep02/rodeo/internal/slots/usecase"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type SlotHandler struct {
	slotSvc *usecase.SlotUsecase
}

func NewSlotHandler(slotSvc *usecase.SlotUsecase) *SlotHandler {
	return &SlotHandler{slotSvc}
}

type CreateInBatchesRes struct {
	Batch []domain.Slot `json:"batch"`
}

func (h *SlotHandler) Create(c *gin.Context) {
	var (
		req        CreateInBatchesRes
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	// 1. recuperar slots desde la request
	authorized_user, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 2. Verificar que sea barbero
	if !authorized_user.IsBarber {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usted no tiene acceso"})
		return
	}

	// 3. Recuperar datos de la request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(400, gin.H{"error": "Objeto invalido"})
		return
	}

	// 4. Si todo fue bien, preparar los objetos con el barber id
	allSlots := []domain.Slot{}
	batchSlots := []domain.Slot{}
	batchLimit := 100

	for _, slot := range req.Batch {
		s := domain.Slot{
			BarberID: authorized_user.ID,
			Start:    slot.Start,
			End:      slot.End,
		}

		batchSlots = append(batchSlots, s)
		allSlots = append(allSlots, s)

		if len(batchSlots) == batchLimit {
			// Insertar batch y actualizar IDs en batchSlots
			h.slotSvc.CreateInBatches(context.Background(), &batchSlots)

			// Actualizar los IDs en allSlots
			copy(allSlots[len(allSlots)-batchLimit:], batchSlots)

			// Reiniciar batchSlots
			batchSlots = []domain.Slot{}
		}
	}

	// Insertar cualquier batch restante
	if len(batchSlots) > 0 {
		log.Println("Batch slot", batchSlots)
		h.slotSvc.CreateInBatches(context.Background(), &batchSlots)
		copy(allSlots[len(allSlots)-len(batchSlots):], batchSlots)
	}

	c.JSON(http.StatusOK, allSlots)
}

func (h *SlotHandler) Update(c *gin.Context) {

}

func (h *SlotHandler) GetByDateRange(c *gin.Context) {

	var (
		idStr        = c.Param("barber")
		startDateStr = c.Param("start")
		endDateStr   = c.Param("end")
		auth_token   = os.Getenv("AUTH_TOKEN")
	)

	if idStr == "" || startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "falta de parametros en la consulta"})
		return
	}

	// 1. recuperar slots desde la request
	_, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 2. Parsing de datos
	parsedId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error en parametros de la consulta"})
		return
	}

	// 3. Parsing de fechas
	startDateParsed, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parseando fecha de inicio en parametros de la consulta"})
		return
	}

	endDateParsed, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parseando fecha de fin en parametros de la consulta"})
		return
	}

	slotRange, err := h.slotSvc.GetByDateRange(c.Request.Context(), uint(parsedId), startDateParsed, endDateParsed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error recuperando los slots"})
		return
	}

	c.JSON(http.StatusOK, slotRange)
}
