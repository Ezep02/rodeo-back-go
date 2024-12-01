package orders

import (
	"context"

	"gorm.io/gorm"
)

type OrderRepository struct {
	Connection *gorm.DB
}

func NewOrderRepository(cnn *gorm.DB) *OrderRepository {

	return &OrderRepository{
		Connection: cnn,
	}
}

func (or *OrderRepository) CreateNewOrder(ctx context.Context, order *Order) (*Order, error) {

	result := or.Connection.WithContext(ctx).Create(order)
	if result.Error != nil {
		return nil, result.Error
	}

	return order, nil
}

func (or *OrderRepository) GetAllOrders(ctx context.Context, limit int, offset int) (*[]Order, error) {

	var orders *[]Order

	result := or.Connection.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at desc").Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func (or *OrderRepository) GetOrderByUserID(ctx context.Context, userID int) (*Order, error) {

	var order_response Order

	result := or.Connection.WithContext(ctx).Last(&order_response)
	if result.Error != nil {
		return nil, result.Error
	}

	return &order_response, nil
}
