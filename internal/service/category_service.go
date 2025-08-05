package service

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/domain"
)

type CategoryService struct {
	authRepo domain.CategoryRepository
}

func NewCategoryService(categoryRepo domain.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) error {
	// 1. Validar que tenga nombre
	if category.Name == "" {
		return errors.New("el nombre es un campo requerido")
	}

	// 3. Crear categoria
	return s.authRepo.Create(ctx, category)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) error {
	// 1. Validar que tenga ID
	if category.ID == 0 {
		return errors.New("el ID de la categoria es un campo requerido")
	}

	// 2. Validar que tenga nombre
	if category.Name == "" {
		return errors.New("el nombre es un campo requerido")
	}

	// 4. Actualizar categoria
	return s.authRepo.Update(ctx, category)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uint) error {
	// 1. Validar que tenga ID
	if id == 0 {
		return errors.New("el ID de la categoria es un campo requerido")
	}

	// 2. Eliminar categoria
	return s.authRepo.Delete(ctx, id)
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]domain.Category, error) {
	// 1. Listar categorias
	categories, err := s.authRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
