package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ezep02/rodeo/internal/orders/models"
	"github.com/ezep02/rodeo/internal/orders/services"
	"github.com/ezep02/rodeo/internal/orders/utils"
	"github.com/gorilla/websocket"

	"github.com/ezep02/rodeo/pkg/jwt"

	"github.com/spf13/viper"
)

// var link string = "https://api.mercadopago.com"

type OrderHandler struct {
	ord_srv *services.OrderService
	ctx     context.Context
}

func NewOrderHandler(orders_srv *services.OrderService) *OrderHandler {
	return &OrderHandler{
		ctx:     context.Background(),
		ord_srv: orders_srv,
	}
}

var (
	mp_access_token string
	auth_token      string
)

func init() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer el archivo .env: %v", err)
	}

	auth_token = viper.GetString("AUTH_TOKEN")
	mp_access_token = viper.GetString("MP_ACCESS_TOKEN")

}

func (orh *OrderHandler) CreateOrderHandler(rw http.ResponseWriter, r *http.Request) {

	var (
		newOrder       models.ServiceOrder
		validatedToken *jwt.VerifyTokenRes
		responseBody   map[string]any
	)

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(rw, "No se pudo parsear correctamente el cuerpo de la peticion", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	cookie, err := r.Cookie(auth_token)
	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	token, err := jwt.VerfiyToken(cookie.Value)
	if err != nil {
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
		return
	}
	validatedToken = token

	// crear token transitorio
	now := time.Now()
	orderToken, err := utils.GenerateOrderToken(models.PendingOrderToken{
		Title:               newOrder.Title,
		Payer_name:          validatedToken.Name,
		Payer_surname:       validatedToken.Surname,
		Barber_id:           newOrder.Barber_id,
		Schedule_day_date:   newOrder.Schedule_day_date,
		Schedule_start_time: newOrder.Schedule_start_time,
		User_id:             int(validatedToken.ID),
		Price:               float64(newOrder.Price),
		ID:                  uint(newOrder.Shift_id),
		Created_at:          &now,
	}, time.Now().Add(10*time.Minute))

	if err != nil {
		http.Error(rw, "Algo salio mal intentando generar el token", http.StatusBadRequest)
		return
	}

	// almacenar transitoriamente el token

	if err := orh.ord_srv.SetOrderToken(orh.ctx, orderToken, models.PendingOrderToken{
		Title:               newOrder.Title,
		Payer_name:          validatedToken.Name,
		Payer_surname:       validatedToken.Surname,
		Barber_id:           newOrder.Barber_id,
		Schedule_day_date:   newOrder.Schedule_day_date,
		Schedule_start_time: newOrder.Schedule_start_time,
		User_id:             int(validatedToken.ID),
		ID:                  uint(newOrder.Shift_id),
		Price:               float64(newOrder.Price),
		Created_at:          &now,
	}); err != nil {
		http.Error(rw, "Algo salió mal intentando setear el token", http.StatusBadRequest)
		return
	}

	var success_url string = fmt.Sprintf("http://localhost:5173/payment/success/token=%s", orderToken)

	request := models.Request{
		BackURLs: models.BackURLs{
			Success: success_url,
			Pending: "http://localhost:8080/payment/pending",
			Failure: "http://localhost:8080/payment/failure",
		},

		Items: []models.Item{
			{
				ID:          newOrder.Service_id,
				Title:       newOrder.Title,
				Quantity:    1,
				UnitPrice:   newOrder.Price,
				Description: newOrder.Description,
			},
		},
		Metadata: models.Metadata{
			Service_duration:    newOrder.Service_duration,
			UserID:              validatedToken.ID,
			Barber_id:           newOrder.Barber_id,
			Created_by_id:       newOrder.Created_by_id,
			Schedule_start_time: newOrder.Schedule_start_time,
			Schedule_day_date:   newOrder.Schedule_day_date,
			Shift_id:            newOrder.Shift_id,
			Email:               validatedToken.Email,
		},
		Payer: models.Payer{
			Name:    validatedToken.Name,
			Surname: validatedToken.Surname,
			Phone: models.Phone{
				Number: validatedToken.Phone_number,
			},
		},

		NotificationURL:    "https://885f-181-16-121-41.ngrok-free.app/order/webhook",
		Expires:            true,
		ExpirationDateFrom: func() *time.Time { now := time.Now(); return &now }(),
		ExpirationDateTo:   func(t time.Time) *time.Time { t = t.Add(30 * 24 * time.Hour); return &t }(*newOrder.Schedule_day_date),
	}

	// Serializa el objeto request a JSON
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		http.Error(rw, "Error parseando los datos", http.StatusInternalServerError)
		return
	}

	// Prepara la solicitud HTTP para crear la preferencia
	apiURL := "https://api.mercadopago.com/checkout/preferences"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonRequest))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Establecer las cabeceras necesarias
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mp_access_token)

	//Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	// Leer la respuesta
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con el cuerpo de la respuesta
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(responseBody)
}

