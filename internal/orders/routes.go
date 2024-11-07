package orders

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func OrderRoutes(r chi.Router, db *gorm.DB) {

	// Inicializar OrderHandler con el servicio
	orderHandler := NewOrderHandler()

	// Rutas del módulo de autenticación
	r.Route("/order", func(r chi.Router) {
		r.Post("/new", orderHandler.CreateOrderHandler)
		r.Post("/webhook", orderHandler.WebHook)
	})
}
