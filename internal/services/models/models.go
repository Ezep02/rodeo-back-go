package models

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Title            string `json:"title" gorm:"type:varchar(150);not null"`
	Description      string `json:"description" gorm:"type:text"`
	Price            int    `json:"price" gorm:"type:decimal(12,0)"`
	Created_by_id    uint   `json:"created_by_id" gorm:"type:int;unsigned"`
	Service_Duration int    `json:"service_duration" gorm:"default:0"`
}

type PopularServices struct {
	Title     string  `json:"title"`
	Total_avg float64 `json:"total_avg"`
}

type ServiceRequest struct {
	Title            string `json:"title" gorm:"type:varchar(150);not null"`
	Description      string `json:"description" gorm:"type:text"`
	Price            int    `json:"price" gorm:"type:decimal(12,0)"`
	Service_Duration int    `json:"service_duration" gorm:"default:0"`
}

type Users struct {
	gorm.Model
	Name         string `gorm:"type:varchar(45);not null" json:"name"`
	Surname      string `gorm:"type:varchar(70);default:null" json:"surname"`
	Email        string `gorm:"type:varchar(255);not null;unique" json:"email"`
	Phone_number string `gorm:"type:varchar(30)" json:"phone_number"`
}
