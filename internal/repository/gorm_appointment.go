package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/domain/appointment"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormAppointmentRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormAppointmentRepo(db *gorm.DB, redis *redis.Client) appointment.AppointmentRepository {
	return &GormAppointmentRepository{db, redis}
}

func (r *GormAppointmentRepository) Create(ctx context.Context, appt *appointment.Appointment) error {

	var (
		cacheKey = fmt.Sprintf("slot-start:%s-end:%s", time.Now(), time.Now().Add(24*time.Hour))
	)

	// Crear los slots en la base de datos
	if err := r.redis.Del(ctx, cacheKey).Err(); err != nil {
		log.Println("Error invalidating cache after product update:", err)
	}

	r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// Crear appointment
		tx.Transaction(func(appt_tx *gorm.DB) error {

			if err := appt_tx.WithContext(ctx).Create(appt).Error; err != nil {
				log.Println("[rolling back new order]")
				appt_tx.Rollback()
				return err
			}
			log.Println("[OK] New Appointment")
			return nil
		})

		// Actualizar el estado del slot
		tx.Transaction(func(slot_tx *gorm.DB) error {

			if err := slot_tx.WithContext(ctx).Model(&domain.Slot{}).Where("id = ?", appt.SlotID).Update("is_booked", true).Error; err != nil {
				log.Println("[rolling back updating slot]")
				slot_tx.Rollback()
				return err
			}
			log.Println("[OK] Updated Slot")
			return nil
		})
		return nil
	})
	return nil
}

func (r *GormAppointmentRepository) GetByID(ctx context.Context, id uint) (*appointment.Appointment, error) {
	var appt appointment.Appointment

	if err := r.db.WithContext(ctx).Preload("Products").Preload("Slot").First(&appt, id).Error; err != nil {

		// Check if the error is a record not found error
		// If so, return a custom error indicating that the appointment was not found
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &appt, nil
}

func (r *GormAppointmentRepository) ListByDateRange(ctx context.Context, start, end time.Time) ([]appointment.Appointment, error) {
	var (
		appts    []appointment.Appointment
		cacheKey = fmt.Sprintf("appointment-start:%s-end:%s", start, end)
	)

	// 1. Recuperar desde el cache
	infoInCache, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(infoInCache), &appts)
		return appts, nil
	}

	// 2. Realizar consulta
	if err := r.db.WithContext(ctx).
		Joins("JOIN slots ON slots.id = appointments.slot_id").
		Where("slots.date >= ? AND slots.date <= ?", start, end).
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "price")
		}).
		Preload("Slot", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Barber", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "surname")
			})
		}).
		Order("slots.date ASC").
		Find(&appts).Error; err != nil {
		return nil, err
	}

	// 3. Guardar nueva informacion
	apptsToByte, err := json.Marshal(appts)
	if err != nil {
		log.Println("Error realizando cache de los productos")
	}

	r.redis.Set(ctx, cacheKey, apptsToByte, 1*time.Minute)
	return appts, nil
}

func (r *GormAppointmentRepository) Update(ctx context.Context, appt *appointment.Appointment, slot_id uint) error {

	r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// Actualizar el registro actual
		tx.Transaction(func(appt_tx *gorm.DB) error {

			if err := appt_tx.WithContext(ctx).Model(&appointment.Appointment{}).Where("id = ?", appt.ID).Updates(map[string]any{
				"SlotID": appt.SlotID,
				"Status": appt.Status,
			}).Error; err != nil {
				return err
			}
			return nil
		})

		// Ocupar el nuevo turno
		tx.Transaction(func(slot_tx *gorm.DB) error {
			if err := slot_tx.WithContext(ctx).Model(&domain.Slot{}).Where("id = ?", appt.SlotID).Update("is_booked", true).Error; err != nil {
				log.Println("[rolling back updating slot]")
				slot_tx.Rollback()
				return err
			}
			log.Println("[OK] Slot is locked")
			return nil

		})

		// Liberar el slot anterior
		tx.Transaction(func(slot_tx *gorm.DB) error {
			if err := slot_tx.WithContext(ctx).Model(&domain.Slot{}).Where("id = ?", slot_id).Update("is_booked", false).Error; err != nil {
				log.Println("[rolling back updating slot]")
				slot_tx.Rollback()
				return err
			}
			log.Println("[OK] Slot is unlocked")
			return nil

		})

		return nil
	})

	return nil
}

func (r *GormAppointmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var appt appointment.Appointment

		// Cargar cita y slot
		if err := tx.Preload("Slot").First(&appt, id).Error; err != nil {
			log.Println("[ERROR] Appointment not found:", err)
			return err
		}

		// Liberar el slot
		if appt.Slot.ID != 0 {
			if err := tx.Model(&appt.Slot).Update("is_booked", false).Error; err != nil {
				log.Println("[ERROR] Could not unlock slot:", err)
				return err
			}
		}

		// Marcar la cita como cancelada
		if err := tx.Model(&appt).Update("status", "cancelled").Error; err != nil {
			log.Println("[ERROR] Could not update status:", err)
			return err
		}

		return nil
	})
}

func (r *GormAppointmentRepository) GetByUserID(ctx context.Context, id uint, offset int) ([]appointment.Appointment, error) {

	var (
		appt        []appointment.Appointment
		userApptKey string = fmt.Sprintf("user-appointments-id:%d-offset:%d", id, offset)
	)

	log.Println("[DEBUG] userApptKey:", userApptKey)

	// 1. Recuperar productos del cache
	infoInCache, err := r.redis.Get(ctx, userApptKey).Result()
	if err == nil {
		json.Unmarshal([]byte(infoInCache), &appt)
		return appt, nil
	}

	// TODO: cachear informacion
	if err := r.db.WithContext(ctx).
		Select("id", "slot_id", "payment_percentage", "created_at", "status").
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "price")
		}).
		Preload("Slot", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "date", "time")
		}).
		Preload("Review", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "appointment_id", "comment", "rating", "created_at")
		}).
		Where("user_id = ?", id).
		Offset(offset).
		Limit(5).
		Order("created_at DESC").
		Find(&appt).Error; err != nil {
		return nil, err
	}

	// 3. Cachear los datos recuperados
	data, err := json.Marshal(appt)
	if err != nil {
		log.Println("Error realizando cache de los productos")
	}
	r.redis.Set(ctx, userApptKey, data, 1*time.Minute)

	return appt, nil
}
