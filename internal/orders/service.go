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

func (or_srv *OrderService) CreateNewOrder(ctx context.Context, payment *Payment) (*Payment, error) {
	return or_srv.OrderRepo.CreateNewOrder(ctx, payment)
}
