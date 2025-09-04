package review

import "time"

// Review
type Review struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	AppointmentID uint   `json:"appointment_id"`
	Rating        int    `gorm:"not null" json:"rating"`
	Comment       string `json:"comment"`
}

// ReviewDetail combina la review, el appointment y el usuario
// Se usa para listar las reviews con informacion adicional en la pagina principal
type ReviewDetail struct {
	ReviewID  uint      `json:"review_id"`
	Rating    uint8     `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`

	AppointmentID     uint   `json:"appointment_id"`
	AppointmentStatus string `json:"appointment_status"`
	ClientName        string `json:"client_name"`
	ClientSurname     string `json:"client_surname"`

	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
