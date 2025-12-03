package services

import "context"

// TODO REEMPLAZAR POR SERVICES
type ServicesRepository interface {
	GetByID(ctx context.Context, id uint) (*Service, error)
	GetTotalPriceByIDs(ctx context.Context, serviceIDs []uint) (float64, error)
	SetBookingServices(ctx context.Context, services []BookingServices) error
}
