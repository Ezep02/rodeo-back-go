package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ezep02/rodeo/internal/schedules/models"
	"github.com/ezep02/rodeo/internal/schedules/services"
	"github.com/ezep02/rodeo/utils"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gorilla/websocket"
)

type ScheduleHandler struct {
	Sch_serv *services.ScheduleService
	Ctx      context.Context
}

// WEBSOCKET
type Peer struct {
	connection []*websocket.Conn // Conexión WebSocket activa
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

func (sch_h *ScheduleHandler) GetBarberSchedulesHandler(rw http.ResponseWriter, r *http.Request) {

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

	if !token.Is_barber {
		http.Error(rw, "Barbero no autorizado", http.StatusUnauthorized)
		return
	}

	limit := chi.URLParam(r, "limit")

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(rw, "Error parseando dato", http.StatusConflict)
		return
	}

	offset := chi.URLParam(r, "offset")

	parsetOffset, err := strconv.Atoi(offset)

	if err != nil {
		http.Error(rw, "Error parseando dato", http.StatusConflict)
		return
	}

	schedulesList, err := sch_h.Sch_serv.GetBarberSchedules(sch_h.Ctx, int(token.ID), parsedLimit, parsetOffset)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(schedulesList)
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

func (sch *ScheduleHandler) BarberSchedulesHandler(rw http.ResponseWriter, r *http.Request) {

	var (
		schedulesToAdd   []models.Schedule
		schedulesRequest models.ScheduleRequest
	)

	// Leer el cuerpo de la solicitud
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "No se puede procesar el cuerpo de la solicitud", http.StatusNoContent)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(b, &schedulesRequest); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
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
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
		return
	}

	if !token.Is_barber {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	// iterar sobre Schedule_add

	for _, schedule := range schedulesRequest.Schedule_add {

		//
		switch strings.ToUpper(schedule.Schedule_status) {
		case "NEW":
			id, err := utils.GenerateRandomID()
			if err != nil {
				// Manejo del error
				fmt.Println("Error generando ID:", err)
				return
			}

			scheduleToCreate := models.Schedule{
				Model: &gorm.Model{
					ID: id,
				},
				Created_by_name:   token.Name,
				Barber_id:         int(token.ID),
				Available:         schedule.Available,
				Start_time:        schedule.Start_time,
				Schedule_day_date: schedule.Schedule_day_date,
			}
			schedulesToAdd = append(schedulesToAdd, scheduleToCreate)
		case "UPDATE":
			log.Println("Schedule", schedule)
		default:
			log.Println("no es new")
		}
	}

	// Create Schedules
	if len(schedulesToAdd) > 0 {
		log.Println("Start creating process ")
		go func() {
			log.Println("elements", schedulesToAdd)
			_, err := sch.Sch_serv.CreateBarberSchedules(sch.Ctx, &schedulesToAdd)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusConflict)
				log.Println("Schedules,", err.Error())
				return
			}
		}()

		// convertir a array de bytes y enviar paquete por websocket
		b, err := json.Marshal(schedulesToAdd)
		if err != nil {
			log.Println("Convertion error to byte")
			http.Error(rw, err.Error(), http.StatusConflict)
			return
		}

		if err := sendUpdatedData(websocket.TextMessage, b); err != nil {
			log.Println("Error al enviar mensaje al cliente:", err.Error())
			http.Error(rw, "Error interno al procesar la orden", http.StatusInternalServerError)
			return
		}
	}

	// Delete schedules
	if len(schedulesRequest.Schedule_delete) > 0 {
		log.Println("Starting Delete process")
		go func() {
			//
			scheduleIDs := make([]int, len(schedulesRequest.Schedule_delete))
			for i, schedule := range schedulesRequest.Schedule_delete {
				scheduleIDs[i] = schedule.ID
			}

			if err := sch.Sch_serv.DeleteSchedules(sch.Ctx, scheduleIDs); err != nil {
				http.Error(rw, err.Error(), http.StatusConflict)
				log.Println("Schedules,", err.Error())
				return
			}
		}()
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("Operacion exitosa")
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

	log.Println("[UPDATE SCHEDULES] Nueva conexión establecida")

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

	// Al cerrar, eliminar la conexión activa
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

// sendUpdatedData envía datos de actualización a todos los peers conectados
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
