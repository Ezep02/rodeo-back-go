package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ezep02/rodeo/internal/auth/models"
	"golang.org/x/oauth2"
)

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
	// token, err := googleOauthConfig.Exchange(context.Background(), code)
	// if err != nil {
	// 	http.Error(rw, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Token obtenido, puedes usarlo para hacer solicitudes autenticadas
	// userInfo, err := GetGoogleUserInfo(token)
	// if err != nil {
	// 	log.Println("Error fetching user info:", err)
	// 	http.Error(rw, "Failed to fetch user info", http.StatusInternalServerError)
	// 	return
	// }

	// parsear el string a tipo uuid
	// parsedID := utils.GenerateDeterministicUint(userInfo.Sub)

	// autenticar creando el token
	// tokenString, err := jwt.GenerateToken(parsedID, false, userInfo.Name, userInfo.Email, "", "", false, time.Now().Add(24*time.Hour))
	// if err != nil {
	// 	http.Error(rw, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// Establece la cookie con el token
	// http.SetCookie(rw, &http.Cookie{
	// 	Name:     auth_token,
	// 	Value:    tokenString,
	// 	Expires:  time.Now().Add(24 * time.Hour * 30),
	// 	Domain:   "",
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteLaxMode,
	// 	Secure:   false,
	// 	Path:     "/",
	// })
	// redireccionar al dashboard

	http.Redirect(rw, r, "http://localhost:5173/dashboard", http.StatusTemporaryRedirect)
}

func GetGoogleUserInfo(token *oauth2.Token) (*models.GoogleUserInfo, error) {

	client := googleOauthConfig.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo models.GoogleUserInfo

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}
