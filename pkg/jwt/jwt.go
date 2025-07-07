package jwt

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/golang-jwt/jwt"
)

type JWTClaim struct {
	ID           uint   `json:"ID"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsAdmin      bool   `json:"is_admin"`
	Surname      string `json:"surname"`
	Phone_number string `json:"phone_number"`
	Is_barber    bool   `json:"is_barber"`
	jwt.StandardClaims
}

type VerifyTokenRes struct {
	ID           uint   `json:"ID"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Is_admin     bool   `json:"is_admin"`
	Surname      string `json:"surname"`
	Phone_number string `json:"phone_number"`
	Is_barber    bool   `json:"is_barber"`
}

type JWTResetPassowrdClaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var TokenKey = []byte("mytokenapikey")

func GenerateToken(user domain.User, expirationTime time.Time) (string, error) {

	claim := JWTClaim{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		IsAdmin:      user.Is_admin,
		Surname:      user.Surname,
		Phone_number: user.Phone_number,
		Is_barber:    user.Is_barber,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(TokenKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func VerfiySessionToken(tokenString string) (*VerifyTokenRes, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Asegura que la firma sea la esperada
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return TokenKey, nil
	})

	if err != nil {
		return nil, errors.New("token couldn't be parse")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verifica cualquier otra cosa que necesites en las reclamaciones
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return nil, errors.New("token has expired")
			}
		}

		// Realiza la conversión del ID desde float64 a uint
		id, ok := claims["ID"].(float64)
		if !ok {
			return nil, errors.New("invalid token: ID claim missing or not a number")
		}

		user := &VerifyTokenRes{
			ID:           uint(id),
			Name:         claims["name"].(string),
			Email:        claims["email"].(string),
			Is_admin:     claims["is_admin"].(bool),
			Surname:      claims["surname"].(string),
			Phone_number: claims["phone_number"].(string),
			Is_barber:    claims["is_barber"].(bool),
		}

		return user, nil
	}

	return nil, errors.New("invalid token")
}

// Crea una cookie de autenticación con el token JWT
func NewAuthTokenCookie(token string) *http.Cookie {

	name := os.Getenv("AUTH_TOKEN")

	return &http.Cookie{
		Name:     name,
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour * 30),
		Domain:   "", // Usa el dominio actual por defecto
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Cambiar a true si se usa HTTPS
		Path:     "/",
	}
}
