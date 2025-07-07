package posts

import (
	"net/http"

	"github.com/ezep02/rodeo/internal/posts/handler"
	"github.com/ezep02/rodeo/internal/posts/repository"
	"github.com/ezep02/rodeo/internal/posts/services"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ServicesRouter(mux *http.ServeMux, db *gorm.DB, redis *redis.Client) {
	// Inicializar componentes
	posts_repository := repository.NewPostsRepository(db, redis)
	posts_service := services.NewPostsService(posts_repository)
	_ = handler.NewOrderHandler(posts_service)

}
