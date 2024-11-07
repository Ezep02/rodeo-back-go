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

	r.Route("/service", func(r chi.Router) {
		r.Post("/new", srv_handler.CreateService)
		r.Get("/all", srv_handler.GetAllServices)
	})

}
