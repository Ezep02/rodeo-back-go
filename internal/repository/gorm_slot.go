package repository

import (
	"context"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"gorm.io/gorm"
)

type GormSlotRepository struct {
	db *gorm.DB
}

func NewGormSlotRepo(db *gorm.DB) domain.SlotRepository {
	return &GormSlotRepository{db}
}

func (r *GormSlotRepository) Create(ctx context.Context, slot *[]domain.Slot) error {
	return r.db.WithContext(ctx).Create(slot).Error
}

func (r *GormSlotRepository) GetByID(ctx context.Context, id uint) (*domain.Slot, error) {
	var slot domain.Slot

	if err := r.db.WithContext(ctx).First(&slot, id).Error; err != nil {
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
	var slot []domain.Slot

	if err := r.db.WithContext(ctx).Where("DATE(date) = ?", date.Format("2006-01-02")).Find(&slot).Error; err != nil {
		return nil, err
	}

	return slot, nil
}

func (r *GormSlotRepository) List(ctx context.Context, offset int) ([]domain.Slot, error) {
	var slot []domain.Slot

	if err := r.db.WithContext(ctx).Where("date >= CURRENT_DATE").Limit(31).Offset(offset).Find(&slot).Error; err != nil {
		return nil, err
	}

	return slot, nil
}
