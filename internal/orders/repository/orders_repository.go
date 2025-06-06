package repository

import (
	"context"
	"log"

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

			if order_err := order_tx.Model(models.Order{}).Create(order).Error; order_err != nil {
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

				log.Println("[SECHEDULE: Error actualizando el estado]")
				return nil
			}
			return nil
		})
		return nil
	})

	return order, nil
}

func (r *OrderRepository) GetBarberPendingOrders(ctx context.Context, barberID int, limit int, offset int) ([]models.BarberPendingOrder, error) {

	var (
		barberPendingOrders []models.BarberPendingOrder
	)

	//redisCacheKey := fmt.Sprintf("barber_pending_orders:barber_id-%d", barberID)

	// if ordersInCache, cacheErr := r.RedisConnection.Get(ctx, redisCacheKey).Result(); ordersInCache != "" && cacheErr == nil {
	// 	// devolver los datos en cache
	// 	json.Unmarshal([]byte(ordersInCache), &barberPendingOrders)
	// 	return barberPendingOrders, nil
	// }

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

	// cachear la informacion
	// if pendingOrderBytes, _ := json.Marshal(barberPendingOrders); pendingOrderBytes != nil {
	// 	r.RedisConnection.Set(ctx, redisCacheKey, pendingOrderBytes, time.Minute*5)
	// 	return barberPendingOrders, nil
	// }

	return barberPendingOrders, nil
}

// obtener los turnos pendientes del cliente
func (r *OrderRepository) GettingCustomerPendingOrders(ctx context.Context, userID int) ([]models.CustomerPendingOrder, error) {
	var (
		customerPendingTurns []models.CustomerPendingOrder
		mp_status            string = "approved"
	)

	//customerOrdersCacheKey := fmt.Sprintf("customer_order:id-%d", userID)

	// if cachedCustomerPendingOrders, cacheErr := r.RedisConnection.Get(ctx, customerOrdersCacheKey).Result(); cachedCustomerPendingOrders != "" && cacheErr == nil {
	// 	json.Unmarshal([]byte(cachedCustomerPendingOrders), &customerPendingTurns)
	// 	log.Println("[cache hit] customer pending orders")
	// 	return customerPendingTurns, nil
	// }

	// No estaba en cache o cache inválida, ir a la DB
	dbErr := r.Connection.WithContext(ctx).Raw(`
		SELECT 
			id,
			title,
			shift_id,
			schedule_day_date,
			schedule_start_time,
			created_at,
			updated_at,
			deleted_at
		FROM orders 
		WHERE 
			user_id = ?
			AND mp_status = ?
			AND deleted_at IS NULL
			AND schedule_day_date >= CURRENT_DATE
			AND EXTRACT(MONTH FROM schedule_day_date) = EXTRACT(MONTH FROM CURRENT_DATE)
		ORDER BY schedule_day_date ASC, schedule_start_time ASC
		LIMIT 5
	`, userID, mp_status).Scan(&customerPendingTurns).Error

	if dbErr != nil {
		return nil, dbErr
	}

	// cachear datos
	// if data, _ := json.Marshal(customerPendingTurns); data != nil {
	// 	r.RedisConnection.Set(ctx, customerOrdersCacheKey, data, 5*time.Minute)
	// 	return customerPendingTurns, nil
	// }

	return customerPendingTurns, nil
}

func (r *OrderRepository) ReschedulingDateTimeOrder(ctx context.Context, schedule models.RescheduleRequest, user_id int) (*models.UpdatedCustomerPendingOrder, error) {

	// Update order and update schedule status as available
	r.Connection.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		tx.Transaction(func(order_tx *gorm.DB) error {

			if dbError := order_tx.Exec(`
				UPDATE orders 
				SET shift_id = ?, schedule_day_date = ?, schedule_start_time = ?, barber_id = ?
				WHERE user_id = ? AND id = ?
			`,
				schedule.Shift_id,
				schedule.Schedule_day_date,
				schedule.Start_time,
				schedule.Barber_id,
				user_id,
				schedule.Order_id,
			).Error; dbError != nil {
				log.Println("[rolling back updating order]")
				order_tx.Rollback()
				return dbError
			}

			return nil
		})

		// update schedule as available
		tx.Transaction(func(updating_status_tx *gorm.DB) error {

			if upating_status_err := updating_status_tx.Exec(`
				UPDATE schedules 
				SET available = ?
				WHERE id = ?
			`, true, schedule.Old_schedule_id).Error; upating_status_err != nil {
				return nil
			}
			log.Println("[status available updated]")
			return nil
		})

		// setting new schedule as not available
		tx.Transaction(func(updating_status_tx *gorm.DB) error {

			if upating_status_err := updating_status_tx.Exec(`
				UPDATE schedules 
				SET available = ?
				WHERE id = ?
			`, false, schedule.Shift_id).Error; upating_status_err != nil {
				return nil
			}
			log.Println("[status not available updated]")
			return nil
		})
		return nil
	})

	return &models.UpdatedCustomerPendingOrder{
		ID:                  schedule.Order_id,
		Title:               schedule.Service_title,
		Schedule_day_date:   schedule.Schedule_day_date,
		Schedule_start_time: schedule.Start_time,
		Shift_id:            schedule.Shift_id,
	}, nil
}

// Coupons
func (r *OrderRepository) GettingCustomerCoupons(ctx context.Context, user_id int) (*[]models.Coupon, error) {

	var available_coupons *[]models.Coupon

	if err := r.Connection.WithContext(ctx).Raw(`
		SELECT * FROM coupons 
		WHERE user_id = ? 
		AND available_to_date >= CURRENT_DATE 
		AND available = ?
		LIMIT 10
	`, user_id, true).Scan(&available_coupons).Error; err != nil {
		return nil, err
	}

	return available_coupons, nil
}

// Get customer previous orders
func (r *OrderRepository) GettingCustomerPreviusOrders(ctx context.Context, user_id int, offset int) (*[]models.CustomerPreviusOrders, error) {
	var previousOrders *[]models.CustomerPreviusOrders

	query := `
		SELECT 
			o.id,
			o.shift_id,
			o.title,
			o.schedule_day_date,
			o.schedule_start_time,
			o.payer_name,
			o.payer_surname,
			o.price,
			COALESCE(r.comment, '') AS comment,
			COALESCE(r.rating, 0) AS rating,
			COALESCE(r.review_status, 0) AS review_status
		FROM orders o
		LEFT JOIN reviews r 
			ON r.order_id = o.id 
			AND r.schedule_id = o.shift_id 
			AND r.user_id = ?
		WHERE o.user_id = ?
		ORDER BY o.schedule_day_date DESC
		LIMIT 5 OFFSET ?
	`

	if err := r.Connection.WithContext(ctx).Raw(query, user_id, user_id, offset).Scan(&previousOrders).Error; err != nil {
		log.Println("Error obteniendo las reviews:", err)
		return nil, err
	}

	return previousOrders, nil
}
