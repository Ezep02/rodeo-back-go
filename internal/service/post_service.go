package service

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/domain"
)

type PostService struct {
	postRepository domain.PostRepository
}

func NewPostService(svc domain.PostRepository) *PostService {
	return &PostService{svc}
}

func (s PostService) Create(ctx context.Context, post *domain.Post) error {

	// 1. Validar que tenga un título
	if post.Title == "" {
		return errors.New("el título es obligatorio")
	}

	return s.postRepository.Create(ctx, post)
}

func (s PostService) List(ctx context.Context, offset int) ([]domain.Post, error) {
	// 1. Llamar al repositorio para obtener la lista de posts
	posts, err := s.postRepository.List(ctx, offset)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s PostService) Update(ctx context.Context, post *domain.Post, post_id uint) error {
	// 1. Validar que el post exista
	_, err := s.postRepository.GetByID(ctx, post_id)
	if err != nil {
		if err == domain.ErrNotFound {
			return errors.New("post no encontrado")
		}
		return err
	}

	// 3. Llamar al repositorio para actualizar el post
	return s.postRepository.Update(ctx, post, post_id)
}

func (s PostService) Delete(ctx context.Context, post_id uint) error {
	// 1. Validar que el post exista
	_, err := s.postRepository.GetByID(ctx, post_id)
	if err != nil {
		if err == domain.ErrNotFound {
			return errors.New("post no encontrado")
		}
		return err
	}

	// 2. Llamar al repositorio para eliminar el post
	return s.postRepository.Delete(ctx, post_id)
}

func (s PostService) GetByID(ctx context.Context, post_id uint) (*domain.Post, error) {
	// 1. Llamar al repositorio para obtener el post por ID
	post, err := s.postRepository.GetByID(ctx, post_id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, errors.New("post no encontrado")
		}
		return nil, err
	}

	return post, nil
}

func (s PostService) Count(ctx context.Context) (int64, error) {
	// 1. Llamar al repositorio para contar los posts
	count, err := s.postRepository.Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}
