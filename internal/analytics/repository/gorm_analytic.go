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

func (r *GormAnalyticRepository) NewClientRate(ctx context.Context) (*analytics.NewClientRate, error) {
	var newClients = &analytics.NewClientRate{}

	//  Contar total de usuarios
	if err := r.db.WithContext(ctx).Table("users").Count(&newClients.TotalCount).Error; err != nil {
		return nil, err
	}

	// Agrupar usuarios por mes (MySQL)
	type monthlyData struct {
		Month      string
		NewClients int
	}

	var monthStats []monthlyData

	if err := r.db.WithContext(ctx).
		Table("users").
		Select("DATE_FORMAT(created_at, '%Y-%m') AS month, COUNT(*) AS new_clients").
		Group("month").
		Order("month ASC").
		Scan(&monthStats).Error; err != nil {
		return nil, err
	}

	// Adaptar al struct final
	for _, m := range monthStats {
		newClients.Data = append(newClients.Data, struct {
			Month      string `json:"month"`
			NewClients int    `json:"new_clients"`
		}{
			Month:      m.Month,
			NewClients: m.NewClients,
		})
	}

	return newClients, nil
}

func (r *GormAnalyticRepository) MonthlyRevenue(ctx context.Context) (*analytics.MonthlyRevenue, error) {
	var monthlyRevenue = &analytics.MonthlyRevenue{}

	// Calcular total global de ingresos aprobados
	if err := r.db.WithContext(ctx).
		Table("payments").
		Select("COALESCE(SUM(amount), 0)").
		Where("status = ?", "aprobado").
		Scan(&monthlyRevenue.TotalRevenue).Error; err != nil {
		return nil, err
	}

	// Agrupar ingresos por mes
	type monthlyData struct {
		Month        string
		TotalRevenue float64
	}

	var monthStats []monthlyData

	if err := r.db.WithContext(ctx).
		Table("payments").
		Select("DATE_FORMAT(paid_at, '%Y-%m') AS month, COALESCE(SUM(amount), 0) AS total_revenue").
		Where("status = ?", "aprobado").
		Group("month").
		Order("month ASC").
		Scan(&monthStats).Error; err != nil {
		return nil, err
	}

	// Adaptar al struct final
	for _, m := range monthStats {
		monthlyRevenue.Data = append(monthlyRevenue.Data, struct {
			Month        string  `json:"month"`
			TotalRevenue float64 `json:"total_revenue"`
		}{
			Month:        m.Month,
			TotalRevenue: m.TotalRevenue,
		})
	}

	return monthlyRevenue, nil
}
