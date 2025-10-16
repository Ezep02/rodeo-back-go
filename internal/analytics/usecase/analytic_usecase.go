package usecase

import (
	"context"

	"github.com/ezep02/rodeo/internal/analytics/domain/analytics"
)

type AnalyticService struct {
	analyticRepo analytics.AnalyticRepository
}

func NewAnalyticService(analyticRepo analytics.AnalyticRepository) *AnalyticService {
	return &AnalyticService{analyticRepo}
}

func (s *AnalyticService) PopularTimeSlot(ctx context.Context) ([]analytics.PopularTimeSlot, error) {
	return s.analyticRepo.PopularTimeSlot(ctx)
}
func (s *AnalyticService) BookingOcupationRate(ctx context.Context) (*analytics.BookingOcupationRate, error) {
	return s.analyticRepo.BookingOcupationRate(ctx)
}
func (s *AnalyticService) MonthBookingCount(ctx context.Context) ([]analytics.MonthBookingCount, error) {
	return s.analyticRepo.MonthBookingCount(ctx)
}
func (s *AnalyticService) WeeklyBookingRate(ctx context.Context) ([]analytics.WeeklyBookingRate, error) {
	return s.analyticRepo.WeeklyBookingRate(ctx)
}
func (s *AnalyticService) NewClientRate(ctx context.Context) ([]analytics.NewClientRate, error) {
	return s.analyticRepo.NewClientRate(ctx)
}
func (s *AnalyticService) MonthlyRevenue(ctx context.Context) ([]analytics.MonthlyRevenue, error) {
	return s.analyticRepo.MonthlyRevenue(ctx)
}
