package delivery

import (
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/booking/delivery/http"
	"github.com/ezep02/rodeo/internal/booking/delivery/sse"
	"github.com/ezep02/rodeo/internal/booking/repository"
	"github.com/ezep02/rodeo/internal/booking/usecases"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewAppointmentRoutes(r *gin.RouterGroup, cnn *gorm.DB, redis *redis.Client) {

	log.Println("[APPOINTMENT ROUTES] Setting up appointment routes")

	// Hub SSE
	sseHub := sse.NewHub()
	sseHandler := sse.NewSSEHandler(sseHub)

	// Respositorios y casos de uso de Cupones
	couponRepo := repository.NewGormCouponRepo(cnn, redis)
	couponSvc := usecases.NewCouponService(couponRepo)

	// Respositorios y casos de uso de Payment
	paymentRepo := repository.NewGormPaymentRepo(cnn, redis)
	paymentSvc := usecases.NewPaymentService(paymentRepo)

	// Respositorios y casos de uso de Bookings
	bookingRepo := repository.NewGormBookingRepo(cnn, redis)
	bookingSvc := usecases.NewBookingService(bookingRepo, paymentRepo, couponRepo)

	// Respositorios y casos de uso de Servicios
	svcRepo := repository.NewGormServiceRepo(cnn, redis)
	serviceSvc := usecases.NewServicesService(svcRepo)

	// Repositorio y casos de usos de Mep
	mepSvc := usecases.NewMepService(bookingRepo, paymentRepo, svcRepo)

	// Job para cancelar las reservas que no fueron pagados aun
	bookingRepo.StartBookingCleanupJob(15 * time.Minute)

	booking := r.Group("/appointment")
	{
		bookingHandler := http.NewBookingHandler(bookingSvc, paymentSvc, couponSvc, serviceSvc)

		booking.GET("/upcoming/:date/:barber", bookingHandler.Upcoming)
		booking.GET("/stats/:id", bookingHandler.StatsByBarberID)
		booking.GET("/all/pending-payment", bookingHandler.AllPendingPayment)
		booking.PUT("/mark-as-paid/:id", bookingHandler.MarkAsPaid)
		booking.PUT("/mark-as-rejected/:id", bookingHandler.MarkAsRejected)

		// Crear una reserva sin mercado pago (creada cuando se la opcion de pago con alias es seleccionada)
		booking.POST("/", bookingHandler.Create)

		// Listado de citas
		booking.GET("/user/:id", bookingHandler.AllByUserId)

		// Reprogramacion de turno
		booking.POST("/user/reschedule", bookingHandler.Reschedule)

		// Cancelacion de turno
		booking.PUT("/user/cancel/:id", bookingHandler.Cancel)
		booking.GET("/user/cancel/verify/:id", bookingHandler.PreviewCancelation)

		// Obtener payment de una reserva
		booking.GET("/payment/:id", bookingHandler.BookingPayment)
	}

	// Rutas de cupones

	// Mercado Pago
	mercado_pago := r.Group("/mercado_pago")
	{
		mepHandler := http.NewMepaHandler(bookingSvc, paymentSvc, couponSvc, serviceSvc, mepSvc)
		mercado_pago.POST("/", mepHandler.CreatePreference)
		mercado_pago.POST("/notification", mepHandler.HandleNotification)
		mercado_pago.POST("/notification/reschedule", mepHandler.RescheduleWithSurcharge)
	}

	// Conexion SSE streaming de datos
	r.GET("/stream", sseHandler.Handle)
}
