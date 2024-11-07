package orders

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	Quantity         int    `json:"quantity"`
	UnitPrice        int    `json:"unit_price"`
	Description      string `json:"description"`
	CategoryID       string `json:"category_id"`
	Service_Duration int    `json:"service_duration"`
}

type Phone struct {
	AreaCode string `json:"area_code"`
	Number   string `json:"number"`
}

type Identification struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

type Address struct {
	StreetName   string `json:"street_name"`
	StreetNumber int    `json:"street_number"`
	ZipCode      string `json:"zip_code"`
}

type Payer struct {
	Email          string         `json:"email"`
	Name           string         `json:"name"`
	Surname        string         `json:"surname"`
	Phone          Phone          `json:"phone"`
	Identification Identification `json:"identification"`
	Address        Address        `json:"address"`
}

type Metadata struct {
	Service_duration int `json:"service_duration"`
	UserID           int `json:"user_id"`
}

type PaymentMethods struct {
	ExcludedPaymentTypes   []string `json:"excluded_payment_types"`
	ExcludedPaymentMethods []string `json:"excluded_payment_methods"`
	Installments           int      `json:"installments"`
	DefaultPaymentMethodID string   `json:"default_payment_method_id"`
}

type BackURLs struct {
	Success string `json:"success"`
	Failure string `json:"failure"`
	Pending string `json:"pending"`
}

type Request struct {
	AutoReturn          string         `json:"auto_return"`
	BackURLs            BackURLs       `json:"back_urls"`
	StatementDescriptor string         `json:"statement_descriptor"`
	BinaryMode          bool           `json:"binary_mode"`
	ExternalReference   string         `json:"external_reference"`
	Items               []Item         `json:"items"`
	Payer               Payer          `json:"payer"`
	Metadata            Metadata       `json:"metadata"`
	PaymentMethods      PaymentMethods `json:"payment_methods"`
	NotificationURL     string         `json:"notification_url"`
	Expires             bool           `json:"expires"`
	ExpirationDateFrom  string         `json:"expiration_date_from"`
	ExpirationDateTo    string         `json:"expiration_date_to"`
}

type Order struct {
	ID               int     `json:"ID"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Price            float64 `json:"price"`
	Created_by_id    int
	Service_Duration int `json:"service_duration"`
}

type PaymentRequest struct {
	*gorm.Model
	UserID    int `json:"user_id" gorm:"not null"`
	ServiceID int `json:"service_id" gorm:"not null"`
	PaymentID int `json:"payment_id" gorm:"not null"`
	Payment
}

// Payment es la estructura que representará los datos del pago.
type Payment struct {
	ID              int64      `json:"id"`
	Title           string     `json:"description"`
	Price           float64    `json:"transaction_amount"`
	ServiceDuration int        `json:"metadata.service_duration"`
	UserID          int        `json:"metadata.user_id"`
	ServiceID       int        `json:"items.0.id"` // Accede al primer item
	PaymentID       int        `json:"payment_id"`
	PayerName       string     `json:"payer.first_name"`
	PayerSurname    string     `json:"payer.last_name"`
	PayerEmail      string     `json:"payer.email"`
	PayerPhone      string     `json:"payer.phone.number"`
	MpOrderID       int64      `json:"order.id"`
	DateApproved    *time.Time `json:"charges_execution_info.date_approved"`
	MpStatus        string     `json:"status"`
	MpStatusDetail  string     `json:"status_detail"`
}
