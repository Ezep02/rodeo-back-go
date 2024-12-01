package orders

import "context"

type OrderService struct {
	OrderRepo *OrderRepository
}

func NewOrderService(ord_srv *OrderRepository) *OrderService {
	return &OrderService{
		OrderRepo: ord_srv,
	}
}

func (or_srv *OrderService) CreateNewOrder(ctx context.Context, order *Order) (*Order, error) {
	return or_srv.OrderRepo.CreateNewOrder(ctx, order)
}

func (or_srv *OrderService) GetOrderService(ctx context.Context, limit int, offset int) (*[]Order, error) {
	return or_srv.OrderRepo.GetAllOrders(ctx, limit, offset)
}

func (or_srv *OrderService) GetOrderByUserID(ctx context.Context, userID int) (*Order, error) {
	return or_srv.OrderRepo.GetOrderByUserID(ctx, userID)
}
