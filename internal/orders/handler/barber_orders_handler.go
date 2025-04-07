package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

func (orh *OrderHandler) GetBarberPendingOrdersHandler(rw http.ResponseWriter, r *http.Request) {

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
		return
	}

	if !token.Is_barber {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	lmt := chi.URLParam(r, "limit")
	off := chi.URLParam(r, "offset")

	limit, _ := strconv.Atoi(lmt)
	offset, _ := strconv.Atoi(off)
	// si todo bien, se solicitan las ordenes
	orders, err := orh.ord_srv.GetOrderService(orh.ctx, int(token.ID), limit, offset)

	if err != nil {
		http.Error(rw, "Error al obtener las ordenes", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(orders)
}
