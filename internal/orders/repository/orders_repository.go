package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/orders/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderRepository struct {
	Connection      *gorm.DB
	RedisConnection *redis.Client
}

func NewOrderRepository(cnn *gorm.DB, redis *redis.Client) *OrderRepository {

	return &OrderRepository{
		Connection:      cnn,
		RedisConnection: redis,
	}
}

func (or *OrderRepository) CreatingNewOrder(ctx context.Context, order *models.Order) (*models.Order, error) {

	// Create order and update schedule status as not available
	or.Connection.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// Set in db a new order
		tx.Transaction(func(order_tx *gorm.DB) error {

			if order_err := order_tx.Create(order).Error; order_err != nil {
				log.Println("[rolling back new order]")
				order_tx.Rollback()
				return order_err
			}
			log.Println("[setting new order]")
			return nil
		})

		// update schedule as not available
		tx.Transaction(func(updating_status_tx *gorm.DB) error {

			if upating_status_err := updating_status_tx.Exec(`
				UPDATE schedules 
				SET available = ?
				WHERE id = ?
			`, false, order.Shift_id).Error; upating_status_err != nil {
				log.Println("[status updated]")
				return nil
			}
			return nil
		})
		return nil
	})

	return order, nil
}

func (r *OrderRepository) GetBarberPendingOrders(ctx context.Context, barberID int, limit int, offset int) ([]models.BarberPendingOrder, error) {

	var barberPendingOrders []models.BarberPendingOrder

	// extraer las ordenes penditenes, cuyo dia se despues del dia actual
	err := r.Connection.WithContext(ctx).Raw(`
		SELECT 
			id,
			title, 
			payer_name,
			payer_surname,
			barber_id,
			schedule_day_date,
			schedule_start_time,
			mp_status,
			price,
			created_at,
			updated_at,
			deleted_at
		FROM orders 
		WHERE 
			barber_id = ? 
			AND schedule_day_date >= CURRENT_DATE 
			AND EXTRACT(MONTH FROM schedule_day_date) = EXTRACT(MONTH FROM CURRENT_DATE)
		LIMIT ? 
		OFFSET ?
	`, barberID, limit, offset).Scan(&barberPendingOrders).Error

	if err != nil {
		return nil, err
	}

	return barberPendingOrders, nil
}

// obtener los turnos pendientes del cliente
func (r *OrderRepository) GettingCustomerPendingOrders(ctx context.Context, userID int) ([]models.CustomerPendingOrder, error) {

	var (
		redisCacheKey        string = "CustomerPendingTurns"
		customerPendingTurns []models.CustomerPendingOrder
	)

	// buscar en cache
	if cachedPendingTurns, err := r.RedisConnection.Get(ctx, redisCacheKey).Result(); err == nil {
		json.Unmarshal([]byte(cachedPendingTurns), &customerPendingTurns)
		return customerPendingTurns, nil
	}

	err := r.Connection.WithContext(ctx).Raw(`
		SELECT 
			id,
			title,
			schedule_day_date,
			schedule_start_time,
			created_at,
			updated_at,
			deleted_at
		FROM orders 
		WHERE 
			user_id = ?
			AND deleted_at IS NULL
			AND schedule_day_date >= CURRENT_DATE
			AND EXTRACT(MONTH FROM schedule_day_date) = EXTRACT(MONTH FROM CURRENT_DATE)
		ORDER BY schedule_day_date ASC, schedule_start_time ASC
		LIMIT 5
	`, userID).Scan(&customerPendingTurns).Error

	println("No estaba, peticionar sql")

	if err != nil {
		return nil, err
	}

	// cachear datos
	data, _ := json.Marshal(customerPendingTurns)
	r.RedisConnection.Set(ctx, redisCacheKey, data, 5*time.Minute)

	return customerPendingTurns, nil
}

func (or *OrderRepository) GetOrderByUserID(ctx context.Context, userID int) (*models.Order, error) {

	var order_response models.Order

	result := or.Connection.WithContext(ctx).Last(&order_response)
	if result.Error != nil {
		return nil, result.Error
	}

	return &order_response, nil
}
