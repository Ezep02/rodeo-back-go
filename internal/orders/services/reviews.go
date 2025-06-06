package services

import (
	"context"

	"github.com/ezep02/rodeo/internal/orders/models"
)

func (s *OrderService) CreateNewReview(ctx context.Context, review_req models.Review) error {
	return s.OrderRepo.CreatingReview(ctx, review_req)
}

func (s *OrderService) GetReviews(ctx context.Context, offset int) (*[]models.ReviewResponse, error) {
	return s.OrderRepo.GettingReviews(ctx, offset)
}
