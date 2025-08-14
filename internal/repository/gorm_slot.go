package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormSlotRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormSlotRepo(db *gorm.DB, redis *redis.Client) domain.SlotRepository {
	return &GormSlotRepository{db, redis}
}

func (r *GormSlotRepository) Create(ctx context.Context, slot *[]domain.Slot) error {

	// Crear los slots en la base de datos

	return r.db.WithContext(ctx).Create(slot).Error
}

func (r *GormSlotRepository) GetByID(ctx context.Context, id uint) (*domain.Slot, error) {
	var slot domain.Slot

	if err := r.db.WithContext(ctx).Preload("Barber", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "surname")
	}).First(&slot, id).Error; err != nil {
		// Check if the error is a record not found error
		// If so, return a custom error indicating that the Product was not found
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &slot, nil
}

func (r *GormSlotRepository) Update(ctx context.Context, slot *[]domain.Slot) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(slot).Error
}

func (r *GormSlotRepository) Delete(ctx context.Context, slot *[]domain.Slot) error {

	tx := r.db.WithContext(ctx)

	for _, s := range *slot {
		if err := tx.Where("date = ? AND time = ?", s.Date, s.Time).Delete(&domain.Slot{}).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *GormSlotRepository) ListByDate(ctx context.Context, date time.Time) ([]domain.Slot, error) {
	var slots []domain.Slot

	err := r.db.WithContext(ctx).
		Preload("Barber"). // ← ¡esto es obligatorio!
		Where("DATE(date) = ?", date.Format("2006-01-02")).
		Find(&slots).Error

	if err != nil {
		return nil, err
	}

	return slots, nil
}

func (r *GormSlotRepository) ListByDateRange(ctx context.Context, start, end time.Time) ([]domain.Slot, error) {
	var (
		slot     []domain.Slot
		cacheKey = fmt.Sprintf("slot-start:%s-end:%s", start, end)
	)

	// 1. Recuperar datos desde cache
	infoInCache, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(infoInCache), &slot); err != nil {
			log.Println("error decodificando slots desde cache")
		}
		return slot, nil
	}

	// 2. Recuperar desde la base de datos
	if err := r.db.WithContext(ctx).
		Where("slots.date >= ? AND slots.date <= ?", start, end).
		Preload("Barber", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "surname")
		}).
		Find(&slot).Error; err != nil {
		return nil, err
	}

	// 3. Cachear nueva informacion
	slotToByte, err := json.Marshal(slot)
	if err != nil {
		log.Println("Error realizando cache de los productos")
	}

	r.redis.Set(ctx, cacheKey, slotToByte, 1*time.Minute)

	return slot, nil
}
