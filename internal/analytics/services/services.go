package services

import (
	"context"

	"github.com/ezep02/rodeo/internal/analytics/models"
	"github.com/ezep02/rodeo/internal/analytics/repository"
)

type AnalyticsServices struct {
	An_repo *repository.AnalyticsRepository
}

func NewAnalyticsServices(an_repo *repository.AnalyticsRepository) *AnalyticsServices {
	return &AnalyticsServices{
		An_repo: an_repo,
	}
}

func (an_srv *AnalyticsServices) GetTotalRevenue(ctx context.Context) (*[]models.Revenue, error) {
	return an_srv.An_repo.GetTotalRevenue(ctx)
}

func (an_srv *AnalyticsServices) GetTotalRegisteredUsers(ctx context.Context) (*int64, error) {
	return an_srv.An_repo.GetTotalRegisteredUsers(ctx)
}

func (an_srv *AnalyticsServices) GetRecivedUsers(ctx context.Context) (*int64, error) {
	return an_srv.An_repo.GetRecivedUsers(ctx)
}

// expense
func (an_srv *AnalyticsServices) NewExpenseSrv(ctx context.Context, exp *models.Expenses) (*models.Expenses, error) {
	return an_srv.An_repo.NewExpense(ctx, exp)
}

func (an_srv *AnalyticsServices) GetExpensesHistorialSrv(ctx context.Context, limit int, offset int) (*[]models.Expenses, error) {
	return an_srv.An_repo.GetExpensesHistorial(ctx, limit, offset)
}

func (an_srv *AnalyticsServices) GetTotalExpenses(ctx context.Context) (*[]models.Expense, error) {
	return an_srv.An_repo.GetTotalExpenses(ctx)
}
