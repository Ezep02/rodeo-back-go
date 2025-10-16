package http

import (
	"github.com/cloudinary/cloudinary-go/v2"

	analyticsRouter "github.com/ezep02/rodeo/internal/analytics/delivery"
	bookingRouter "github.com/ezep02/rodeo/internal/auth/delivery"
	apptRouter "github.com/ezep02/rodeo/internal/booking/delivery"
	calendarRouter "github.com/ezep02/rodeo/internal/calendar/delivery"
	catalogRouter "github.com/ezep02/rodeo/internal/catalog/delivery"
	slotRouter "github.com/ezep02/rodeo/internal/slots/delivery"
	userRouter "github.com/ezep02/rodeo/internal/users/delivery"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, cloud *cloudinary.Cloudinary, redis *redis.Client) *gin.Engine {

	r := gin.Default()

	// Middleware de CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api/v1")

	// Inicializa los controladores y rutas
	bookingRouter.NewAuthRoutes(api, db)
	apptRouter.NewAppointmentRoutes(api, db, redis)
	analyticsRouter.NewAnalyticsRoutes(api, db, redis)
	calendarRouter.NewCalendarRouter(api, db)
	userRouter.NewUserRouter(api, db, redis, cloud)
	userRouter.NewCloudRouter(api, db, redis, cloud)
	catalogRouter.NewCatalogRoutes(api, db, redis)
	slotRouter.NewSlotRouter(api, db, redis)

	return r
}
