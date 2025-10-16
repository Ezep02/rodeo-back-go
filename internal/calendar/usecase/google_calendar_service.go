package usecase

import (
	"context"

	"github.com/ezep02/rodeo/internal/calendar/domain"
)

type CalendarService struct {
	calendarRepo domain.CalendarRepository
}

func NewCalendarService(calendarRepo domain.CalendarRepository) *CalendarService {
	return &CalendarService{calendarRepo}
}

func (s *CalendarService) SaveToken(ctx context.Context, userID uint, token *domain.GoogleCalendarToken) error {
	return s.calendarRepo.SaveToken(ctx, userID, token)
}

func (s *CalendarService) GetToken(ctx context.Context, userID uint) (*domain.GoogleCalendarToken, error) {
	return s.calendarRepo.GetToken(ctx, userID)
}

func (s *CalendarService) AssignBarberCalendar(ctx context.Context, calendarId string, userId uint) error {
	return s.calendarRepo.AssignBarberCalendar(ctx, calendarId, userId)
}
