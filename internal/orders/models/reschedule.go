package models

import (
	"time"
)

type RescheduleRequest struct {
	Old_schedule_id   int
	Order_id          int
	Shift_id          int
	Barber_id         int
	Service_title     string
	Start_time        string
	Schedule_day_date *time.Time
}
