package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	auth_token = viper.GetString("AUTH_TOKEN")
)

type Srvs_Handler struct {
	Srvs_Service *Srv_Service
	Ctx          context.Context
}

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

func NewServiceHandler(srv_service *Srv_Service) *Srvs_Handler {
	return &Srvs_Handler{
		Srvs_Service: srv_service,
		Ctx:          context.Background(),
	}
}

// create service handler
func (h *Srvs_Handler) CreateService(rw http.ResponseWriter, r *http.Request) {

	var srv ServiceRequest

	if err := json.NewDecoder(r.Body).Decode(&srv); err != nil {
		http.Error(rw, "No se pudo parsear correctamente el cuerpo de la peticion", http.StatusBadRequest)
		log.Printf("[Error] %s", err.Error())
		return
	}

	defer r.Body.Close()

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		log.Printf("[TOKEN] no se pudo verificar el token, %s", err.Error())
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	newSrvReq := Service{
		Title:            srv.Title,
		Created_by_id:    token.ID,
		Description:      srv.Description,
		Price:            srv.Price,
		Service_Duration: srv.Service_Duration,
		Preview_url:      srv.Preview_url,
	}

	newSrv, err := h.Srvs_Service.CreateService(h.Ctx, &newSrvReq)

	if err != nil {
		log.Printf("[Create Req] No se pudo crear el servicio %s", err.Error())
		http.Error(rw, "Error al crear el servicio", http.StatusBadRequest)
		return
	}

	msg, err := json.Marshal(newSrv)
	if err != nil {
		log.Println("Error al parsear la informacion")
		http.Error(rw, "Error parseando el mensaje", http.StatusExpectationFailed)
		return
	}

	// Enviar el mensaje al cliente específico
	err = sendUpdatedData(websocket.TextMessage, msg)
	if err != nil {
		log.Println("Error al enviar mensaje al cliente:", err.Error())
		http.Error(rw, "Error interno al procesar la orden", http.StatusInternalServerError)
		return
	}

	// 3. notification push en el client view

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("Servicio creado correctamente")
}

// Servicios para mostrar al cliente
func (h *Srvs_Handler) GetServices(rw http.ResponseWriter, r *http.Request) {

	limit := chi.URLParam(r, "limit")
	offset := chi.URLParam(r, "offset")

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(rw, "Error parseando parametro", http.StatusBadRequest)
		return
	}

	parsedOffset, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(rw, "Error parseando parametro", http.StatusBadRequest)
		return
	}

	services, err := h.Srvs_Service.GetServices(h.Ctx, parsedLimit, parsedOffset)

	if err != nil {
		http.Error(rw, "Algo salio mal al intentar obtener los servicios", http.StatusBadRequest)
		return
	}

	// si todo bien, devolves el servicio creado
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(services)
}

func (h *Srvs_Handler) GetBarberServices(rw http.ResponseWriter, r *http.Request) {
	limit := chi.URLParam(r, "limit")
	offset := chi.URLParam(r, "offset")
	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(rw, "Error parseando parametro", http.StatusBadRequest)
		return
	}

	parsedOffset, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(rw, "Error parseando parametro", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Is_barber {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}
	service, err := h.Srvs_Service.GetBarberServices(h.Ctx, parsedLimit, parsedOffset, int(token.ID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	// si todo bien, devolves el servicio creado
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(service)
}

func (h *Srvs_Handler) UpdateServices(rw http.ResponseWriter, r *http.Request) {
	var srv Service

	if err := json.NewDecoder(r.Body).Decode(&srv); err != nil {
		http.Error(rw, "No se pudo parsear correctamente el cuerpo de la peticion", http.StatusBadRequest)
		log.Printf("[Error] %s", err.Error())
		return
	}

	defer r.Body.Close()

	cookie, err := r.Cookie("auth_token")

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, "No se pudo verificar el token", http.StatusUnauthorized)
		return
	}

	if !token.Is_barber {
		log.Println(token.Is_barber)
		http.Error(rw, "Solamente un barbero puede actualizar esta informacion", http.StatusUnauthorized)
		return
	}

	srv_id := chi.URLParam(r, "id")

	values := Service{
		Model: gorm.Model{
			ID: srv.ID,
		},
		Title:            srv.Title,
		Description:      srv.Description,
		Price:            srv.Price,
		Created_by_id:    srv.Created_by_id,
		Service_Duration: srv.Service_Duration,
	}

	updatedService, err := h.Srvs_Service.UpdateService(h.Ctx, &values, srv_id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	// log.Println("updated service", updatedService)

	msg, err := json.Marshal(updatedService)
	if err != nil {
		log.Println("Error al parsear la informacion")
		http.Error(rw, "Error parseando el mensaje", http.StatusExpectationFailed)
		return
	}

	// Enviar el mensaje al cliente específico
	err = sendUpdatedData(websocket.TextMessage, msg)
	if err != nil {
		log.Println("Error al enviar mensaje al cliente:", err.Error())
		http.Error(rw, "Error interno al procesar la orden", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("Servicio correctamente actualizado")
}

// Delete service by ID
func (h *Srvs_Handler) DeleteServiceByID(rw http.ResponseWriter, r *http.Request) {

	srv_id := chi.URLParam(r, "id")

	cookie, err := r.Cookie("auth_token")

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		log.Printf("[TOKEN] no se pudo verificar el token, %s", err.Error())
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Is_barber {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	parsedID, err := strconv.Atoi(srv_id)

	if err != nil {
		http.Error(rw, "Error parseando el service id", http.StatusConflict)
		return
	}

	if err := h.Srvs_Service.DeleteServiceByID(h.Ctx, parsedID); err != nil {
		http.Error(rw, "No se pudo completar la eliminacion, vuelva a intentarlo", http.StatusExpectationFailed)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("Eliminado correctamente")

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

	log.Println("[UPDATE SERVICES] Nueva conexion establecida ")

	// Leer mensajes del cliente y reenviar directamente al peer
	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}

		// Reenviar el mensaje al peer
		err = sendUpdatedData(messageType, msg)
		if err != nil {
			break
		}
	}

	// Al cerrar, eliminar la conexión activa
	peer.mu.Lock()
	peer.connection = nil
	peer.mu.Unlock()
	log.Println("[UPDATE SERVICES] Conexión cerrada")
}

// sendMessageToPeer envia los datos de actualizacion del viewer
func sendUpdatedData(messageType int, msg []byte) error {

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
