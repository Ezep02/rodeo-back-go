package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ezep02/rodeo/internal/booking/domain/coupon"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormCouponRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormCouponRepo(db *gorm.DB, redis *redis.Client) coupon.CouponRepository {
	return &GormCouponRepository{db: db, redis: redis}
}

func (r *GormCouponRepository) Create(ctx context.Context, c *coupon.Coupon) error {
	if c == nil {
		return errors.New("coupon es nil")
	}
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *GormCouponRepository) GetByCode(ctx context.Context, code string) (*coupon.Coupon, error) {
	var c coupon.Coupon
	if err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *GormCouponRepository) GetByUserID(ctx context.Context, userID uint) ([]coupon.Coupon, error) {
	var coupons []coupon.Coupon
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&coupons).Error; err != nil {
		return nil, err
	}
	return coupons, nil
}

func (r *GormCouponRepository) UpdateStatus(ctx context.Context, code string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&coupon.Coupon{}).
		Where("code = ?", code).
		Updates(map[string]any{
			"is_available": false,
			"used_at":      now,
		}).Error
}
