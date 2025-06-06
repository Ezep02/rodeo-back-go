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

type BackURLs struct {
	Failure string `json:"failure"`
	Pending string `json:"pending"`
	Success string `json:"success"`
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
	User_id             int        `json:"User_id"`
	Payer_name          string     `json:"Payer_name"`
	Payer_surname       string     `json:"Payer_surname"`
	Payer_email         string     `json:"Payer_email"`
	Payer_phone_number  string     `json:"Payer_phone_number"`
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
	Price               float64    `json:"price"`
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
	Transaction_type    string     `json:"transaction_type"`
}
