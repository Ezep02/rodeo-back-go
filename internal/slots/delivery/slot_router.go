package delivery

import (
	"log"

	"github.com/ezep02/rodeo/internal/slots/delivery/http"
	"github.com/ezep02/rodeo/internal/slots/repository"
	"github.com/ezep02/rodeo/internal/slots/usecase"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewSlotRouter(r *gin.RouterGroup, db *gorm.DB, redis *redis.Client) {

	log.Println("[SLOT ROUTES] Setting up slot routes")

	// Repositio u casos de uso de claudinary
	slotRepo := repository.NewGormSlotsRepo(db, redis)
	slotSvc := usecase.NewSlotUsecase(slotRepo)

	// Rutas de usuario
	slot := r.Group("/slot")
	{
		slotHandler := http.NewSlotHandler(slotSvc)
		slot.POST("/", slotHandler.Create)
		slot.PUT("/:id", slotHandler.Update)
		slot.GET("/range/:start/:end/:barber", slotHandler.GetByDateRange)
	}
}
