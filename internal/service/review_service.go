package service

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/domain/review"
)

type ReviewService struct {
	revRepo  review.ReviewRepository
	apptRepo domain.AppointmentRepository
}

func NewReviewService(revRepo review.ReviewRepository, apptRepo domain.AppointmentRepository) *ReviewService {
	return &ReviewService{revRepo, apptRepo}
}

func (s *ReviewService) Create(ctx context.Context, review *review.Review) error {

	// 1. Verificar la existencia del appointment
	if _, err := s.apptRepo.GetByID(ctx, review.AppointmentID); err != nil {
		return errors.New("no es posible crear una reseña de una cita que no existe")
	}

	// 2. Verficar que tenga un comentario
	if review.Comment == "" {
		return errors.New("debe contener un comentario")
	}

	return s.revRepo.Create(ctx, review)
}

func (s *ReviewService) List(ctx context.Context) ([]review.ReviewDetail, error) {
	return s.revRepo.List(ctx)
}

func (s *ReviewService) ListByProductID(ctx context.Context, productID uint) ([]review.Review, error) {
	return s.revRepo.ListByProductID(ctx, productID)
}

func (s *ReviewService) ListByUserID(ctx context.Context, userID uint, offset int) ([]review.ReviewDetail, error) {
	if offset < 0 {
		offset = 0
	}

	return s.revRepo.ListByUserID(ctx, userID, offset)
}
