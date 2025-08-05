package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(apptService *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: apptService}
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	CategoryID  uint    `json:"category_id"`
	Preview_url string  `json:"preview_url"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Description string  `json:"description"`
	CategoryID  uint    `json:"category_id"`
	Preview_url string  `json:"preview_url"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Crear producto
	product := &domain.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		PreviewUrl:  req.Preview_url,
	}

	if err := h.svc.Create(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Producto creado exitosamente",
		"product": product,
	})
}

func (h *ProductHandler) List(c *gin.Context) {
	Products, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching Products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": Products,
		"total":    len(Products),
	})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	prod, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prod)
}

func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 3. Actualizar la cita
	product := &domain.Product{
		ID:          uint(id),
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		PreviewUrl:  req.Preview_url,
		UpdatedAt:   time.Now(),
	}

	if err := h.svc.Update(c.Request.Context(), uint(id), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Producto actualizado correctamente",
		"product": product,
	})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, "eliminado exitosamente")
}

func (h *ProductHandler) Popular(c *gin.Context) {

	popular, err := h.svc.Popular(c.Request.Context())

	if err == domain.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "solicitud exitorisa",
		"popular": popular,
	})
}
