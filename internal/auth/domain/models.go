package domain

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(45);not null" json:"name"`
	Surname        string    `gorm:"type:varchar(70);default:null" json:"surname"`
	Password       string    `gorm:"type:varchar(70);not null" json:"password"`
	Email          string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	Phone_number   string    `gorm:"type:varchar(30)" json:"phone_number"`
	Is_admin       bool      `gorm:"default:false" json:"is_admin"`
	Is_barber      bool      `gorm:"default:false" json:"is_barber"`
	LastNameChange time.Time `json:"last_name_change"`
	Username       string    `gorm:"type:varchar(45);not null;unique" json:"username"`
	Avatar         string    `json:"avatar"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
