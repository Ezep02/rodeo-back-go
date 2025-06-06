package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ezep02/rodeo/internal/schedules/models"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gorilla/websocket"
)

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

		switch strings.ToUpper(schedule.Schedule_status) {
		case "NEW":
			scheduleToCreate := models.Schedule{
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
			continue
		}
	}

	// Create Schedules
	if len(schedulesToAdd) > 0 {
		log.Println("Start creating process ")
		go func() {
			log.Println("elements", schedulesToAdd)
			data, err := sch.Sch_serv.CreateBarberSchedules(sch.Ctx, &schedulesToAdd)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusConflict)
				log.Println("Schedules,", err.Error())
				return
			}
			// convertir a array de bytes y enviar paquete por websocket
			b, err := json.Marshal(data)
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
		}()
	}

	// Delete schedules
	if len(schedulesRequest.Schedule_delete) > 0 {
		log.Println("Starting Delete process")
		go func() {
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

	rw.Header().Add("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("OK")

}

func (sch_h *ScheduleHandler) GetBarberSchedulesHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)
	if err != nil {
		http.Error(w, "Error al verificar el token", http.StatusBadRequest)
		return
	}

	if !token.Is_barber {
		http.Error(w, "Barbero no autorizado", http.StatusUnauthorized)
		return
	}

	// Ruta esperada: /services/barber/{limit}/{offset}
	path := strings.TrimPrefix(r.URL.Path, "/schedules/barber/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.Error(w, "Missing limit or offset", http.StatusBadRequest)
		return
	}

	limit := parts[0]
	offset := parts[1]

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Error parseando dato", http.StatusConflict)
		return
	}

	parsetOffset, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Error parseando dato", http.StatusConflict)
		return
	}

	schedulesList, err := sch_h.Sch_serv.GetBarberSchedules(sch_h.Ctx, int(token.ID), parsedLimit, parsetOffset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedulesList)
}
