package http

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ezep02/rodeo/internal/catalog/domain/categorie"
	"github.com/ezep02/rodeo/internal/catalog/usecase"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categorySvc *usecase.CategoryService
}

func NewCategorieHandler(categorySvc *usecase.CategoryService) *CategoryHandler {
	return &CategoryHandler{categorySvc}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var (
		req        categorie.Categorie
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	// 1, validar que el usuario sea un barbero o admin
	existing, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !existing.IsBarber || !existing.IsAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usted no tiene permiso suficiente"})
		return
	}

	// 2. Recuperar datos de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	req_constructor := &categorie.Categorie{
		Name:       req.Name,
		PreviewURL: req.PreviewURL,
	}

	if err := h.categorySvc.CreateCategory(c.Request.Context(), req_constructor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req_constructor)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		req        categorie.Categorie
		idStr      = c.Param("id")
	)

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	// 1, validar que el usuario sea un barbero o admin
	existing, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !existing.IsBarber || !existing.IsAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usted no tiene permiso suficiente"})
		return
	}

	// 2. Parsing de datos
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	req_constructor := &categorie.Categorie{
		Name:       req.Name,
		PreviewURL: req.PreviewURL,
	}

	if err := h.categorySvc.UpdateCategory(c.Request.Context(), uint(id), req_constructor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req_constructor)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		idStr      = c.Param("id")
	)

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	// 1, validar que el usuario sea un barbero o admin
	existing, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !existing.IsBarber || !existing.IsAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usted no tiene permiso suficiente"})
		return
	}

	// 2. Parsing de datos
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	if err := h.categorySvc.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categoria eliminada correctamente"})
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {

	// var (
	// 	offsetStr = c.Param("offset")
	// )

	// 2. Parsing de datos
	// offset, err := strconv.Atoi(offsetStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Offset invalido"})
	// 	return
	// }

	categories, err := h.categorySvc.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
