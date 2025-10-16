package media

import "context"

// TODO REEMPLAZAR POR SERVICES
type MediaRepository interface {
	Create(ctx context.Context, data *Medias) error
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, id uint, data *Medias) error
	ListByServiceId(ctx context.Context, id uint) ([]Medias, error)
}
