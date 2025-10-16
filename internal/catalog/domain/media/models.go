package media

import "time"

type Medias struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ServiceID uint64    `json:"service_id" gorm:"not null;index"`
	URL       string    `json:"url" gorm:"type:text;not null"`
	Type      string    `json:"type" gorm:"type:enum('image','video');default:'image'"`
	Position  int       `json:"position" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
