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

func (or_srv *OrderService) CreateNewOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	return or_srv.OrderRepo.CreateNewOrder(ctx, order)
}

func (or_srv *OrderService) GetOrderService(ctx context.Context, barberID int, limit int, offset int) ([]models.BarberPendingOrder, error) {
	return or_srv.OrderRepo.GetBarberPendingOrders(ctx, barberID, limit, offset)
}

func (or_srv *OrderService) GetOrderByUserID(ctx context.Context, userID int) (*models.Order, error) {
	return or_srv.OrderRepo.GetOrderByUserID(ctx, userID)
}

func (or_srv *OrderService) GetOrdersHistorial(ctx context.Context, userID int, limit int, offset int) (*[]models.Order, error) {
	return or_srv.OrderRepo.GetOrdersHistorial(ctx, userID, limit, offset)
}
