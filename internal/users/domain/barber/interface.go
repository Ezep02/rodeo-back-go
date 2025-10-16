package barber

import "context"

type BarberRepository interface {
	GetByID(ctx context.Context, id uint) (*Barber, error)
	List(ctx context.Context) ([]BarberWithUser, error)
}
