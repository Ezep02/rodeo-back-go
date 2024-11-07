package auth

import "context"

type AuthService struct {
	AuthRepo *AuthRepository
}

func NewAuthService(auth_r *AuthRepository) *AuthService {

	return &AuthService{
		AuthRepo: auth_r,
	}
}

func (s *AuthService) RegisterUserServ(ctx context.Context, user *User) (*User, error) {
	return s.AuthRepo.RegisterUser(ctx, user)
}

func (s *AuthService) LoginUserServ(ctx context.Context, user *LogUserReq) (*User, error) {
	return s.AuthRepo.LoginUser(ctx, user)
}
