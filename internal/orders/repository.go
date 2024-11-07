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

func (or *OrderRepository) CreateNewOrder(ctx context.Context, payment *Payment) (*Payment, error) {

	result := or.Connection.WithContext(ctx).Create(payment)

	if result.Error != nil {
		return nil, result.Error
	}

	return payment, nil
}
