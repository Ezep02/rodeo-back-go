package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/booking/domain/payments"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormPaymentRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormPaymentRepo(db *gorm.DB, redis *redis.Client) payments.PaymentRepository {
	return &GormPaymentRepository{db: db, redis: redis}
}

func (r *GormPaymentRepository) Create(ctx context.Context, p *payments.Payment) error {
	if p == nil {
		return errors.New("payment es nil")
	}
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *GormPaymentRepository) GetByID(ctx context.Context, paymentID uint) (*payments.Payment, error) {
	var p payments.Payment
	if err := r.db.WithContext(ctx).
		Preload("Booking"). // opcional, si querés traer datos del booking
		Where("id = ?", paymentID).
		First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *GormPaymentRepository) GetByBookingID(ctx context.Context, bookingID uint) ([]payments.Payment, error) {
	var payments []payments.Payment
	if err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *GormPaymentRepository) UpdateStatus(ctx context.Context, paymentID uint, status string, paidAt *time.Time) error {
	if status == "" {
		return errors.New("status no puede ser vacío")
	}
	updates := map[string]any{"status": status}
	if paidAt != nil {
		updates["paid_at"] = *paidAt
	}
	return r.db.WithContext(ctx).
		Model(&payments.Payment{}).
		Where("id = ?", paymentID).
		Updates(updates).Error
}

func (r *GormPaymentRepository) Update(ctx context.Context, p *payments.Payment) error {
	if p == nil {
		return errors.New("payment es nil")
	}
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *GormPaymentRepository) MarkAsPaid(ctx context.Context, paymentID uint, mpPaymentID string) error {

	// return r.db.WithContext(ctx).Save(p).Error
	updates := map[string]any{
		"status":          "aprobado",
		"mercado_pago_id": mpPaymentID,
		"paid_at":         time.Now(),
	}

	if err := r.db.WithContext(ctx).Model(&payments.Payment{}).Where("id = ?", paymentID).Updates(updates).Error; err != nil {
		log.Println("Error updating payments:", err)
		return err
	}

	return nil
}
