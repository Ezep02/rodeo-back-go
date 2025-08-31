package domain

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(45);not null" json:"name"`
	Surname        string    `gorm:"type:varchar(70);default:null" json:"surname"`
	Password       string    `gorm:"type:varchar(70);not null" json:"password"`
	Email          string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	Phone_number   string    `gorm:"type:varchar(30)" json:"phone_number"`
	Is_admin       bool      `gorm:"default:false" json:"is_admin"`
	Is_barber      bool      `gorm:"default:false" json:"is_barber"`
	LastNameChange time.Time `json:"last_name_change"`
	Username       string    `gorm:"type:varchar(45);not null;unique" json:"username"`
	Avatar         string    `json:"avatar"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

//
type Slot struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Date     time.Time `gorm:"not null" json:"date"`
	Time     string    `gorm:"not null" json:"time"`
	IsBooked bool      `gorm:"default:false" json:"is_booked"`
	BarberID uint      `json:"barber_id"`
	Barber   User      `gorm:"foreignKey:BarberID;references:ID" json:"barber"`
}

type Appointment struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ClientName        string    `gorm:"size:100;not null" json:"client_name"`
	ClientSurname     string    `gorm:"size:100;not null" json:"client_surname"`
	SlotID            uint      `json:"slot_id"`
	UserID            uint      `gorm:"foreignKey:UserID;references:ID" json:"user_id"`
	Slot              Slot      `gorm:"foreignKey:SlotID;references:ID" json:"slot"`
	PaymentPercentage int       `gorm:"not null;default:0" json:"payment_percentage"`
	Status            string    `gorm:"size:100;not null;default:'active'" json:"status"`
	Products          []Product `gorm:"many2many:appointment_products;" json:"products"`
	Review            *Review   `gorm:"foreignKey:AppointmentID;" json:"review"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Product struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Name              string     `gorm:"size:100;not null" json:"name"`
	Description       string     `gorm:"size:255" json:"description"`
	Price             float64    `gorm:"not null" json:"price"`
	CategoryID        uint       `gorm:"not null" json:"category_id"` // <- Este es clave
	Category          *Category  `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
	RatingSum         int        `gorm:"default:0" json:"rating_sum"`
	NumberOfReviews   int        `gorm:"default:0" json:"number_of_reviews"`
	PromotionDiscount int        `json:"promotion_discount"` // porcentaje de descuento
	PromotionEndDate  *time.Time `json:"promotion_end_date"`
	HasPromotion      bool       `gorm:"default:false" json:"has_promotion"`
	PreviewUrl        string     `json:"preview_url"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Color     string    `gorm:"size:7" json:"color"` // #RRGGBB
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Review
type Review struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	AppointmentID uint   `json:"appointment_id"`
	Rating        int    `gorm:"not null" json:"rating"`
	Comment       string `json:"comment"`
}

// Analiticas de la franja horaria mas popular
type PopularTimeSlot struct {
	Time     string `json:"time"`
	Bookings int    `json:"bookings"`
}

// Analiticas de la tasa de ocupacion de los slots por mes
type BookingOcupationRate struct {
	Month   string  `json:"month"`
	Occ_pct float64 `json:"ocuppancy_percentage"`
}

// Analiticas de numero de citas por mes
type MonthBookingCount struct {
	Month             string `json:"month"`
	TotalAppointments int    `json:"total_appointments"`
}

// Analiticas del promedio de citas por semana
type WeeklyBookingRate struct {
	Week                string `json:"week"`
	AppointmentThisWeek int    `json:"appointment_this_week"`
}

// Analiticas de nuevos clientes por mes
type NewClientRate struct {
	Month      string `json:"month"`
	NewClients int    `json:"new_clients"`
}

// Analiticas del total de ingresos por mes
type MonthlyRevenue struct {
	Month        string  `json:"month"`
	TotalRevenue float64 `json:"total_revenue"`
}

// Coupon
type Coupon struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Code               string    `gorm:"type:varchar(12);not null" json:"code"`
	UserID             uint      `gorm:"foreingkey:UserID;references:ID" json:"user_id"`
	DiscountPercentage float64   `json:"discount_percentage"`
	IsAvailable        bool      `gorm:"default:true" json:"is_available"`
	CreatedAt          time.Time `json:"created_at"`
	ExpireAt           time.Time `json:"expire_at"`
}

// Information
type BarberInformation struct {
	Member           int     `json:"member"`
	Promedy          float64 `json:"promedy"`
	TotalAppointment int     `json:"total_appointment"`
}

// Posts
type Post struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"foreignKey:UserID;references:ID" json:"user_id"`
	Title       string    `gorm:"type:varchar(12);not null" json:"title"`
	PreviewUrl  string    `json:"preview_url"`
	Description string    `json:"description"`
	IsPublished bool      `gorm:"default:true" json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
