package repository

import (
	"context"
	"log"
	"strings"

	"github.com/ezep02/rodeo/internal/orders/models"
	"gorm.io/gorm"
)

func (r *OrderRepository) CreatingRefund(ctx context.Context, refund models.RefundRequest) (*models.UpdatedCustomerPendingOrder, error) {

	var (
		expected_status string = "canceled"
	)

	r.Connection.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// ACTUALIZAR ORDER A CANCELED
		tx.Transaction(func(refund_tx *gorm.DB) error {
			if err := refund_tx.Exec(`UPDATE orders SET mp_status=? WHERE id = ?`, expected_status, refund.Order_id).Error; err != nil {
				log.Println("[ERROR CANCELANDO ORDEN]:", err)
				return err
			}
			log.Println("[ESTADO DE ORDEN CANCELADO]")
			return nil
		})

		// ACTUALIZAR DISPONIBILIDAD SCHEDULE
		tx.Transaction(func(schedule_tx *gorm.DB) error {

			if schedule_tx_err := schedule_tx.Exec(`UPDATE schedules SET available = ? WHERE id = ?`, true, refund.Shift_id).Error; schedule_tx_err != nil {
				return schedule_tx_err
			}
			log.Println("[ESTADO DE SCHEDULE DISPONIBLE]")

			return nil
		})

		return nil
	})

	return &models.UpdatedCustomerPendingOrder{
		ID:                  refund.Order_id,
		Title:               refund.Title,
		Schedule_day_date:   &refund.Schedule_day_date,
		Shift_id:            refund.Shift_id,
		Schedule_start_time: refund.Schedule_start_time,
	}, nil
}

func (r *OrderRepository) CreatingCoupon(ctx context.Context, coupon models.Coupon) (models.Coupon, error) {

	r.Connection.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Set in db a new order
		tx.Transaction(func(coupon_tx *gorm.DB) error {

			if coupon_tx_err := coupon_tx.Model(models.Coupon{}).Create(coupon).Error; coupon_tx_err != nil {
				log.Println("[ERROR CREANDO CUPON ROLLING BACK]")
				coupon_tx.Rollback()
				return coupon_tx_err
			}
			log.Println("[setting new order]")
			return nil
		})

		return nil
	})

	return coupon, nil
}

// corrobora que el estado de la orden este cancelado o no
func (r *OrderRepository) CheckingOrderStatus(ctx context.Context, order_id int) (bool, error) {
	var (
		expected_status string = "canceled"
		response_status string
	)

	err := r.Connection.WithContext(ctx).Raw(`SELECT mp_status FROM orders WHERE id = ?`, order_id).Scan(&response_status).Error

	if err != nil {
		return false, err
	}

	if strings.Compare(response_status, expected_status) == 0 {
		return true, nil
	}

	return false, nil
}
