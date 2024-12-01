package orders

import (
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
	Service_duration int  `json:"service_duration"`
	UserID           uint `json:"user_id"`
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

type ServiceOrder struct {
	ID               int     `json:"ID"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Price            float64 `json:"price"`
	Created_by_id    int
	Service_Duration int `json:"service_duration"`
}

// type PaymentRequest struct {
// 	*gorm.Model
// 	UserID    int `json:"user_id" gorm:"not null"`
// 	ServiceID int `json:"service_id" gorm:"not null"`
// 	PaymentID int `json:"payment_id" gorm:"not null"`
// 	Payment
// }

// Payment es la estructura que representará los datos del pago.
type Order struct {
	*gorm.Model
	Title            string
	Price            string
	Service_Duration int
	User_id          int
	Service_id       string
	Payment_id       int
	Payer_name       string
	Payer_surname    string
	Email            string
	Payer_phone      string
	Mp_order_id      int64
	Date_approved    string
	Mp_status        string
	Mp_status_detail string
}

type PaymentResponse struct {
	AccountsInfo           interface{}   `json:"accounts_info"`
	AcquirerReconciliation []interface{} `json:"acquirer_reconciliation"`
	AdditionalInfo         struct {
		AuthenticationCode interface{} `json:"authentication_code"`
		AvailableBalance   interface{} `json:"available_balance"`
		IpAddress          string      `json:"ip_address"`
		Items              []ItemRes   `json:"items"`
		NsuProcessadora    interface{} `json:"nsu_processadora"`
		Payer              PayerRes    `json:"payer"`
	} `json:"additional_info"`
	AuthorizationCode     interface{}          `json:"authorization_code"`
	BinaryMode            bool                 `json:"binary_mode"`
	BrandID               interface{}          `json:"brand_id"`
	BuildVersion          string               `json:"build_version"`
	CallForAuthorizeID    interface{}          `json:"call_for_authorize_id"`
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
	DifferentialPricingID interface{}          `json:"differential_pricing_id"`
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
	Refunds               []interface{}        `json:"refunds"`
	Status                string               `json:"status"`
	StatusDetail          string               `json:"status_detail"`
	TransactionAmount     float64              `json:"transaction_amount"`
	TransactionDetails    TransactionDetails   `json:"transaction_details"`
}

type ItemRes struct {
	CategoryID  interface{} `json:"category_id"`
	Description interface{} `json:"description"`
	ID          string      `json:"id"`
	PictureURL  interface{} `json:"picture_url"`
	Quantity    string      `json:"quantity"`
	Title       string      `json:"title"`
	UnitPrice   string      `json:"unit_price"`
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
	ClientID      int           `json:"client_id"`
	DateCreated   string        `json:"date_created"`
	ID            string        `json:"id"`
	LastUpdated   string        `json:"last_updated"`
	Metadata      MetadataRes   `json:"metadata"`
	Name          string        `json:"name"`
	RefundCharges []interface{} `json:"refund_charges"`
	Type          string        `json:"type"`
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
	AcquirerReference        interface{} `json:"acquirer_reference"`
	ExternalResourceURL      interface{} `json:"external_resource_url"`
	FinancialInstitution     interface{} `json:"financial_institution"`
	InstallmentAmount        float64     `json:"installment_amount"`
	NetReceivedAmount        float64     `json:"net_received_amount"`
	OverpaidAmount           float64     `json:"overpaid_amount"`
	PayableDeferralPeriod    interface{} `json:"payable_deferral_period"`
	PaymentMethodReferenceID interface{} `json:"payment_method_reference_id"`
	TotalPaidAmount          float64     `json:"total_paid_amount"`
}
