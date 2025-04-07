package schedules

import (
	"github.com/ezep02/rodeo/internal/schedules/handler"
	"github.com/ezep02/rodeo/internal/schedules/repository"
	"github.com/ezep02/rodeo/internal/schedules/services"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SchedulesRoutes(r chi.Router, db *gorm.DB) {
	// Inicializar AuthRepository con la conexion a la DB

	sch_Repo := repository.NewSchedulesRepository(db)
	// Inicializar AuthService con el repositorio
	sch_Serv := services.NewOrderService(sch_Repo)
	// Inicializar AuthHandler con el servicio
	sch_Handler := handler.NewSchedulHandler(sch_Serv)

	// Rutas del módulo de autenticación
	r.Route("/schedules", func(r chi.Router) {
		r.Post("/", sch_Handler.BarberSchedulesHandler) // Success, TODO: manejar la respuestas con un websocket
		r.Get("/{limit}/{offset}", sch_Handler.GetBarberSchedulesHandler)
		r.Get("/{limit}/{offset}", sch_Handler.GetAvailableSchedulesHandler)
		r.HandleFunc("/updates", handler.HandleConnection)
	})
}
