package models

import "time"

type Refund struct {
	User_id             int       `json:"user_id"`
	Shift_id            int       `json:"shift_id"`
	Service_id          int       `json:"service_id"`
	Payer_name          string    `json:"payer_name"`
	Payer_surname       string    `json:"payer_surname"`
	Refund_percentaje   int       `json:"refund_percenataje"`
	Total_price         float64   `json:"total_price"`
	Schedule_day_date   time.Time `json:"schedule_day_date"`
	Schedule_start_time string    `json:"schedule_start_time"`
}
