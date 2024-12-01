package models

import (
	"gorm.io/gorm"
)

// Schedule model: Representa la tabla "schedules"
type ScheduleResponse struct {
	*gorm.Model
	User_id       int     `gorm:"not null"`
	Schedule_type string  `gorm:"type:enum('Semanal','Personalizado');not null"`
	Start_date    string  `gorm:"not null"`
	End_date      string  `gorm:"null"`
	Day           string  `gorm:"type:enum('Lunes','Martes','Miercoles','Jueves','Viernes','Sabado','Domingo');not null"`
	Shifts        []Shift `gorm:"foreignKey:Schedule_id"`
}

// Shift model: Representa la tabla "shifts"
type Shift struct {
	*gorm.Model
	Schedule_id uint   `json:"Schedule_id" gorm:"not null"`
	Day         string `gorm:"type:enum('Lunes','Martes','Miercoles','Jueves','Viernes','Sabado','Domingo');not null"`
	Start_time  string `gorm:"not null"`
}

type Schedule struct {
	*gorm.Model
	User_id       int    `gorm:"not null"`
	Schedule_type string `gorm:"type:enum('Semanal','Personalizado');not null"`
	Schedule_day  string `gorm:"type:enum('Lunes','Martes','Miercoles','Jueves','Viernes','Sabado','Domingo');not null"`
	Start_date    string `gorm:"not null"`
	End_date      string `gorm:"null"`
}

type ShiftRequest struct {
	Date struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	Day              string `gorm:"type:enum('Lunes','Martes','Miercoles','Jueves','Viernes','Sabado','Domingo');not null"`
	DistributionType string `gorm:"type:enum('Semanal', 'Personalizado')"`
	Shifts           []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"shifts"`
}

// UPDATE MODELS
type ScheduleModifyDay struct {
	Day      string `json:"Day"`
	ShiftAdd []struct {
		CreatedAt   string `json:"CreatedAt,omitempty"`
		Day         string `json:"Day,omitempty"`
		DeletedAt   string `json:"DeletedAt,omitempty"` // Se usa puntero para representar null
		ID          int    `json:"ID"`
		ScheduleID  int    `json:"Schedule_id,omitempty"`
		Start_Time  string `json:"Start_time"`
		UpdatedAt   string `json:"UpdatedAt,omitempty"`
		ShiftStatus string `json:"Shift_status"`
	} `json:"Shift_add"`
	ShiftsDelete []struct {
		ID *int `json:"ID"`
	} `json:"Shifts_delete"`
	Date struct {
		Start string  `json:"start"`         // Fecha de inicio en formato "YYYY-MM-DD"
		End   *string `json:"end,omitempty"` // Fecha de fin en formato "YYYY-MM-DD", puede ser vacío o null
	} `json:"Date"`

	ID               int
	DistributionType string `json:"DistributionType"`
	ScheduleStatus   string `json:"ScheduleStatus"`
}
