package http

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ezep02/rodeo/internal/catalog/domain/media"
	"github.com/ezep02/rodeo/internal/catalog/usecase"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	mediaSvc *usecase.MediaService
}

func NewMediaHandler(mediaSvc *usecase.MediaService) *MediaHandler {
	return &MediaHandler{mediaSvc}
}

func (h *MediaHandler) SetMedia(c *gin.Context) {
	var (
		req        media.Medias
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

	// 2. Recuperar datos de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	req_constructor := &media.Medias{
		ServiceID: id,
		URL:       req.URL,
		Type:      req.Type,
	}

	if err := h.mediaSvc.Create(c.Request.Context(), req_constructor); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req_constructor)
}

func (h *MediaHandler) Update(c *gin.Context) {
	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		promoIdStr = c.Param("id")
		req        media.Medias
	)

	// 1. Validar sesion del usuario
	existing, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if !existing.IsAdmin || !existing.IsBarber {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usted no tiene autorizacion"})
		return
	}

	// 2. Recuperar los datos de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fue posible recuperar los datos de la consulta"})
		return
	}

	// 3. Parsear el id
	parsedMediaId, err := strconv.ParseUint(promoIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parseando datos"})
		return
	}

	// 4. Construir promocion
	req_constructor := &media.Medias{
		URL:       req.URL,
		ID:        parsedMediaId,
		ServiceID: req.ID,
	}

	if err := h.mediaSvc.Update(c.Request.Context(), uint(parsedMediaId), req_constructor); err != nil {
		log.Println("Error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req_constructor)
}

func (h *MediaHandler) DeleteMedia(c *gin.Context) {
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

	if err := h.mediaSvc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categoria eliminada correctamente"})
}

func (h *MediaHandler) ListByServiceId(c *gin.Context) {

	var (
		idStr = c.Param("id")
	)

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	// 2. Parsing de datos
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offset invalido"})
		return
	}

	categories, err := h.mediaSvc.ListByServiceId(c.Request.Context(), uint(id))
	if err != nil {
		log.Println("Err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
