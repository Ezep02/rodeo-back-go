package orders

import (
	"net/http"

	"github.com/ezep02/rodeo/internal/orders/handler"
	"github.com/ezep02/rodeo/internal/orders/repository"
	"github.com/ezep02/rodeo/internal/orders/services"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func OrderRoutes(mux *http.ServeMux, db *gorm.DB, redis *redis.Client) {

	// Iniciar OrderRepository
	order_repo := repository.NewOrderRepository(db, redis)
	// Iniciar OrderService
	order_srv := services.NewOrderService(order_repo)
	// Inicializar OrderHandler con el servicio
	orderHandler := handler.NewOrderHandler(order_srv)

	// Rutas del módulo de ordenes
	mux.HandleFunc("/order/new", orderHandler.CreateOrderHandler)
	mux.HandleFunc("/order/webhook", orderHandler.WebHook)
	mux.HandleFunc("/order/pending/", orderHandler.GetBarberPendingOrdersHandler)
	mux.HandleFunc("/order/notification", handler.HandleConnection)

	// Rutas ordenes de clientes
	mux.HandleFunc("/order/customer", orderHandler.CustomerPendingOrderHandler)
	mux.HandleFunc("/order/customer/success", orderHandler.GetSuccessPaymentHandler)
	mux.HandleFunc("/order/customer/refund", orderHandler.CreateRefundHandler)
	mux.HandleFunc("/order/customer/reschedule", orderHandler.CreateReschedule)
	mux.HandleFunc("/order/customer/coupons", orderHandler.GetCouponsHandler)
	mux.HandleFunc("/order/customer/notification", handler.HandleConnection)
	mux.HandleFunc("/order/customer/notification-coupon", handler.HandleConnection)
	mux.HandleFunc("/order/customer/previous/", orderHandler.GetCustomerPreviousOrdersHandler)

	// Rutas reviews de los clientes
	mux.HandleFunc("/review/new", orderHandler.SetReviewHandler)
	mux.HandleFunc("/review/all/", orderHandler.GetReviewsHandler)
}
