package promotions

import "time"

type Promotion struct {
	ID        uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ServiceID uint64     `json:"service_id" gorm:"not null;index"`
	Discount  float64    `json:"discount" gorm:"not null"`
	Type      string     `json:"type" gorm:"type:enum('percentage','fixed');default:'percentage'"`
	StartDate time.Time  `json:"start_date" gorm:"autoCreateTime"`
	EndDate   *time.Time `json:"end_date"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
