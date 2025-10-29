package http

import (
	"net/http"
	"os"

	"github.com/ezep02/rodeo/internal/analytics/usecase"
	"github.com/ezep02/rodeo/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type AnalyticHandler struct {
	svc *usecase.AnalyticService
}

func NewAnalyticHandler(analyticSvc *usecase.AnalyticService) *AnalyticHandler {
	return &AnalyticHandler{svc: analyticSvc}
}

func (h *AnalyticHandler) MonthlyRevenue(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	// 1. Verificar que el usuario sea un admin
	authenticated, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !authenticated.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "usted no tiene autorizacion",
		})
		return
	}

	// 1. Analiticas del total de ingresos por mes
	MonthlyRevenue, err := h.svc.MonthlyRevenue(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fue posible recuperar la informacion",
		})
		return
	}

	c.JSON(http.StatusOK, MonthlyRevenue)
}

func (h *AnalyticHandler) NewClientRate(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	// 1. Verificar que el usuario sea un admin
	authenticated, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !authenticated.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "usted no tiene autorizacion",
		})
		return
	}

	// 1.  Analiticas de nuevos clientes por mes
	clientRate, err := h.svc.NewClientRate(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fue posible recuperar la informacion",
		})
		return
	}

	c.JSON(http.StatusOK, clientRate)
}
