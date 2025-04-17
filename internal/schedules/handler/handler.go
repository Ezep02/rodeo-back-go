package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/ezep02/rodeo/internal/schedules/services"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"

	"github.com/gorilla/websocket"
)

type ScheduleHandler struct {
	Sch_serv *services.ScheduleService
	Ctx      context.Context
}

// WEBSOCKET
type Peer struct {
	connection []*websocket.Conn
	mu         sync.Mutex
}

// variables globales
var (
	auth_token = viper.GetString("AUTH_TOKEN")
)

// Crear una instancia global del peer
var peer Peer

// Configuracion del upgrader de WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true }, // Permitir todas las conexiones
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewSchedulHandler(sch_srv *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		Sch_serv: sch_srv,
		Ctx:      context.Background(),
	}
}

func (sch_h *ScheduleHandler) GetAvailableSchedulesHandler(rw http.ResponseWriter, r *http.Request) {

	limit := chi.URLParam(r, "limit")
	offset := chi.URLParam(r, "offset")

	parsedLimit, err := strconv.Atoi(limit)

	if err != nil {
		log.Println("Parsing error")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	parsedOffset, err := strconv.Atoi(offset)

	if err != nil {
		log.Println("Parsing error")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	availableSchedules, err := sch_h.Sch_serv.GetAvailableSchedules(sch_h.Ctx, parsedLimit, parsedOffset)

	if err != nil {
		log.Println("error searching schedules")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(availableSchedules)
}

// HandleConnection gestiona una conexión WebSocket
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

	log.Println("[UPDATE SCHEDULES] Nueva conexion establecida")

	// Leer mensajes del cliente y reenviar directamente al peer
	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error leyendo mensaje:", err.Error())
			break
		}

		// Reenviar el mensaje al peer
		err = sendUpdatedData(messageType, msg)
		if err != nil {
			log.Println("Error enviando datos actualizados:", err.Error())
			break
		}
	}

	// Al cerrar, eliminar la conexion activa
	removeConnection(ws)
	log.Println("[UPDATE SCHEDULES] Conexión cerrada")
}

// removeConnection elimina una conexión específica del peer
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

// sendUpdatedData envia datos de actualización a todos los peers conectados
func sendUpdatedData(messageType int, msg []byte) error {
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
