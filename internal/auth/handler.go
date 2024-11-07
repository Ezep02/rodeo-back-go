package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/ezep02/rodeo/utils"
)

type AuthHandler struct {
	AuthServ *AuthService
	ctx      context.Context
}

func NewAuthHandler(authServ *AuthService) *AuthHandler {
	return &AuthHandler{
		AuthServ: authServ,
		ctx:      context.Background(),
	}
}

func (h *AuthHandler) RegisterUserHandler(rw http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "No se puede procesar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Deserializar el cuerpo en un objeto User
	var user User
	if err := json.Unmarshal(b, &user); err != nil {
		http.Error(rw, "Error al deserializar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// hash de la contraseña (encriptacion)
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	user.Password = hash
	user.Name = strings.ToLower(user.Name)

	// Registrar el usuario utilizando el servicio
	registeredUser, err := h.AuthServ.RegisterUserServ(h.ctx, &user)
	if err != nil {
		http.Error(rw, "Error registrando usuario", http.StatusInternalServerError)
		return
	}

	// Devolver el usuario registrado como respuesta
	rw.WriteHeader(http.StatusCreated)

	// Serializar el usuario registrado en JSON y enviarlo como respuesta
	response, err := json.Marshal(registeredUser)
	if err != nil {
		http.Error(rw, "Error al serializar la respuesta", http.StatusInternalServerError)
		return
	}

	// si el registro fue exitoso, se crea un token
	tokenString, err := jwt.GenerateToken(user.ID, user.Is_admin, user.Name, user.Email, user.Surname, user.Phone_number)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Establece la cookie con el token
	http.SetCookie(rw, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   "", // Usa el dominio actual por defecto
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Cambiar a true si se usa HTTPS
		Path:     "/",
	})

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(response)
}

func (h *AuthHandler) LoginUserHandler(rw http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Couldn't parse request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var loggedUserReq LogUserReq

	if err := json.Unmarshal(b, &loggedUserReq); err != nil {
		http.Error(rw, "Error al deserializar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	loggedUser, err := h.AuthServ.LoginUserServ(h.ctx, &loggedUserReq)
	if err != nil {
		http.Error(rw, "Error al iniciar sesion del usuario", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(loggedUser)
	if err != nil {
		http.Error(rw, "Error al serializar la respuesta", http.StatusInternalServerError)
		return
	}

	// si el registro fue exitoso, se crea un token
	tokenString, err := jwt.GenerateToken(loggedUser.ID, loggedUser.Is_admin, loggedUser.Name, loggedUser.Email, loggedUser.Surname, loggedUser.Phone_number)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Establece la cookie con el token
	http.SetCookie(rw, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   "", // Usa el dominio actual por defecto
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Cambiar a true si se usa HTTPS
		Path:     "/",
	})

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(response)
}

func (h *AuthHandler) VerifyTokenHandler(rw http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("auth_token")

	if err != nil {
		http.Error(rw, "No token provided", http.StatusUnauthorized)
		return
	}
	// Validar el token
	tokenString := cookie.Value

	user, err := jwt.VerfiyToken(tokenString)
	if err != nil {
		http.Error(rw, "Token indalido o expirado", http.StatusUnauthorized)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(user)

}

func (h *AuthHandler) LogoutSession(w http.ResponseWriter, r *http.Request) {

	c := http.Cookie{
		Name:     "auth_token",
		MaxAge:   -1,
		Path:     "/",  // Asegúrate de que el Path coincida con el de la cookie original
		HttpOnly: true, // Evita que la cookie sea accesible desde JavaScript
		Secure:   true, // Solo permite que se envíe por HTTPS
	}

	http.SetCookie(w, &c)
	w.WriteHeader(http.StatusOK)
}
