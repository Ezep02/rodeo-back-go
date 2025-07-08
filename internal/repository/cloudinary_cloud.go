package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/ezep02/rodeo/internal/domain"
	"github.com/redis/go-redis/v9"
)

type CloudinaryRepository struct {
	cloud *cloudinary.Cloudinary
	redis *redis.Client
}

func NewCloudinaryCloudRepo(cloud *cloudinary.Cloudinary, redis *redis.Client) domain.CloudinaryRepository {
	return &CloudinaryRepository{cloud, redis}
}

func (r *CloudinaryRepository) List(ctx context.Context) ([]api.BriefAssetResult, error) {
	var (
		resourceKey string = "cloudinaryImagesKey"
		imgInCache  []api.BriefAssetResult
	)

	// 1. Verificar si ya están en el cache
	dataInCache, err := r.redis.Get(ctx, resourceKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(dataInCache), &imgInCache); err == nil {
			return imgInCache, nil
		}
	}

	// 2. Consultar Cloudinary
	res, err := r.cloud.Admin.Assets(ctx, admin.AssetsParams{
		AssetType:  "image",
		MaxResults: 10,
	})
	if err != nil {
		return nil, err
	}

	// 3. Cachear los datos recuperados
	data, err := json.Marshal(res.Assets)
	if err != nil {
		log.Println("Error serializando imágenes para cache:", err)
	} else {
		err = r.redis.Set(ctx, resourceKey, data, 10*time.Minute).Err()
		if err != nil {
			log.Println("Error guardando en cache:", err)
		}
	}

	return res.Assets, nil
}

func (r *CloudinaryRepository) Video(ctx context.Context) ([]api.BriefAssetResult, error) {
	var (
		resourceKey string = "cloudinaryVideoKey"
		vidInCache  []api.BriefAssetResult
	)

	// 1. Verificar si ya están en el cache
	dataInCache, err := r.redis.Get(ctx, resourceKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(dataInCache), &vidInCache); err == nil {
			return vidInCache, nil
		}
	}

	resp, err := r.cloud.Admin.AssetsByAssetFolder(ctx, admin.AssetsByAssetFolderParams{
		AssetFolder: "rodeo_video_container",
		MaxResults:  3,
	})

	if err != nil {
		return nil, err
	}

	// 3. Cachear los datos recuperados
	data, err := json.Marshal(resp.Assets)
	if err != nil {
		log.Println("Error serializando imágenes para cache:", err)
	} else {
		err = r.redis.Set(ctx, resourceKey, data, 30*time.Minute).Err()
		if err != nil {
			log.Println("Error guardando en cache:", err)
		}
	}

	return resp.Assets, nil
}
