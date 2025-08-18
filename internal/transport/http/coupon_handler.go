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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontro cookie de sesion"})
		return
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

func (h *CouponHandler) GetByUserID(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		userIDStr  = c.Param("id")
	)

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 1. convertir userID a int
	id, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	// 2. Recuperar cookie de sesion
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontro cookie de sesion"})
		return
	}

	// 3. Validar la sesion
	user, err := jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	// 4. Validar que el usuario sea el mismo que el del ID
	if user.ID != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para acceder a estos cupones"})
		return
	}

	// 5. Recuperar cupones por userID
	coupons, err := h.svc.GetByUserID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coupons": coupons,
	})
}

func (h *CouponHandler) GetByCode(c *gin.Context) {

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coupon code is required"})
		return
	}

	coupon, err := h.svc.GetByCode(c.Request.Context(), code)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coupon": coupon,
	})
}
