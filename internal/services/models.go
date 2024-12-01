package services

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Title            string `json:"title" gorm:"type:varchar(150);not null"`
	Description      string `json:"description" gorm:"type:text"`
	Price            int    `json:"price" gorm:"type:decimal(12,0)"`
	Created_by_id    uint   `json:"created_by_id" gorm:"type:int;unsigned"`
	Service_Duration int    `json:"service_duration" gorm:"default:0"`
}
