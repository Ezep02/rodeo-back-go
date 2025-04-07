package utils

import (
	"errors"
	"time"

	"github.com/ezep02/rodeo/internal/orders/models"
	"github.com/golang-jwt/jwt"
)

type JWTOrderClaim struct {
	ID                  uint       `json:"ID"`
	Title               string     `json:"title"`
	Payer_name          string     `json:"payer_name"`
	Payer_surname       string     `json:"payer_surname"`
	Barber_id           int        `json:"barber_id"`
	Schedule_day_date   *time.Time `json:"schedule_day_date"`
	Schedule_start_time string     `json:"schedule_start_time"`
	User_id             int        `json:"user_id"`
	Price               float64    `json:"price"`
	Created_at          *time.Time `json:"Created_at"`
	jwt.StandardClaims
}

type VerifyOrderTokenRes struct {
	ID                  uint       `json:"ID"`
	Title               string     `json:"title"`
	Payer_name          string     `json:"payer_name"`
	Payer_surname       string     `json:"payer_surname"`
	Barber_id           int        `json:"barber_id"`
	User_id             int        `json:"user_id"`
	Schedule_day_date   *time.Time `json:"schedule_day_date"`
	Schedule_start_time string     `json:"schedule_start_time"`
	Price               float64    `json:"price"`
	Created_at          *time.Time `json:"Created_at"`
}

var TokenKey = []byte("mytokenapikey")

func GenerateOrderToken(order models.PendingOrderToken, expirationTime time.Time) (string, error) {

	claim := JWTOrderClaim{
		ID:                  order.ID,
		Title:               order.Title,
		Payer_name:          order.Payer_name,
		Payer_surname:       order.Payer_surname,
		Barber_id:           order.Barber_id,
		Schedule_day_date:   order.Schedule_day_date,
		Schedule_start_time: order.Schedule_start_time,
		User_id:             order.User_id,
		Price:               order.Price,
		Created_at:          order.Created_at,
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

func ValidateOrderToken(signedString string) error {

	token, err := jwt.ParseWithClaims(
		signedString,
		&JWTOrderClaim{},
		func(t *jwt.Token) (any, error) {
			return TokenKey, nil
		},
	)

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*JWTOrderClaim)
	if !ok {
		return errors.New("couldn't parse claims or token is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return errors.New("token expired")
	}
	return nil
}

func VerfiyToken(tokenString string) (*VerifyOrderTokenRes, error) {
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

		user := &VerifyOrderTokenRes{
			ID:                  claims["ID"].(uint),
			Title:               claims["title"].(string),
			Payer_name:          claims["email"].(string),
			Payer_surname:       claims["payer_surname"].(string),
			Barber_id:           claims["barber_id"].(int),
			User_id:             claims["user_id"].(int),
			Schedule_day_date:   claims["schedule_day_date"].(*time.Time),
			Schedule_start_time: claims["schedule_start_time"].(string),
			Price:               claims["price"].(float64),
			Created_at:          claims["Created_at"].(*time.Time),
		}

		return user, nil
	}

	return nil, errors.New("invalid token")
}
