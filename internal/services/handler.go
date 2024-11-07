package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/ezep02/rodeo/pkg/jwt"
)

var (
	TOKEN = os.Getenv("SECRET_TOKEN")
)

type Srvs_Handler struct {
	Srvs_Service *Srv_Service
	Ctx          context.Context
}

func NewServiceHandler(srv_service *Srv_Service) *Srvs_Handler {
	return &Srvs_Handler{
		Srvs_Service: srv_service,
		Ctx:          context.Background(),
	}
}

// create service handler
func (h *Srvs_Handler) CreateService(rw http.ResponseWriter, r *http.Request) {

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
		log.Printf("[TOKEN] no se pudo verificar el token, %s", err.Error())
		return
	}

	// obtener el user ID
	srv.Created_by_id = token.ID

	newSrv, err := h.Srvs_Service.CreateService(h.Ctx, &srv)

	if err != nil {
		log.Printf("[Create Req] No se pudo crear el servicio %s", err.Error())
		return
	}

	// si todo bien, devolves el servicio creado
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(newSrv)
}

// get all services handler
func (h *Srvs_Handler) GetAllServices(rw http.ResponseWriter, r *http.Request) {

	services, err := h.Srvs_Service.GetServices(h.Ctx)

	if err != nil {
		http.Error(rw, "Algo salio mal al intentar obtener los servicios", http.StatusBadRequest)
		return
	}

	// si todo bien, devolves el servicio creado
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(services)
}
