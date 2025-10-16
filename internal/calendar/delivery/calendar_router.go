package delivery

import (
	"log"

	"github.com/ezep02/rodeo/internal/calendar/delivery/http"
	"github.com/ezep02/rodeo/internal/calendar/repository"
	"github.com/ezep02/rodeo/internal/calendar/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewCalendarRouter(r *gin.RouterGroup, db *gorm.DB) {

	log.Println("[CALENDAR ROUTER] Setting up calendar routes")

	// Iniciar repositorio y servicio de calendar
	calendarRepo := repository.NewGormCalendarRepo(db)
	calendarSvc := usecase.NewCalendarService(calendarRepo)

	calendar := r.Group("/calendar")
	{
		calendarHandler := http.NewGoogleCalendarHandler(calendarSvc)
		calendar.GET("/google-calendar/login", calendarHandler.GoogleCalendarLogin)
		calendar.GET("/google-calendar/callback", calendarHandler.GoogleCalendarCallback)
		calendar.GET("/google-calendar/verify-status", calendarHandler.GoogleCalendarVerify)
		calendar.POST("/new", calendarHandler.Create)
	}

}
