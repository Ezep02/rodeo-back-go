package http

import (
	"github.com/ezep02/rodeo/internal/middleware"
	"github.com/ezep02/rodeo/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	appSvc *service.AppointmentService,
	prodSvc *service.ProductService,
	authSvc *service.AuthService,
	slotSvc *service.SlotService,
	revSvc *service.ReviewService,
	analyticSvc *service.AnalyticService,
	couponSvc *service.CouponService,
	infoSvc *service.InformationService,
) *gin.Engine {

	r := gin.Default()

	// Middleware de CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API V1
	v1 := r.Group("/api/v1")
	{

		// Rutas de autenticacion
		auth := v1.Group("/auth")
		{
			authHandler := NewAuthHandler(authSvc)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/logout", authHandler.Logout)
			auth.GET("/verify", authHandler.VerifySession)
			auth.GET("/google", authHandler.GoogleAuth)
			auth.GET("/callback", authHandler.CallbackHandler)
		}

		// Rutas de Appointment
		appts := v1.Group("/appointments")
		{
			apptHandler := NewAppointmentHandler(appSvc, couponSvc)
			appts.POST("/", apptHandler.Create)
			appts.GET("/", apptHandler.List)
			appts.GET("/:id", apptHandler.GetByID)
			appts.PUT("/:id", apptHandler.Update)
			appts.POST("/surcharge", apptHandler.Surcharge)
			appts.DELETE("/:id", apptHandler.Cancel)

			// Rutas especificas de los usuarios
			appts.GET("/user/:id", apptHandler.GetByUserID)
		}

		mercado_pago := v1.Group("/mercado_pago")
		{
			mepHandler := NewMepaHandler(prodSvc, appSvc, slotSvc)
			mercado_pago.POST("/", mepHandler.CreatePreference)
			mercado_pago.POST("/surcharge", mepHandler.CreateSurchargePreference)
			mercado_pago.GET("/:token", mepHandler.GetPayment)
		}

		// Rutas de Product
		products := v1.Group("/products")
		{
			prodHandler := NewProductHandler(prodSvc)
			products.GET("/", prodHandler.List)
			products.POST("/", prodHandler.Create)
			products.GET("/:id", prodHandler.GetByID)
			products.PUT("/:id", prodHandler.Update)
			products.DELETE("/:id", prodHandler.Delete)
			products.GET("/popular", prodHandler.Popular)
		}

		// Rutas de Slots
		slots := v1.Group("/slots")
		{
			slotHandler := NewSlotHandler(slotSvc)
			slots.GET("/:offset", slotHandler.List)
			slots.GET("/date/:id", slotHandler.ListByDate)
			slots.POST("/", slotHandler.Create)
			slots.DELETE("/", slotHandler.Delete)
			slots.PUT("/", slotHandler.Update)
		}

		// Rutas de reviews
		reviews := v1.Group("/reviews")
		{
			reviewHandler := NewReviewHandler(revSvc)
			reviews.POST("/", reviewHandler.Create)
			reviews.GET("/", reviewHandler.List)
		}

		// Rutas de analiticas
		analytics := v1.Group("/analytics").Use(middleware.AuthorizeAdmin())
		{
			analyticHandler := NewAnalyticHandler(analyticSvc)
			analytics.GET("/booking-rate", analyticHandler.BookingOcupationRate)
			analytics.GET("/booking-count", analyticHandler.MonthBookingCount)
			analytics.GET("/booking-weekly-rate", analyticHandler.WeeklyBookingRate)
			analytics.GET("/month-revenue", analyticHandler.MonthlyRevenue)
			analytics.GET("/client-rate", analyticHandler.NewClientRate)
			analytics.GET("/slot-popular-time", analyticHandler.PopularTimeSlot)
		}

		info := v1.Group("/info")
		{
			infoHandler := NewInfoHandler(infoSvc)
			info.GET("/", infoHandler.Information)
		}

		// Rutas de cupones
		coupons := v1.Group("/coupons")
		{
			couponHandler := NewCouponHandler(couponSvc)
			coupons.POST("/", couponHandler.Create)
		}

		// Rutas de cloudinary
		cloudinary := v1.Group("/cloudinary")
		{

			cloudinary.GET("/images", GetCloudinaryImages)
		}
	}

	return r
}
