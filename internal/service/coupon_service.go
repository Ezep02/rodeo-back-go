package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
)

type CouponService struct {
	couponRepo domain.CouponRepository
}

func NewCouponService(repo domain.CouponRepository) *CouponService {
	return &CouponService{repo}
}

func (s *CouponService) Create(ctx context.Context, coupon *domain.Coupon) error {

	// 1. validar que tenga codigo
	if coupon.Code == "" {
		return errors.New("sin codigo")
	}

	// 2. validar que la fecha de expiracion no sea en el pasado
	if coupon.ExpireAt.Before(time.Now()) {
		return errors.New("su cupon ya caduco")
	}

	// 3. validar que el cupon este disponible
	if !coupon.IsAvailable {
		return errors.New("cupon no disponible")
	}

	// 4. verificar que el codigo del id no este en uso
	// TODO

	return s.couponRepo.Create(ctx, coupon)
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// GenerateCoupon generates a random coupon code of the given length.
func (s *CouponService) GenerateCoupon(length int) (string, error) {
	coupon := make([]byte, length)
	for i := range coupon {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		coupon[i] = charset[randomIndex.Int64()]
	}
	return string(coupon), nil
}

func (s *CouponService) GetByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	return s.couponRepo.GetByCode(ctx, code)
}

func (s *CouponService) GetByUserID(ctx context.Context, userID uint) ([]domain.Coupon, error) {
	return s.couponRepo.GetByUserID(ctx, userID)
}

func (s *CouponService) UpdateStatus(ctx context.Context, code string) error {
	return s.couponRepo.UpdateStatus(ctx, code)
}
