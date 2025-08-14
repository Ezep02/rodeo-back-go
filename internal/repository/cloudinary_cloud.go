package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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

var (
	rodeo_video_container = "rodeo_video_container"
	rodeo_img_container   = "rodeo_img_container"
)

func (r *CloudinaryRepository) List(ctx context.Context, next_cursor string) ([]api.BriefAssetResult, string, error) {
	var (
		resourceKey string = fmt.Sprintf("cloudinaryImagesKey:next-cursor-%s", next_cursor)
		imgInCache  []api.BriefAssetResult
	)

	//1. Verificar si ya están en el cache
	dataInCache, err := r.redis.Get(ctx, resourceKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(dataInCache), &imgInCache); err == nil {
			return imgInCache, next_cursor, nil
		}
	}

	// 2. Consultar Cloudinary
	res, err := r.cloud.Admin.AssetsByAssetFolder(ctx, admin.AssetsByAssetFolderParams{
		AssetFolder: rodeo_img_container,
		MaxResults:  10,
		NextCursor:  next_cursor,
	})
	if err != nil {
		return nil, next_cursor, err
	}

	// 3. Cachear los datos recuperados
	data, err := json.Marshal(res.Assets)
	if err != nil {
		log.Println("Error serializando imágenes para cache:", err)
	} else {
		err = r.redis.Set(ctx, resourceKey, data, 3*time.Minute).Err()
		if err != nil {
			log.Println("Error guardando en cache:", err)
		}
	}

	return res.Assets, res.NextCursor, nil
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
		AssetFolder: rodeo_video_container,
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

func (r *CloudinaryRepository) Upload(ctx context.Context, file io.Reader, filename string) error {
	_, err := r.cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:    filename,
		AssetFolder: rodeo_img_container,
	})

	if err != nil {
		log.Println("Error uploading file to Cloudinary:", err)
		return err
	}

	log.Println("File uploaded successfully to Cloudinary")
	return nil
}
