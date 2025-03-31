package services

import (
	"context"

	"github.com/ezep02/rodeo/internal/analytics/models"
	"github.com/ezep02/rodeo/internal/analytics/repository"
)

type Analytics_service struct {
	Analytics_repository *repository.Analytics_repository
}

func NewAnalyticServices(analytics_repo *repository.Analytics_repository) *Analytics_service {
	return &Analytics_service{
		Analytics_repository: analytics_repo,
	}
}

// Servicio para obtener el total de ingresos en el mes, y el promedio en comparacion al mes anterior
func (s *Analytics_service) ObtainMonthlyRevenueAndAvgComparedToLastMonth(ctx context.Context) (float64, float64, error) {
	return s.Analytics_repository.GetMonthlyRevenueAndAvgComparedToLastMonth(ctx)
}

// Servicio para obtener el total de citas reservadas en el mes, y el promedio en comparacion al mes anterior
func (s *Analytics_service) ObtainMonthlyAppointmentsAndAvgComparedToLastMonth(ctx context.Context) (int, float64, error) {
	return s.Analytics_repository.GetMonthlyAppointmentsAndAvgComparedToLastMonth(ctx)
}

// Servicio para obtener el numero de nuevos clientes en el mes, y el promedio en comparacion al mes anterior
func (s *Analytics_service) ObtainMonthlyNewCustomersAndAvgComparedToLastMonth(ctx context.Context) (int, float64, error) {
	return s.Analytics_repository.GetMonthlyNewCustomersAndAvgComparedToLastMonth(ctx)
}

// Servicio para obtener el numero de cancelaciones en el mes, y el promedio en comparacion al mes anterior
func (s *Analytics_service) ObtainMonthlyCancellationsAndAvgComparedToLastMonth(ctx context.Context) (int, float64, error) {
	return s.Analytics_repository.GetMonthlyCancellationsAndAvgComparedToLastMonth(ctx)
}

// Servicio para obtener un listado de los meses del año con su respectivo monto total de ingresos
func (s *Analytics_service) ObtainCurrentYearMonthlyRevenue(ctx context.Context) ([]models.CurrentYearMonthlyRevenue, error) {
	return s.Analytics_repository.GetCurrentYearMonthlyRevenue(ctx)
}

// Servicio para obtener un listado de los meses del año con su respectivo monto total de ingresos
func (s *Analytics_service) ObtainMonthlyPopularServices(ctx context.Context) ([]models.MonthlyPopularService, error) {
	return s.Analytics_repository.GetMonthlyPopularServices(ctx)
}

// Servicio para obtener un listado de clientes con mas frecuencia y el monto abonado
func (s *Analytics_service) ObtainFrequentCustomers(ctx context.Context) ([]models.FrequentCustomer, error) {
	return s.Analytics_repository.GetFrequentCustomers(ctx)
}

// Servicio para obtener un listado mes a mes del total de cortes realizados por el barbero
func (s *Analytics_service) ObtainYearlyBarberHaircuts(ctx context.Context, barberID int) ([]models.MonthlyHaircuts, error) {
	return s.Analytics_repository.GetYearlyBarberHaircuts(ctx, barberID)
}
