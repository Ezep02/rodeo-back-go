package services

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Title            string  `json:"title" gorm:"type:varchar(150);not null"`
	Description      string  `json:"description" gorm:"type:text"`
	Price            float64 `json:"price" gorm:"type:decimal(12,2)"`
	Created_by_id    int
	Service_Duration int `json:"service_duration" gorm:"default:0"`
}
