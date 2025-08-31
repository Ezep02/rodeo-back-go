package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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
	var (
		user     domain.User
		cacheKey = fmt.Sprintf("user:%d", id)
	)
	// 1. Intentar obtener el usuario desde Redis
	cachedUser, err := r.redis.Get(ctx, cacheKey).Result()

	if err == nil {
		json.Unmarshal([]byte(cachedUser), &user)
		return &user, nil
	}

	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	// 3. Almacenar el usuario en Redis para futuras solicitudes
	data, _ := json.Marshal(user)
	r.redis.Set(ctx, cacheKey, data, time.Hour)

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

	var (
		updates = map[string]any{
			"name":         user.Name,
			"surname":      user.Surname,
			"email":        user.Email,
			"phone_number": user.Phone_number,
		}

		cacheKey = fmt.Sprintf("user:%d", user.ID)
	)

	// Eliminar el usuario en cache
	if err := r.redis.Del(ctx, cacheKey); err != nil {
		log.Println("Error deleting user from cache:", err)
	}

	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
		log.Println("Error updating user:", err)
		return errors.New("error actualizando el usuario")
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

func (r *GormUserRepository) UpdateUsername(ctx context.Context, new_username string, id uint) error {
	var (
		cacheKey = fmt.Sprintf("user:%d", id)
		updates  = map[string]any{
			"username":         new_username,
			"last_name_change": time.Now(),
		}
	)

	// Eliminar el usuario en cache
	if err := r.redis.Del(ctx, cacheKey).Err(); err != nil {
		log.Println("Error deleting user from cache:", err)
	}

	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Println("Error updating username:", err)
		return err
	}

	return nil
}

func (r *GormUserRepository) UpdateAvatar(ctx context.Context, avatar string, id uint) error {
	var (
		cacheKey = fmt.Sprintf("user:%d", id)
	)

	// Eliminar el usuario en cache
	_, err := r.redis.Del(ctx, cacheKey).Result()
	if err != nil {
		log.Println("Error deleting user from cache:", err)
	}

	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("avatar", avatar).Error; err != nil {
		log.Println("Error updating avatar:", err)
		return err
	}

	return nil
}
