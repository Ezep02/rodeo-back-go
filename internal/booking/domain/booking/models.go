package booking

import "time"

type Booking struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	SlotID         uint    `gorm:"not null" json:"slot_id"`
	ClientID       uint    `gorm:"not null" json:"client_id"`
	Status         string  `gorm:"type:enum('pendiente_pago','confirmado','cancelado','rechazado','completado', 'reprogramado');default:'pendiente_pago';not null" json:"status"`
	TotalAmount    float64 `gorm:"type:decimal(10,2);default:0" json:"total_amount"`
	CouponCode     *string `gorm:"size:12" json:"coupon_code"`
	DiscountAmount float64 `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	GoogleEventID  *string `gorm:"size:255" json:"google_event_id"`

	Client          User             `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE" json:"client"`
	Slot            Slot             `gorm:"foreignKey:SlotID;constraint:OnDelete:CASCADE" json:"slot"`
	BookingServices []BookingService `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE" json:"services"`

	ExpiresAt *time.Time `gorm:"default:null" json:"expires_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type User struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(45);not null" json:"name"`
	Surname     string `gorm:"type:varchar(70)" json:"surname"`
	Email       string `gorm:"type:varchar(255);not null;unique" json:"email"`
	PhoneNumber string `gorm:"type:varchar(30)" json:"phone_number"`
	Username    string `gorm:"type:varchar(45);not null;unique" json:"username"`
	Avatar      string `json:"avatar"`
}

type Slot struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	BarberID uint      `json:"barber_id"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

type BookingService struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	BookingID uint `gorm:"not null" json:"booking_id"`
	ServiceID uint `gorm:"not null" json:"service_id"`

	Service Service `gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE" json:"service"`
	Booking Booking `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE" json:"booking"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Service struct {
	ID         uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	PreviewURL string  `gorm:"type:text" json:"preview_url"`
	Name       string  `gorm:"size:100;not null" json:"name"`
	Price      float64 `gorm:"type:decimal(10,2);not null" json:"price"`
}

type BookingStats struct {
	TotalBookings     int64   `json:"total_bookings"`
	PendingBookings   int64   `json:"pending_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	CanceledBookings  int64   `json:"canceled_bookings"`
	ExpectedRevenue   float64 `json:"expected_revenue"`
}

// Devuelve informacion sobre el estado de la reprogramacion
type RescheduleResponse struct {
	RequiresPayment bool    `json:"requires_payment"`
	Amount          float64 `json:"amount,omitempty"`
	Percentage      int     `json:"percentage,omitempty"`
	InitPoint       string  `json:"init_point,omitempty"`
	Free            bool    `json:"free"`
	Reprogrammed    bool    `json:"reprogrammed"`
	Message         string  `json:"message"`
}

// Devuelve informacion sobre el estado de la cancelacion
type CancelationResponse struct {
	RequiresCoupon bool   `json:"requires_coupon"`          // se generara un cupon como devolucion?
	CouponPercent  int    `json:"coupon_percent,omitempty"` // 25, 50, 75, etc.
	LosesDeposit   bool   `json:"loses_deposit,omitempty"`  // indica si pierde la seña
	Canceled       bool   `json:"canceled"`                 // si la cancelacion fue efectuada
	Message        string `json:"message"`                  // explicacion para el usuario
}
