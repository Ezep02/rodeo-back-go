package helpers

import (
	"fmt"
	"time"

	"github.com/ezep02/rodeo/internal/orders/models"
)

func BuildOrderFromWebhook(root map[string]any) (*models.Order, error) {
	// helpers para extraer datos seguros
	getString := func(m map[string]any, key string) string {
		if val, ok := m[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		return ""
	}

	getFloat := func(m map[string]any, key string) float64 {
		if val, ok := m[key]; ok {
			if f, ok := val.(float64); ok {
				return f
			}
		}
		return 0
	}

	// obtener metadata y item
	metadata, _ := root["metadata"].(map[string]any)
	additionalInfo, _ := root["additional_info"].(map[string]any)
	items, _ := additionalInfo["items"].([]any)
	item := map[string]any{}

	if len(items) > 0 {
		item = items[0].(map[string]any)
	}

	payer, _ := additionalInfo["payer"].(map[string]any)

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
		Price:               int(getFloat(item, "unit_price")),
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