func getMap(m map[string]any, key string) map[string]any {
	val, ok := m[key].(map[string]any)
	if !ok {
		log.Printf("Campo '%s' no es un objeto", key)
		return map[string]any{}
	}
	return val
}

func getSliceMap(m map[string]any, key string, index int) map[string]any {
	slice, ok := m[key].([]any)
	if !ok || len(slice) <= index {
		log.Printf("Campo '%s' no es un slice válido", key)
		return map[string]any{}
	}
	val, ok := slice[index].(map[string]any)
	if !ok {
		log.Printf("Elemento %d de '%s' no es un objeto", index, key)
		return map[string]any{}
	}
	return val
}

func getString(m map[string]any, key string) string {
	val, _ := m[key].(string)
	return val
}

func getFloat(m map[string]any, key string) float64 {
	val, _ := m[key].(float64)
	return val
}

func (orh *OrderHandler) WebHook(rw http.ResponseWriter, r *http.Request) {
	var bodyPayment map[string]any
	if err := json.NewDecoder(r.Body).Decode(&bodyPayment); err != nil {
		log.Println("Error decodificando el cuerpo del webhook:", err)
		http.Error(rw, "JSON inválido", http.StatusBadRequest)
		return
	}

	data, ok := bodyPayment["data"].(map[string]any)
	if !ok {
		log.Println("Campo 'data' no encontrado o malformado")
		http.Error(rw, "Campo 'data' inválido", http.StatusBadRequest)
		return
	}

	paymentID := fmt.Sprintf("%v", data["id"])
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creando solicitud GET:", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mp_access_token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("Error haciendo solicitud a MercadoPago:", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var root map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&root); err != nil {
		log.Println("Error decodificando respuesta de MercadoPago:", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	additionalInfo := getMap(root, "additional_info")
	metadata := getMap(root, "metadata")
	payer := getMap(additionalInfo, "payer")
	item := getSliceMap(additionalInfo, "items", 0)

	price := int(getFloat(item, "unit_price"))
	itemID := int(getFloat(item, "id"))
	serviceID := fmt.Sprintf("%d", itemID)

	scheduleDateStr := getString(metadata, "schedule_day_date")
	scheduleDate, err := time.Parse(time.RFC3339, scheduleDateStr)
	if err != nil {
		log.Println("Error parseando fecha de turno:", err)
		http.Error(rw, "Fecha inválida", http.StatusBadRequest)
		return
	}

	newOrder, err := orh.ord_srv.CreateNewOrder(orh.ctx, &models.Order{
		Title:               getString(item, "title"),
		Price:               price,
		Service_duration:    int(getFloat(metadata, "service_duration")),
		User_id:             int(getFloat(metadata, "user_id")),
		Schedule_start_time: getString(metadata, "schedule_start_time"),
		Email:               getString(metadata, "email"),
		Service_id:          serviceID,
		Description:         getString(root, "description"),
		Payer_name:          getString(payer, "first_name"),
		Payer_surname:       getString(payer, "last_name"),
		Date_approved:       getString(root, "date_approved"),
		Mp_status:           getString(root, "status"),
		Barber_id:           int(getFloat(metadata, "barber_id")),
		Schedule_day_date:   &scheduleDate,
		Created_by_id:       int(getFloat(metadata, "created_by_id")),
		Shift_id:            int(getFloat(metadata, "shift_id")),
	})

	if err != nil {
		http.Error(rw, "No se creó la orden", http.StatusBadRequest)
		return
	}

	msgBytes, err := json.Marshal(newOrder)
	if err != nil {
		log.Println("Error al serializar la orden:", err)
		http.Error(rw, "Error al procesar la orden", http.StatusInternalServerError)
		return
	}

	if err := sendMessageToPeer(websocket.TextMessage, msgBytes); err != nil {
		log.Println("Error enviando mensaje al cliente:", err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("ok")
}

func (orh *OrderHandler) Failure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Failure")
}

func (orh *OrderHandler) Pending(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Pending")
}
