package services

type Service struct {
	ID    uint64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Price float64 `json:"price" gorm:"type:decimal(10,2);not null"`
}
