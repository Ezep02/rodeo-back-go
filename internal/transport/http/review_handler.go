package http

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ezep02/rodeo/internal/domain/review"

	"github.com/ezep02/rodeo/internal/service"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	svc *service.ReviewService
}

func NewReviewHandler(revService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{revService}
}

type CreateReviewReq struct {
	AppointmentID uint   `json:"appointment_id"`
	Comment       string `json:"comment"`
	Rating        int    `json:"rating"`
}

func (h *ReviewHandler) Create(c *gin.Context) {

	var (
		req CreateReviewReq
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Crear producto
	review := &review.Review{
		AppointmentID: req.AppointmentID,
		Rating:        req.Rating,
		Comment:       req.Comment,
	}

	if err := h.svc.Create(c.Request.Context(), review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Gracias por participar",
		"review":  review,
	})
}

func (h *ReviewHandler) List(c *gin.Context) {

	review, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "algo no fue bien recuperando la reseñas",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reseñas recuperadas correctamente",
		"review":  review,
	})
}

func (h *ReviewHandler) ListByProductID(c *gin.Context) {
	productIDStr := c.Param("id")

	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "product ID is required",
		})
		return
	}

	// Convertir productID a entero
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	reviews, err := h.svc.ListByProductID(c.Request.Context(), uint(productID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "algo no fue bien recuperando las reseñas del producto",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reseñas del producto recuperadas correctamente",
		"reviews": reviews,
	})
}

func (h *ReviewHandler) ListByUserID(c *gin.Context) {
	// Validar session de usuario
	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		userIDStr  = c.Param("id")
		offsetStr  = c.Param("offset")
	)

	if userIDStr == "" || offsetStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user ID and offset are required",
		})
		return
	}

	// 1. Recuperar cookies
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	// 2. Validar la cookie
	_, err = jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	// Convertir offeset e userID a enteros
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	// Realizar consulta
	reviews, err := h.svc.ListByUserID(c.Request.Context(), uint(userID), int(offset))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "algo no fue bien recuperando las reseñas del usuario",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reseñas del usuario recuperadas correctamente",
		"reviews": reviews,
	})
}
