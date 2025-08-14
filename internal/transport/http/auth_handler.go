package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
