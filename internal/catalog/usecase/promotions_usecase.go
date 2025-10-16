package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ezep02/rodeo/internal/catalog/domain/promotions"
)

type PromoService struct {
	promoRepo promotions.PromoRepository
}

func NewPromoService(promoRepo promotions.PromoRepository) *PromoService {
	return &PromoService{promoRepo}
}

func (s *PromoService) Create(ctx context.Context, data *promotions.Promotion) error {

	if data.ServiceID == 0 {
		return errors.New("error recupeando servicio")
	}

	if data.StartDate.IsZero() {
		return errors.New("el parametro fecha de inicio es requerido")
	}

	if data.EndDate.IsZero() {
		return errors.New("el parametro fecha de fin es requerido")
	}

	if time.Now().After(*data.EndDate) {
		// handle case where end date is in the past
		return errors.New("la fecha de fin debe ser posterior a la fecha actual")
	}

	return s.promoRepo.Create(ctx, data)
}

func (s *PromoService) ListByServiceId(ctx context.Context, id uint, offset int) ([]promotions.Promotion, error) {
	return s.promoRepo.ListByServiceId(ctx, id, offset)
}

func (s *PromoService) Update(ctx context.Context, id uint, data *promotions.Promotion) error {
	if id == 0 {
		return errors.New("error recupeando servicio")
	}

	if data.StartDate.IsZero() {
		return errors.New("el parametro fecha de inicio es requerido")
	}

	if data.EndDate.IsZero() {
		return errors.New("el parametro fecha de fin es requerido")
	}

	if time.Now().After(*data.EndDate) {
		// handle case where end date is in the past
		return errors.New("la fecha de fin debe ser posterior a la fecha actual")
	}

	return s.promoRepo.Update(ctx, id, data)
}

func (s *PromoService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("el id es un parametro requerido")
	}
	return s.promoRepo.Delete(ctx, id)
}
