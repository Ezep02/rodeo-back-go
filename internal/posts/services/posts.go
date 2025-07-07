package services

import "github.com/ezep02/rodeo/internal/posts/repository"

type PostsService struct {
	Posts_repository *repository.PostsRepository
}

func NewPostsService(post_service *repository.PostsRepository) *PostsService {
	return &PostsService{
		Posts_repository: post_service,
	}
}
