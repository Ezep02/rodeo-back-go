package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ezep02/rodeo/pkg/jwt"
)

func (orh *OrderHandler) GetBarberPendingOrdersHandler(w http.ResponseWriter, r *http.Request) {

	// Validar el token
	cookie, err := r.Cookie(auth_token)
	if err != nil {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(w, "Error al verificar el token", http.StatusBadRequest)
		return
	}

	if !token.Is_barber {
		http.Error(w, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	// Obtener limit y offset desde la url
	path := strings.TrimPrefix(r.URL.Path, "/order/pending/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.Error(w, "Missing limit or offset", http.StatusBadRequest)
		return
	}

	limit := parts[0]
	offset := parts[1]

	// parsing
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

	// si todo bien, se solicitan las ordenes
	orders, err := orh.ord_srv.GetOrderService(orh.ctx, int(token.ID), parsedLimit, parsetOffset)

	if err != nil {
		http.Error(w, "Error al obtener las ordenes", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
