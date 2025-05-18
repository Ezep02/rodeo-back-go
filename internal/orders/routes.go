package orders

import (
	"github.com/ezep02/rodeo/internal/orders/handler"
	"github.com/ezep02/rodeo/internal/orders/repository"
	"github.com/ezep02/rodeo/internal/orders/services"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func OrderRoutes(r chi.Router, db *gorm.DB, redis *redis.Client) {
	// Iniciar OrderRepository
	order_repo := repository.NewOrderRepository(db, redis)
	// Iniciar OrderService
	order_srv := services.NewOrderService(order_repo)
	// Inicializar OrderHandler con el servicio
	orderHandler := handler.NewOrderHandler(order_srv)

	// Rutas del módulo de ordenes
	r.Route("/order", func(r chi.Router) {
		r.Post("/new", orderHandler.CreateOrderHandler)
		r.Post("/webhook", orderHandler.WebHook)
		r.Get("/pending/{limit}/{offset}", orderHandler.GetBarberPendingOrdersHandler)
		r.HandleFunc("/notification", handler.HandleConnection)
	})

	r.Route("/order/customer", func(r chi.Router) {
		r.Get("/", orderHandler.CustomerPendingOrderHandler)
		r.Post("/success", orderHandler.GetSuccessPaymentHandler)
		r.Post("/refund", orderHandler.CreateRefundHandler)
		r.Post("/reschedule", orderHandler.CreateReschedule)
		r.Get("/coupons", orderHandler.GetCouponsHandler)
		r.HandleFunc("/notification", handler.HandleConnection)
		r.HandleFunc("/notification-coupon", handler.HandleConnection)
	})
}
