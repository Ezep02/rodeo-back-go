package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/orders/models"
	"github.com/ezep02/rodeo/internal/orders/services"
	"github.com/ezep02/rodeo/internal/orders/utils"

	"github.com/ezep02/rodeo/pkg/jwt"

	"github.com/gorilla/websocket"
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
		AutoReturn: "approved",
		BackURLs: models.BackURLs{
			Success: success_url,
			Pending: "http://localhost:8080/payment/pending",
			Failure: "http://localhost:8080/payment/failure",
		},
		StatementDescriptor: "TestStore",
		BinaryMode:          false,
		ExternalReference:   "IWD1238971",
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
		},
		Payer: models.Payer{
			Email:   validatedToken.Email,
			Name:    validatedToken.Name,
			Surname: validatedToken.Surname,
			Phone: models.Phone{
				Number: validatedToken.Phone_number,
			},
		},
		PaymentMethods: models.PaymentMethods{
			ExcludedPaymentTypes:   []string{},
			ExcludedPaymentMethods: []string{},
			Installments:           12,
			DefaultPaymentMethodID: "account_money",
		},
		NotificationURL:    "https://f8ce-181-16-121-41.ngrok-free.app/order/webhook",
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

// WebHook maneja las solicitudes de webhook de Mercado Pago.
func (orh *OrderHandler) WebHook(rw http.ResponseWriter, r *http.Request) {

	var (
		bodyData map[string]any
		payment  models.PaymentResponse
	)

	// Leer el cuerpo de la solicitud
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error leyendo el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo JSON en un mapa de interfaces
	if err := json.Unmarshal(body, &bodyData); err != nil {
		http.Error(rw, "Error decoding request body", http.StatusBadRequest)
		return
	}

	data, ok := bodyData["data"].(map[string]any)
	if !ok {
		http.Error(rw, "Error: 'data' field is missing or invalid", http.StatusBadRequest)
		return
	}

	idStr, ok := data["id"].(string)
	if !ok {
		http.Error(rw, "Error: 'id' field is missing or invalid", http.StatusBadRequest)
		return
	}

	// Convertir el ID de cadena a entero
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Error: 'id' is not a valid integer", http.StatusBadRequest)
		return
	}

	// URL de la API de pagos, reemplazando :id con el ID real
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%d", id)

	// Crear la solicitud HTTP GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Establecer las cabeceras necesarias
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mp_access_token)

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		http.Error(rw, "Error leyendo el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo JSON en el objeto Payment
	if err := json.Unmarshal(b, &payment); err != nil {
		http.Error(rw, "Error decoding request body", http.StatusBadRequest)
		return
	}

	newOrder, err := orh.ord_srv.CreateNewOrder(orh.ctx, &models.Order{
		Title:               payment.AdditionalInfo.Items[0].Title,
		Price:               payment.AdditionalInfo.Items[0].UnitPrice,
		Service_duration:    payment.Metadata.Service_duration,
		User_id:             int(payment.Metadata.UserID),
		Service_id:          payment.AdditionalInfo.Items[0].ID,
		Payment_id:          payment.ID,
		Payer_name:          payment.AdditionalInfo.Payer.FirstName,
		Payer_surname:       payment.AdditionalInfo.Payer.LastName,
		Email:               payment.PayerInfo.Email,
		Mp_order_id:         int64(payment.ID),
		Date_approved:       payment.DateApproved,
		Mp_status:           payment.Status,
		Barber_id:           payment.Metadata.Barber_id,
		Schedule_day_date:   payment.Metadata.Schedule_day_date,
		Created_by_id:       payment.Metadata.Created_by_id,
		Schedule_start_time: payment.Metadata.Schedule_start_time,
		Shift_id:            payment.Metadata.Shift_id,
	})

	if err != nil {
		http.Error(rw, "no se creo la orden", http.StatusBadRequest)
		return
	}

	// Convertir `newOrder` a bytes y enviar al cliente
	msgBytes, err := json.Marshal(newOrder)
	if err != nil {
		log.Println("Error al serializar la orden:", err.Error())
		http.Error(rw, "Error interno al procesar la orden", http.StatusInternalServerError)
		return
	}

	// Enviar el mensaje al cliente específico
	err = sendMessageToPeer(websocket.TextMessage, msgBytes)
	if err != nil {
		log.Println("Error al enviar mensaje al cliente:", err.Error())
		return
	}

	// Respuesta exitosa
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
