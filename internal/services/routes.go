package services

import (
	"net/http"

	handler "github.com/ezep02/rodeo/internal/services/handlers"
	"github.com/ezep02/rodeo/internal/services/repository"
	"github.com/ezep02/rodeo/internal/services/services"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ServicesRouter(mux *http.ServeMux, db *gorm.DB, redis *redis.Client) {
	// Inicializar componentes
	srv_repository := repository.NewServiceRepository(db, redis)
	srv_service := services.NewSrvRepository(srv_repository)
	srv_handler := handler.NewServiceHandler(srv_service)

	// Handlers estáticos
	mux.HandleFunc("/services/new", srv_handler.CreateService)
	mux.HandleFunc("/services/popular-services", srv_handler.GetPopularServices)
	mux.HandleFunc("/services/notification-update", handler.HandleConnection)

	// Handlers con parámetros (deben parsear desde la URL dentro del handler)
	mux.HandleFunc("/services/", srv_handler.GetServices)              // para /services/{limit}/{offset}
	mux.HandleFunc("/services/barber/", srv_handler.GetBarberServices) // para /services/barber/{limit}/{offset}
	mux.HandleFunc("/services/update/", srv_handler.UpdateServices)    // para /services/update/{id}
	mux.HandleFunc("/services/delete/", srv_handler.DeleteServiceByID) // para /services/{id} en DELETE
}
