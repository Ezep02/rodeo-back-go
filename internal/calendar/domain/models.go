package domain

import (
	"time"
)

type GoogleCalendarToken struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uint      `gorm:"not null;index"` // FK al usuario
	AccessToken  string    `gorm:"type:text;not null"`
	RefreshToken string    `gorm:"type:text;not null"`
	Expiry       time.Time `gorm:"not null"`
	TokenType    string    `gorm:"type:varchar(50);not null;default:'Bearer'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
