package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/booking/domain/booking"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormBookingRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormBookingRepo(db *gorm.DB, redis *redis.Client) booking.BookingRepository {
	return &GormBookingRepository{db: db, redis: redis}
}

func (r *GormBookingRepository) Create(ctx context.Context, b *booking.Booking) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *GormBookingRepository) UpdateStatus(ctx context.Context, bookingID uint, status string) error {

	return r.db.WithContext(ctx).
		Model(&booking.Booking{}).
		Where("id = ?", bookingID).
		Update("status", status).Error
}

func (r *GormBookingRepository) Update(ctx context.Context, b *booking.Booking) error {
	return r.db.WithContext(ctx).Save(b).Error
}

func (r *GormBookingRepository) GetByID(ctx context.Context, bookingID uint) (*booking.Booking, error) {
	var b booking.Booking
	if err := r.db.WithContext(ctx).
		Preload("Services").
		Preload("Payments").
		Where("id = ?", bookingID).
		First(&b).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// Cuando se paga la cita, se marca como confirmada para que no sea cancelada
func (r *GormBookingRepository) MarkAsPaid(ctx context.Context, bookingID uint) error {
	return r.db.WithContext(ctx).Model(&booking.Booking{}).Where("id = ?", bookingID).Update("status", "confirmado").Error
}

// Devuelve las proximas citas dado un id de barbero
func (r *GormBookingRepository) Upcoming(ctx context.Context, barberID uint, date time.Time, status string) ([]booking.Booking, error) {
	var (
		//bookingKey = fmt.Sprintf("booking-barberID:%d-status:%s", barberID, status)
		bookings []booking.Booking
	)

	// Intentar recuperar del cache
	// if infoInCache, err := r.redis.Get(ctx, bookingKey).Result(); err == nil {
	// 	if err := json.Unmarshal([]byte(infoInCache), &bookings); err == nil {
	// 		return bookings, nil
	// 	}
	// }

	startOfDay := date
	endOfDay := date.Add(24 * time.Hour)

	query := r.db.WithContext(ctx).
		Preload("Client").
		Preload("Slot").
		Preload("Services").
		Joins("JOIN slots s ON s.id = bookings.slot_id").
		Where("s.barber_id = ?", barberID).
		Where("s.start >= ? AND s.start < ?", startOfDay, endOfDay).
		Order("s.start ASC")

	// Filtrar por status si viene
	if status != "" {
		query = query.Where("bookings.status = ?", status)
	}

	// Ejecutar consulta
	if err := query.Find(&bookings).Error; err != nil {
		return nil, err
	}

	// Guardar cache
	// if data, err := json.Marshal(bookings); err == nil {
	// 	r.redis.Set(ctx, bookingKey, data, 3*time.Minute)
	// }

	return bookings, nil
}

// Devuelve las estadisticas de cada barbero
func (r *GormBookingRepository) StatsByBarberID(ctx context.Context, barberID uint) (*booking.BookingStats, error) {

	var (
		statsKey string = fmt.Sprintf("booking-stats:%d", barberID)
		stats           = &booking.BookingStats{}
	)

	//1. Intentar recuperar del cache
	if infoInCache, err := r.redis.Get(ctx, statsKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(infoInCache), &stats); err == nil {
			return stats, nil
		}
	}

	// Query base
	baseQuery := r.db.WithContext(ctx).
		Table("bookings b").
		Joins("JOIN slots s ON s.id = b.slot_id").
		Where("s.barber_id = ?", barberID)

	// Total bookings
	if err := baseQuery.Count(&stats.TotalBookings).Error; err != nil {
		return nil, err
	}

	// Confirmadas
	if err := baseQuery.
		Where("b.status = ?", "confirmado").
		Count(&stats.PendingBookings).Error; err != nil {
		return nil, err
	}

	// Canceladas
	if err := baseQuery.
		Where("b.status = ?", "cancelado").
		Count(&stats.CanceledBookings).Error; err != nil {
		return nil, err
	}

	// Completadas
	if err := baseQuery.
		Where("b.status = ?", "completado").
		Count(&stats.CompletedBookings).Error; err != nil {
		return nil, err
	}

	// Ingresos estimados
	revenueQuery := r.db.WithContext(ctx).
		Table("bookings b").
		Joins("JOIN slots s ON s.id = b.slot_id").
		Where("s.barber_id = ?", barberID)

	if err := revenueQuery.
		Select("COALESCE(SUM(b.total_amount), 0)").
		Where("b.status = ?", "confirmado").
		Scan(&stats.ExpectedRevenue).Error; err != nil {
		return nil, err
	}

	// Guardar cache
	if data, err := json.Marshal(stats); err == nil {
		r.redis.Set(ctx, statsKey, data, 3*time.Minute)
	} else {
		log.Println("Error cacheando estadisticas:", err)
	}

	return stats, nil
}

func (r *GormBookingRepository) AllPendingPayment(ctx context.Context) ([]booking.Booking, error) {
	var bookings []booking.Booking

	if err := r.db.WithContext(ctx).
		Preload("Client", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "surname", "avatar")
		}).
		Preload("Slot", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "start", "end")
		}).
		Where("status = ? AND expires_at IS NULL", "pendiente_pago").
		Find(&bookings).Error; err != nil {
		return nil, err
	}

	return bookings, nil
}

// Proceso en segundo plano para eliminar los bookings que no fueron abonados
func (r *GormBookingRepository) StartBookingCleanupJob(interval time.Duration) {

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			now := time.Now()
			if err := r.db.
				Where("status = ? AND expires_at < ?", "pendiente_pago", now).
				Delete(&booking.Booking{}).Error; err != nil {
				log.Println("Error cancelando bookings expirados:", err)
			}
			log.Println("[CONSULTING EXPIRED BOOKINGS]")
		}
	}()
}
