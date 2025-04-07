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
	Email          string         `json:"email"`
	Name           string         `json:"name"`
	Surname        string         `json:"surname"`
	Phone          Phone          `json:"phone"`
	Identification Identification `json:"identification"`
	Address        Address        `json:"address"`
}

type Metadata struct {
	Service_duration    int        `json:"service_duration"`
	UserID              uint       `json:"user_id"`
	Barber_id           int        `json:"Barber_id"`
	Created_by_id       int        `json:"Created_by_id"`
	Schedule_start_time string     `json:"Schedule_start_time"`
	Schedule_day_date   *time.Time `json:"Schedule_day_date"`
	Shift_id            int
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
	ExpirationDateFrom  *time.Time     `json:"expiration_date_from"`
	ExpirationDateTo    *time.Time     `json:"expiration_date_to"`
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
	Price               string     `json:"price"`
	User_id             int        `json:"user_id"`
	Service_id          string     `json:"service_id"`
	Payment_id          int        `json:"payment_id"`
	Payer_name          string     `json:"payer_name"`
	Payer_surname       string     `json:"payer_surname"`
	Email               string     `json:"email"`
	Payer_phone         string     `json:"payer_phone"`
	Mp_order_id         int64      `json:"mp_order_id"`
	Date_approved       string     `json:"date_approved"`
	Mp_status           string     `json:"mp_status"`
	Mp_status_detail    string     `json:"mp_status_datail"`
	Barber_id           int        `json:"barber_id"`
	Created_by_id       int        `json:"created_by_id"`
	Shift_id            int        `json:"shift_id"`
	Schedule_day_date   *time.Time `json:"schedule_day_date"`
	Service_duration    int        `json:"service_duration"`
	Schedule_start_time string     `json:"schedule_start_time"`
}

type PaymentResponse struct {
	AccountsInfo           any   `json:"accounts_info"`
	AcquirerReconciliation []any `json:"acquirer_reconciliation"`
	AdditionalInfo         struct {
		AuthenticationCode any       `json:"authentication_code"`
		AvailableBalance   any       `json:"available_balance"`
		IpAddress          string    `json:"ip_address"`
		Items              []ItemRes `json:"items"`
		NsuProcessadora    any       `json:"nsu_processadora"`
		Payer              PayerRes  `json:"payer"`
	} `json:"additional_info"`
	AuthorizationCode     any                  `json:"authorization_code"`
	BinaryMode            bool                 `json:"binary_mode"`
	BrandID               any                  `json:"brand_id"`
	BuildVersion          string               `json:"build_version"`
	CallForAuthorizeID    any                  `json:"call_for_authorize_id"`
	Captured              bool                 `json:"captured"`
	Card                  Card                 `json:"card"`
	ChargesDetails        []ChargeDetail       `json:"charges_details"`
	ChargesExecutionInfo  ChargesExecutionInfo `json:"charges_execution_info"`
	CollectorID           int                  `json:"collector_id"`
	CurrencyID            string               `json:"currency_id"`
	DateApproved          string               `json:"date_approved"`
	DateCreated           string               `json:"date_created"`
	DateLastUpdated       string               `json:"date_last_updated"`
	Description           string               `json:"description"`
	DifferentialPricingID any                  `json:"differential_pricing_id"`
	ExternalReference     string               `json:"external_reference"`
	FeeDetails            []FeeDetail          `json:"fee_details"`
	ID                    int                  `json:"id"`
	Installments          int                  `json:"installments"`
	LiveMode              bool                 `json:"live_mode"`
	Metadata              Metadata             `json:"metadata"`
	MoneyReleaseDate      string               `json:"money_release_date"`
	MoneyReleaseStatus    string               `json:"money_release_status"`
	NotificationURL       string               `json:"notification_url"`
	OperationType         string               `json:"operation_type"`
	Order                 OrderRes             `json:"order"`
	PayerInfo             PayerInfo            `json:"payer"`
	PaymentMethod         PaymentMethod        `json:"payment_method"`
	PaymentMethodID       string               `json:"payment_method_id"`
	PaymentTypeID         string               `json:"payment_type_id"`
	ProcessingMode        string               `json:"processing_mode"`
	Refunds               []any                `json:"refunds"`
	Status                string               `json:"status"`
	StatusDetail          string               `json:"status_detail"`
	TransactionAmount     float64              `json:"transaction_amount"`
	TransactionDetails    TransactionDetails   `json:"transaction_details"`
}

type ItemRes struct {
	CategoryID  any    `json:"category_id"`
	Description any    `json:"description"`
	ID          string `json:"id"`
	PictureURL  any    `json:"picture_url"`
	Quantity    string `json:"quantity"`
	Title       string `json:"title"`
	UnitPrice   string `json:"unit_price"`
}

type PayerRes struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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

// REFOUND
// Define una estructura que coincida con la respuesta
type RefundResponse struct {
	ID                   int        `json:"id"`
	PaymentID            int        `json:"payment_id"`
	Amount               float64    `json:"amount"`
	Metadata             []struct{} `json:"metadata"`
	Source               []Source   `json:"source"`
	DateCreated          string     `json:"date_created"`
	UniqueSequenceNumber any        `json:"unique_sequence_number"`
	RefundMode           string     `json:"refund_mode"`
	AdjustmentAmount     float64    `json:"adjustment_amount"`
	Status               int        `json:"status"`
	Reason               any        `json:"reason"`
	Label                []struct{} `json:"label"`
	PartitionDetails     []struct{} `json:"partition_details"`
}

type Source struct {
	Name struct {
		EN string `json:"en"`
		PT string `json:"pt"`
		ES string `json:"es"`
	} `json:"name"`
	ID   string `json:"id"`
	Type string `json:"type"`
}
