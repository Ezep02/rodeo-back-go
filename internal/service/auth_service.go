package service

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/domain"
)

type AuthService struct {
	authRepo domain.AuthRepository
}

func NewAuthService(authRepo domain.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) Register(ctx context.Context, user *domain.User) error {

	// 1. Validar que tenga email
	if user.Email == "" {
		return errors.New("el gmail es un campo requerido")
	}

	// 2. Validar que tenga nombre
	if user.Name == "" {
		return errors.New("el nombre es un campo requerido")
	}

	// 3. Validar que tenga apellido
	if user.Surname == "" {
		return errors.New("el apellido es un campo requerido")
	}

	// 4. Verificar contraseña
	if user.Password == "" {
		return errors.New("la contraseña es un campo requerido")
	}

	// 5. Verificar que el usuario no este registrado
	existing, _ := s.authRepo.Login(ctx, user.Email)

	if existing != nil {
		return errors.New("al parecer ya existe un usuario registrado con ese gmail")
	}

	// 6. Crear registro
	return s.authRepo.Register(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email string) (*domain.User, error) {

	// 1. Validar que el email exista
	if email == "" {
		return nil, errors.New("el email es un campo requerido")
	}

	return s.authRepo.Login(ctx, email)
}

func (s *AuthService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.authRepo.GetByID(ctx, id)
}
