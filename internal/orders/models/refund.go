package models

import "time"

type RefundRequest struct {
	Order_id            int       `json:"order_id"`
	Shift_id            int       `json:"shift_id"`
	Refund_percentaje   float64   `json:"refund_percenataje"`
	Schedule_day_date   time.Time `json:"schedule_day_date"`
	Schedule_start_time string    `json:"schedule_start_time"`
	Refund_type         string    `json:"refund_type"`
	Title               string    `json:"title"`
}

// DISCOUNT CARD

type Coupon struct {
	Code            string     `db:"code" json:"code"`
	UserID          int        `db:"user_id" json:"user_id"`
	DiscountPercent float64    `db:"discount_percent" json:"discount_percent"`
	Available       bool       `db:"available" json:"available"`
	Used            bool       `db:"used" json:"used"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	AvailableToDate time.Time  `db:"available_to_date" json:"available_to_date"`
	UsedAt          *time.Time `db:"used_at" json:"used_at,omitempty"`
	Coupon_type     string     `db:"coupon_type" json:"type"`
}
