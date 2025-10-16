package usecase

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/catalog/domain/categorie"
)

type CategoryService struct {
	categorieRepo categorie.CategorieRepository
}

func NewCategorieService(categoryRepo categorie.CategorieRepository) *CategoryService {
	return &CategoryService{categoryRepo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, categorie *categorie.Categorie) error {
	// 1. Validar que tenga nombre
	if categorie.Name == "" {
		return errors.New("el nombre es un campo requerido")
	}

	// 3. Crear categoria
	return s.categorieRepo.Create(ctx, categorie)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id uint, data *categorie.Categorie) error {
	// 1. Validar que tenga ID
	if id == 0 {
		return errors.New("el ID de la categoria es un campo requerido")
	}

	// 2. Validar que tenga nombre
	if data.Name == "" {
		return errors.New("el nombre es un campo requerido")
	}

	// 4. Actualizar categoria
	return s.categorieRepo.Update(ctx, id, data)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uint) error {
	// 1. Validar que tenga ID
	if id == 0 {
		return errors.New("el ID es un campo requerido")
	}

	// 2. Eliminar categoria
	return s.categorieRepo.Delete(ctx, id)
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]categorie.Categorie, error) {
	// 1. Listar categorias
	categories, err := s.categorieRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id uint) (*categorie.Categorie, error) {
	// 1. Validar que tenga ID
	if id == 0 {
		return nil, errors.New("el ID de la categoria es un campo requerido")
	}

	// 2. Obtener categoria por ID
	category, err := s.categorieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}
