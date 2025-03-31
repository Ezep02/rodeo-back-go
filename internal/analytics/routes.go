package analytics

import (
	"github.com/ezep02/rodeo/internal/analytics/handler"
	"github.com/ezep02/rodeo/internal/analytics/repository"
	"github.com/ezep02/rodeo/internal/analytics/services"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AnalyticsRoutes(r chi.Router, db *gorm.DB, redisClient *redis.Client) {

	analyticsService := services.NewAnalyticServices(repository.NewAnalyticsRepository(db, redisClient))
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	r.Route("/analytics", func(r chi.Router) {
		r.Get("/revenue", analyticsHandler.GetMonthlyRevenueAndAvgHandler)
		r.Get("/appointments", analyticsHandler.GetMonthlyAppointmentsAndAvgHandler)
		r.Get("/customers", analyticsHandler.GetMonthlyNewCustomersAndAvgHandler)
		r.Get("/cancellations", analyticsHandler.GetMonthlyCancellationsAndAvgHandler)
		r.Get("/revenue/current-year", analyticsHandler.GetCurrentYearMonthlyRevenueHandler)
		r.Get("/popular-services", analyticsHandler.GetMonthlyPopularServicesHandler)
		r.Get("/frequent-customers", analyticsHandler.GetFrequentCustomersHandler)
	})

	r.Route("/barber", func(r chi.Router) {
		r.Get("/yearly-haircut", analyticsHandler.GetYearlyBarberHaircuts)
	})

}
