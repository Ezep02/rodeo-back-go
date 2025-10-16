package usecase

import (
	"context"
	"errors"

	"github.com/ezep02/rodeo/internal/catalog/domain/media"
)

type MediaService struct {
	mediaRepo media.MediaRepository
}

func NewMediaService(mediaRepo media.MediaRepository) *MediaService {
	return &MediaService{mediaRepo}
}

func (uc *MediaService) Create(ctx context.Context, data *media.Medias) error {

	if data.URL == "" {
		return errors.New("no fue posible encontrar la imagen seleccionada")
	}

	if data.ServiceID == 0 {
		return errors.New("error recupeando servicio")
	}

	return uc.mediaRepo.Create(ctx, data)
}

func (uc *MediaService) Delete(ctx context.Context, id uint) error {

	if id == 0 {
		return errors.New("error recupeando servicio")
	}

	return uc.mediaRepo.Delete(ctx, id)
}

func (s *MediaService) Update(ctx context.Context, id uint, data *media.Medias) error {

	if id == 0 {
		return errors.New("error recupeando servicio")
	}

	return s.mediaRepo.Update(ctx, id, data)
}

func (uc *MediaService) ListByServiceId(ctx context.Context, id uint) ([]media.Medias, error) {

	if id == 0 {
		return nil, errors.New("error recupeando servicio")
	}

	return uc.mediaRepo.ListByServiceId(ctx, id)
}
