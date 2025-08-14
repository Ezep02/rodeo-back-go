package http

import (
	"github.com/ezep02/rodeo/internal/service"
	"github.com/ezep02/rodeo/internal/transport/sse"
	"github.com/ezep02/rodeo/internal/transport/ws"
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
	cloudinarySvc *service.CloudinaryService,
	postSvc *service.PostService,
	categorySvc *service.CategoryService,
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

	// Hub websocket
	wsHub := ws.NewHub()
	wsHandler := ws.NewWSHandler(wsHub)

	// Hub SSE
	sseHub := sse.NewHub()
	sseHandler := sse.NewSSEHandler(sseHub)

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
			apptHandler := NewAppointmentHandler(appSvc, couponSvc, sseHandler, slotSvc)
			appts.POST("/", apptHandler.Create)
			appts.GET("/page/:start/:end", apptHandler.ListByDateRange)
			appts.GET("/:id", apptHandler.GetByID)
			appts.PUT("/:id", apptHandler.Update)
			appts.POST("/surcharge", apptHandler.Surcharge)
			appts.POST("/cancel/:id", apptHandler.Cancel)
			appts.POST("/reminder/:id", apptHandler.Reminder)

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
			prodHandler := NewProductHandler(prodSvc, categorySvc)
			products.GET("/", prodHandler.List)
			products.POST("/", prodHandler.Create)
			products.GET("/:id", prodHandler.GetByID)
			products.PUT("/:id", prodHandler.Update)
			products.DELETE("/:id", prodHandler.Delete)
			products.GET("/popular", prodHandler.Popular)
			products.GET("/promotion", prodHandler.Promotion)
		}

		// Rutas de Slots
		slots := v1.Group("/slots")
		{
			slotHandler := NewSlotHandler(slotSvc)
			slots.GET("/page/:start/:end", slotHandler.ListByDateRange) // cambiar por page offset
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
			reviews.GET("/product/:id", reviewHandler.ListByProductID)
			reviews.GET("/user/:id/page/:offset", reviewHandler.ListByUserID)
		}

		// Rutas de analiticas
		analytics := v1.Group("/analytics")
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
			cloudinaryHandler := NewCloudinaryHandler(cloudinarySvc)
			cloudinary.GET("/images", cloudinaryHandler.Images)
			cloudinary.GET("/video", cloudinaryHandler.Video)
			cloudinary.POST("/upload", cloudinaryHandler.Upload)
		}

		// Rutas de Posts
		posts := v1.Group("/posts")
		{
			postHandler := NewPostHandler(postSvc)
			posts.GET("/page/:offset", postHandler.List)
			posts.GET("/:id", postHandler.GetByID)
			posts.POST("/", postHandler.Create)
			posts.PUT("/:id", postHandler.Update)
			posts.DELETE("/:id", postHandler.Delete)
			posts.GET("/count", postHandler.Count)
		}

		// Rutas de categorias
		categories := v1.Group("/categories")
		{
			categoryHandler := NewCategoryHandler(categorySvc)
			categories.POST("/", categoryHandler.CreateCategory)
			categories.PUT("/", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
			categories.GET("/", categoryHandler.ListCategories)
		}

		// Conexion websocket
		v1.GET("/ws", wsHandler.HandleWS)

		// Conexion SSE streaming de datos
		v1.GET("/stream", sseHandler.Handle)
	}

	return r
}
