package services

import (
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ServicesRouter(r chi.Router, db *gorm.DB, redis *redis.Client) {

	// Inicializar srv Repository
	srv_repository := NewServiceRepository(db, redis)

	// Inicializar srv Service
	srv_service := NewSrvRepository(srv_repository)

	// Inicializar srv Handler
	srv_handler := NewServiceHandler(srv_service)

	r.Route("/services", func(r chi.Router) {
		r.Get("/{limit}/{offset}", srv_handler.GetServices)
		r.Post("/new", srv_handler.CreateService)
		r.Get("/barber/{limit}/{offset}", srv_handler.GetBarberServices)
		r.Put("/update/{id}", srv_handler.UpdateServices)
		r.Delete("/{id}", srv_handler.DeleteServiceByID)
		r.HandleFunc("/notification-update", HandleConnection)
	})

}
