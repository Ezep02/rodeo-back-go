package repository

import (
	"context"
	"time"

	"github.com/ezep02/rodeo/internal/analytics/models"
	"gorm.io/gorm"
)

type AnalyticsRepository struct {
	Connection *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{
		Connection: db,
	}
}

func (an_r *AnalyticsRepository) GetTotalRegisteredUsers(ctx context.Context) (*int64, error) {

	var totalUsers int64
	result := an_r.Connection.WithContext(ctx).Model(&models.User{}).Count(&totalUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	return &totalUsers, nil
}

func (an_r *AnalyticsRepository) GetRecivedUsers(ctx context.Context) (*int64, error) {

	var totalUsers int64

	result := an_r.Connection.WithContext(ctx).Model(&models.Schedule{}).Where("schedule_day_date <= curdate() AND available = ?", false).Count(&totalUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	return &totalUsers, nil
}

func (an_r *AnalyticsRepository) GetTotalRevenue(ctx context.Context) (*[]models.Revenue, error) {

	var revenue *[]models.Revenue
	query := `
		SELECT
			STR_TO_DATE(DATE_FORMAT(schedule_day_date, '%Y-%m-01'), '%Y-%m-%d') AS month_date,
			SUM(price) AS total_revenue
		FROM
			orders
		WHERE
			deleted_at IS NULL
		GROUP BY
			STR_TO_DATE(DATE_FORMAT(schedule_day_date, '%Y-%m-01'), '%Y-%m-%d')
		ORDER BY
			month_date ASC
	`
	result := an_r.Connection.WithContext(ctx).Raw(query).Scan(&revenue)

	if result.Error != nil {
		return nil, result.Error
	}

	return revenue, nil
}

func (an_r *AnalyticsRepository) NewExpense(ctx context.Context, exp *models.Expenses) (*models.Expenses, error) {

	result := an_r.Connection.WithContext(ctx).Create(exp)

	if result.Error != nil {
		return nil, result.Error
	}

	return exp, nil
}

func (an_r *AnalyticsRepository) GetExpensesHistorial(ctx context.Context, limit int, offset int) (*[]models.Expenses, error) {
	var expensesList *[]models.Expenses

	result := an_r.Connection.WithContext(ctx).Offset(offset).Limit(limit).Find(&expensesList)

	if result.Error != nil {
		return nil, result.Error
	}

	return expensesList, nil
}

func (an_r *AnalyticsRepository) GetTotalExpenses(ctx context.Context) (*[]models.Expense, error) {
	var totalExpenses *[]models.Expense

	currentYear := time.Now().Year()

	// Realiza la sumatoria del campo amount para el año actual
	query := `
	SELECT
		STR_TO_DATE(DATE_FORMAT(created_at, '%Y-%m-01'), '%Y-%m-%d') AS month_date,
		SUM(amount) AS total_expense
	FROM
		expenses
	WHERE
		deleted_at IS NULL
		AND YEAR(created_at) = ?  -- Filtro para el año actual
	GROUP BY
		STR_TO_DATE(DATE_FORMAT(created_at, '%Y-%m-01'), '%Y-%m-%d')
	ORDER BY
		month_date ASC
	`
	result := an_r.Connection.WithContext(ctx).Raw(query, currentYear).Scan(&totalExpenses)

	if result.Error != nil {
		return nil, result.Error
	}

	return totalExpenses, nil
}
