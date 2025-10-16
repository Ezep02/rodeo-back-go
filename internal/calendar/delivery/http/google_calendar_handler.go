package http

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ezep02/rodeo/internal/calendar/domain"
	"github.com/ezep02/rodeo/internal/calendar/usecase"
	googleauth "github.com/ezep02/rodeo/pkg/google_auth"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendarHandler struct {
	calendarService *usecase.CalendarService
}

func NewGoogleCalendarHandler(calendarSvc *usecase.CalendarService) *GoogleCalendarHandler {
	return &GoogleCalendarHandler{calendarSvc}
}

func (h *GoogleCalendarHandler) GoogleCalendarLogin(c *gin.Context) {
	googleOauthConfig := googleauth.CreateGoogleAuthConfig(
		[]string{calendar.CalendarScope},
		"http://localhost:9090/api/v1/calendar/google-calendar/callback",
	)

	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

func (h *GoogleCalendarHandler) GoogleCalendarCallback(c *gin.Context) {
	ctx := context.Background()
	authToken := os.Getenv("AUTH_TOKEN")

	googleOauthConfig := googleauth.CreateGoogleAuthConfig(
		[]string{calendar.CalendarScope},
		"http://localhost:9090/api/v1/calendar/google-calendar/callback",
	)

	// Recibir el "code" de Google
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se encontró el código en la respuesta"})
		return
	}

	// Intercambiar el code por tokens
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al intercambiar el code por token: " + err.Error()})
		return
	}

	// Validar sesión del usuario mediante cookie
	cookie, err := c.Cookie(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}
	user, err := jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
		return
	}

	// Intentar obtener el refresh token actual (si ya existe en DB)
	storedToken, _ := h.calendarService.GetToken(c.Request.Context(), user.ID)

	refresh := token.RefreshToken
	if refresh == "" && storedToken != nil {
		// Mantener el refresh token viejo
		refresh = storedToken.RefreshToken
	}

	// Guardar token en DB sin perder el refresh
	err = h.calendarService.SaveToken(c.Request.Context(), user.ID, &domain.GoogleCalendarToken{
		UserID:       user.ID,
		AccessToken:  token.AccessToken,
		RefreshToken: refresh,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo guardar el token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Google Calendar conectado correctamente"})
}

func (h *GoogleCalendarHandler) GoogleCalendarVerify(c *gin.Context) {
	authToken := os.Getenv("AUTH_TOKEN")

	// Validar sesión del usuario mediante cookie
	cookie, err := c.Cookie(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	user, err := jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
		return
	}

	// Verificar si el usuario tiene un token de Google Calendar
	storedToken, err := h.calendarService.GetToken(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"calendar_is_active": false})
		return
	}

	// Inicializar googleOauthConfig
	googleOauthConfig := googleauth.CreateGoogleAuthConfig(
		[]string{calendar.CalendarScope},
		"http://localhost:9090/api/v1/calendar/google-calendar/callback",
	)

	// Si el token existe, crear un tokenSource
	tokenSource := googleOauthConfig.TokenSource(context.Background(), &oauth2.Token{
		AccessToken:  storedToken.AccessToken,
		RefreshToken: storedToken.RefreshToken,
		TokenType:    storedToken.TokenType,
		Expiry:       storedToken.Expiry,
	})

	newToken, err := tokenSource.Token()

	if err != nil {
		log.Println("[DEBUG NEW TOKEN]", newToken)
		log.Println("[DEBUG NEW TOKEN ERROR]", err.Error())
		c.JSON(http.StatusOK, gin.H{"calendar_is_active": false})
		return
	}

	// Opcional: si hay cambios, actualizar la DB
	if newToken.AccessToken != storedToken.AccessToken {
		if err := h.calendarService.SaveToken(c.Request.Context(), user.ID, &domain.GoogleCalendarToken{
			UserID:       user.ID,
			AccessToken:  newToken.AccessToken,
			RefreshToken: newToken.RefreshToken,
			Expiry:       newToken.Expiry,
			TokenType:    newToken.TokenType,
		}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo actualizar el token: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"calendar_is_active": true})
}

func (h *GoogleCalendarHandler) Create(c *gin.Context) {

	var (
		authToken = os.Getenv("AUTH_TOKEN")
	)

	// Validar sesión del usuario mediante cookie
	cookie, err := c.Cookie(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	user, err := jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
		return
	}

	if !user.IsBarber {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	// Si el token existe, crear un tokenSource

	savedToken, err := h.calendarService.GetToken(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usuario no encontrado"})
		return
	}

	googleOauthConfig := googleauth.CreateGoogleAuthConfig(
		[]string{calendar.CalendarScope},

		"http://localhost:9090/api/v1/calendar/google-calendar/callback",
	)

	client := googleOauthConfig.Client(context.Background(), &oauth2.Token{
		AccessToken:  savedToken.AccessToken,
		TokenType:    savedToken.TokenType,
		RefreshToken: savedToken.RefreshToken,
		Expiry:       savedToken.Expiry,
	})

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Println("error creando servicio calendario:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "error creando servicio calendario"})
		return
	}

	cal := &calendar.Calendar{
		Summary:  "El Rodeo",
		TimeZone: "America/Argentina/Buenos_Aires",
	}

	createdCal, err := srv.Calendars.Insert(cal).Do()
	if err != nil {
		log.Println("error creando calenadrio", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "error creando calendario"})
		return
	}

	// Asignar calendario al usuario
	go func(calendar_id string) {

		ctx := context.Background()
		if err := h.calendarService.AssignBarberCalendar(ctx, calendar_id, user.ID); err != nil {
			log.Println("Algo no fue bien asignando el calendario al barbero")
			return
		}

	}(createdCal.Id)

	c.JSON(http.StatusOK, gin.H{
		"calendar_id": createdCal.Id,
	})
}
