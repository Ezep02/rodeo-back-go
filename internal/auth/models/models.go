package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `gorm:"type:varchar(45);not null" json:"name"`
	Surname      string `gorm:"type:varchar(70);default:null" json:"surname"`
	Password     string `gorm:"type:varchar(70);not null" json:"password"`
	Email        string `gorm:"type:varchar(255);not null;unique" json:"email"`
	Phone_number string `gorm:"type:varchar(30)" json:"phone_number"`
	Is_admin     bool   `gorm:"default:false" json:"is_admin"`
	Is_barber    bool   `gorm:"default:false" json:"is_barber"`
}

type LogUserReq struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserResetPassowrdReq struct {
	New_password string `json:"new_password"`
	Token        string `json:"token"`
}
