package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func (b *BookingHandler) AllByUserId(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		idStr      = c.Param("id")
	)

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fue posible recuperar el id del usuario"})
		return
	}

	// verifiar sesion
	if _, err := jwt.VerifyUserSession(c, auth_token); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// parsear el user id
	parsedId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fue posible recuperar el id del usuario"})
		return
	}

	// recuperar datos
	list, err := b.bookingSvc.GetByUserID(c.Request.Context(), uint(parsedId), 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fue posible recuperar los datos del usuario"})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (b *BookingHandler) Reschedule(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		reqBody    struct {
			BookingID uint `json:"booking_id"`
			NewSlotID uint `json:"new_slot_id"`
		}
	)

	// 1. Verificar sesion del usuario
	if _, err := jwt.VerifyUserSession(c, auth_token); err != nil {
		fmt.Printf("[Error verificando sesion] %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 2, Verificar id del booking
	// 2. Parsear request
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 3.
	res, err := b.bookingSvc.Reschedule(c.Request.Context(), reqBody.BookingID, reqBody.NewSlotID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
