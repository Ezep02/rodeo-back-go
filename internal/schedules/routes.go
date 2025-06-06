package schedules

import (
	"net/http"

	"github.com/ezep02/rodeo/internal/schedules/handler"
	"github.com/ezep02/rodeo/internal/schedules/repository"
	"github.com/ezep02/rodeo/internal/schedules/services"
	"gorm.io/gorm"
)

func SchedulesRoutes(mux *http.ServeMux, db *gorm.DB) {

	sch_Repo := repository.NewSchedulesRepository(db)
	sch_Serv := services.NewOrderService(sch_Repo)
	sch_Handler := handler.NewSchedulHandler(sch_Serv)

	mux.HandleFunc("/schedules", sch_Handler.BarberSchedulesHandler)
	mux.HandleFunc("/schedules/barber/", sch_Handler.GetBarberSchedulesHandler)
	mux.HandleFunc("/schedules/available/", sch_Handler.GetAvailableSchedulesHandler)
	mux.HandleFunc("/schedules/updates", handler.HandleConnection)
}
