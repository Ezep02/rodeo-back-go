package repository

import (
	"context"
	"fmt"

	"github.com/ezep02/rodeo/internal/auth/models"
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

func (r *AuthRepository) RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {

	result := r.Connection.WithContext(ctx).Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &models.User{
		Model:        user.Model,
		Name:         user.Name,
		Surname:      user.Surname,
		Email:        user.Email,
		Phone_number: user.Phone_number,
		Is_admin:     user.Is_admin,
	}, nil
}

func (r *AuthRepository) SearchUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User

	result := r.Connection.WithContext(ctx).Where("email = ?", email).Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// Funcion encargada de restablecer la contraseña del usuario
func (r *AuthRepository) UpdateUserPassword(ctx context.Context, userID int, newPassword string) error {

	// Actualizar la contraseña en la base de datos
	if err := r.Connection.WithContext(ctx).Exec(`
		UPDATE users 
		SET password = ? 
		WHERE id = ?
	`, newPassword, userID).Error; err != nil {
		return fmt.Errorf("error al actualizar la contraseña: %v", err)
	}

	return nil
}
