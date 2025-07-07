package http

import (
	"net/http"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/service"
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
	review := &domain.Review{
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
		"message": "Gracias por participar",
		"review":  review,
	})
}
