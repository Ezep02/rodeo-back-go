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
	Code              string     `json:"code"`
	UserID            int        `json:"user_id"`
	Refunded_order_id int        `json:"refunded_order_id"`
	DiscountPercent   float64    `json:"discount_percent"`
	Available         bool       `json:"available"`
	Used              bool       `json:"used"`
	CreatedAt         time.Time  `json:"created_at"`
	AvailableToDate   time.Time  `json:"available_to_date"`
	UsedAt            *time.Time `json:"used_at"`
	Coupon_type       string     `json:"coupon_type"`
	Transaction_type  string     `json:"transaction_type"`
}
