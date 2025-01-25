package models

import (
	"time"

	"gorm.io/gorm"
)

type Revenue struct {
	MonthDate    time.Time `json:"month_start_date"`
	TotalRevenue float64   `json:"total_revenue"`
}

type Expense struct {
	MonthDate    time.Time `json:"month_start_date"`
	TotalExpense float64   `json:"total_expense"`
}

type User struct {
	gorm.Model
	Name         string `gorm:"type:varchar(45);not null" json:"name"`
	Surname      string `gorm:"type:varchar(70);default:null" json:"surname"`
	Password     string `gorm:"type:varchar(70);not null" json:"password"`
	Email        string `gorm:"type:varchar(255);not null;unique" json:"email"`
	Phone_number string `gorm:"type:varchar(30)" json:"phone_number"`
	Is_admin     bool   `gorm:"default:false" json:"is_admin"`
	Is_barber    bool   `gorm:"default:false" json:"is_barber"`
}

type Schedule struct {
	*gorm.Model
	Created_by_name   string `gorm:"not null"`
	Barber_id         int    `gorm:"not null"`
	Available         bool
	Start_time        string     `gorm:"not null"`
	Schedule_day_date *time.Time `gorm:"not null"`
}

type Expenses struct {
	*gorm.Model
	Created_by_name string `gorm:"not null"`
	Admin_id        int    `gorm:"not null"`
	Title           string `json:"title" gorm:"type:varchar(150);not null"`
	Description     string `json:"description" gorm:"type:text"`
	Amount          int    `json:"amount" gorm:"type:decimal(12,0);not null; default 0"`
}

type ExpenseRequest struct {
	Created_by_name string `gorm:"not null"`
	Admin_id        int    `gorm:"not null"`
	Title           string `json:"title" gorm:"type:varchar(150);not null"`
	Description     string `json:"description" gorm:"type:text"`
	Amount          int    `json:"amount" gorm:"type:decimal(12,0);not null; default 0"`
}
