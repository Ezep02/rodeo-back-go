package services

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func ServicesRouter(r chi.Router, db *gorm.DB) {

	// Inicializar srv Repository
	srv_repository := NewServiceRepository(db)

	// Inicializar srv Service
	srv_service := NewSrvRepository(srv_repository)

	// Inicializar srv Handler
	srv_handler := NewServiceHandler(srv_service)

	r.Route("/services", func(r chi.Router) {
		r.Post("/new", srv_handler.CreateService)
		r.Get("/{limit}/{offset}", srv_handler.GetServices)
		r.Get("/barber/{limit}/{offset}", srv_handler.GetBarberServices)
		r.Put("/update/{id}", srv_handler.UpdateServices)
		r.Delete("/{id}", srv_handler.DeleteServiceByID)
		r.HandleFunc("/notification-update", HandleConnection)
	})

}
