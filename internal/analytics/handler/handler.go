package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/ezep02/rodeo/internal/analytics/models"
	"github.com/ezep02/rodeo/internal/analytics/services"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

type AnalyticsHandler struct {
	An_Srv *services.AnalyticsServices
	Ctx    context.Context
}

var (
	auth_token = "auth_token"
)

func NewAnalyticsHandler(an_srv *services.AnalyticsServices) *AnalyticsHandler {
	return &AnalyticsHandler{
		An_Srv: an_srv,
		Ctx:    context.Background(),
	}
}

func (an_handler *AnalyticsHandler) GetTotalUsers(rw http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Is_admin {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	totalUsers, err := an_handler.An_Srv.GetTotalRegisteredUsers(an_handler.Ctx)

	if err != nil {
		http.Error(rw, "Ocurrio un error al intentar obtener los usuarios", http.StatusUnauthorized)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(totalUsers)
}

func (an_handler *AnalyticsHandler) GetRevicedTotalUsers(rw http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Is_admin {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	totalUsers, err := an_handler.An_Srv.GetRecivedUsers(an_handler.Ctx)

	if err != nil {
		http.Error(rw, "Ocurrio un error al intentar obtener los usuarios", http.StatusUnauthorized)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(totalUsers)
}

func (an_handler *AnalyticsHandler) GetRevenues(rw http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Is_admin {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	revenue, err := an_handler.An_Srv.GetTotalRevenue(an_handler.Ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(revenue)
}

func (an_handler *AnalyticsHandler) NewExpense(rw http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Is_admin {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Couldn't parse request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var expenseReq *models.ExpenseRequest

	if err := json.Unmarshal(b, &expenseReq); err != nil {
		http.Error(rw, "Error al deserializar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	expenseCreated, err := an_handler.An_Srv.NewExpenseSrv(an_handler.Ctx, &models.Expenses{
		Created_by_name: token.Name,
		Admin_id:        int(token.ID),
		Title:           expenseReq.Title,
		Description:     expenseReq.Description,
		Amount:          expenseReq.Amount,
	})

	if err != nil {
		http.Error(rw, "Error agregando gasto al listado", http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(expenseCreated)
}

func (an_handler *AnalyticsHandler) GetExpensesList(rw http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Is_admin {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	// recuperar parametros

	limitParam := chi.URLParam(r, "limit")
	offsetParam := chi.URLParam(r, "offset")

	// parsear
	parsedLmt, err := strconv.Atoi(limitParam)
	if err != nil {
		http.Error(rw, "Error de parseo", http.StatusUnauthorized)
		return
	}

	parsedOffset, err := strconv.Atoi(offsetParam)
	if err != nil {
		http.Error(rw, "Error de parseo", http.StatusUnauthorized)
		return
	}

	expensesList, err := an_handler.An_Srv.GetExpensesHistorialSrv(an_handler.Ctx, parsedLmt, parsedOffset)
	if err != nil {
		http.Error(rw, "Algo salio mal itentando recuperar el listado de gastos", http.StatusUnauthorized)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(expensesList)
}

func (an_handler *AnalyticsHandler) GetTotalExpensesCount(rw http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(auth_token)

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value
	token, err := jwt.VerfiyToken(tokenString)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Is_admin {
		http.Error(rw, "Usuario no autorizado", http.StatusUnauthorized)
		return
	}

	totalexp, err := an_handler.An_Srv.GetTotalExpenses(an_handler.Ctx)

	if err != nil {
		http.Error(rw, "Ocurrio un error al intentar obtener los usuarios", http.StatusUnauthorized)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(totalexp)
}
