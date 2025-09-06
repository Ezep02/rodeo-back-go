package review

import "context"

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	Update(ctx context.Context, review *Review) error
	Delete(ctx context.Context, id, user_id uint) error
	List(ctx context.Context, offset int) ([]ReviewDetail, error)
	ListByProductID(ctx context.Context, productID uint) ([]Review, error)
	ListByUserID(ctx context.Context, userID uint, offset int) ([]ReviewDetail, error)
	ReviewRatingStats(ctx context.Context) (*ReviewRatingStats, error)
	GetByID(ctx context.Context, id uint) (*Review, error)
}
