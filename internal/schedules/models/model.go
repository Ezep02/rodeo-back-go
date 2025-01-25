package models

import (
	"time"

	"gorm.io/gorm"
)

// Objeto de respuesta
type ScheduleResponse struct {
	*gorm.Model
	Created_by_name   string
	Barber_id         int
	Available         bool
	Schedule_day_date *time.Time
	Start_time        string
}

type Schedule struct {
	*gorm.Model
	Created_by_name   string `gorm:"not null"`
	Barber_id         int    `gorm:"not null"`
	Available         bool
	Start_time        string     `gorm:"not null"`
	Schedule_day_date *time.Time `gorm:"not null"`
}

type CutsQuantity struct {
	Barber_id         int    `json:"Barber_id"`
	Schedule_day_date string `json:"Schedule_day_date"`
	Quantity          int    `json:"Quantity"`
}

type ScheduleRequest struct {
	Schedule_add []struct {
		gorm.Model
		Created_by_name   string
		Barber_id         int
		Available         bool
		Schedule_day_date *time.Time
		Start_time        string
		Schedule_status   string
	} `json:"Schedule_add"`

	Schedule_delete []struct {
		ID int
	} `json:"Schedule_delete"`
}
