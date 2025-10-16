package http

import (
	"net/http"

	"github.com/ezep02/rodeo/internal/analytics/usecase"

	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	infoSvc *usecase.InformationService
}

func NewInfoHandler(infoSvc *usecase.InformationService) *InfoHandler {
	return &InfoHandler{infoSvc}
}

func (h *InfoHandler) Information(c *gin.Context) {

	// 1. Consultar analiticas de la tasa de ocupacion de los slots por mes
	information, err := h.infoSvc.Information(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fue posible recuperar la informacion",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "solicitud exitosa",
		"info":    information,
	})
}
