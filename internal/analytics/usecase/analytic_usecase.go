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

func (s *AnalyticService) NewClientRate(ctx context.Context) (*analytics.NewClientRate, error) {
	return s.analyticRepo.NewClientRate(ctx)
}
func (s *AnalyticService) MonthlyRevenue(ctx context.Context) (*analytics.MonthlyRevenue, error) {
	return s.analyticRepo.MonthlyRevenue(ctx)
}
