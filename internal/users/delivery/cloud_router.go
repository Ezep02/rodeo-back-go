package delivery

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/ezep02/rodeo/internal/users/delivery/http"
	"github.com/ezep02/rodeo/internal/users/repository"
	"github.com/ezep02/rodeo/internal/users/usecase"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewCloudRouter(r *gin.RouterGroup, db *gorm.DB, redis *redis.Client, cloudConfig *cloudinary.Cloudinary) {

	log.Println("[CLOUD ROUTES] Setting up cloud routes")

	// Repositio u casos de uso de claudinary
	claudinaryRepo := repository.NewCloudinaryCloudRepo(cloudConfig, redis)
	cloudinarySvc := usecase.NewCloudService(claudinaryRepo)

	// Rutas de usuario
	cloudinary := r.Group("/cloudinary")
	{
		cloudinaryHandler := http.NewCloudinaryHandler(cloudinarySvc)
		cloudinary.GET("/images", cloudinaryHandler.Images)
		cloudinary.GET("/video", cloudinaryHandler.Video)
		cloudinary.POST("/upload", cloudinaryHandler.Upload)

	}
}
