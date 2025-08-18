package service

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/ezep02/rodeo/internal/domain"
)

type CloudinaryService struct {
	cloudRepo domain.CloudinaryRepository
}

func NewCloudService(cloudRepo domain.CloudinaryRepository) *CloudinaryService {
	return &CloudinaryService{cloudRepo}
}

func (s *CloudinaryService) List(ctx context.Context, next_cursor string) ([]api.BriefAssetResult, string, error) {
	result, nexCursor, err := s.cloudRepo.List(ctx, next_cursor)
	if err != nil {
		return nil, nexCursor, err
	}

	return result, nexCursor, nil
}

func (s *CloudinaryService) Video(ctx context.Context) ([]api.BriefAssetResult, error) {
	return s.cloudRepo.Video(ctx)
}

func (s *CloudinaryService) Upload(ctx context.Context, file io.Reader, filename string) error {
	return s.cloudRepo.Upload(ctx, file, filename)
}
