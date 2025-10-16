package repository

import (
	"context"

	"github.com/ezep02/rodeo/internal/users/domain/barber"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormBarberRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormBarberRepo(db *gorm.DB, redis *redis.Client) barber.BarberRepository {
	return &GormBarberRepository{db, redis}
}

func (r *GormBarberRepository) GetByID(ctx context.Context, id uint) (*barber.Barber, error) {
	var barber barber.Barber

	if err := r.db.WithContext(ctx).Where("user_id = ?", id).First(&barber).Error; err != nil {
		return nil, err
	}

	return &barber, nil
}

func (r *GormBarberRepository) List(ctx context.Context) ([]barber.BarberWithUser, error) {
	var barbers []barber.BarberWithUser

	if err := r.db.WithContext(ctx).
		Table("barbers b").
		Select(`b.id, b.user_id, b.calendar_id,
                u.id as user_id, u.name as user_name, u.surname as user_surname, u.avatar as user_avatar, u.username as user_username`).
		Joins("JOIN users u ON u.id = b.user_id").
		Find(&barbers).Error; err != nil {
		return nil, err
	}

	return barbers, nil
}
