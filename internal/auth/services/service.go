package services

import (
	"context"

	"github.com/ezep02/rodeo/internal/auth/models"
	"github.com/ezep02/rodeo/internal/auth/repository"
)

type AuthService struct {
	AuthRepo *repository.AuthRepository
}

func NewAuthService(auth_r *repository.AuthRepository) *AuthService {

	return &AuthService{
		AuthRepo: auth_r,
	}
}

func (s *AuthService) RegisterUserServ(ctx context.Context, user *models.User) (*models.User, error) {
	return s.AuthRepo.RegisterUser(ctx, user)
}

func (s *AuthService) SearchUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.AuthRepo.SearchUserByEmail(ctx, email)
}

func (s *AuthService) UpdateUserPasswordServ(ctx context.Context, userID int, newPassword string) error {
	return s.AuthRepo.UpdateUserPassword(ctx, userID, newPassword)
}
