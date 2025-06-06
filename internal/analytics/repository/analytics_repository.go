package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ezep02/rodeo/internal/analytics/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Analytics_repository struct {
	Connection      *gorm.DB
	RedisConnection *redis.Client
}

func NewAnalyticsRepository(Db_cnn *gorm.DB, Redis_cnn *redis.Client) *Analytics_repository {
	return &Analytics_repository{
		Connection:      Db_cnn,
		RedisConnection: Redis_cnn,
	}
}

var (
	statusApproved string = "approved"
)

// Obtener el total de ingresos en el mes, y el promedio en comparacion al mes anterior
func (r *Analytics_repository) GetMonthlyRevenueAndAvgComparedToLastMonth(ctx context.Context) (float64, float64, error) {

	var (
		totalRevenue         float64
		avgComparedLastMonth float64
		monthlyRevenueAndAvg models.MonthlyRevenue
		redisCacheKey        string = "RevenueAndAvg"
	)

	// Verificar que los datos no esten cacheados
	if cachedRevenueAndAvg, err := r.RedisConnection.Get(ctx, redisCacheKey).Result(); err == nil {
		json.Unmarshal([]byte(cachedRevenueAndAvg), &monthlyRevenueAndAvg)
		return monthlyRevenueAndAvg.Total_month_revenue, monthlyRevenueAndAvg.Avg_compared_last_month, nil
	}

	// Iniciar una transaccion
	r.Connection.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// Obtener el total de ingresos del mes actual
		tx.Transaction(func(txSum *gorm.DB) error {
			if sumErr := txSum.Raw(`
				SELECT SUM(price) AS Revenue FROM orders 
				WHERE MONTH(schedule_day_date) = MONTH(CURRENT_DATE)`).Scan(&totalRevenue); sumErr != nil {
				return sumErr.Error
			}

			return nil
		})

		// Obtener el promedio de ingresos del mes actual comparado al anterior
		tx.Transaction(func(txAvg *gorm.DB) error {
			if avgErr := txAvg.Raw(`
				SELECT
					COALESCE(
						AVG(CASE WHEN MONTH(schedule_day_date) = MONTH(CURRENT_DATE) THEN COALESCE(price, 0) END) -
						AVG(CASE WHEN MONTH(schedule_day_date) = MONTH(DATE_SUB(CURDATE(), INTERVAL 1 MONTH)) THEN COALESCE(price, 0) END),
						0
					) AS diferencia_promedios
				FROM orders
				WHERE mp_status = ?`, statusApproved).Scan(&avgComparedLastMonth); avgErr != nil {
				return avgErr.Error
			}
			return nil
		})

		return nil
	})

	// se se consulto DB, entonces cachear el nuevo resultado
	data, _ := json.Marshal(
		models.MonthlyRevenue{
			Total_month_revenue:     totalRevenue,
			Avg_compared_last_month: avgComparedLastMonth})

	r.RedisConnection.Set(ctx, redisCacheKey, data, 5*time.Minute)
	// Devuelve los resultados
	return totalRevenue, avgComparedLastMonth, nil
}

// Obtener el total de citas reservadas en el mes
func (r *Analytics_repository) GetMonthlyAppointmentsAndAvgComparedToLastMonth(ctx context.Context) (int, float64, error) {

	var (
		total_month_appointments int
		avg_compared_last_month  float64
		redisCacheKey            string = "MonthlyAppointmentsAndAvg"
		appointmensAndAvg        models.MonthlyAppointmens
	)

	if cachedAppointmentsAndAvg, err := r.RedisConnection.Get(ctx, redisCacheKey).Result(); err == nil {
		json.Unmarshal([]byte(cachedAppointmentsAndAvg), &appointmensAndAvg)
		return appointmensAndAvg.Total_month_appointments, appointmensAndAvg.Avg_compared_last_month, nil
	}

	r.Connection.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// Contar el total de clientes atentidos en el mes
		tx.Transaction(func(txCountAppointments *gorm.DB) error {
			if txCountAppointmentsErr := txCountAppointments.Raw(`
				SELECT COUNT(*) FROM schedules 
				WHERE MONTH(schedule_day_date) = MONTH(CURRENT_DATE) AND available = ?`, true).
				Scan(&total_month_appointments); txCountAppointmentsErr != nil {
				return txCountAppointments.Error
			}
			return nil
		})

		// calcular promedio de turnos en promedio a el mes anterior
		tx.Transaction(func(txAppointmentsAvg *gorm.DB) error {

			if txAppointmentsAvgErr := txAppointmentsAvg.Raw(`
				SELECT 
					(COALESCE(
						(SELECT AVG(count_appointments) 
						FROM (SELECT COUNT(*) AS count_appointments 
							FROM schedules 
							WHERE MONTH(schedule_day_date) = MONTH(CURRENT_DATE) AND available = ?
							GROUP BY schedule_day_date) AS current_month_avg), 0
					) - 
					COALESCE(
						(SELECT AVG(count_appointments) 
						FROM (SELECT COUNT(*) AS count_appointments 
							FROM schedules 
							WHERE MONTH(schedule_day_date) = MONTH(DATE_SUB(CURDATE(), INTERVAL 1 MONTH)) AND available = ?
							GROUP BY schedule_day_date) AS last_month_avg), 0)) AS diferencia_promedios
				`, false, false).Scan(&avg_compared_last_month); txAppointmentsAvgErr != nil {
				return txAppointmentsAvgErr.Error
			}
			return nil
		})

		return nil
	})

	//cachear la respuesta
	data, _ := json.Marshal(
		models.MonthlyAppointmens{
			Total_month_appointments: total_month_appointments,
			Avg_compared_last_month:  avg_compared_last_month})

	r.RedisConnection.Set(ctx, redisCacheKey, data, 5*time.Minute)
	return total_month_appointments, avg_compared_last_month, nil
}

