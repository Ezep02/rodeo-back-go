package promotions

import "context"

// TODO REEMPLAZAR POR SERVICES
type PromoRepository interface {
	Create(ctx context.Context, data *Promotion) error
	ListByServiceId(ctx context.Context, id uint, page int) ([]Promotion, error)
	Update(ctx context.Context, id uint, data *Promotion) error
	Delete(ctx context.Context, id uint) error

	// GetByID(ctx context.Context, id uint) (*Promotion, error)
}
