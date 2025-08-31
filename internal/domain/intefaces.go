package domain

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
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
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *Appointment) error
	Update(ctx context.Context, appointment *Appointment, slot_id uint) error
	Delete(ctx context.Context, id uint) error
	ListByDateRange(ctx context.Context, start, end time.Time) ([]Appointment, error)
	GetByID(ctx context.Context, id uint) (*Appointment, error)
	GetByUserID(ctx context.Context, id uint) ([]Appointment, error)
}

// TODO REEMPLAZAR POR SERVICES
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset int) ([]Product, error)
	GetByID(ctx context.Context, id uint) (*Product, error)
	Popular(ctx context.Context) ([]Product, error)
	Promotion(ctx context.Context) ([]Product, error)
}

type SlotRepository interface {
	Create(ctx context.Context, slot *[]Slot) error
	Update(ctx context.Context, slot *[]Slot) error
	Delete(ctx context.Context, slot *[]Slot) error
	GetByID(ctx context.Context, id uint) (*Slot, error)
	ListByDateRange(ctx context.Context, start, end time.Time) ([]Slot, error)
	ListByDate(ctx context.Context, date time.Time) ([]Slot, error)
}

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	Update(ctx context.Context, review *Review) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]Appointment, error)
	ListByProductID(ctx context.Context, productID uint) ([]Review, error)
	ListByUserID(ctx context.Context, userID uint, offset int) ([]Appointment, error)
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
	GetByCode(ctx context.Context, code string) (*Coupon, error)
	GetByUserID(ctx context.Context, userID uint) ([]Coupon, error)
	UpdateStatus(ctx context.Context, code string) error
}

type InformationRepository interface {
	BarberInformation(ctx context.Context) (*BarberInformation, error)
}

type CloudinaryRepository interface {
	List(ctx context.Context, next_cursor string) ([]api.BriefAssetResult, string, error)
	Video(ctx context.Context) ([]api.BriefAssetResult, error)
	Upload(ctx context.Context, file io.Reader, filename string) error
	UploadAvatar(ctx context.Context, file io.Reader, filename string) (string, error)
}

type PostRepository interface {
	List(ctx context.Context, offset int) ([]Post, error)
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, post *Post, post_id uint) error
	Delete(ctx context.Context, post_id uint) error
	GetByID(ctx context.Context, id uint) (*Post, error)
	Count(ctx context.Context) (int64, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]Category, error)
	GetByID(ctx context.Context, id uint) (*Category, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*User, error)
	Update(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdatePassword(ctx context.Context, user *User) error
	UpdateUsername(ctx context.Context, new_username string, id uint) error
	UpdateAvatar(ctx context.Context, avatar string, id uint) error
}
