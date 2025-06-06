package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ezep02/rodeo/internal/services/services"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var (
	auth_token = viper.GetString("AUTH_TOKEN")
)

type Srvs_Handler struct {
	Srvs_Service *services.Srv_Service
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

func NewServiceHandler(srv_service *services.Srv_Service) *Srvs_Handler {
	return &Srvs_Handler{
		Srvs_Service: srv_service,
		Ctx:          context.Background(),
	}
}

// Servicios para mostrar al cliente
func (h *Srvs_Handler) GetServices(w http.ResponseWriter, r *http.Request) {
	// Ruta esperada: /services/{limit}/{offset}
	path := strings.TrimPrefix(r.URL.Path, "/services/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.Error(w, "Missing limit or offset", http.StatusBadRequest)
		return
	}

	limit := parts[0]
	offset := parts[1]
	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Error parseando parametro", http.StatusBadRequest)
		return
	}

	parsedOffset, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Error parseando parametro", http.StatusBadRequest)
		return
	}

	services, err := h.Srvs_Service.GetServices(h.Ctx, parsedLimit, parsedOffset)

	if err != nil {
		http.Error(w, "Algo salio mal al intentar obtener los servicios", http.StatusBadRequest)
		return
	}

	// si todo bien, devolves el servicio creado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}

// GetPopularServices obtiene los servicios populares
func (h *Srvs_Handler) GetPopularServices(rw http.ResponseWriter, r *http.Request) {
	services, err := h.Srvs_Service.GetPopularServices(h.Ctx)
	if err != nil {
		http.Error(rw, "Error al obtener los servicios populares", http.StatusBadRequest)
		return
	}
	// si todo bien, devolves el servicio creado

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(services)

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
