package analytics

import (
	"net/http"

	"github.com/ezep02/rodeo/internal/analytics/handler"
	"github.com/ezep02/rodeo/internal/analytics/repository"
	"github.com/ezep02/rodeo/internal/analytics/services"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AnalyticsRoutes(mux *http.ServeMux, db *gorm.DB, redisClient *redis.Client) {
	analyticsService := services.NewAnalyticServices(repository.NewAnalyticsRepository(db, redisClient))
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Analytics
	mux.HandleFunc("/analytics/revenue", analyticsHandler.GetMonthlyRevenueAndAvgHandler)
	mux.HandleFunc("/analytics/appointments", analyticsHandler.GetMonthlyAppointmentsAndAvgHandler)
	mux.HandleFunc("/analytics/customers", analyticsHandler.GetMonthlyNewCustomersAndAvgHandler)
	mux.HandleFunc("/analytics/cancellations", analyticsHandler.GetMonthlyCancellationsAndAvgHandler)
	mux.HandleFunc("/analytics/revenue/current-year", analyticsHandler.GetCurrentYearMonthlyRevenueHandler)
	mux.HandleFunc("/analytics/popular-services", analyticsHandler.GetMonthlyPopularServicesHandler)
	mux.HandleFunc("/analytics/frequent-customers", analyticsHandler.GetFrequentCustomersHandler)

	// Barber
	mux.HandleFunc("/barber/yearly-haircut", analyticsHandler.GetYearlyBarberHaircuts)
}
