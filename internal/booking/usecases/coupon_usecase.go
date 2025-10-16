package usecases

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/booking/domain/coupon"
)

type CouponService struct {
	couponRepo coupon.CouponRepository
}

// Constructor
func NewCouponService(couponRepo coupon.CouponRepository) *CouponService {
	return &CouponService{couponRepo: couponRepo}
}

func (s *CouponService) CreateCoupon(ctx context.Context, c *coupon.Coupon) error {
	if c == nil {
		return errors.New("coupon es nil")
	}
	return s.couponRepo.Create(ctx, c)
}

func (s *CouponService) GetCouponByCode(ctx context.Context, code string) (*coupon.Coupon, error) {
	if code == "" {
		return nil, errors.New("code no puede ser vacío")
	}
	return s.couponRepo.GetByCode(ctx, code)
}

func (s *CouponService) GetCouponsByUserID(ctx context.Context, userID uint) ([]coupon.Coupon, error) {
	if userID == 0 {
		return nil, errors.New("userID no puede ser cero")
	}
	return s.couponRepo.GetByUserID(ctx, userID)
}

func (s *CouponService) MarkCouponAsUsed(ctx context.Context, code string) error {
	if code == "" {
		return errors.New("code no puede ser vacío")
	}
	return s.couponRepo.UpdateStatus(ctx, code)
}
