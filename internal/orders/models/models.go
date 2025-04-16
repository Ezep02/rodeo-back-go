package models

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
	Name           string         `json:"name"`
	Surname        string         `json:"surname"`
	Phone          Phone          `json:"phone"`
	Identification Identification `json:"identification"`
	Address        Address        `json:"address"`
}

type Metadata struct {
	Service_duration    int        `json:"service_duration"`
	UserID              uint       `json:"user_id"`
	Barber_id           int        `json:"barber_id"`
	Created_by_id       int        `json:"created_by_id"`
	Schedule_start_time string     `json:"schedule_start_time"`
	Schedule_day_date   *time.Time `json:"schedule_day_date"`
	Shift_id            int        `json:"shift_id"`
	Email               string     `json:"email"`
	Service_id          int        `json:"service_id"`
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
	BackURLs           BackURLs   `json:"back_urls"`
	Items              []Item     `json:"items"`
	Payer              Payer      `json:"payer"`
	Metadata           Metadata   `json:"metadata"`
	NotificationURL    string     `json:"notification_url"`
	Expires            bool       `json:"expires"`
	ExpirationDateFrom *time.Time `json:"expiration_date_from"`
	ExpirationDateTo   *time.Time `json:"expiration_date_to"`
}

type ServiceOrder struct {
	Barber_id           int        `json:"Barber_id"`
	Created_by_id       int        `json:"Created_by_id"`
	Description         string     `json:"Description"`
	Price               int        `json:"Price"`
	Service_duration    int        `json:"Service_duration"`
	Service_id          int        `json:"Service_id"`
	Title               string     `json:"Title"`
	Schedule_start_time string     `json:"Schedule_start_time"`
	Schedule_day_date   *time.Time `json:"Schedule_day_date"`
	Shift_id            int        `json:"Shift_id"`
}

type Order struct {
	*gorm.Model
	Title               string     `json:"title"`
	Price               int        `json:"price"`
	User_id             int        `json:"user_id"`
	Service_id          int        `json:"service_id"`
	Payer_name          string     `json:"payer_name"`
	Payer_surname       string     `json:"payer_surname"`
	Description         string     `json:"description"`
	Email               string     `json:"email"`
	Payer_phone         string     `json:"payer_phone"`
	Date_approved       string     `json:"date_approved"`
	Mp_status           string     `json:"mp_status"`
	Barber_id           int        `json:"barber_id"`
	Created_by_id       int        `json:"created_by_id"`
	Shift_id            int        `json:"shift_id"`
	Schedule_day_date   *time.Time `json:"schedule_day_date"`
	Service_duration    int        `json:"service_duration"`
	Schedule_start_time string     `json:"schedule_start_time"`
}

type Card struct{}

type ChargeDetail struct {
	Accounts struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"accounts"`
	Amounts struct {
		Original float64 `json:"original"`
		Refunded float64 `json:"refunded"`
	} `json:"amounts"`
	ClientID      int         `json:"client_id"`
	DateCreated   string      `json:"date_created"`
	ID            string      `json:"id"`
	LastUpdated   string      `json:"last_updated"`
	Metadata      MetadataRes `json:"metadata"`
	Name          string      `json:"name"`
	RefundCharges []any       `json:"refund_charges"`
	Type          string      `json:"type"`
}

type ChargesExecutionInfo struct {
	InternalExecution struct {
		Date        string `json:"date"`
		ExecutionID string `json:"execution_id"`
	} `json:"internal_execution"`
}

type FeeDetail struct {
	Amount   float64 `json:"amount"`
	FeePayer string  `json:"fee_payer"`
	Type     string  `json:"type"`
}

type MetadataRes struct {
	UserID          int `json:"user_id"`
	ServiceDuration int `json:"service_duration"`
}

type OrderRes struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type PayerInfo struct {
	Email          string `json:"email"`
	ID             string `json:"id"`
	Identification struct {
		Number string `json:"number"`
		Type   string `json:"type"`
	} `json:"identification"`
}

type PaymentMethod struct {
	ID       string `json:"id"`
	IssuerID string `json:"issuer_id"`
	Type     string `json:"type"`
}

type TransactionDetails struct {
	AcquirerReference        any     `json:"acquirer_reference"`
	ExternalResourceURL      any     `json:"external_resource_url"`
	FinancialInstitution     any     `json:"financial_institution"`
	InstallmentAmount        float64 `json:"installment_amount"`
	NetReceivedAmount        float64 `json:"net_received_amount"`
	OverpaidAmount           float64 `json:"overpaid_amount"`
	PayableDeferralPeriod    any     `json:"payable_deferral_period"`
	PaymentMethodReferenceID any     `json:"payment_method_reference_id"`
	TotalPaidAmount          float64 `json:"total_paid_amount"`
}
