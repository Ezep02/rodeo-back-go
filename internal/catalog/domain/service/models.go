package service

import "time"

type Service struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	BarberID    uint64    `json:"barber_id" gorm:"not null"`
	PreviewURL  string    `json:"preview_url" gorm:"type:text"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Medias     []Medias    `json:"medias" gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
	Categories []Category  `json:"categories" gorm:"many2many:service_categories"`
	Promotions []Promotion `json:"promotions" gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
}

type Medias struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ServiceID uint64    `json:"service_id" gorm:"not null;index"`
	URL       string    `json:"url" gorm:"type:text;not null"`
	Type      string    `json:"type" gorm:"type:enum('image','video');default:'image'"`
	Position  int       `json:"position" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type Category struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:100;not null;unique"`
	Color     string    `json:"color" gorm:"size:7"` // opcional
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

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

// Estadisticas de los servicios
type ServiceStats struct {
	TotalServices      int64 `json:"total_services"`
	PendingJourneys    int64 `json:"pending_journeys"`
	InProgressJourneys int64 `json:"in_progress_journeys"`
	CompletedJourneys  int64 `json:"completed_journeys"`
	RetiredJourneys    int64 `json:"retired_journeys"`
}
