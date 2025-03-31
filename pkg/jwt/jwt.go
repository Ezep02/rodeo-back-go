package jwt

import (
	"errors"
	"time"

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

func GenerateToken(user_id uint, isAdmin bool, name string, email string, surname string, phone_number string, isBarber bool, expirationTime time.Time) (string, error) {

	claim := JWTClaim{
		ID:           user_id,
		Name:         name,
		Email:        email,
		IsAdmin:      isAdmin,
		Surname:      surname,
		Phone_number: phone_number,
		Is_barber:    isBarber,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Hora de expiración en formato Unix
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(TokenKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func ValidateToken(signedString string) error {

	token, err := jwt.ParseWithClaims(
		signedString,
		&JWTClaim{},
		func(t *jwt.Token) (any, error) {
			return TokenKey, nil
		},
	)

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return errors.New("couldn't parse claims or token is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return errors.New("token expired")
	}
	return nil
}

func VerfiyToken(tokenString string) (*VerifyTokenRes, error) {
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
