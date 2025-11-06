package http

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ezep02/rodeo/internal/users/usecase"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type BarberHandler struct {
	barberSvc *usecase.BarberService
}

func NewBarberHandler(barberSvc *usecase.BarberService) *BarberHandler {
	return &BarberHandler{barberSvc}
}

func (h *BarberHandler) GetByID(c *gin.Context) {
	var (
		idStr = c.Param("id")
	)

	// 1. Vlidar el id
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// 2. Parsear el id a uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El id debe ser valido"})
		return
	}

	// 3. Obtener el barber
	barber, err := h.barberSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Barber not found"})
		return
	}

	c.JSON(http.StatusOK, barber)
}

func (h *BarberHandler) List(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	if _, err := jwt.VerifyUserSession(c, auth_token); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	// 1. Obtener el barber
	barber, err := h.barberSvc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Barber not found"})
		return
	}

	c.JSON(http.StatusOK, barber)
}
