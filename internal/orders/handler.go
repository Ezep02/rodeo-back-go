package orders

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ezep02/rodeo/pkg/jwt"
)

// var link string = "https://api.mercadopago.com"

type OrderHandler struct {
	ctx context.Context
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		ctx: context.Background(),
	}
}

func (orh *OrderHandler) CreateOrderHandler(rw http.ResponseWriter, r *http.Request) {

	var newOrder Order

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(rw, "No se pudo parsear correctamente el cuerpo de la peticion", http.StatusBadRequest)
		log.Printf("[Error] %s", err.Error())
		return
	}

	defer r.Body.Close()

	// Si el pago es aceptado, se crea la orden
	accessToken := "APP_USR-196506190136225-092022-41af146cb6426644ccd360b92edc7ef6-1432087693"
	cookie, err := r.Cookie("auth_token")

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)
	if err != nil {
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
	}

	request := Request{
		AutoReturn: "approved",
		BackURLs: BackURLs{
			Success: "http://localhost:5173/dashboard",
			Pending: "http://localhost:8080/payment/pending",
			Failure: "http://localhost:8080/payment/failure",
		},
		StatementDescriptor: "TestStore",
		BinaryMode:          false,
		ExternalReference:   "IWD1238971",
		Items: []Item{
			{
				ID:          newOrder.ID,
				Title:       newOrder.Title,
				Quantity:    1,
				UnitPrice:   int(newOrder.Price),
				Description: newOrder.Description,
			},
		},
		Metadata: Metadata{
			Service_duration: newOrder.Service_Duration,
			UserID:           token.ID,
		},
		Payer: Payer{
			Email:   token.Email,
			Name:    token.Name,
			Surname: token.Surname,
			Phone: Phone{
				Number: token.Phone_number,
			},
		},
		PaymentMethods: PaymentMethods{
			ExcludedPaymentTypes:   []string{},
			ExcludedPaymentMethods: []string{},
			Installments:           12,
			DefaultPaymentMethodID: "account_money",
		},
		NotificationURL:    "https://67c0-181-16-123-233.ngrok-free.app/order/webhook",
		Expires:            true,
		ExpirationDateFrom: "2024-01-01T12:00:00.000-04:00",
		ExpirationDateTo:   "2024-12-31T12:00:00.000-04:00",
	}

	// Serializa el objeto request a JSON
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
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
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con el cuerpo de la respuesta
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(responseBody)
}

// WebHook maneja las solicitudes de webhook de Mercado Pago.
func (orh *OrderHandler) WebHook(rw http.ResponseWriter, r *http.Request) {

	// Leer el cuerpo de la solicitud
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error leyendo el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo JSON en un mapa de interfaces
	var bodyData map[string]interface{}
	if err := json.Unmarshal(body, &bodyData); err != nil {
		http.Error(rw, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Acceder al campo "data.id"
	data, ok := bodyData["data"].(map[string]interface{})
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

	// Token de acceso para la API de Mercado Pago
	accessToken := "APP_USR-196506190136225-092022-41af146cb6426644ccd360b92edc7ef6-1432087693"

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
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta para verificar su contenido
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reiniciar el cuerpo para decodificarlo después
	resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	// Leer el cuerpo de la respuesta
	var preferenceData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&preferenceData); err != nil {
		http.Error(rw, "Error decoding response body", http.StatusInternalServerError)
		return
	}

	log.Printf("[PAYMENT BODY] %+v", responseBody)

	// si todo fue bien se crea la orden
	var paymentResponse Payment
	if err := json.Unmarshal(responseBody, &paymentResponse); err != nil {
		http.Error(rw, "Error decoding response body", http.StatusInternalServerError)
		return
	}

	log.Printf("[PAYMENT BODY] %+v", paymentResponse)

	// Respuesta exitosa
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]string{"message": "webhook received"})
}

func (orh *OrderHandler) Success(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Success")
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
