package service

import (
	"context"

	"github.com/ezep02/rodeo/internal/domain"
)

type AnalyticService struct {
	analyticRepo domain.AnalyticRepository
}

func NewAnalyticService(analyticRepo domain.AnalyticRepository) *AnalyticService {
	return &AnalyticService{analyticRepo}
}

func (s *AnalyticService) PopularTimeSlot(ctx context.Context) ([]domain.PopularTimeSlot, error) {
	return s.analyticRepo.PopularTimeSlot(ctx)
}
func (s *AnalyticService) BookingOcupationRate(ctx context.Context) (*domain.BookingOcupationRate, error) {
	return s.analyticRepo.BookingOcupationRate(ctx)
}
func (s *AnalyticService) MonthBookingCount(ctx context.Context) ([]domain.MonthBookingCount, error) {
	return s.analyticRepo.MonthBookingCount(ctx)
}
func (s *AnalyticService) WeeklyBookingRate(ctx context.Context) ([]domain.WeeklyBookingRate, error) {
	return s.analyticRepo.WeeklyBookingRate(ctx)
}
func (s *AnalyticService) NewClientRate(ctx context.Context) ([]domain.NewClientRate, error) {
	return s.analyticRepo.NewClientRate(ctx)
}
func (s *AnalyticService) MonthlyRevenue(ctx context.Context) ([]domain.MonthlyRevenue, error) {
	return s.analyticRepo.MonthlyRevenue(ctx)
}
