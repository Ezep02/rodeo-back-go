package auth

import (
	"context"

	"gorm.io/gorm"
)

type AuthRepository struct {
	Connection *gorm.DB
}

func NewAuthRepository(Connection *gorm.DB) *AuthRepository {

	return &AuthRepository{
		Connection: Connection,
	}
}

func (r *AuthRepository) RegisterUser(ctx context.Context, user *User) (*User, error) {

	result := r.Connection.WithContext(ctx).Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &User{
		Model:        user.Model,
		Name:         user.Name,
		Surname:      user.Surname,
		Email:        user.Email,
		Phone_number: user.Phone_number,
		Is_admin:     user.Is_admin,
	}, nil
}

func (r *AuthRepository) LoginUser(ctx context.Context, user *LogUserReq) (*User, error) {

	var loggedUser User

	result := r.Connection.WithContext(ctx).Where("email = ?", user.Email).First(&loggedUser)

	if result.Error != nil {
		return nil, result.Error
	}

	return &User{
		Model:        loggedUser.Model,
		Name:         loggedUser.Name,
		Surname:      loggedUser.Surname,
		Email:        loggedUser.Email,
		Phone_number: loggedUser.Phone_number,
		Is_admin:     loggedUser.Is_admin,
		Is_barber:    loggedUser.Is_barber,
	}, nil
}
