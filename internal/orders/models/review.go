package models

import (
	"time"
)

type ReviewData struct {
	Schedule_id int    `json:"schedule_id"`
	Order_id    int    `json:"order_id"`
	Comment     string `json:"comment"`
	Rating      int    `json:"rating"`
}

type Review struct {
	ReviewData
	User_id int `json:"user_id"`
}

type ReviewResponse struct {
	Review
	Title               string    `json:"title"`
	Schedule_day_date   time.Time `json:"schedule_day_date"`
	Schedule_start_time string    `json:"schedule_start_time"`
	Payer_name          string    `json:"payer_name"`
	Payer_surname       string    `json:"payer_surname"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
