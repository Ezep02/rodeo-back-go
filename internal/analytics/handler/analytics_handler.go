package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ezep02/rodeo/internal/analytics/models"
	"github.com/ezep02/rodeo/internal/analytics/services"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/spf13/viper"
)

type Analytics_handler struct {
	Analytics_service *services.Analytics_service
	Ctx               context.Context
}

var (
	auth_token string
)

func init() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer el archivo .env: %v", err)
	}

	auth_token = viper.GetString("AUTH_TOKEN")
}

func NewAnalyticsHandler(analytics_srv *services.Analytics_service) *Analytics_handler {
	return &Analytics_handler{
		Analytics_service: analytics_srv,
		Ctx:               context.Background(),
	}
}

// Obtener el total de ingresos en el mes, y el promedio en comparacion al mes anterior
func (h *Analytics_handler) GetMonthlyRevenueAndAvgHandler(rw http.ResponseWriter, r *http.Request) {

	totalRevenue, avg, err := h.Analytics_service.ObtainMonthlyRevenueAndAvgComparedToLastMonth(h.Ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.MonthlyRevenue{
		Total_month_revenue:     totalRevenue,
		Avg_compared_last_month: avg,
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response)
}

// Obtener el total de citas en el mes, y el promedio en comparacion al mes anterior
func (h *Analytics_handler) GetMonthlyAppointmentsAndAvgHandler(rw http.ResponseWriter, r *http.Request) {

	month_appointments, avg_compared_last_month, err := h.Analytics_service.ObtainMonthlyAppointmentsAndAvgComparedToLastMonth(h.Ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.MonthlyAppointmens{
		Total_month_appointments: month_appointments,
		Avg_compared_last_month:  avg_compared_last_month,
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response)

}

// Obtener el numero de nuevos clientes nuevo en el mes, y el promedio en comparacion del mes anterior
func (h *Analytics_handler) GetMonthlyNewCustomersAndAvgHandler(rw http.ResponseWriter, r *http.Request) {

	month_new_users, avg_new_users_compared_last_month, err := h.Analytics_service.ObtainMonthlyNewCustomersAndAvgComparedToLastMonth(h.Ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.MonthlyNewCustomers{
		Total_month_new_users:   month_new_users,
		Avg_compared_last_month: avg_new_users_compared_last_month,
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response)

}

// Obtener el total de cancelaciones/devoluciones en el mes, y el promedio en comparacion del mes anterior
func (h *Analytics_handler) GetMonthlyCancellationsAndAvgHandler(rw http.ResponseWriter, r *http.Request) {

}

func (h *Analytics_handler) GetCurrentYearMonthlyRevenueHandler(rw http.ResponseWriter, r *http.Request) {

	monthlyRevenueList, err := h.Analytics_service.ObtainCurrentYearMonthlyRevenue(h.Ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(monthlyRevenueList)
}

// Obterner un listado con el top de servicios elegidos en el mes
func (h *Analytics_handler) GetMonthlyPopularServicesHandler(rw http.ResponseWriter, r *http.Request) {

	topServicesList, err := h.Analytics_service.ObtainMonthlyPopularServices(h.Ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(topServicesList)
}

// Obtener un listado con un top de clientes recurrentes con el monto total abonado
func (h *Analytics_handler) GetFrequentCustomersHandler(rw http.ResponseWriter, r *http.Request) {

	topFrequentCustomers, err := h.Analytics_service.ObtainFrequentCustomers(h.Ctx)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(topFrequentCustomers)
}

// Obtener un listado mes a mes del total de cortes realizados por el barbero
func (h *Analytics_handler) GetYearlyBarberHaircuts(rw http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)
	if err != nil {
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
		return
	}

	if !token.Is_barber {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	BarberYearlyHaircuts, err := h.Analytics_service.ObtainYearlyBarberHaircuts(h.Ctx, int(token.ID))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(BarberYearlyHaircuts)
}
