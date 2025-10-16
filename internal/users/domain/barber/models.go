package barber

import "time"

type Barber struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CalendarID string    `gorm:"type:varchar(255)" json:"calendar_id"`
	UserID     uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type BarberWithUser struct {
	ID       uint   `gorm:"column:id" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Surname  string `gorm:"column:surname" json:"surname"`
	Avatar   string `gorm:"column:avatar" json:"avatar"`
	Username string `gorm:"column:username" json:"username"`
}
