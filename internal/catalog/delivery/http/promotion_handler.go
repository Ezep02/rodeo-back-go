package http

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ezep02/rodeo/internal/catalog/domain/promotions"
	"github.com/ezep02/rodeo/internal/catalog/usecase"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type PromoHandler struct {
	promoSvc *usecase.PromoService
}

func NewPromoHandler(promoSvc *usecase.PromoService) *PromoHandler {
	return &PromoHandler{promoSvc}
}

type CreatePromoReq struct {
	ServiceId uint                 `json:"id"`
	Data      promotions.Promotion `json:"data"`
}

func (h *PromoHandler) Create(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		req        CreatePromoReq
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

	// 2 parsing de datos
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "algo no fue bien recuperando los datos de la consulta"})
		return
	}

	// 3. crear consulta
	req_constructor := &promotions.Promotion{
		ServiceID: uint64(req.ServiceId),
		Discount:  req.Data.Discount,
		Type:      req.Data.Type,
		StartDate: req.Data.StartDate,
		EndDate:   req.Data.EndDate,
	}

	if err := h.promoSvc.Create(c.Request.Context(), req_constructor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req_constructor)
}

func (h *PromoHandler) ListByServiceId(c *gin.Context) {
	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		svcIdStr   = c.Param("id")
		offsetStr  = c.Param("offset")
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

	// 2 Parsing de datos
	parsedSvcId, err := strconv.ParseUint(svcIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parseando datos"})
		return
	}

	parsedOffset, err := strconv.ParseUint(offsetStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parseando datos"})
		return
	}

	list, err := h.promoSvc.ListByServiceId(c.Request.Context(), uint(parsedSvcId), int(parsedOffset))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error recuperando informacion"})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *PromoHandler) Update(c *gin.Context) {
	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		promoIdStr = c.Param("id")
		req        promotions.Promotion
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
	parsedPromoId, err := strconv.ParseUint(promoIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parseando datos"})
		return
	}

	// 4. Construir promocion
	req_constructor := &promotions.Promotion{
		Discount:  req.Discount,
		Type:      req.Type,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		ID:        parsedPromoId,
	}

	if err := h.promoSvc.Update(c.Request.Context(), uint(parsedPromoId), req_constructor); err != nil {
		log.Println("Error actualizando promocion", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req_constructor)
}

func (h *PromoHandler) Delete(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		promoIdStr = c.Param("id")
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

	// 2. Parsear id
	parsedPromoId, err := strconv.ParseUint(promoIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parseando datos"})
		return
	}

	if err := h.promoSvc.Delete(c.Request.Context(), uint(parsedPromoId)); err != nil {
		log.Println("[DEBUG]: error eliminando promocion", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "promocion eliminada correctamente"})
}
