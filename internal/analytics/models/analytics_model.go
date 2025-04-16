package models

import (
	"time"
)

type MonthlyRevenue struct {
	Total_month_revenue     float64 `json:"Total_month_revenue"`
	Avg_compared_last_month float64 `json:"Avg_compared_last_month"`
}

type MonthlyAppointmens struct {
	Total_month_appointments int     `json:"Total_month_appointments"`
	Avg_compared_last_month  float64 `json:"Avg_compared_last_month"`
}

type MonthlyNewCustomers struct {
	Total_month_new_users   int     `json:"Total_month_new_users"`
	Avg_compared_last_month float64 `json:"Avg_compared_last_month"`
}

type CurrentYearMonthlyRevenue struct {
	Month         string  `json:"Month"`
	Month_revenue float64 `json:"Month_revenue"`
}

type MonthlyPopularService struct {
	Service_name  string `json:"Service_name"`
	Service_count int    `json:"Service_count"`
}

type FrequentCustomer struct {
	Customer_name    string    `json:"Customer_name"`
	Customer_surname string    `json:"Customer_surname"`
	Visits_count     int       `json:"Visits_count"`
	Total_spent      float64   `json:"Total_spent"`
	Last_visit       time.Time `json:"Last_visit"`
}

type MonthlyHaircuts struct {
	Month          int `json:"Month"`
	Total_haircuts int `json:"Total_haircuts"`
}
