package models

import (
	"time"
)

type RescheduleRequest struct {
	Old_schedule_id   int        `json:"old_schedule_id"`
	Order_id          int        `json:"order_id"`
	Shift_id          int        `json:"shift_id"`
	Barber_id         int        `json:"barber_id"`
	Service_title     string     `json:"service_title"`
	Start_time        string     `json:"start_time"`
	Schedule_day_date *time.Time `json:"schedule_day_date"`
}
