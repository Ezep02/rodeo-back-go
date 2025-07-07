package domain

import (
	"context"
	"errors"
	"time"
)

var (

	// ErrNotFound is returned when an entity is not found
	ErrNotFound = errors.New("entity not found")
	// ErrAlreadyExists is returned when an entity already exists
	ErrAlreadyExists = errors.New("entity already exists")
)

type AuthRepository interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
}

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *Appointment) error
	Update(ctx context.Context, appointment *Appointment, slot_id uint) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]Appointment, error)
	GetByID(ctx context.Context, id uint) (*Appointment, error)
	GetByUserID(ctx context.Context, id uint) ([]Appointment, error)
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id uint) (*Product, error)
	Popular(ctx context.Context) ([]Product, error)
}

type SlotRepository interface {
	Create(ctx context.Context, slot *[]Slot) error
	Update(ctx context.Context, slot *[]Slot) error
	Delete(ctx context.Context, slot *[]Slot) error
	GetByID(ctx context.Context, id uint) (*Slot, error)
	List(ctx context.Context, offset int) ([]Slot, error)
	ListByDate(ctx context.Context, date time.Time) ([]Slot, error)
}

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	Update(ctx context.Context, review *Review) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]Appointment, error)
}

type AnalyticRepository interface {
	PopularTimeSlot(ctx context.Context) ([]PopularTimeSlot, error)
	BookingOcupationRate(ctx context.Context) (*BookingOcupationRate, error)
	MonthBookingCount(ctx context.Context) ([]MonthBookingCount, error)
	WeeklyBookingRate(ctx context.Context) ([]WeeklyBookingRate, error)
	NewClientRate(ctx context.Context) ([]NewClientRate, error)
	MonthlyRevenue(ctx context.Context) ([]MonthlyRevenue, error)
}

type CouponRepository interface {
	Create(ctx context.Context, coupon *Coupon) error
	// Update(ctx context.Context, id uint) (*Coupon, error)
	// ListAll(ctx context.Context) ([]Coupon, error)
	// GetByCode(ctx context.Context, code string) (*Coupon, error)
}

type InformationRepository interface {
	BarberInformation(ctx context.Context) (*BarberInformation, error)
}
