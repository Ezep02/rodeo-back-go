package analytics

import "context"

type AnalyticRepository interface {
	PopularTimeSlot(ctx context.Context) ([]PopularTimeSlot, error)
	BookingOcupationRate(ctx context.Context) (*BookingOcupationRate, error)
	MonthBookingCount(ctx context.Context) ([]MonthBookingCount, error)
	WeeklyBookingRate(ctx context.Context) ([]WeeklyBookingRate, error)
	NewClientRate(ctx context.Context) ([]NewClientRate, error)
	MonthlyRevenue(ctx context.Context) ([]MonthlyRevenue, error)
}
