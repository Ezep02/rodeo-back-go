package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ezep02/rodeo/internal/auth/models"
	"github.com/ezep02/rodeo/internal/auth/services"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/ezep02/rodeo/utils"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	AuthServ *services.AuthService
	ctx      context.Context
}

func NewAuthHandler(authServ *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthServ: authServ,
		ctx:      context.Background(),
	}
}

var (
	clientID          string
	redirectURI       string
	clientSecret      string
	auth_token        string
	googleOauthConfig *oauth2.Config
)

func init() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer el archivo .env: %v", err)
	}

	auth_token = viper.GetString("AUTH_TOKEN")
	clientID = viper.GetString("CLIENT_ID")
	clientSecret = viper.GetString("CLIENT_SECRET")
	redirectURI = viper.GetString("REDIRECT_URI")

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		log.Fatalf("Error: Las variables de entorno requeridas no están configuradas")
	}

	googleOauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
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
	var user models.User
	if err := json.Unmarshal(b, &user); err != nil {
		http.Error(rw, "Error al deserializar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// verificar que el email no este en uso
	u, err := h.AuthServ.SearchUserByEmail(h.ctx, user.Email)
	if err != nil {
		http.Error(rw, "Algo salio mal intentando registrar al usuario", http.StatusBadRequest)
		return
	}

	if u.Email != "" {
		http.Error(rw, "El email ingresado ya se encuentra en uso", http.StatusBadRequest)
		return
	}

	// hash de la contraseña (encriptacion)
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = hash
	user.Name = strings.ToLower(user.Name)

	// Registrar el usuario utilizando el servicio
	registeredUser, err := h.AuthServ.RegisterUserServ(h.ctx, &user)
	if err != nil {
		http.Error(rw, "Error registrando usuario", http.StatusInternalServerError)
		return
	}

	// Serializar el usuario registrado en JSON y enviarlo como respuesta
	response, err := json.Marshal(registeredUser)
	if err != nil {
		http.Error(rw, "Error al serializar la respuesta", http.StatusInternalServerError)
		return
	}

	// si el registro fue exitoso, se crea un token
	tokenString, err := jwt.GenerateToken(user.ID, user.Is_admin, user.Name, user.Email, user.Surname, user.Phone_number, user.Is_barber, time.Now().Add(24*time.Hour))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Establece la cookie con el token
	http.SetCookie(rw, &http.Cookie{
		Name:     auth_token,
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour * 30),
		Domain:   "", // Usa el dominio actual por defecto
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Cambiar a true si se usa HTTPS
		Path:     "/",
	})

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response)
}

func (h *AuthHandler) LoginUserHandler(rw http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Couldn't parse request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var loggedUserReq models.LogUserReq

	if err := json.Unmarshal(b, &loggedUserReq); err != nil {
		http.Error(rw, "Error al deserializar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	user, err := h.AuthServ.SearchUserByEmail(h.ctx, loggedUserReq.Email)
	if err != nil {
		http.Error(rw, "Error al iniciar sesion del usuario", http.StatusInternalServerError)
		return
	}

	if err := utils.HashCompare(user.Password, loggedUserReq.Password); err != nil {
		http.Error(rw, "Contraseña incorrecta, volve a intentarlo", http.StatusInternalServerError)
		return
	}

	// si el registro fue exitoso, se crea un token
	tokenString, err := jwt.GenerateToken(user.ID, user.Is_admin, user.Name, user.Email, user.Surname, user.Phone_number, user.Is_barber, time.Now().Add(24*time.Hour))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Establece la cookie con el token
	http.SetCookie(rw, &http.Cookie{
		Name:     auth_token,
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour * 30),
		Domain:   "", // Usa el dominio actual por defecto
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Cambiar a true si se usa HTTPS
		Path:     "/",
	})

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(user)
}

func (h *AuthHandler) VerifyTokenHandler(rw http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(auth_token)

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
		Name:     auth_token,
		MaxAge:   -1,
		Path:     "/",  // Asegúrate de que el Path coincida con el de la cookie original
		HttpOnly: true, // Evita que la cookie sea accesible desde JavaScript
		Secure:   true, // Solo permite que se envíe por HTTPS
	}

	http.SetCookie(w, &c)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) ResetUserPassword(rw http.ResponseWriter, r *http.Request) {

	var (
		userData models.UserResetPassowrdReq
	)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Couldn't parse request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err := json.Unmarshal(b, &userData); err != nil {
		http.Error(rw, "Error al deserializar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	token, err := jwt.VerfiyToken(userData.Token)

	if err != nil {
		http.Error(rw, "Token indalido o expirado", http.StatusUnauthorized)
		return
	}

	// hash de la contraseña (encriptacion)
	hash, err := utils.HashPassword(userData.New_password)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.AuthServ.UpdateUserPasswordServ(h.ctx, int(token.ID), hash); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("contraseña actualizada correctamente")
}
