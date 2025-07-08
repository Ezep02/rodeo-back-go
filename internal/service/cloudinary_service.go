package service

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/ezep02/rodeo/internal/domain"
)

type CloudinaryService struct {
	cloudRepo domain.CloudinaryRepository
}

func NewCloudService(cloudRepo domain.CloudinaryRepository) *CloudinaryService {
	return &CloudinaryService{cloudRepo}
}

func (s *CloudinaryService) List(ctx context.Context) ([]api.BriefAssetResult, error) {
	return s.cloudRepo.List(ctx)
}

func (s *CloudinaryService) Video(ctx context.Context) ([]api.BriefAssetResult, error) {
	return s.cloudRepo.Video(ctx)
}
