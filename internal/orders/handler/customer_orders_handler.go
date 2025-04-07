package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ezep02/rodeo/pkg/jwt"
)

func (h *OrderHandler) CustomerPendingOrderHandler(rw http.ResponseWriter, r *http.Request) {
	var (
		validatedToken *jwt.VerifyTokenRes
	)

	if cookie, err := r.Cookie(auth_token); err == nil {
		token, err := jwt.VerfiyToken(cookie.Value)
		if err != nil {
			http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
			return
		}
		validatedToken = token
	} else {
		http.Error(rw, "Error al verificar el token", http.StatusBadRequest)
		return
	}

	pendingOrders, err := h.ord_srv.GetCustomerPendingOrder(h.ctx, int(validatedToken.ID))
	log.Println("[200 ORDERS]", pendingOrders)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(pendingOrders)
}

// Obtener turno pendiente REFACTORIZAR
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

func (h *OrderHandler) GetSuccessPaymentHandler(rw http.ResponseWriter, r *http.Request) {
	// Definir estructura correcta
	var requestData struct {
		Token string `json:"token"`
	}

	// Decodificar JSON en la estructura
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(rw, "No se pudo parsear correctamente el cuerpo de la petición", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Obtener la orden por token
	order, err := h.ord_srv.GetOrderByToken(h.ctx, requestData.Token)
	if err != nil {
		http.Error(rw, "No se pudo obtener la orden", http.StatusBadRequest)
		return
	}

	// Responder con JSON
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(order)
}
