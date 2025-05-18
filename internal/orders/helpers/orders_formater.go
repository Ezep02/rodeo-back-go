package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/orders/models"
)

func BuildOrderFromWebhook(root map[string]any) (*models.Order, error) {

	// obtener metadata y item
	metadata, _ := root["metadata"].(map[string]any)
	additionalInfo, _ := root["additional_info"].(map[string]any)
	items, _ := additionalInfo["items"].([]any)
	item := map[string]any{}

	if len(items) > 0 {
		item = items[0].(map[string]any)
	}

	payer, _ := additionalInfo["payer"].(map[string]any)

	// helpers para extraer datos seguros
	getString := func(m map[string]any, key string) string {
		if val, ok := m[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		return ""
	}

	// log.Println("ROOT", root)
	getFloat := func(m map[string]any, key string) float64 {
		if val, ok := m[key]; ok {
			switch v := val.(type) {
			case float64:
				return v
			case int:
				return float64(v)
			case int64:
				return float64(v)
			case string:
				res, _ := strconv.Atoi(v)
				return float64(res)
			case json.Number:
				f, _ := v.Float64()
				return f
			}
		}
		return 0
	}

	// convertir fecha
	scheduleDateStr := getString(metadata, "schedule_day_date")
	var scheduleDate time.Time
	if scheduleDateStr != "" {
		t, err := time.Parse(time.RFC3339, scheduleDateStr)
		if err == nil {
			scheduleDate = t
		}
	}

	// construir la orden
	order := &models.Order{
		Title:               getString(item, "title"),
		Price:               getFloat(item, "unit_price"),
		Service_duration:    int(getFloat(metadata, "service_duration")),
		User_id:             int(getFloat(metadata, "user_id")),
		Schedule_start_time: getString(metadata, "schedule_start_time"),
		Email:               getString(metadata, "email"),
		Service_id:          int(getFloat(metadata, "service_id")),
		Description:         getString(root, "description"),
		Payer_name:          getString(payer, "first_name"),
		Payer_surname:       getString(payer, "last_name"),
		Date_approved:       getString(root, "date_approved"),
		Mp_status:           getString(root, "status"),
		Barber_id:           int(getFloat(metadata, "barber_id")),
		Schedule_day_date:   &scheduleDate,
		Created_by_id:       int(getFloat(metadata, "created_by_id")),
		Shift_id:            int(getFloat(metadata, "shift_id")),
	}

	return order, nil
}

func BuildOrderPreference(service_order models.ServiceOrder, orderToken string) (models.Request, error) {

	var success_url string = fmt.Sprintf("http://localhost:5173/payment/success/token=%s", orderToken)

	return models.Request{
		BackURLs: models.BackURLs{
			Success: success_url,
			Pending: "http://localhost:8080/payment/pending",
			Failure: "http://localhost:8080/payment/failure",
		},

		Items: []models.Item{
			{
				ID:          service_order.Service_id,
				Title:       service_order.Title,
				Quantity:    1,
				UnitPrice:   service_order.Price,
				Description: service_order.Description,
			},
		},
		Metadata: models.Metadata{
			UserID:              uint(service_order.User_id),
			Barber_id:           service_order.Barber_id,
			Service_id:          service_order.Service_id,
			Created_by_id:       service_order.Created_by_id,
			Shift_id:            service_order.Shift_id,
			Email:               service_order.Payer_email,
			Service_duration:    service_order.Service_duration,
			Schedule_start_time: service_order.Schedule_start_time,
			Schedule_day_date:   service_order.Schedule_day_date,
		},
		Payer: models.Payer{
			Name:    service_order.Payer_name,
			Surname: service_order.Payer_surname,
			Phone: models.Phone{
				Number: service_order.Payer_phone_number,
			},
		},
		NotificationURL: "https://2da0-181-16-121-41.ngrok-free.app/order/webhook",

		Expires:            true,
		ExpirationDateFrom: func() *time.Time { now := time.Now(); return &now }(),
		ExpirationDateTo:   func(t time.Time) *time.Time { t = t.Add(30 * 24 * time.Hour); return &t }(*service_order.Schedule_day_date),
	}, nil
}
