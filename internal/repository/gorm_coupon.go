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
