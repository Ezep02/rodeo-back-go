package services

import (
	"context"

	"github.com/ezep02/rodeo/internal/orders/models"
	"github.com/ezep02/rodeo/internal/orders/repository"
)

type OrderService struct {
	OrderRepo *repository.OrderRepository
}

func NewOrderService(ord_srv *repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepo: ord_srv,
	}
}

func (s *OrderService) CreateNewOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	return s.OrderRepo.CreatingNewOrder(ctx, order)
}

func (s *OrderService) GetOrderService(ctx context.Context, barberID int, limit int, offset int) ([]models.BarberPendingOrder, error) {
	return s.OrderRepo.GetBarberPendingOrders(ctx, barberID, limit, offset)
}

// obtener ordenes pendientes del cliente
func (s *OrderService) GetCustomerPendingOrder(ctx context.Context, userID int) ([]models.CustomerPendingOrder, error) {
	return s.OrderRepo.GettingCustomerPendingOrders(ctx, userID)
}

// almacenar el token transitoriamente para realizar el pago
func (s *OrderService) SetOrderToken(ctx context.Context, token string, order models.PendingOrderToken) error {
	return s.OrderRepo.SavingOrderToken(ctx, token, order)
}

func (s *OrderService) GetOrderByToken(ctx context.Context, token string) (models.PendingOrderToken, error) {
	return s.OrderRepo.SearchingOrderToken(ctx, token)
}
