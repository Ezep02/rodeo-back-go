package service

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/domain"
)

type ProductService struct {
	prodRepo domain.ProductRepository
}

func NewProductService(prodRepo domain.ProductRepository) *ProductService {
	return &ProductService{prodRepo}
}

func (s *ProductService) Create(ctx context.Context, product *domain.Product) error {

	// 1. Validar que el producto tenga un nombre
	if product.Name == "" {
		return errors.New("el producto debe tener un nombre")
	}

	// 2. Validar que el producto tenga un precio positivo
	if product.Price <= 0 {
		return errors.New("el producto debe tener un precio mayor o igual a cero")
	}

	// 3. Validar que el producto no exista
	existing, err := s.prodRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, p := range existing {
		if p.Name == product.Name {
			return errors.New("el producto ya existe")
		}
	}

	// 4. Crear el producto
	return s.prodRepo.Create(ctx, product)
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	return s.prodRepo.GetByID(ctx, id)
}

func (s *ProductService) ListAll(ctx context.Context) ([]domain.Product, error) {
	return s.prodRepo.List(ctx)
}

func (s *ProductService) Update(ctx context.Context, id uint, updatedProd *domain.Product) error {
	// 1. Validar que el producto tenga un nombre
	if updatedProd.Name == "" {
		return errors.New("el producto debe tener un nombre")
	}

	// 2. Validar que el producto tenga un precio positivo
	if updatedProd.Price <= 0 {
		return errors.New("el producto debe tener un precio mayor o igual a cero")
	}

	// 3. Verificar si existe un producto con el mismo nombre
	allProducts, err := s.prodRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, p := range allProducts {
		// Si el nombre es igual y el ID es diferente, significa que ya existe otro
		if p.Name == updatedProd.Name && p.ID != id {
			return errors.New("ya existe un producto con ese nombre")
		}
	}

	// 4. Manetener el id original
	updatedProd.ID = id

	return s.prodRepo.Update(ctx, updatedProd)
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	// 1. Verificar que el producto exista
	_, err := s.prodRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Eliminar el producto
	return s.prodRepo.Delete(ctx, id)
}

func (s *ProductService) Popular(ctx context.Context) ([]domain.Product, error) {
	return s.prodRepo.Popular(ctx)
}
