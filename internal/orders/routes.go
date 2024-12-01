package orders

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func OrderRoutes(r chi.Router, db *gorm.DB) {
	// Iniciar OrderRepository
	order_repo := NewOrderRepository(db)
	// Iniciar OrderService
	order_srv := NewOrderService(order_repo)
	// Inicializar OrderHandler con el servicio
	orderHandler := NewOrderHandler(order_srv)

	// Rutas del módulo de autenticacion
	r.Route("/order", func(r chi.Router) {
		r.Post("/new", orderHandler.CreateOrderHandler)
		r.Post("/webhook", orderHandler.WebHook)
		r.Get("/list/{limit}/{offset}", orderHandler.GetOrders)
		r.Get("/success", orderHandler.Success)
		r.HandleFunc("/notification", HandleConnection)
	})
}
