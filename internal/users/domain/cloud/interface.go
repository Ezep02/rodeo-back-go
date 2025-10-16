package cloud

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2/api"
)

type CloudinaryRepository interface {
	List(ctx context.Context, next_cursor string) ([]api.BriefAssetResult, string, error)
	Video(ctx context.Context) ([]api.BriefAssetResult, error)
	Upload(ctx context.Context, file io.Reader, filename string) error
	UploadAvatar(ctx context.Context, file io.Reader, filename string) (string, error)
}
