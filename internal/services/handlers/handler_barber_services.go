package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ezep02/rodeo/internal/services/models"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func (h *Srvs_Handler) CreateService(rw http.ResponseWriter, r *http.Request) {

	var srv models.ServiceRequest

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

	newSrvReq := models.Service{
		Title:            srv.Title,
		Created_by_id:    token.ID,
		Description:      srv.Description,
		Price:            srv.Price,
		Service_Duration: srv.Service_Duration,
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

func (h *Srvs_Handler) GetBarberServices(w http.ResponseWriter, r *http.Request) {
	// Ruta esperada: /services/barber/{limit}/{offset}
	path := strings.TrimPrefix(r.URL.Path, "/services/barber/")
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

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Is_barber {
		http.Error(w, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}
	service, err := h.Srvs_Service.GetBarberServices(h.Ctx, parsedLimit, parsedOffset, int(token.ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// si todo bien, devolves el servicio creado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(service)
}

func (h *Srvs_Handler) UpdateServices(w http.ResponseWriter, r *http.Request) {
	var srv models.Service

	if err := json.NewDecoder(r.Body).Decode(&srv); err != nil {
		http.Error(w, "No se pudo parsear correctamente el cuerpo de la peticion", http.StatusBadRequest)
		log.Printf("[Error] %s", err.Error())
		return
	}

	defer r.Body.Close()

	cookie, err := r.Cookie("auth_token")

	if err != nil {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(w, "No se pudo verificar el token", http.StatusUnauthorized)
		return
	}

	if !token.Is_barber {
		log.Println(token.Is_barber)
		http.Error(w, "Solamente un barbero puede actualizar esta informacion", http.StatusUnauthorized)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/services/update/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 {
		http.Error(w, "Missing limit or offset", http.StatusBadRequest)
		return
	}

	srv_id := parts[0]

	values := models.Service{
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// log.Println("updated service", updatedService)

	msg, err := json.Marshal(updatedService)
	if err != nil {
		log.Println("Error al parsear la informacion")
		http.Error(w, "Error parseando el mensaje", http.StatusExpectationFailed)
		return
	}

	// Enviar el mensaje al cliente específico
	err = sendUpdatedData(websocket.TextMessage, msg)
	if err != nil {
		log.Println("Error al enviar mensaje al cliente:", err.Error())
		http.Error(w, "Error interno al procesar la orden", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Servicio correctamente actualizado")
}

// Delete service by ID
func (h *Srvs_Handler) DeleteServiceByID(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("auth_token")

	if err != nil {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		log.Printf("[TOKEN] no se pudo verificar el token, %s", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Is_barber {
		http.Error(w, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/services/delete/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}
	srv_id := parts[0]

	parsedLimit, err := strconv.Atoi(srv_id)

	if err != nil {
		http.Error(w, "Error parseando el service id", http.StatusConflict)
		return
	}

	if err := h.Srvs_Service.DeleteServiceByID(h.Ctx, parsedLimit); err != nil {
		http.Error(w, "No se pudo completar la eliminacion, vuelva a intentarlo", http.StatusExpectationFailed)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Eliminado correctamente")
}
