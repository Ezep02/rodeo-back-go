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
	"time"

	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
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

var (
	mp_access_token = viper.GetString("MP_ACCESS_TOKEN")
	auth_token      = viper.GetString("AUTH_TOKEN")
)

func (orh *OrderHandler) CreateOrderHandler(rw http.ResponseWriter, r *http.Request) {

	var newOrder ServiceOrder

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(rw, "No se pudo parsear correctamente el cuerpo de la peticion", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Si el pago es aceptado, se crea la orden
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
				ID:          newOrder.Service_id,
				Title:       newOrder.Title,
				Quantity:    1,
				UnitPrice:   newOrder.Price,
				Description: newOrder.Description,
			},
		},
		Metadata: Metadata{
			Service_duration:    newOrder.Service_duration,
			UserID:              token.ID,
			Barber_id:           newOrder.Barber_id,
			Created_by_id:       newOrder.Created_by_id,
			Schedule_start_time: newOrder.Schedule_start_time,
			Schedule_day_date:   newOrder.Schedule_day_date,
			Shift_id:            newOrder.Shift_id,
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
		NotificationURL:    "https://3bfe-181-16-120-185.ngrok-free.app/order/webhook",
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
	rw.WriteHeader(http.StatusOK)
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
		log.Println("aqui ocurre", err.Error())
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
	json.NewEncoder(rw).Encode(newOrder)
}

func (orh *OrderHandler) Success(rw http.ResponseWriter, r *http.Request) {

	// Validar el token
	cookie, err := r.Cookie(auth_token)
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
	cookie, err := r.Cookie(auth_token)
	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
		return
	}

	if !token.Is_admin {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
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

// Refound maneja la solicitud de reembolso
func (orh *OrderHandler) Refound(rw http.ResponseWriter, r *http.Request) {

	parsedID := chi.URLParam(r, "id")

	// Construye la URL
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s/refunds", parsedID)

	// Crea la solicitud HTTP
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error al crear la solicitud:", err)
		return
	}
	fmt.Println("Access Token:", mp_access_token)

	// Configura las cabeceras
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mp_access_token))
	req.Header.Set("X-Idempotency-Key", "77e1c83b-7bb0-437b-bc50-a7a58e5660ac")

	// Envía la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al enviar la solicitud:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("res", resp)

	// Lee y procesa el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error al leer la respuesta: %v", err), http.StatusInternalServerError)
		return
	}

	// Procesa la respuesta
	var refund RefundResponse
	err = json.Unmarshal(body, &refund)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error al analizar la respuesta JSON: %v", err), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(refund)
}

// Obtener turno pendiente
func (orh *OrderHandler) GetPendingOrder(rw http.ResponseWriter, r *http.Request) {

	// Validar el token
	cookie, err := r.Cookie(auth_token)
	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
	}

	nextOrder, err := orh.ord_srv.GetOrderByUserID(orh.ctx, int(token.ID))
	if err != nil {
		log.Println("No fue posible recuperar la orden")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(nextOrder)
}

// Obtener historial de turnos
func (orh *OrderHandler) GetOrderHistorial(rw http.ResponseWriter, r *http.Request) {

	lim := chi.URLParam(r, "limit")
	off := chi.URLParam(r, "offset")

	parsedLim, err := strconv.Atoi(lim)
	if err != nil {
		http.Error(rw, "No se pudo parsear el limite", http.StatusUnauthorized)
		return
	}

	parsedOff, err := strconv.Atoi(off)
	if err != nil {
		http.Error(rw, "No se pudo parsear el offset", http.StatusUnauthorized)
		return
	}

	// Validar el token
	cookie, err := r.Cookie(auth_token)
	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
		return
	}

	historial, err := orh.ord_srv.GetOrdersHistorial(orh.ctx, int(token.ID), parsedLim, parsedOff)
	if err != nil {
		http.Error(rw, "Algo fallo al recuperar el historial", http.StatusUnauthorized)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(historial)
}

// WEBSOCKET
// Peer estructura para manejar una conexión peer-to-peer
type Peer struct {
	connection []*websocket.Conn
	mu         sync.Mutex
}

var peer Peer

// Configuración del upgrader de WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true }, // Permitir todas las conexiones (ajusta según sea necesario)
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleConnection gestiona conexión WebSocket
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
	peer.connection = append(peer.connection, ws)
	peer.mu.Unlock()

	log.Println("[ORDERS] Nueva conexión establecida")

	// Leer mensajes del cliente y reenviar directamente al peer
	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error leyendo mensaje:", err.Error())
			break
		}

		// Reenviar el mensaje al peer
		err = sendMessageToPeer(messageType, msg)
		if err != nil {
			log.Println("Error enviando datos actualizados:", err.Error())
			break
		}
	}

	// Al cerrar, eliminar la conexión activa
	removeConnection(ws)
	log.Println("[ORDERS] Conexión cerrada")
}

// removeConnection elimina una conexiones cerradas
func removeConnection(conn *websocket.Conn) {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	for i, c := range peer.connection {
		if c == conn {
			peer.connection = append(peer.connection[:i], peer.connection[i+1:]...)
			break
		}
	}
}

// sendUpdatedData envía mediante ws
func sendMessageToPeer(messageType int, msg []byte) error {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	if len(peer.connection) == 0 {
		log.Println("No hay peers conectados para recibir el mensaje schedules")
		return nil
	}

	var activeConnections []*websocket.Conn
	for _, conn := range peer.connection {
		err := conn.WriteMessage(messageType, msg)
		if err != nil {
			log.Println("Error al enviar mensaje al peer:", err.Error())
			conn.Close() // Cerrar la conexión fallida
			continue     // Omitir esta conexión en la lista activa
		}
		activeConnections = append(activeConnections, conn)
	}

	// Actualizar la lista de conexiones activas
	peer.connection = activeConnections
	return nil
}
