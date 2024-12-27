package auth

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r chi.Router, db *gorm.DB) {
	// Inicializar AuthRepository con la conexion a la DB
	authRepo := NewAuthRepository(db)
	// Inicializar AuthService con el repositorio
	authServ := NewAuthService(authRepo)
	// Inicializar AuthHandler con el servicio
	authHandler := NewAuthHandler(authServ)

	// Rutas del módulo de autenticación
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.RegisterUserHandler)
		r.Post("/login", authHandler.LoginUserHandler)
		r.HandleFunc("/google", GoogleAuth)        // Inicia el flujo OAuth2
		r.HandleFunc("/callback", CallbackHandler) // Recibe el código de autorización
		r.Get("/verify-token", authHandler.VerifyTokenHandler)
		r.Get("/logout", authHandler.LogoutSession)
	})
}
