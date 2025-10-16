package repository

import (
	"context"

	"github.com/ezep02/rodeo/internal/auth/domain"
	"gorm.io/gorm"
)

type GormAuthRepository struct {
	db *gorm.DB
}

func NewGormAuthRepo(db *gorm.DB) domain.AuthRepository {
	return &GormAuthRepository{db}
}

func (r *GormAuthRepository) Register(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *GormAuthRepository) Login(ctx context.Context, email string) (*domain.User, error) {

	var user domain.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {

		// Check if the error is a record not found error
		// If so, return a custom error indicating that the appointment was not found
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}
	return &user, nil
}

func (r *GormAuthRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User

	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
	}

	return &user, nil
}

func (r *GormAuthRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {

	var user domain.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {

		// Check if the error is a record not found error
		// If so, return a custom error indicating that the appointment was not found
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}
	return &user, nil
}
