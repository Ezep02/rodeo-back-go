package repository

import (
	"context"
	"log"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormUserRepo(db *gorm.DB, redis *redis.Client) domain.UserRepository {
	return &GormUserRepository{db, redis}
}

func (r *GormUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User

	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		// Check if the error is a record not found error
		// If so, return a custom error indicating that the User was not found
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		// Check if the error is a record not found error
		// If so, return a custom error indicating that the User was not found
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) Update(ctx context.Context, user *domain.User) error {
	updates := map[string]interface{}{
		"name":         user.Name,
		"surname":      user.Surname,
		"email":        user.Email,
		"phone_number": user.Phone_number,
	}

	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
		log.Println("Error updating user:", err)
		return err
	}

	return nil
}

func (r *GormUserRepository) UpdatePassword(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).Update("password", user.Password).Error; err != nil {
		log.Println("Error updating user password:", err)
		return err
	}

	return nil
}
