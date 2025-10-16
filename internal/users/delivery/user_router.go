package delivery

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/ezep02/rodeo/internal/users/delivery/http"
	"github.com/ezep02/rodeo/internal/users/repository"
	"github.com/ezep02/rodeo/internal/users/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(r *gin.RouterGroup, db *gorm.DB, redis *redis.Client, cloudConfig *cloudinary.Cloudinary) {

	log.Println("[USER ROUTES] Setting up user routes")

	userRepo := repository.NewGormUserRepo(db, redis)
	userSvc := usecase.NewUserService(userRepo)

	// Repositio u casos de uso de claudinary
	claudinaryRepo := repository.NewCloudinaryCloudRepo(cloudConfig, redis)
	cloudinarySvc := usecase.NewCloudService(claudinaryRepo)

	// Rutas de usuario
	users := r.Group("/users")
	{
		userHandler := http.NewUserHandler(userSvc, cloudinarySvc)
		users.PUT("/:id", userHandler.Update) // TODO: Actualizar datos del usuario
		users.GET("/:id", userHandler.GetByID)
		users.GET("/info", userHandler.UserInfo)
		users.PUT("/username/:id", userHandler.UpdateUsername)
		users.PUT("/password/:id", userHandler.UpdatePassword)
		users.POST("/avatar", userHandler.UploadAvatar)
	}

	// Repositorio y caso de uso de barberos
	barberRepo := repository.NewGormBarberRepo(db, redis)
	barberSvc := usecase.NewBarberService(barberRepo)

	// Rutas de barberos
	barbers := r.Group("/barbers")
	{
		barberHandler := http.NewBarberHandler(barberSvc)
		barbers.GET("/:id", barberHandler.GetByID)
		barbers.GET("/all", barberHandler.List)
	}
}
