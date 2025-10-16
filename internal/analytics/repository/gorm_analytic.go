package repository

import (
	"context"

	"github.com/ezep02/rodeo/internal/analytics/domain/analytics"

	"gorm.io/gorm"
)

type GormAnalyticRepository struct {
	db *gorm.DB
}

func NewGormAnalyticRepo(db *gorm.DB) analytics.AnalyticRepository {
	return &GormAnalyticRepository{db}
}

func (r *GormAnalyticRepository) PopularTimeSlot(ctx context.Context) ([]analytics.PopularTimeSlot, error) {
	var popTime []analytics.PopularTimeSlot

	// err := r.db.WithContext(ctx).
	// 	Model(&analytics.Slot{}).
	// 	Select("time, COUNT(*) AS bookings").
	// 	Where("is_booked = ?", true).
	// 	Group("time").
	// 	Order("bookings DESC").
	// 	Limit(10).
	// 	Scan(&popTime).Error

	// if err != nil {
	// 	return nil, err
	// }

	return popTime, nil
}

func (r *GormAnalyticRepository) BookingOcupationRate(ctx context.Context) (*analytics.BookingOcupationRate, error) {

	var bookingRates *analytics.BookingOcupationRate

	// err := r.db.WithContext(ctx).
	// 	Model(&domain.Slot{}).
	// 	Select(`
	// 		DATE_FORMAT(date, '%M %Y') AS month,
	// 		ROUND(SUM(CASE WHEN is_booked = TRUE THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2) AS occ_pct
	// 	`).
	// 	Group("month").
	// 	Order("STR_TO_DATE(month, '%M %Y')").
	// 	Scan(&bookingRates).Error

	// if err != nil {
	// 	return nil, err
	// }
	return bookingRates, nil
}

func (r *GormAnalyticRepository) MonthBookingCount(ctx context.Context) ([]analytics.MonthBookingCount, error) {
	var monthBooking []analytics.MonthBookingCount

	query := `
		SELECT 
			DATE_FORMAT(s.date, '%m-%Y') AS month, 
			COUNT(*) AS total_appointments
		FROM appointments a
		JOIN slots s ON s.id = a.slot_id
		GROUP BY month
		ORDER BY STR_TO_DATE(month, '%m-%Y')
	`

	err := r.db.WithContext(ctx).Raw(query).Scan(&monthBooking).Error
	if err != nil {
		return nil, err
	}
	return monthBooking, nil
}

func (r *GormAnalyticRepository) WeeklyBookingRate(ctx context.Context) ([]analytics.WeeklyBookingRate, error) {
	var weeklyRate []analytics.WeeklyBookingRate

	// err := r.db.WithContext(ctx).Model(&appointment.Appointment{}).
	// 	Select(`DATE_FORMAT(created_at, '%Y-%u-%m') AS week, COUNT(*) AS appointment_this_week`).
	// 	Group("week").
	// 	Order("week").
	// 	Scan(&weeklyRate).Error

	// if err != nil {
	// 	return nil, err
	// }

	return weeklyRate, nil
}

func (r *GormAnalyticRepository) NewClientRate(ctx context.Context) ([]analytics.NewClientRate, error) {

	var newClients []analytics.NewClientRate

	// err := r.db.WithContext(ctx).Model(&analytics.User{}).
	// 	Select(`DATE_FORMAT(created_at, '%Y-%m') AS month, COUNT(*) AS new_clients`).
	// 	Where("is_barber = ?", false).
	// 	Group("DATE_FORMAT(created_at, '%Y-%m')").
	// 	Order("month").
	// 	Scan(&newClients).Error

	// if err != nil {
	// 	return nil, err
	// }

	return newClients, nil
}

func (r *GormAnalyticRepository) MonthlyRevenue(ctx context.Context) ([]analytics.MonthlyRevenue, error) {

	var monthlyRevenue []analytics.MonthlyRevenue

	err := r.db.WithContext(ctx).
		Raw(`
            SELECT 
                DATE_FORMAT(a.created_at, '%Y-%m-01') AS month,
                SUM(p.price) AS total_revenue
            FROM appointments a
            JOIN appointment_products ap ON ap.appointment_id = a.id
            JOIN products p ON p.id = ap.product_id
            GROUP BY month
            ORDER BY month
        `).
		Scan(&monthlyRevenue).Error

	if err != nil {
		return nil, err
	}

	return monthlyRevenue, nil
}
