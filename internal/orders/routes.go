package orders

import (
	"github.com/ezep02/rodeo/internal/orders/handler"
	"github.com/ezep02/rodeo/internal/orders/repository"
	"github.com/ezep02/rodeo/internal/orders/services"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func OrderRoutes(r chi.Router, db *gorm.DB) {
	// Iniciar OrderRepository
	order_repo := repository.NewOrderRepository(db)
	// Iniciar OrderService
	order_srv := services.NewOrderService(order_repo)
	// Inicializar OrderHandler con el servicio
	orderHandler := handler.NewOrderHandler(order_srv)

	// Rutas del módulo de autenticacion
	r.Route("/order", func(r chi.Router) {
		r.Post("/new", orderHandler.CreateOrderHandler)
		r.Post("/webhook", orderHandler.WebHook)
		r.Get("/pending/{limit}/{offset}", orderHandler.GetBarberPendingOrdersHandler)
		r.Get("/success", orderHandler.Success)
		// r.Get("/pending", orderHandler.GetPendingOrder)
		r.Post("/refound/{id}/{amount}", orderHandler.Refound)
		r.Get("/historial/{limit}/{offset}", orderHandler.GetOrderHistorial)
		r.HandleFunc("/notification", handler.HandleConnection)
	})
}
