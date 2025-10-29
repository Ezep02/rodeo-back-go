package analytics

import "context"

type AnalyticRepository interface {
	NewClientRate(ctx context.Context) (*NewClientRate, error)
	MonthlyRevenue(ctx context.Context) (*MonthlyRevenue, error)
}
