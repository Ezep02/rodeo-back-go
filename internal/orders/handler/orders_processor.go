package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ezep02/rodeo/internal/orders/helpers"
	"github.com/ezep02/rodeo/internal/orders/models"
	"github.com/ezep02/rodeo/internal/orders/services"
	"github.com/ezep02/rodeo/internal/orders/utils"
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
		newOrder     models.ServiceOrder
		responseBody map[string]any
	)

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(rw, "No se pudo parsear correctamente el cuerpo de la peticion", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// crear token transitorio
	now := time.Now()
	orderToken, err := utils.GenerateOrderToken(models.PendingOrderToken{
		Title:               newOrder.Title,
		Payer_name:          newOrder.Payer_name,
		Payer_surname:       newOrder.Payer_surname,
		Barber_id:           newOrder.Barber_id,
		Schedule_day_date:   newOrder.Schedule_day_date,
		Schedule_start_time: newOrder.Schedule_start_time,
		User_id:             newOrder.User_id,
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
		Payer_name:          newOrder.Payer_name,
		Payer_surname:       newOrder.Payer_surname,
		Barber_id:           newOrder.Barber_id,
		Schedule_day_date:   newOrder.Schedule_day_date,
		Schedule_start_time: newOrder.Schedule_start_time,
		User_id:             newOrder.User_id,
		ID:                  uint(newOrder.Shift_id),
		Price:               float64(newOrder.Price),
		Created_at:          &now,
	}); err != nil {
		http.Error(rw, "Algo salió mal intentando setear el token", http.StatusBadRequest)
		return
	}

	request, err := helpers.BuildOrderPreference(newOrder, orderToken)
	if err != nil {
		http.Error(rw, "Algo salio mal intentando crear la preferencia", http.StatusBadRequest)
		return
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

func (orh *OrderHandler) WebHook(rw http.ResponseWriter, r *http.Request) {
	var bodyPayment map[string]any
	if err := json.NewDecoder(r.Body).Decode(&bodyPayment); err != nil {
		log.Println("Error decodificando el cuerpo del webhook:", err)
		http.Error(rw, "JSON inválido", http.StatusBadRequest)
		return
	}

	data, ok := bodyPayment["data"].(map[string]any)
	if !ok {
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

	order, err := helpers.BuildOrderFromWebhook(root)
	if err != nil {
		log.Println("Error formateando respuesta de MercadoPago:", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	newOrder, err := orh.ord_srv.CreateNewOrder(orh.ctx, order)

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

	// Send real time data using sse

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusAccepted)
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
