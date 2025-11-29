package payments

import "time"

type Payment struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	BookingID     uint       `gorm:"not null" json:"booking_id"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Type          string     `gorm:"type:enum('total','parcial);default:'total';not null" json:"type"`
	Method        string     `gorm:"type:enum('mercadopago','efectivo','tarjeta','transferencia');not null" json:"method"`
	Status        string     `gorm:"type:enum('pendiente','aprobado','rechazado','reembolsado');default:'pendiente';not null" json:"status"`
	MercadoPagoID *string    `gorm:"size:255" json:"mercado_pago_id"`
	PaymentURL    *string    `gorm:"type:text" json:"payment_url"`
	PaidAt        *time.Time `gorm:"default:null" json:"paid_at"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
