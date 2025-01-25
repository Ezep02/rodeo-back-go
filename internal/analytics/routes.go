package analytics

import (
	"github.com/ezep02/rodeo/internal/analytics/handler"
	"github.com/ezep02/rodeo/internal/analytics/repository"
	"github.com/ezep02/rodeo/internal/analytics/services"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func AnalyticsRoutes(r chi.Router, db *gorm.DB) {

	analyticsRepo := repository.NewAnalyticsRepository(db)

	analyticsServices := services.NewAnalyticsServices(analyticsRepo)

	analyticsHandler := handler.NewAnalyticsHandler(analyticsServices)

	// Rutas del módulo de autenticación
	r.Route("/analytics", func(r chi.Router) {
		r.Get("/", analyticsHandler.GetRevenues)
		r.Get("/users", analyticsHandler.GetTotalUsers)
		r.Get("/recived-clients", analyticsHandler.GetRevicedTotalUsers)

	})

	r.Route("/analytics/expense", func(r chi.Router) {
		r.Post("/new", analyticsHandler.NewExpense)
		r.Get("/historial/{limit}/{offset}", analyticsHandler.GetExpensesList)
		r.Get("/total", analyticsHandler.GetTotalExpensesCount)
	})
}