// Obtener el numero de nuevos clientes nuevo en el mes, y el promedio en comparacion del mes anterior
func (r *Analytics_repository) GetMonthlyNewCustomersAndAvgComparedToLastMonth(ctx context.Context) (int, float64, error) {

	var (
		monthlyNewCustomers           int
		customersAvgComparedLastMonth float64
		redisCacheKey                 string = "MonthlyCustomersAndAvg"
		monthlyCustomersAndAvg        models.MonthlyNewCustomers
	)

	if cachedCustomersAndAvg, err := r.RedisConnection.Get(ctx, redisCacheKey).Result(); err == nil {
		json.Unmarshal([]byte(cachedCustomersAndAvg), &monthlyCustomersAndAvg)
		return monthlyCustomersAndAvg.Total_month_new_users, monthlyCustomersAndAvg.Avg_compared_last_month, nil
	}

	r.Connection.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		tx.Transaction(func(txCustomers *gorm.DB) error {

			if txCustomersErrs := txCustomers.Raw(`
				SELECT COUNT(*) FROM users 
				WHERE MONTH(created_at) = MONTH(CURRENT_DATE)`).Scan(&monthlyNewCustomers); txCustomersErrs != nil {
				return txCustomers.Error
			}
			return nil
		})

		tx.Transaction(func(txCustomersAvg *gorm.DB) error {

			if txCustomersAvgErrs := txCustomersAvg.Raw(`
				SELECT 
					(COALESCE(
						(SELECT AVG(count_new_users) 
						FROM (SELECT COUNT(*) AS count_new_users 
							FROM users 
							WHERE MONTH(created_at) = MONTH(CURRENT_DATE)
							GROUP BY created_at) AS current_month_avg), 0
					) - 
					COALESCE(
						(SELECT AVG(count_new_users) 
						FROM (SELECT COUNT(*) AS count_new_users 
							FROM users 
							WHERE MONTH(created_at) = MONTH(DATE_SUB(CURDATE(), INTERVAL 1 MONTH))
							GROUP BY created_at) AS last_month_avg), 0)) AS diferencia_promedios
			`).Scan(&customersAvgComparedLastMonth); txCustomersAvgErrs != nil {
				return txCustomersAvgErrs.Error
			}
			return nil
		})

		return nil
	})

	// volver a cachear
	data, _ := json.Marshal(
		models.MonthlyNewCustomers{
			Total_month_new_users:   monthlyNewCustomers,
			Avg_compared_last_month: customersAvgComparedLastMonth,
		})

	r.RedisConnection.Set(ctx, redisCacheKey, data, 5*time.Minute)

	return monthlyNewCustomers, customersAvgComparedLastMonth, nil
}

// Obtener el total de cancelaciones/devoluciones en el mes, y el promedio en comparacion del mes anterior TODO
func (r *Analytics_repository) GetMonthlyCancellationsAndAvgComparedToLastMonth(ctx context.Context) (int, float64, error) {
	return 0, 0, nil
}

// Obtener el listado de ingresos por mes
func (r *Analytics_repository) GetCurrentYearMonthlyRevenue(ctx context.Context) ([]models.CurrentYearMonthlyRevenue, error) {

	var (
		currentYearMonthlyRevenue []models.CurrentYearMonthlyRevenue
		redisCacheKey             string = "YearMonthlyRevenue"
	)

	// Verificar que los datos no esten cacheados
	if cachedRevenueAndAvg, err := r.RedisConnection.Get(ctx, redisCacheKey).Result(); err == nil {
		json.Unmarshal([]byte(cachedRevenueAndAvg), &currentYearMonthlyRevenue)
		return currentYearMonthlyRevenue, nil
	}

	err := r.Connection.WithContext(ctx).Raw(`
        SELECT EXTRACT(MONTH FROM schedule_day_date) AS month, SUM(price) AS month_revenue 
        FROM orders 
        WHERE EXTRACT(YEAR FROM schedule_day_date) = EXTRACT(YEAR FROM CURRENT_DATE) 
        AND mp_status = ?
        GROUP BY month 
        ORDER BY month
    `, statusApproved).Scan(&currentYearMonthlyRevenue).Error

	if err != nil {
		return nil, err
	}

	// Si se realizo la consulta, volver a cachear
	data, _ := json.Marshal(currentYearMonthlyRevenue)
	r.RedisConnection.Set(ctx, redisCacheKey, data, 20*time.Minute)

	return currentYearMonthlyRevenue, nil
}

