package coupon

import "time"

// Coupon
type Coupon struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Code               string    `gorm:"type:varchar(12);not null" json:"code"`
	UserID             uint      `gorm:"foreingkey:UserID;references:ID" json:"user_id"`
	DiscountPercentage float64   `json:"discount_percentage"`
	IsAvailable        bool      `gorm:"default:true" json:"is_available"`
	CreatedAt          time.Time `json:"created_at"`
	ExpireAt           time.Time `json:"expire_at"`
}
