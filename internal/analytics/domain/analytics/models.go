package analytics

// Analiticas de la franja horaria mas popular
type PopularTimeSlot struct {
	Time     string `json:"time"`
	Bookings int    `json:"bookings"`
}

// Analiticas de la tasa de ocupacion de los slots por mes
type BookingOcupationRate struct {
	Month   string  `json:"month"`
	Occ_pct float64 `json:"ocuppancy_percentage"`
}

// Analiticas de numero de citas por mes
type MonthBookingCount struct {
	Month             string `json:"month"`
	TotalAppointments int    `json:"total_appointments"`
}

// Analiticas del promedio de citas por semana
type WeeklyBookingRate struct {
	Week                string `json:"week"`
	AppointmentThisWeek int    `json:"appointment_this_week"`
}

// Analiticas de nuevos clientes por mes
type NewClientRate struct {
	Month      string `json:"month"`
	NewClients int    `json:"new_clients"`
}

// Analiticas del total de ingresos por mes
type MonthlyRevenue struct {
	Month        string  `json:"month"`
	TotalRevenue float64 `json:"total_revenue"`
}
