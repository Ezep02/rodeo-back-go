package coupon

import "context"

type CouponRepository interface {
	Create(ctx context.Context, coupon *Coupon) error
	GetByCode(ctx context.Context, code string) (*Coupon, error)
	GetByUserID(ctx context.Context, id uint) ([]Coupon, error)
	UpdateStatus(ctx context.Context, code string) error
}
