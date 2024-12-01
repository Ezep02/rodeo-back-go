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
	"sync"

	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

// var link string = "https://api.mercadopago.com"

type OrderHandler struct {
	ord_srv *OrderService
	ctx     context.Context
}

func NewOrderHandler(orders_srv *OrderService) *OrderHandler {
	return &OrderHandler{
		ctx:     context.Background(),
		ord_srv: orders_srv,
	}
}

func (orh *OrderHandler) CreateOrderHandler(rw http.ResponseWriter, r *http.Request) {

	var newOrder ServiceOrder

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
			Success: "http://localhost:5173/payment/success",
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
		NotificationURL:    "https://2fd0-181-16-121-113.ngrok-free.app/order/webhook",
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

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		http.Error(rw, "Error leyendo el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	// Decodificar el cuerpo JSON en el objeto Payment
	var payment PaymentResponse
	if err := json.Unmarshal(b, &payment); err != nil {
		http.Error(rw, "Error decoding request body", http.StatusBadRequest)
		return
	}

	newOrder, err := orh.ord_srv.CreateNewOrder(orh.ctx, &Order{
		Title:            payment.AdditionalInfo.Items[0].Title,
		Price:            payment.AdditionalInfo.Items[0].UnitPrice,
		Service_Duration: payment.Metadata.Service_duration,
		User_id:          int(payment.Metadata.UserID),
		Service_id:       payment.AdditionalInfo.Items[0].ID,
		Payment_id:       payment.ID,
		Payer_name:       payment.AdditionalInfo.Payer.FirstName,
		Payer_surname:    payment.AdditionalInfo.Payer.LastName,
		Email:            payment.PayerInfo.Email,
		Mp_order_id:      int64(payment.ID),
		Date_approved:    payment.DateApproved,
		Mp_status:        payment.Status,
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
	}

	// Respuesta exitosa
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(newOrder)
}

func (orh *OrderHandler) Success(rw http.ResponseWriter, r *http.Request) {

	// Validar el token
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value

	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
	}

	order, err := orh.ord_srv.GetOrderByUserID(orh.ctx, int(token.ID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(order)
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

func (orh *OrderHandler) GetOrders(rw http.ResponseWriter, r *http.Request) {

	// Validar el token
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
	}

	if !token.Is_admin {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
	}

	lmt := chi.URLParam(r, "limit")
	off := chi.URLParam(r, "offset")

	limit, _ := strconv.Atoi(lmt)
	offset, _ := strconv.Atoi(off)
	// si todo bien, se solicitan las ordenes
	orders, err := orh.ord_srv.GetOrderService(orh.ctx, limit, offset)

	if err != nil {
		http.Error(rw, "Error al obtener las ordenes", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(orders)
}

// WEBSOCKET
// Peer estructura para manejar una conexión peer-to-peer
type Peer struct {
	connection *websocket.Conn // Conexión WebSocket activa
	mu         sync.Mutex      // Mutex para concurrencia en la conexión
}

// Crear una instancia global del peer
var peer Peer

// Configuración del upgrader de WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true }, // Permitir todas las conexiones (ajusta según sea necesario)
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleConnection gestiona una conexión WebSocket P2P
func HandleConnection(rw http.ResponseWriter, r *http.Request) {
	// Actualizar a WebSocket
	ws, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println("Error al actualizar la conexión:", err.Error())
		return
	}
	defer ws.Close()

	// Registrar la conexión como la conexión activa
	peer.mu.Lock()
	peer.connection = ws
	peer.mu.Unlock()

	log.Println("Nueva conexión P2P establecida")

	// Leer mensajes del cliente y reenviar directamente al peer
	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}

		// Reenviar el mensaje al peer
		err = sendMessageToPeer(messageType, msg)
		if err != nil {
			break
		}
	}

	// Al cerrar, eliminar la conexión activa
	peer.mu.Lock()
	peer.connection = nil
	peer.mu.Unlock()
	log.Println("Conexión P2P cerrada")
}

// sendMessageToPeer envía un mensaje al peer conectado
func sendMessageToPeer(messageType int, msg []byte) error {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	if peer.connection == nil {
		log.Println("No hay peer conectado para recibir el mensaje")
		return nil
	}

	// Enviar el mensaje
	err := peer.connection.WriteMessage(messageType, msg)
	if err != nil {
		log.Println("Error al enviar mensaje al peer:", err.Error())
		return err
	}

	return nil
}
