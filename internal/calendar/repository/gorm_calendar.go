package repository

import (
	"context"

	"github.com/ezep02/rodeo/internal/calendar/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormCalendarRepository struct {
	db *gorm.DB
}

func NewGormCalendarRepo(db *gorm.DB) domain.CalendarRepository {
	return &GormCalendarRepository{db}
}

func (r *GormCalendarRepository) SaveToken(ctx context.Context, userID uint, token *domain.GoogleCalendarToken) error {
	token.UserID = userID

	// Upsert: si ya existe un token para user_id, actualiza los campos
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"access_token", "refresh_token", "expiry", "token_type", "updated_at"}),
	}).Create(token).Error
}

func (r *GormCalendarRepository) GetToken(ctx context.Context, userID uint) (*domain.GoogleCalendarToken, error) {
	var token domain.GoogleCalendarToken
	if err := r.db.Where("user_id = ?", userID).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *GormCalendarRepository) AssignBarberCalendar(ctx context.Context, calendarId string, userId uint) error {
	return r.db.WithContext(ctx).
		Exec("UPDATE barbers SET calendar_id = ? WHERE user_id = ?", calendarId, userId).
		Error
}
