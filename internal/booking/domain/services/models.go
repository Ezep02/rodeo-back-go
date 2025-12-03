package services

type Service struct {
	ID    uint64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Price float64 `json:"price" gorm:"type:decimal(10,2);not null"`
}

// Modelo que relaciona los servicios con el bookings
type BookingServices struct {
	ID        uint `json:"id"`
	BookingID uint `json:"booking_id"`
	ServiceID uint `json:"service_id"`
}
