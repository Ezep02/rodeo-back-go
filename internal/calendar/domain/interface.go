package domain

import "context"

type CalendarRepository interface {
	SaveToken(ctx context.Context, userID uint, token *GoogleCalendarToken) error
	GetToken(ctx context.Context, userID uint) (*GoogleCalendarToken, error)
	AssignBarberCalendar(ctx context.Context, calendarId string, userId uint) error
}
