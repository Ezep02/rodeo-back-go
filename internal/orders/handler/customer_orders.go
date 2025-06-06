package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ezep02/rodeo/internal/orders/helpers"
	"github.com/ezep02/rodeo/internal/orders/models"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gorilla/websocket"
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
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(pendingOrders)
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

// Refaund
func (h *OrderHandler) CreateRefundHandler(w http.ResponseWriter, r *http.Request) {

	var (
		requestBody *models.RefundRequest
	)

	// obtener los datos del refaund
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Error recibiendo los datos", http.StatusBadRequest)
		return
	}

	// validar session
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

	// validar que la orden ya este cancelada
	if isCanceled, err := h.ord_srv.CheckOrderStatus(h.ctx, requestBody.Order_id); isCanceled || err != nil {
		http.Error(w, "Algo salio mal al intentar cancelar el turno", http.StatusBadRequest)
		return
	}

	// cacelar la orden y liberar turno, devuelve el turno a liberar
	available_schedule, err := h.ord_srv.NewRefound(h.ctx, *requestBody)
	if err != nil {
		http.Error(w, "Algo no fue bien intentando cancelar el turno", http.StatusBadRequest)
		return
	}

	parsed_coupon, err := helpers.CouponFormater(*requestBody, int(token.ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		// GENERAR CUPON, DEVOLVER MEDIANTE WS

		available_coupon, err := h.ord_srv.GenerateCoupon(h.ctx, parsed_coupon)
		if err != nil {
			http.Error(w, "Algo no fue bien intentando cancelar el turno", http.StatusBadRequest)
			return
		}

		// enviar contenido ws
		msg, err := json.Marshal(available_coupon)
		if err != nil {
			http.Error(w, "Error preparando entrega de datos", http.StatusBadRequest)
			return
		}

		if err := sendMessageToPeer(websocket.TextMessage, msg); err != nil {
			http.Error(w, "Error durante la entrega de datos", http.StatusBadRequest)
			return
		}
	}()

	// Comunicar al dashboard
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(available_schedule)
}

// Reschedule
func (h *OrderHandler) CreateReschedule(rw http.ResponseWriter, r *http.Request) {

	var (
		schedule       models.RescheduleRequest
		validatedToken *jwt.VerifyTokenRes
	)

	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(rw, "No se pudo parsear correctamente el cuerpo de la petición", http.StatusBadRequest)
		return
	}

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

	// actualizar los datos de la orden
	updated_pending_order, err := h.ord_srv.UpdateScheduleOrder(h.ctx, schedule, int(validatedToken.ID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// enviar contenido ws
	msg, err := json.Marshal(updated_pending_order)
	if err != nil {
		http.Error(rw, "Error preparando entrega de datos", http.StatusBadRequest)
		return
	}

	if err := sendMessageToPeer(websocket.TextMessage, msg); err != nil {
		http.Error(rw, "Error durante la entrega de datos", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(updated_pending_order)
}

// Coupons
func (h *OrderHandler) GetCouponsHandler(rw http.ResponseWriter, r *http.Request) {

	// validar session
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

	available_coupons, err := h.ord_srv.GetCustomerCoupons(h.ctx, int(token.ID))
	if err != nil {
		http.Error(rw, "Error obteniendo cupones", http.StatusUnauthorized)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(available_coupons)
}

// Obtener las ordenes previas, historial del cliente, ademas con su review si es que tiene
func (h *OrderHandler) GetCustomerPreviousOrdersHandler(w http.ResponseWriter, r *http.Request) {

	// validar session
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

	// Obtener limit y offset desde la url
	path := strings.TrimPrefix(r.URL.Path, "/order/customer/previous/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 {
		http.Error(w, "Missing limit or offset", http.StatusBadRequest)
		return
	}

	offset := parts[0]

	// parsing
	parsedOffset, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Error parseando dato", http.StatusConflict)
		return
	}

	previus_orders, err := h.ord_srv.GetCustomerPreviusOrders(h.ctx, int(token.ID), parsedOffset)
	if err != nil {
		http.Error(w, "Error obteniendo las ordenes anteriores", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(previus_orders)

}
