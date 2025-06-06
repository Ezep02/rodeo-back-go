package repository

import (
	"context"
	"log"

	"github.com/ezep02/rodeo/internal/orders/models"
)

func (r *OrderRepository) CreatingReview(ctx context.Context, review_req models.Review) error {

	log.Println("[REVIEW]:", review_req)
	if err := r.Connection.WithContext(ctx).Model(models.Review{}).Create(review_req).Error; err != nil {
		log.Println("Error creando review", err)
		return err
	}

	return nil
}

// Reviews de la pagina principal
func (r *OrderRepository) GettingReviews(ctx context.Context, offset int) (*[]models.ReviewResponse, error) {
	var reviews *[]models.ReviewResponse

	err := r.Connection.WithContext(ctx).Raw(`
		SELECT 
			r.order_id,
			r.schedule_id,
			r.user_id,
			r.rating,
			r.review_status,
			r.comment,
			r.created_at,
			o.title,
			o.schedule_day_date,
			o.schedule_start_time,
			o.payer_name,
			o.payer_surname
		FROM reviews r
		INNER JOIN orders o ON r.order_id = o.id  
		ORDER BY r.created_at DESC
		LIMIT 5 OFFSET ?
	`, offset).Scan(&reviews).Error

	if err != nil {
		return nil, err
	}

	return reviews, nil
}
