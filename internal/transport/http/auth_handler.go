package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/service"
	googleauth "github.com/ezep02/rodeo/pkg/google_auth"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/ezep02/rodeo/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc}
}

type RegisterUserRequest struct {
	Name         string `json:"name" binding:"required"`
	Surname      string `json:"surname" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Phone_number string `json:"phone_number"`
	IsAdmin      bool   `json:"is_admin"`
	IsBarber     bool   `json:"is_barber"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GoogleUserInfoReq struct {
	Sub           string `json:"sub"` // ID unico del usuario
	Email         string `json:"email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	VerifiedEmail bool   `json:"verified_email"`
}

func (h *AuthHandler) Register(c *gin.Context) {

	var (
		req RegisterUserRequest
	)

	// 1. Obtener datos de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 2. Encriptar contraseña
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	// 3. Contruir consulta
	user := domain.User{
		Name:         req.Name,
		Surname:      req.Surname,
		Password:     hash,
		Phone_number: req.Phone_number,
		Email:        req.Email,
		Is_admin:     req.IsAdmin,
		Is_barber:    req.IsBarber,
	}

	// 4. Registrar usuario
	if err := h.svc.Register(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 5. Recuperar usuario
	existing, err := h.svc.GetByEamil(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "se registro, pero no fue posible recuperar al usuario"})
		return
	}

	// 6. Crear token de sesion
	tokenStr, err := jwt.GenerateToken(*existing, time.Now().Add(24*time.Hour*30))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error creando token de sesion"})
		return
	}

	// 7. Establece la cookie con el token
	httpCookie := jwt.NewAuthTokenCookie(tokenStr)
	http.SetCookie(c.Writer, httpCookie)

	c.JSON(http.StatusOK, gin.H{
		"message": "operacion exitosa",
		"user":    existing,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {

	var (
		req LoginUserRequest
	)

	// 1. Obtener datos de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 2. Obtener usuario
	existing, err := h.svc.Login(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usuario no registrado"})
		return
	}

	// 3. Comparar contraseña
	if err := utils.HashCompare(existing.Password, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contraseña incorrecta volve a intentarlo"})
		return
	}

	// 4. Crear token de sesion
	tokenStr, err := jwt.GenerateToken(*existing, time.Now().Add(24*time.Hour*30))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error creando token de sesion"})
		return
	}

	// 5. Establece la cookie con el token
	httpCookie := jwt.NewAuthTokenCookie(tokenStr)
	http.SetCookie(c.Writer, httpCookie)

	c.JSON(http.StatusOK, gin.H{
		"message": "inicio de sesion exitoso",
		"user":    existing,
	})
}

func (h *AuthHandler) VerifySession(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	// 1. Recuperar cookies
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	// 2. Validar la cookie
	user, err := jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	// 3. Devolver usuario autenticado
	c.JSON(http.StatusOK, gin.H{
		"message": "usuario autenticado correctamente",
		"user":    user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	// 1. Establecer fecha de vencimiento
	cookie := http.Cookie{
		Name:     auth_token,
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	// 2. Establer expiracion de cookie
	http.SetCookie(c.Writer, &cookie)
	c.JSON(http.StatusOK, gin.H{"message": "sesion cerrada correctamente"})
}

func (h *AuthHandler) GoogleAuth(c *gin.Context) {

	// 1. Crear configuracion basica de Google Auth
	googleOauthConfig := googleauth.CreateGoogleAuthConfig()

	// 2. Generar URL necesaria para redirigir al usuario a Google para autenticación
	googleAuthURL := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	c.Redirect(http.StatusTemporaryRedirect, googleAuthURL)
}

func (h *AuthHandler) CallbackHandler(c *gin.Context) {

	// 1. Leer el código de autorización del query string
	code := c.Request.URL.Query().Get("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se encontro el codigo en la respuesta"})
		return
	}

	// 2. Crear configuracion basica de Google Auth
	googleOauthConfig := googleauth.CreateGoogleAuthConfig()

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible intercambiar el token"})
		return
	}

	// 3. Recuperar la informacion del usuario
	userInfo, err := GetGoogleUserInfo(token, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no fue posible recuperar los datos del usuario"})
		return
	}

	// 4. Verificar la existencia el usuario
	existingUser, err := h.svc.Login(c.Request.Context(), userInfo.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "al parecer el usuario no esta registrado"})
		return
	}

	// 5. Creando token de sesion
	tokenString, err := jwt.GenerateToken(*existingUser, time.Now().Add(24*time.Hour))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error al crear el token de sesion"})
		return
	}

	// 7. Crear cookie de sesion
	cookie := jwt.NewAuthTokenCookie(tokenString)

	// 8. Establece la cookie
	http.SetCookie(c.Writer, cookie)

	// 9. Redireccionar al dashboard
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/")
}

func GetGoogleUserInfo(token *oauth2.Token, c *gin.Context) (*GoogleUserInfoReq, error) {
	var userInfo GoogleUserInfoReq

	// 1. Crear configuracion basica de Google Auth
	googleOauthConfig := googleauth.CreateGoogleAuthConfig()

	// 2. Crear un cliente http
	client := googleOauthConfig.Client(c, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}

// TODO: Integrar handlers para recuperar contraseña, y envio de token por mail

// Reset de contraseña
type UserEmailRes struct {
	Email string `json:"email"`
}

func (h *AuthHandler) SendResetPasswordEmail(c *gin.Context) {

	var (
		auth_token        = os.Getenv("AUTH_TOKEN")
		smtpHost   string = "smtp.gmail.com"
		//smtpPort = "587"
		sender   string   = "epereyra443@gmail.com"
		password string   = "cubrrxzypaskawzc"
		to       []string = []string{"pereyraezequiel15617866@outlook.es"}
		req      UserEmailRes
	)

	// 1. Recuperar cookies
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	// 2. Validar la cookie
	if _, err := jwt.VerfiySessionToken(cookie); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	// 3. Obtener ID del usuario desde el parametro de la ruta
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario no proporcionado"})
		return
	}

	// 4. Parsear el ID a entero
	paredID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario invalido"})
		return
	}

	// 5. Obtener datos de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.svc.GetByID(c.Request.Context(), uint(paredID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error al recuperar el usuario"})
		return
	}

	// 6. crear un token utilizando los datos de user
	tokenString, err := jwt.GenerateToken(*user, time.Now().Add(15*time.Minute))
	if err != nil {
		log.Println("Error creando token:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "algo no fue bien al enviar el email"})
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Si el correo es válido, recibirás un email con instrucciones.",
	})
}

type UserResetPassowrdReq struct {
	New_password string `json:"new_password"`
	Token        string `json:"token"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {

	var (
		req UserResetPassowrdReq
	)

	// 1. Obtener datos de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := jwt.VerfiySessionToken(req.Token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	// hash de la contraseña (encriptacion)
	_, err = utils.HashPassword(req.New_password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contraseña invalida"})
		return
	}

	// if err := h.AuthServ.UpdateUserPasswordServ(c., int(token.ID), hash); err != nil {
	// 	http.Error(rw, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Contraseña actualizada correctamente",
	})
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
		req        domain.User
	)

	// 1. Recuperar cookies
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	// 2. Validar la cookie
	if _, err := jwt.VerfiySessionToken(cookie); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	// 3. Obtener ID del usuario desde el parametro de la ruta
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario no proporcionado"})
		return
	}

	// 4. Parsear el ID a entero
	_, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario invalido"})
		return
	}

	// 5. Obtener datos de la consulta
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 4. Actualizar usuario
	// if err := h.svc.UpdateUser(c.Request.Context(), &req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "no fue posible actualizar el usuario"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "usuario actualizado correctamente",
		"user":    req,
	})
}
