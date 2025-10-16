package http

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ezep02/rodeo/internal/catalog/domain/service"
	"github.com/ezep02/rodeo/internal/catalog/usecase"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	svc *usecase.ServicesService
}

func NewServiceHandler(apptService *usecase.ServicesService) *ServiceHandler {
	return &ServiceHandler{apptService}
}

func (h *ServiceHandler) Create(c *gin.Context) {
	var (
		req        service.Service
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	// Verificar sesion del cliente
	existing, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !existing.IsBarber || !existing.IsAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usted no tiene permiso suficiente"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req_constructor := &service.Service{
		BarberID:    uint64(existing.ID),
		PreviewURL:  req.PreviewURL,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := h.svc.Create(c, req_constructor); err != nil {
		log.Println("Error creando servicio", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "algo no fue bien creando el servicio"})
		return
	}

	c.JSON(http.StatusOK, req_constructor)
}

func (h *ServiceHandler) List(c *gin.Context) {
	offsetStr := c.Param("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offset invalido"})
		return
	}

	list, err := h.svc.ListAll(c.Request.Context(), offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching Products"})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	prod, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prod)
}

func (h *ServiceHandler) Update(c *gin.Context) {

	var (
		req        service.Service
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

	// 2. Recuperar informacion de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	// 3. Crear objeto
	req_constructor := &service.Service{
		PreviewURL:  req.PreviewURL,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		IsActive:    req.IsActive,
		ID:          id,
	}

	if err := h.svc.Update(c, uint(id), req_constructor); err != nil {
		log.Println("Error actualizando servicio", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "algo no fue bien actualizando el servicio"})
		return
	}

	c.JSON(http.StatusOK, req_constructor)
}

func (h *ServiceHandler) Delete(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		srvIdStr   = c.Param("id")
	)

	if srvIdStr == "" {
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

	// 2. parsing del id
	id, err := strconv.ParseUint(srvIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, "servicio eliminado exitosamente")
}

func (h *ServiceHandler) Stats(c *gin.Context) {

	var (
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

	stats, err := h.svc.Stats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "algo no fue bien recuperando las estadisticas"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *ServiceHandler) Popular(c *gin.Context) {

	popular, err := h.svc.Popular(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "solicitud exitorisa",
		"popular": popular,
	})
}

func (h *ServiceHandler) AddCategories(c *gin.Context) {
	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		req        []uint
		srvIdStr   = c.Param("id")
	)

	if srvIdStr == "" {
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

	// 2. Recuperar informacion de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 3. parsing del id
	id, err := strconv.ParseUint(srvIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	if err := h.svc.AddCategories(c.Request.Context(), uint(id), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "categorias agregadas correctamente"})
}

func (h *ServiceHandler) RemoveCategories(c *gin.Context) {
	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		req        []uint
		srvIdStr   = c.Param("id")
	)

	if srvIdStr == "" {
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

	// 2. Recuperar informacion de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 3. parsing del id
	id, err := strconv.ParseUint(srvIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	if err := h.svc.RemoveCategories(c.Request.Context(), uint(id), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "categorias actualizadas correctamente"})
}
