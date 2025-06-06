package models

import "gorm.io/gorm"

type CustomerServices struct {
	gorm.Model
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Price            int     `json:"price"`
	Service_Duration int     `json:"service_duration"`
	Category         string  `json:"category"`
	Rating           float64 `json:"rating"`
	Reviews          int     `json:"reviews_count"`
	Preview_url      string  `json:"preview_url"`
	Created_by_id    int     `json:"created_by_id"`
}
