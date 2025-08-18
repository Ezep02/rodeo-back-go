package repository

import (
	"context"

	"github.com/ezep02/rodeo/internal/domain"
	"gorm.io/gorm"
)

type GormCouponRepository struct {
	db *gorm.DB
}

func NewGormCouponRepo(db *gorm.DB) domain.CouponRepository {
	return &GormCouponRepository{db}
}

func (r *GormCouponRepository) Create(ctx context.Context, coupon *domain.Coupon) error {

	return r.db.WithContext(ctx).Create(coupon).Error
}

func (r *GormCouponRepository) GetByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	var coupon domain.Coupon

	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&coupon).Error; err != nil {
		return nil, err
	}

	return &coupon, nil
}

func (r *GormCouponRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Coupon, error) {
	var coupons []domain.Coupon

	if err := r.db.WithContext(ctx).Where("user_id = ? AND expire_at > NOW()", userID).Find(&coupons).Error; err != nil {
		return nil, err
	}

	return coupons, nil
}
