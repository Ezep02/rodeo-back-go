package domain

import "time"

// Modelo enviado por el barbero para luego generar los horarios
type Slot struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BarberID  uint      `json:"barber_id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Modelo que devuelve si el horario esta ocupado o no
type SlotWithStatus struct {
	ID        uint      `json:"id"`
	BarberID  uint      `json:"barber_id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	IsBooked  bool      `json:"is_booked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
