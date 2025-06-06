package auth

import (
	"net/http"

	"github.com/ezep02/rodeo/internal/auth/handlers"
	"github.com/ezep02/rodeo/internal/auth/repository"
	"github.com/ezep02/rodeo/internal/auth/services"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(mux *http.ServeMux, db *gorm.DB) {
	// Inicializar AuthRepository con la conexion a la DB
	authRepo := repository.NewAuthRepository(db)
	// Inicializar AuthService con el repositorio
	authServ := services.NewAuthService(authRepo)
	// Inicializar AuthHandler con el servicio
	authHandler := handlers.NewAuthHandler(authServ)

	mux.HandleFunc("/auth/auth", authHandler.RegisterUserHandler)
	mux.HandleFunc("/auth/login", authHandler.LoginUserHandler)

	mux.HandleFunc("/auth/google", handlers.GoogleAuth)
	mux.HandleFunc("/auth/callback", handlers.CallbackHandler)

	mux.HandleFunc("/auth/verify-token", authHandler.VerifyTokenHandler)
	mux.HandleFunc("/auth/logout", authHandler.LogoutSession)

	mux.HandleFunc("/auth/send-email", authHandler.SendResetUserPasswordEmailHandler)
	mux.HandleFunc("/auth/reset-password", authHandler.ResetUserPassword)
}
