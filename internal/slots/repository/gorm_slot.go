package repository

import (
	"context"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/slots/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormSlotRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGormSlotsRepo(db *gorm.DB, redis *redis.Client) domain.SlotRepository {
	return &GormSlotRepository{db, redis}
}

func (r *GormSlotRepository) CreateInBatches(ctx context.Context, slot *[]domain.Slot) error {

	var (
		batchSize = 100
	)

	return r.db.WithContext(ctx).CreateInBatches(slot, batchSize).Error
}

func (r *GormSlotRepository) Update(ctx context.Context, slot *domain.Slot, slot_id uint) error {
	return nil
}

func (r *GormSlotRepository) ListByDateRange(ctx context.Context, barber_id uint, start, end time.Time) ([]domain.SlotWithStatus, error) {

	var (
		slotList []domain.SlotWithStatus
		//cacheKey    = fmt.Sprintf("barber:%d-slot-start:%s-end:%s", barber_id, start, end)
		parsedStart = start.Truncate(24 * time.Hour)
		parsedEnd   = end.Truncate(24 * time.Hour).Add(24 * time.Hour) // incluye todo el último día
	)

	// 1. Recuperar datos desde cache
	// infoInCache, err := r.redis.Get(ctx, cacheKey).Result()
	// if err == nil {
	// 	if err := json.Unmarshal([]byte(infoInCache), &slotList); err != nil {
	// 		log.Println("error decodificando slots desde cache")
	// 	}
	// 	return slotList, nil
	// }

	// 2. Recuperar desde la base de datos
	if err := r.db.WithContext(ctx).
		Table("slots").
		Select(`
		slots.id,
		slots.barber_id,
		slots.start,
		slots.end,
		CASE 
			WHEN b.id IS NULL THEN FALSE 
			WHEN b.status IN ('pendiente_pago', 'confirmado', 'completado', 'reprogramado') THEN TRUE
			ELSE FALSE 
		END AS is_booked
	`).
		Joins(`
		LEFT JOIN bookings b
		ON b.slot_id = slots.id
	`).
		Where("slots.barber_id = ? AND slots.start BETWEEN ? AND ?", barber_id, parsedStart, parsedEnd).
		Scan(&slotList).Error; err != nil {
		log.Println("List by range err", err)
		return nil, err
	}

	// 3. Cachear nueva informacion
	// slotToByte, err := json.Marshal(slotList)
	// if err != nil {
	// 	log.Println("Error realizando cache de los productos")
	// }

	// r.redis.Set(ctx, cacheKey, slotToByte, 1*time.Minute)

	return slotList, nil
}
