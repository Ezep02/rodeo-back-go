package repository

import (
	"context"
	"time"

	"github.com/ezep02/rodeo/internal/orders/models"
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

func (or *OrderRepository) CreateNewOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	result := or.Connection.WithContext(ctx).Create(order)
	if result.Error != nil {
		return nil, result.Error
	}

	return order, nil
}

func (r *OrderRepository) GetBarberPendingOrders(ctx context.Context, barberID int, limit int, offset int) ([]models.BarberPendingOrder, error) {

	var barberPendingOrders []models.BarberPendingOrder

	// extraer las ordenes penditenes, cuyo dia se despues del dia actual
	err := r.Connection.WithContext(ctx).Raw(`
		SELECT 
			title, 
			payer_name,
			payer_surname,
			barber_id,
			schedule_day_date,
			schedule_start_time,
			mp_status as status,
		FROM orders 
		WHERE 
			barber_id = ? 
			AND schedule_day_date > CURRENT_DATE 
			AND EXTRACT(DAY FROM schedule_day_date) = EXTRACT(MONTH FROM CURRENT_DATE)
		LIMIT ? 
		OFFSET ?
	`, barberID, limit, offset).Scan(&barberPendingOrders).Error

	if err != nil {
		return nil, err
	}

	return barberPendingOrders, nil
}

func (or *OrderRepository) GetOrderByUserID(ctx context.Context, userID int) (*models.Order, error) {

	var order_response models.Order

	result := or.Connection.WithContext(ctx).Last(&order_response)
	if result.Error != nil {
		return nil, result.Error
	}

	return &order_response, nil
}

// obtiene una lista de los turnos realizados
func (or *OrderRepository) GetOrdersHistorial(ctx context.Context, userID int, limit int, offset int) (*[]models.Order, error) {
	var ordersList *[]models.Order

	// Obtener la hora actual
	currentTime := time.Now()

	// Construir la consulta
	result := or.Connection.WithContext(ctx).
		Where("user_id = ? AND date < ?", userID, currentTime).
		Order("date desc").
		Limit(limit).
		Offset(offset).
		Find(&ordersList)

	if result.Error != nil {
		return nil, result.Error
	}

	return ordersList, nil
}
