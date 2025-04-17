package auth

import (
	"github.com/ezep02/rodeo/internal/auth/handlers"
	"github.com/ezep02/rodeo/internal/auth/repository"
	"github.com/ezep02/rodeo/internal/auth/services"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r chi.Router, db *gorm.DB) {
	// Inicializar AuthRepository con la conexion a la DB
	authRepo := repository.NewAuthRepository(db)
	// Inicializar AuthService con el repositorio
	authServ := services.NewAuthService(authRepo)
	// Inicializar AuthHandler con el servicio
	authHandler := handlers.NewAuthHandler(authServ)

	// Rutas del módulo de autenticación
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.RegisterUserHandler)
		r.Post("/login", authHandler.LoginUserHandler)
		r.HandleFunc("/google", handlers.GoogleAuth)
		r.HandleFunc("/callback", handlers.CallbackHandler)
		r.Get("/verify-token", authHandler.VerifyTokenHandler)
		r.Get("/logout", authHandler.LogoutSession)
		r.Post("/send-email", authHandler.SendResetUserPasswordEmailHandler)
		r.Post("/reset-password", authHandler.ResetUserPassword)
	})
}
