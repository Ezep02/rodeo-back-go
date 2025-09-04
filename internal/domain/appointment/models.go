package appointment

import (
	"time"

	"github.com/ezep02/rodeo/internal/domain"

	"github.com/ezep02/rodeo/internal/domain/review"
)

type Appointment struct {
	ID                uint             `gorm:"primaryKey" json:"id"`
	ClientName        string           `gorm:"size:100;not null" json:"client_name"`
	ClientSurname     string           `gorm:"size:100;not null" json:"client_surname"`
	SlotID            uint             `json:"slot_id"`
	UserID            uint             `gorm:"foreignKey:UserID;references:ID" json:"user_id"`
	Slot              domain.Slot      `gorm:"foreignKey:SlotID;references:ID" json:"slot"`
	PaymentPercentage int              `gorm:"not null;default:0" json:"payment_percentage"`
	Status            string           `gorm:"size:100;not null;default:'active'" json:"status"`
	Products          []domain.Product `gorm:"many2many:appointment_products;" json:"products"`
	Review            *review.Review   `gorm:"foreignKey:AppointmentID;" json:"review"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}
