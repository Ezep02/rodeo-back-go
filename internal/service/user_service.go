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

	if user.Surname == "" {
		return errors.New("apellido requerido")
	}

	// verificar existencia del usuario
	if _, err := s.userRepo.GetByID(ctx, user.ID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return errors.New("usuario no encontrado")
		}
	}

	// si quiere cambiar el email, verificar que no exista otro usuario con ese email
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	if existingUser != nil && existingUser.ID != user.ID {
		return errors.New("ya existe un usuario con ese email")
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

func (s *UserService) UpdateUsername(ctx context.Context, new_username string, id uint) error {

	if id == 0 {
		return errors.New("id de usuario invalido")
	}

	if new_username == "" {
		return errors.New("nombre requerido")
	}

	// verificar existencia del usuario
	if _, err := s.userRepo.GetByID(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return errors.New("usuario no encontrado")
		}
	}

	return s.userRepo.UpdateUsername(ctx, new_username, id)
}

func (s *UserService) UpdateAvatar(ctx context.Context, avatar string, id uint) error {

	if id == 0 {
		return errors.New("id de usuario invalido")
	}

	if avatar == "" {
		return errors.New("avatar requerido")
	}
	// verificar existencia del usuario
	if _, err := s.userRepo.GetByID(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return errors.New("usuario no encontrado")
		}
	}

	return s.userRepo.UpdateAvatar(ctx, avatar, id)
}
