package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PostsRepository struct {
	Connection      *gorm.DB
	RedisConnection *redis.Client
}

func NewPostsRepository(DATABASE *gorm.DB, REDIS *redis.Client) *PostsRepository {
	return &PostsRepository{
		Connection:      DATABASE,
		RedisConnection: REDIS,
	}
}

func (r *PostsRepository) CreatingPost() {

}
