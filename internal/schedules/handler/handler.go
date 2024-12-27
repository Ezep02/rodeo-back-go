package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ezep02/rodeo/internal/schedules/models"
	"github.com/ezep02/rodeo/internal/schedules/services"
	"github.com/go-chi/chi/v5"

	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ScheduleHandler struct {
	Sch_serv *services.ScheduleService
	Ctx      context.Context
}

// WEBSOCKET
type Peer struct {
	connection *websocket.Conn // Conexión WebSocket activa
	mu         sync.Mutex
}

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

func (sch_h *ScheduleHandler) CreateNewSchedule(rw http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "No se puede procesar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var scheduleReq []models.ShiftRequest
	if err := json.Unmarshal(b, &scheduleReq); err != nil {
		log.Printf("Error deserializando JSON: %v", err)
		http.Error(rw, "Error al deserializar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("auth_token")

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

	if !token.Is_barber {
		http.Error(rw, "Barbero no autorizado", http.StatusUnauthorized)
		return
	}

	scheduleList, err := sch_h.Sch_serv.CreateNewSchedulService(sch_h.Ctx, &scheduleReq, int(token.ID), token.Name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	msgBytes, err := json.Marshal(scheduleList)
	if err != nil {
		http.Error(rw, "Error interno al procesar la orden", http.StatusInternalServerError)
		return
	}

	err = sendUpdatedData(websocket.TextMessage, msgBytes)

	if err != nil {
		log.Println("Error al enviar mensaje al cliente:", err.Error())
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(scheduleList)
}

func (sch_h *ScheduleHandler) GetSchedules(rw http.ResponseWriter, r *http.Request) {

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

	if !token.Is_barber {
		http.Error(rw, "Barbero no autorizado", http.StatusUnauthorized)
		return
	}

	schedulesList, err := sch_h.Sch_serv.GetSchedules(sch_h.Ctx, int(token.ID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(schedulesList)
}

func (sch *ScheduleHandler) UpdateSchedules(rw http.ResponseWriter, r *http.Request) {

	// Leer el cuerpo de la solicitud
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "No se puede procesar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var updateSchedule []models.ScheduleResponse

	if err := json.Unmarshal(b, &updateSchedule); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

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

	if !token.Is_barber {
		http.Error(rw, "Barbero no autorizado", http.StatusUnauthorized)
		return
	}

	var newShift []models.Shift
	var updateShift []models.Shift
	var updateScheduleList []models.Schedule

	// CREACION DE NUEVOS SHIFTS
	for _, shift_add := range updateSchedule {
		// Check si cambio el schedule
		if shift_add.ScheduleStatus == "UPDATE" {

			scheduleUpdateRequest := models.Schedule{
				Model: &gorm.Model{
					ID:        uint(shift_add.ID),
					UpdatedAt: time.Now(),
				},
				Barber_id:    int(token.ID),
				Schedule_day: shift_add.Day,
				Start_date:   shift_add.Start,
				End_date:     shift_add.End,
			}
			updateScheduleList = append(updateScheduleList, scheduleUpdateRequest)
		}

		if len(shift_add.ShiftAdd) > 0 {
			for _, add := range shift_add.ShiftAdd {

				switch add.ShiftStatus {
				case "UPDATE":

					shiftUpdateRequest := models.Shift{
						Model: &gorm.Model{
							ID:        uint(add.ID),
							UpdatedAt: time.Now(),
						},
						Day:             shift_add.Day,
						Schedule_id:     uint(shift_add.ID),
						Start_time:      add.Start_time,
						Available:       add.Available,
						Created_by_name: add.Created_by_name,
					}
					updateShift = append(updateShift, shiftUpdateRequest)
				case "NEW":

					// Crear un objeto para nuevos shifts
					shiftRequest := models.Shift{
						Day:             shift_add.Day,
						Schedule_id:     uint(shift_add.ID),
						Start_time:      add.Start_time,
						Created_by_name: token.Name,
						Available:       true,
					}
					newShift = append(newShift, shiftRequest)
				case "NOT CHANGE":
					continue
				default:
					// Manejo opcional para valores desconocidos en ShiftStatus
					log.Printf("ShiftStatus desconocido: %s", add.ShiftStatus)
				}
			}
		}
	}
	newShiftChan := make(chan error)

	go func() {
		if len(newShift) > 0 {
			// Crear nuevos shifts
			newShiftRes, err := sch.Sch_serv.CreateNewShift(sch.Ctx, &newShift)
			if err != nil {
				http.Error(rw, "error creando los nuevos turnos", http.StatusBadRequest)
				return
			}

			// Comprobamos que newShiftRes no sea nil y que tenga elementos
			if newShiftRes != nil && len(*newShiftRes) > 0 {
				counter := 0

				// Recorremos el schedule para actualizar los shifts
				for _, sch := range updateSchedule {
					for indx, shf := range sch.ShiftAdd {
						if shf.ShiftStatus == "NEW" {
							// Obtener el shift de la respuesta
							shift := (*newShiftRes)[counter]

							// Mapeamos el modelo de GORM Shift a la estructura esperada en ShiftAdd
							sch.ShiftAdd[indx] = models.Shift{
								Model:           shift.Model,
								Schedule_id:     uint(sch.ID),
								Day:             sch.Day,
								Start_time:      shift.Start_time,
								Available:       shift.Available,
								Created_by_name: shift.Created_by_name,
							}
							// Incrementar el contador
							counter++
						}
					}
				}

			} else {
				log.Println("No se crearon nuevos shifts o la respuesta fue nil")
			}

			log.Println("newShiftRes:", newShiftRes)
			// DEVOLVER AL FRONT
			msgBytes, err := json.Marshal(updateSchedule)
			if err != nil {
				log.Println("Error al serializar la orden:", err.Error())
				http.Error(rw, "Error interno al procesar la orden", http.StatusInternalServerError)
				return
			}

			err = sendUpdatedData(websocket.TextMessage, msgBytes)
			if err != nil {
				log.Println("Error al enviar mensaje al cliente:", err.Error())
			}

			newShiftChan <- err
		} else {
			newShiftChan <- nil
		}
	}()

	// ELIMINACIONES DE SHIFTS
	del_ID := []int{}

	for _, shift_del := range updateSchedule {

		if len(shift_del.ShiftsDelete) > 0 {

			for _, del := range shift_del.ShiftsDelete {
				del_ID = append(del_ID, *del.ID)
			}
		}
	}
	deleteShiftChan := make(chan error)

	go func() {
		if len(del_ID) > 0 {
			err := sch.Sch_serv.DeleteShifts(sch.Ctx, del_ID)
			deleteShiftChan <- err
		} else {
			deleteShiftChan <- nil
		}
	}()

	// ACTUALIZACIONES DE SHIFTS

	updateShiftChan := make(chan error)

	go func() {
		if len(updateShift) > 0 {

			updatedShiftList, err := sch.Sch_serv.UpdateShiftList(sch.Ctx, &updateShift)
			log.Println("UPDATED:", updatedShiftList)

			updateShiftChan <- err
		} else {
			updateShiftChan <- nil
		}
	}()

	// ACTUALIZACIONES DE SCHEDULES
	updateScheduleChan := make(chan error)
	go func() {

		if len(updateScheduleList) > 0 {

			updateSchedule, err := sch.Sch_serv.UpdateScheduleList(sch.Ctx, int(token.ID), &updateScheduleList)
			log.Println("UPDATED SCHEDULES:", updateSchedule)

			updateScheduleChan <- err
		} else {
			updateShiftChan <- nil
		}
	}()

	if len(newShift) == 0 {
		msgBytes, err := json.Marshal(updateSchedule)
		if err != nil {
			log.Println("Error al serializar la orden:", err.Error())
			http.Error(rw, "Error interno al procesar la orden", http.StatusInternalServerError)
			return
		}

		err = sendUpdatedData(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Println("Error al enviar mensaje al cliente:", err.Error())
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(newShift)
}

func (sch *ScheduleHandler) GetBarberSchedules(rw http.ResponseWriter, r *http.Request) {

	barberID := chi.URLParam(r, "id")

	parsedID, err := strconv.Atoi(barberID)

	if err != nil {
		http.Error(rw, "No se encontro un barbero con ese id", http.StatusInternalServerError)
		return
	}

	barberSchedulesList, err := sch.Sch_serv.GetSchedules(sch.Ctx, parsedID)

	if err != nil {
		http.Error(rw, "Error al obtener los horarios", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(barberSchedulesList)
}

func (sch *ScheduleHandler) UpdateShiftStatus(rw http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	parsedID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(rw, "Error parseando el parametro id", http.StatusBadRequest)
		return
	}

	updatedShift, err := sch.Sch_serv.UpdateShiftByID(sch.Ctx, parsedID)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(updatedShift)

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

	log.Println("[UPDATE SCHEDULES] Nueva conexion establecida ")

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
	log.Println("[UPDATE SCHEDULES] Conexión cerrada")
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
