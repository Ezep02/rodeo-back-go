package delivery

import (
	"log"

	"github.com/ezep02/rodeo/internal/analytics/delivery/http"
	"github.com/ezep02/rodeo/internal/analytics/repository"
	"github.com/ezep02/rodeo/internal/analytics/usecase"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewAnalyticsRoutes(r *gin.RouterGroup, cnn *gorm.DB, redis *redis.Client) {

	log.Println("[ANALYTICS ROUTES] Setting up analytics routes")

	// Aquí irían las rutas relacionadas con la analítica
	analyticRepo := repository.NewGormAnalyticRepo(cnn)
	analyticSvc := usecase.NewAnalyticService(analyticRepo)

	analytics := r.Group("/analytics")
	{
		analyticHandler := http.NewAnalyticHandler(analyticSvc)
		analytics.GET("/month-revenue", analyticHandler.MonthlyRevenue)
		analytics.GET("/client-rate", analyticHandler.NewClientRate)
	}

	// Rutas de informacion de la barberia
	infoRepo := repository.NewGormInfoRepo(cnn, redis)
	infoSvc := usecase.NewInfoService(infoRepo)

	info := r.Group("/info")
	{
		infoHandler := http.NewInfoHandler(infoSvc)
		info.GET("/", infoHandler.Information)
	}
}
