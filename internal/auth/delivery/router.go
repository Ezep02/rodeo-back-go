package delivery

import (
	"log"

	"github.com/ezep02/rodeo/internal/auth/delivery/http"
	"github.com/ezep02/rodeo/internal/auth/repository"
	"github.com/ezep02/rodeo/internal/auth/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewAuthRoutes(r *gin.RouterGroup, cnn *gorm.DB) {

	log.Println("[AUTH ROUTES] Setting up authentication routes")

	//
	authRepo := repository.NewGormAuthRepo(cnn)
	authSvc := usecase.NewAuthService(authRepo)

	auth := r.Group("/auth")
	{
		authHandler := http.NewAuthHandler(authSvc)
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/logout", authHandler.Logout)
		auth.GET("/verify", authHandler.VerifySession)
		auth.GET("/google", authHandler.GoogleAuth)
		auth.GET("/callback", authHandler.CallbackHandler)
		auth.POST("/send-email", authHandler.SendResetPasswordEmail)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.PUT("/update-user/:id", authHandler.UpdateUser)
	}
}
