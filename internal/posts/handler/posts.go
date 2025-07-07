package handler

import (
	"context"

	"github.com/ezep02/rodeo/internal/posts/services"
)

type PostsHandler struct {
	ord_srv *services.PostsService
	ctx     context.Context
}

func NewOrderHandler(posts_srv *services.PostsService) *PostsHandler {
	return &PostsHandler{
		ctx:     context.Background(),
		ord_srv: posts_srv,
	}
}
