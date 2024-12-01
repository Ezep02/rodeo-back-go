package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ezep02/rodeo/internal/schedules/models"
	"github.com/ezep02/rodeo/internal/schedules/services"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ScheduleHandler struct {
	Sch_serv *services.ScheduleService
	Ctx      context.Context
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

func NewSchedulHandler(sch_srv *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		Sch_serv: sch_srv,
		Ctx:      context.Background(),
	}
}

func (sch_h *ScheduleHandler) CreateNewSchedule(rw http.ResponseWriter, r *http.Request) {

	// Leer el cuerpo de la solicitud
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "No se puede procesar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Deserializar el JSON recibido
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
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)
	if err != nil {
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
	}

	scheduleList, err := sch_h.Sch_serv.Sch_repo.CreateNewSchedul(sch_h.Ctx, &scheduleReq, int(token.ID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
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

	schedulesList, err := sch_h.Sch_serv.Sch_repo.GetSchedules(sch_h.Ctx, int(token.ID))
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

	var updateSchedule []models.ScheduleModifyDay

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
				User_id:       int(token.ID),
				Schedule_type: shift_add.DistributionType,
				Schedule_day:  shift_add.Day,
				Start_date:    shift_add.Date.Start,
				End_date:      *shift_add.Date.End,
			}
			updateScheduleList = append(updateScheduleList, scheduleUpdateRequest)
		}

		if len(shift_add.ShiftAdd) > 0 {
			for _, add := range shift_add.ShiftAdd {

				switch add.ShiftStatus {
				case "UPDATE":
					if err != nil {
						log.Println("Error parsing CreatedAt:", err)
						return
					}

					shiftUpdateRequest := models.Shift{
						Model: &gorm.Model{
							ID:        uint(add.ID),
							UpdatedAt: time.Now(),
						},
						Day:         shift_add.Day,
						Schedule_id: uint(shift_add.ID),
						Start_time:  add.Start_Time,
					}
					updateShift = append(updateShift, shiftUpdateRequest)
				case "NEW":
					// Crear un objeto para nuevos shifts
					shiftRequest := models.Shift{
						Day:         shift_add.Day,
						Schedule_id: uint(shift_add.ID),
						Start_time:  add.Start_Time,
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
							sch.ShiftAdd[indx] = struct {
								CreatedAt   string `json:"CreatedAt,omitempty"`
								Day         string `json:"Day,omitempty"`
								DeletedAt   string `json:"DeletedAt,omitempty"`
								ID          int    `json:"ID"`
								ScheduleID  int    `json:"Schedule_id,omitempty"`
								Start_Time  string `json:"Start_time"`
								UpdatedAt   string `json:"UpdatedAt,omitempty"`
								ShiftStatus string `json:"Shift_status"`
							}{
								CreatedAt:   shift.CreatedAt.String(),
								Day:         shift.Day,
								DeletedAt:   shift.DeletedAt.Time.String(),
								ID:          int(shift.ID),
								ScheduleID:  int(shift.Schedule_id),
								Start_Time:  shift.Start_time,
								UpdatedAt:   shift.UpdatedAt.String(),
								ShiftStatus: "NOT CHANGE",
							}
							// Incrementar el contador
							counter++
						}
					}
				}

			} else {
				log.Println("No se crearon nuevos shifts o la respuesta fue nil")
			}

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
