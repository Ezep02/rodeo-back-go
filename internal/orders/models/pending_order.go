package models

import (
	"time"

	"gorm.io/gorm"
)

type BarberPendingOrder struct {
	*gorm.Model
	Title               string     `json:"title"`
	Payer_name          string     `json:"payer_name"`
	Payer_surname       string     `json:"payer_surname"`
	Barber_id           int        `json:"barber_id"`
	Schedule_day_date   *time.Time `json:"schedule_day_date"`
	Schedule_start_time string     `json:"schedule_start_time"`
	Mp_status           string     `json:"mp_status"`
	Price               float64    `json:"price"`
}

type PendingOrderToken struct {
	ID                  uint       `json:"ID"`
	Title               string     `json:"title"`
	Payer_name          string     `json:"payer_name"`
	Payer_surname       string     `json:"payer_surname"`
	Barber_id           int        `json:"barber_id"`
	User_id             int        `json:"user_id"`
	Schedule_day_date   *time.Time `json:"schedule_day_date"`
	Schedule_start_time string     `json:"schedule_start_time"`
	Price               float64    `json:"price"`
	Created_at          *time.Time `json:"Created_at"`
}

type UpdatedCustomerPendingOrder struct {
	ID                  int        `json:"ID"`
	Title               string     `json:"title"`
	Schedule_day_date   *time.Time `json:"schedule_day_date"`
	Shift_id            int        `json:"shift_id"`
	Schedule_start_time string     `json:"schedule_start_time"`
}

type CustomerPendingOrder struct {
	*gorm.Model
	Shift_id            int        `json:"shift_id"`
	Title               string     `json:"title"`
	Schedule_day_date   *time.Time `json:"schedule_day_date"`
	Schedule_start_time string     `json:"schedule_start_time"`
}

type CustomerPreviusOrders struct {
	CustomerPendingOrder
	Price         float64 `json:"price"`
	Comment       string  `json:"comment"`
	Rating        int     `json:"rating"`
	Review_status bool    `json:"review_status"`
	Payer_name    string  `json:"payer_name"`
	Payer_surname string  `json:"payer_surname"`
}
