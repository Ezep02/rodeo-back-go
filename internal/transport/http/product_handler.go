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
	svc         *service.ProductService
	categorySvc *service.CategoryService
}

func NewProductHandler(apptService *service.ProductService, categorySvc *service.CategoryService) *ProductHandler {
	return &ProductHandler{apptService, categorySvc}
}

type CreateProductRequest struct {
	Name              string  `json:"name"`
	Price             float64 `json:"price"`
	Description       string  `json:"description"`
	CategoryID        uint    `json:"category_id"`
	Preview_url       string  `json:"preview_url"`
	PromotionDiscount int     `json:"promotion_discount"`
	PromotionEndDate  string  `json:"promotion_end_date"`
	HasPromotion      bool    `json:"has_promotion"`
}

type UpdateProductRequest struct {
	Name              string  `json:"name" binding:"required"`
	Price             float64 `json:"price" binding:"required"`
	Description       string  `json:"description"`
	CategoryID        uint    `json:"category_id"`
	Preview_url       string  `json:"preview_url"`
	PromotionDiscount int     `json:"promotion_discount"`
	PromotionEndDate  string  `json:"promotion_end_date"`
	HasPromotion      bool    `json:"has_promotion"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var (
		req             CreateProductRequest
		formatedEndDate *time.Time
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Crear producto
	if req.HasPromotion {
		endDate, err := time.Parse("2006-01-02", req.PromotionEndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "promotion end date invalid"})
			return
		}

		formatedEndDate = &endDate
	}

	product := &domain.Product{
		Name:              req.Name,
		Price:             req.Price,
		Description:       req.Description,
		CategoryID:        req.CategoryID,
		PreviewUrl:        req.Preview_url,
		PromotionDiscount: req.PromotionDiscount,
		PromotionEndDate:  formatedEndDate,
		HasPromotion:      req.HasPromotion,
	}

	if err := h.svc.Create(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. Asociar categoria
	if req.CategoryID != 0 {

		category, err := h.categorySvc.GetCategoryByID(c.Request.Context(), req.CategoryID)
		if err != nil {
			if err == domain.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "categoria no encontrada"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		product.Category = &domain.Category{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			Color:     category.Color,
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Producto creado exitosamente",
		"product": product,
	})
}

func (h *ProductHandler) List(c *gin.Context) {
	offsetStr := c.Param("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offset invalido"})
		return
	}

	Products, err := h.svc.ListAll(c.Request.Context(), offset)
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

	var (
		req             UpdateProductRequest
		formatedEndDate *time.Time
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Crear producto

	if req.HasPromotion {
		endDate, err := time.Parse("2006-01-02", req.PromotionEndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "promotion end date invalid"})
			return
		}
		formatedEndDate = &endDate
	}

	// 3. Actualizar la cita
	log.Println("Updating product with ID:", req.HasPromotion)
	product := &domain.Product{
		ID:                uint(id),
		Name:              req.Name,
		Price:             req.Price,
		Description:       req.Description,
		CategoryID:        req.CategoryID,
		PreviewUrl:        req.Preview_url,
		PromotionDiscount: req.PromotionDiscount,
		PromotionEndDate:  formatedEndDate,
		HasPromotion:      req.HasPromotion,
		UpdatedAt:         time.Now(),
	}

	if err := h.svc.Update(c.Request.Context(), uint(id), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.CategoryID != 0 {
		category, err := h.categorySvc.GetCategoryByID(c.Request.Context(), req.CategoryID)
		if err != nil {
			if err == domain.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "categoria no encontrada"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		product.Category = &domain.Category{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			Color:     category.Color,
		}
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

func (h *ProductHandler) Promotion(c *gin.Context) {
	promotion, err := h.svc.Promotion(c.Request.Context())

	if err == domain.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "solicitud exitosa",
		"promotion": promotion,
	})
}
