package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ezep02/rodeo/internal/orders/models"
	"github.com/ezep02/rodeo/pkg/jwt"
)

// Reviews
func (h *OrderHandler) SetReviewHandler(w http.ResponseWriter, r *http.Request) {

	var (
		request_body models.ReviewData
	)

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

	if err := json.NewDecoder(r.Body).Decode(&request_body); err != nil {
		log.Println("ERROR decodificando review data")
		http.Error(w, "Error decodificando review data", http.StatusBadRequest)
		return
	}

	if err = h.ord_srv.CreateNewReview(h.ctx, models.Review{
		ReviewData: request_body,
		User_id:    int(token.ID),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Reseña creada con éxito.")
}

func (h *OrderHandler) GetReviewsHandler(w http.ResponseWriter, r *http.Request) {

	// Obtener limit y offset desde la url
	path := strings.TrimPrefix(r.URL.Path, "/review/all/")
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

	reviews, err := h.ord_srv.GetReviews(h.ctx, parsedOffset)
	if err != nil {
		http.Error(w, "Error recuperando reviews", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reviews)
}
