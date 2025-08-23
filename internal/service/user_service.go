package service

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/domain"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{userRepo}
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*domain.User, error) {

	if id == 0 {
		return nil, errors.New("id de usuario invalido")
	}

	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "" {
		return nil, errors.New("email requerido")
	}

	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) Update(ctx context.Context, user *domain.User) error {

	if user.ID == 0 {
		return errors.New("id de usuario invalido")
	}

	if user.Email == "" {
		return errors.New("email requerido")
	}

	if user.Name == "" {
		return errors.New("nombre requerido")
	}

	return s.userRepo.Update(ctx, user)
}

func (s *UserService) UpdatePassword(ctx context.Context, user *domain.User) error {
	if user.ID == 0 {
		return errors.New("id de usuario invalido")
	}

	if user.Password == "" {
		return errors.New("contraseña requerida")
	}

	return s.userRepo.UpdatePassword(ctx, user)
}