// Obtener una lista del top servicios populares del mes y cantidad de veces que fue elejido
func (r *Analytics_repository) GetMonthlyPopularServices(ctx context.Context) ([]models.MonthlyPopularService, error) {
	var (
		monthlyPopularServices []models.MonthlyPopularService
		redisCacheKey          string = "PopularServices"
	)

	if cachedPopularServices, err := r.RedisConnection.Get(ctx, redisCacheKey).Result(); err == nil {
		json.Unmarshal([]byte(cachedPopularServices), &monthlyPopularServices)
		return monthlyPopularServices, nil
	}

	err := r.Connection.WithContext(ctx).Raw(`
		SELECT title AS Service_name, COUNT(*) AS Service_count 
		FROM orders 
		WHERE mp_status = ? 
		AND EXTRACT(MONTH FROM schedule_day_date) = EXTRACT(MONTH FROM CURRENT_DATE) 
		GROUP BY Service_name 
		ORDER BY Service_count DESC
		LIMIT 3
	`, statusApproved).Scan(&monthlyPopularServices).Error

	if err != nil {
		return nil, err
	}

	log.Println("monthlyPopularServices", monthlyPopularServices)
	// Volver a cachear datos
	data, _ := json.Marshal(monthlyPopularServices)
	r.RedisConnection.Set(ctx, redisCacheKey, data, 5*time.Minute)

	return monthlyPopularServices, nil
}

// Obtener una lista de clientes recurrentes y el total de dinero abonado
func (r *Analytics_repository) GetFrequentCustomers(ctx context.Context) ([]models.FrequentCustomer, error) {

	var (
		TopFrequentCustomers []models.FrequentCustomer
		redisCacheKey        string = "FrequentCustomers"
	)

	if cachedFrequentCustomers, err := r.RedisConnection.Get(ctx, redisCacheKey).Result(); err == nil {
		json.Unmarshal([]byte(cachedFrequentCustomers), &TopFrequentCustomers)
		return TopFrequentCustomers, nil
	}

	// obtener lo usuarios que mas plata gastaron y cuando fue su ultima visita

	err := r.Connection.WithContext(ctx).Raw(`
		SELECT 
			payer_name AS Customer_name,
			payer_surname AS Customer_surname,
			COUNT(user_id) AS Visits_count,
			SUM(price) AS Total_spent,
			MAX(created_at) AS Last_visit
		FROM orders
		GROUP BY payer_name, payer_surname
		ORDER BY Total_spent DESC
		LIMIT 10
	`).Scan(&TopFrequentCustomers).Error

	if err != nil {
		return nil, err
	}

	// Volver a cachear datos
	data, _ := json.Marshal(TopFrequentCustomers)
	r.RedisConnection.Set(ctx, redisCacheKey, data, 5*time.Minute)

	return TopFrequentCustomers, nil
}

// Obtener la cantidad total de cortes realizados por un barbero comparativo mes a mes
func (r *Analytics_repository) GetYearlyBarberHaircuts(ctx context.Context, barberID int) ([]models.MonthlyHaircuts, error) {

	var (
		yearlyHairCuts []models.MonthlyHaircuts
		redisCacheKey  string = "YearlyBarberHaircuts"
	)

	if cachedYearlyBarberHaircuts, err := r.RedisConnection.Get(ctx, redisCacheKey).Result(); err == nil {
		json.Unmarshal([]byte(cachedYearlyBarberHaircuts), &yearlyHairCuts)
		return yearlyHairCuts, nil
	}

	err := r.Connection.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(MONTH FROM schedule_day_date) AS Month, COUNT(*) AS Total_haircuts 
        FROM schedules 
        WHERE EXTRACT(YEAR FROM schedule_day_date) = EXTRACT(YEAR FROM CURRENT_DATE) 
        AND available = ? AND barber_id = ?
        GROUP BY month 
        ORDER BY month
	`, false, barberID).Scan(&yearlyHairCuts).Error

	if err != nil {
		return nil, err
	}

	// cachear los datos
	data, _ := json.Marshal(yearlyHairCuts)
	r.RedisConnection.Set(ctx, redisCacheKey, data, 5*time.Minute)

	return yearlyHairCuts, nil
}
