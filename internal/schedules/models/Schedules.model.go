package models

import (
	"time"

	"gorm.io/gorm"
)

// UPDATE MODELS
type ScheduleResponse struct {
	*gorm.Model
	Day       string `json:"Day"`
	Barber_id int    `gorm:"not null"`

	ShiftAdd     []Shift `json:"Shift_add"`
	ShiftsDelete []struct {
		ID *int `json:"ID"`
	} `json:"Shifts_delete"`

	Start *time.Time `json:"start"`         // Fecha de inicio en formato "YYYY-MM-DD"
	End   *time.Time `json:"end,omitempty"` // Fecha de fin en formato "YYYY-MM-DD", puede ser vacío o null

	ID               int
	DistributionType string `json:"DistributionType"`
	ScheduleStatus   string `json:"ScheduleStatus"`
}

// Shift model: Representa la tabla "shifts"
type Shift struct {
	*gorm.Model
	Schedule_id     uint   `json:"Schedule_id" gorm:"not null"`
	Day             string `gorm:"type:enum('Lunes','Martes','Miercoles','Jueves','Viernes','Sabado','Domingo');not null"`
	Start_time      string `gorm:"not null"`
	Available       bool   `gorm:"not null;default:true"`
	Created_by_name string `gorm:"not null"`
	ShiftStatus     string `json:"Shift_status"`
}

type Schedule struct {
	*gorm.Model
	Barber_id    int        `gorm:"not null"`
	Schedule_day string     `gorm:"type:enum('Lunes','Martes','Miercoles','Jueves','Viernes','Sabado','Domingo');not null"`
	Start_date   *time.Time `gorm:"not null"`
	End_date     *time.Time `gorm:"null"`
}

type ShiftRequest struct {
	Start     *time.Time `json:"start"`
	End       *time.Time `json:"end"`
	Day       string     `gorm:"type:enum('Lunes','Martes','Miercoles','Jueves','Viernes','Sabado','Domingo');not null"`
	Shift_add []struct {
		Start           string `json:"start"`
		End             string `json:"end"`
		Available       string `json:"available"`
		Created_by_name string `json:"created_by_name"`
	} `json:"shift_add"`
}
