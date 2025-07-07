package service

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/domain"
)

type ReviewService struct {
	revRepo  domain.ReviewRepository
	apptRepo domain.AppointmentRepository
}

func NewReviewService(revRepo domain.ReviewRepository, apptRepo domain.AppointmentRepository) *ReviewService {
	return &ReviewService{revRepo, apptRepo}
}

func (s *ReviewService) Create(ctx context.Context, review *domain.Review) error {

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

func (s *ReviewService) List(ctx context.Context) ([]domain.Appointment, error) {
	return s.revRepo.List(ctx)
}
