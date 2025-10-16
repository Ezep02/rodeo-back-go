package delivery

import (
	"log"

	"github.com/ezep02/rodeo/internal/catalog/delivery/http"
	"github.com/ezep02/rodeo/internal/catalog/repository"
	"github.com/ezep02/rodeo/internal/catalog/usecase"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewCatalogRoutes(r *gin.RouterGroup, cnn *gorm.DB, redis *redis.Client) {

	log.Println("[SERVICES ROUTES] Setting up services routes")

	// Repositorio y caso de uso de Servicios
	serviceRepo := repository.NewGormServiceRepo(cnn, redis)
	serviceSvc := usecase.NewServiceService(serviceRepo)

	// Repositorio y caso de uso de Categorias
	categorieRepo := repository.NewGormCategorieRepo(cnn, redis)
	categorieSvc := usecase.NewCategorieService(categorieRepo)

	// Repositorio y case de uso de Promociones
	promoRepo := repository.NewGormPromoRepo(cnn, redis)
	promoSvc := usecase.NewPromoService(promoRepo)

	// Repositorio y caso de uso de medias
	mediaRepo := repository.NewGormMediaRepo(cnn, redis)
	mediaSvc := usecase.NewMediaService(mediaRepo)

	// Rutas de Product
	services := r.Group("/services")
	{
		svcHandler := http.NewServiceHandler(serviceSvc)
		services.GET("/page/:offset", svcHandler.List)
		services.POST("/", svcHandler.Create)
		services.GET("/:id", svcHandler.GetByID)
		services.PUT("/:id", svcHandler.Update)
		services.DELETE("/:id", svcHandler.Delete)
		services.GET("/popular", svcHandler.Popular)
		services.GET("/stats", svcHandler.Stats)
		services.POST("/categories/:id/add", svcHandler.AddCategories)
		services.POST("/categories/:id/remove", svcHandler.RemoveCategories)

	}

	// TODO: Aqui crear los endpoint necesarios para realizar operaciones crud para PROMOCIONES
	promo := r.Group("/promotion")
	{
		promoHandler := http.NewPromoHandler(promoSvc)
		promo.POST("/", promoHandler.Create)
		promo.GET("/page/:id/:offset", promoHandler.ListByServiceId)
		promo.PUT("/:id", promoHandler.Update)
		promo.DELETE("/:id", promoHandler.Delete)
	}

	// Rutas de Category
	categories := r.Group("/categories")
	{
		categorieHandler := http.NewCategorieHandler(categorieSvc)
		categories.POST("/", categorieHandler.CreateCategory)
		categories.PUT("/:id", categorieHandler.UpdateCategory)
		categories.DELETE("/:id", categorieHandler.DeleteCategory)
		categories.GET("/", categorieHandler.ListCategories)
	}

	medias := r.Group("/media")
	{
		mediaHandler := http.NewMediaHandler(mediaSvc)
		medias.POST("/:id", mediaHandler.SetMedia)
		medias.PUT("/:id", mediaHandler.Update)
		medias.DELETE("/:id", mediaHandler.DeleteMedia)
		medias.GET("/:id", mediaHandler.ListByServiceId)
	}
}
