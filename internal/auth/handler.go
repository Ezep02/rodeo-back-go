package auth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/ezep02/rodeo/utils"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
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
	var user User
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

	var loggedUserReq LogUserReq

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

func GoogleAuth(rw http.ResponseWriter, r *http.Request) {

	// Generar una URL para redirigir al usuario a Google para autenticación
	authURL := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(rw, r, authURL, http.StatusTemporaryRedirect)
}

func CallbackHandler(rw http.ResponseWriter, r *http.Request) {
	// Leer el código de autorización del query string
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(rw, "Code not found", http.StatusBadRequest)
		return
	}

	// Intercambiar el código de autorización por un token
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(rw, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Token obtenido, puedes usarlo para hacer solicitudes autenticadas
	userInfo, err := GetGoogleUserInfo(token)
	if err != nil {
		log.Println("Error fetching user info:", err)
		http.Error(rw, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}

	// parsear el string a tipo uuid
	parsedID := utils.GenerateDeterministicUint(userInfo.Sub)

	// autenticar creando el token
	tokenString, err := jwt.GenerateToken(parsedID, false, userInfo.Name, userInfo.Email, "", "", false, time.Now().Add(24*time.Hour))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Establece la cookie con el token
	http.SetCookie(rw, &http.Cookie{
		Name:     auth_token,
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour * 30),
		Domain:   "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Path:     "/",
	})
	// redireccionar al dashboard

	http.Redirect(rw, r, "http://localhost:5173/dashboard", http.StatusTemporaryRedirect)
}

func GetGoogleUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {

	client := googleOauthConfig.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}

func (h *AuthHandler) SendResetUserPasswordEmailHandler(rw http.ResponseWriter, r *http.Request) {

	var (
		smtpHost string = "smtp.gmail.com"
		// smtpPort := "587"
		sender   string   = "epereyra443@gmail.com"
		password string   = "cubrrxzypaskawzc"
		to       []string = []string{"pereyraezequiel15617866@outlook.es"}
		u        UserEmail
	)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Couldn't parse request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err := json.Unmarshal(b, &u); err != nil {
		http.Error(rw, "Error al deserializar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	user, err := h.AuthServ.SearchUserByEmail(h.ctx, u.Email)

	if err != nil {
		http.Error(rw, "Si el correo es válido, recibirás un email con instrucciones.", http.StatusInternalServerError)
		return
	}

	// crear un token utilizando los datos de user
	tokenString, err := jwt.GenerateToken(user.ID, user.Is_admin, user.Name, user.Email, user.Surname, user.Phone_number, user.Is_barber, time.Now().Add(15*time.Minute))

	if err != nil {
		http.Error(rw, "[Creacion token] Algo salio mal, vuelve a intentarlo mas tarde", http.StatusInternalServerError)
		return
	}

	// Autenticación con el servidor
	auth := smtp.PlainAuth("", sender, password, smtpHost)

	// Crear conexión segura con TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// Establecer conexión con el servidor SMTP
	conn, err := tls.Dial("tcp", smtpHost+":465", tlsConfig) // Usa puerto 465 para TLS directo
	if err != nil {
		log.Fatal("Error en conexión TLS:", err)
	}
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Fatal("Error creando cliente SMTP:", err)
	}

	// Autenticarse
	if err = client.Auth(auth); err != nil {
		log.Fatal("Error en autenticación:", err)
	}

	// Configurar el remitente y destinatario
	if err = client.Mail(sender); err != nil {
		log.Fatal(err)
	}
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			log.Fatal(err)
		}
	}

	// Escribir el mensaje
	wc, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf("Subject: 🔐 Recupera tu contraseña\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"<html><body>"+
		"<h2>🔐 Recuperación de contraseña</h2>"+
		"<p>Hola,</p>"+
		"<p>Has solicitado restablecer tu contraseña. Haz clic en el botón de abajo:</p>"+
		"<a href='http://localhost:5173/auth/recover/token=%s' "+
		"style='display:inline-block;background-color:#007bff;color:#ffffff;padding:10px 20px;text-decoration:none;border-radius:5px;'>Restablecer contraseña</a>"+
		"<p>Si no solicitaste esto, ignora este mensaje.</p>"+
		"<p>Saludos,<br>Equipo de Soporte</p>"+
		"</body></html>", tokenString)

	_, err = wc.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Cerrar conexión
	client.Quit()

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode("ok")
}

func (h *AuthHandler) ResetUserPassword(rw http.ResponseWriter, r *http.Request) {

	var (
		userData UserResetPassowrdReq
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
