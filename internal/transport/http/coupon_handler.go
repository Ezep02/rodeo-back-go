package http

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/service"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	svc *service.CouponService
}

func NewCouponHandler(svc *service.CouponService) *CouponHandler {
	return &CouponHandler{svc}
}

type CreateCouponReq struct {
	DiscountPercentage float64 `json:"discount_percentage"`
}

func (h *CouponHandler) Create(c *gin.Context) {

	var (
		req        CreateCouponReq
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	if auth_token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan variables de entorno"})
		return
	}

	// 1. Capturar informacion
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 2. Validar sesion, y recuperar user_id

	// 3. Recuperar id de cliente si es que existe desde la session
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

	// 2. Generar codigo unico
	code, err := h.svc.GenerateCoupon(12)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating code"})
		return
	}
	// 3. Crear modelo
	log.Println("New code", code)

	newCoupon := domain.Coupon{
		Code:               code,
		UserID:             user.ID,
		DiscountPercentage: req.DiscountPercentage,
		IsAvailable:        true,
		CreatedAt:          time.Now(),
		ExpireAt:           time.Now().Add(time.Hour * 24 * 7),
	}

	if err := h.svc.Create(c.Request.Context(), &newCoupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Algo no fue bien creando el cupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cupon creado con exito",
		"coupon":  newCoupon,
	})
}
