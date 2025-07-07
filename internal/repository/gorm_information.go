package repository

import (
	"context"

	"github.com/ezep02/rodeo/internal/domain"
	"gorm.io/gorm"
)

type GormInfoRepository struct {
	db *gorm.DB
}

func NewGormInfoRepo(db *gorm.DB) domain.InformationRepository {
	return &GormInfoRepository{db}
}

func (r *GormInfoRepository) BarberInformation(ctx context.Context) (*domain.BarberInformation, error) {

	var info *domain.BarberInformation

	if err := r.db.WithContext(ctx).Raw(`
		SELECT 
			(SELECT COUNT(*) FROM users) AS member,
			(SELECT COUNT(*) FROM appointments) AS total_appointment,
			COALESCE((SELECT AVG(rating) FROM reviews), 0) AS promedy
		`).Scan(&info).Error; err != nil {
		return nil, err
	}

	return info, nil
}
