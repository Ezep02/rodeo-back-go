package information

import "context"

type InformationRepository interface {
	BarberInformation(ctx context.Context) (*BarberInformation, error)
}
