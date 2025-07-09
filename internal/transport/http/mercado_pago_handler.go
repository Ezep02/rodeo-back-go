package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/service"
	custom_jwt "github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type CreatePreferenceRequest struct {
	CustomerName      string `json:"customer_name"`
	CustomerSurname   string `json:"customer_surname"`
	Date              string `json:"date"`
	Time              string `json:"time"`
	Products          []uint `json:"products"`
	SlotID            uint   `json:"slotID"`
	PaymentPercentage uint8  `json:"payment_percentage"`
}

type CreateSurchargePrefReq struct {
	OldSlotId      uint    `json:"old_slot_id"`
	NewSlotId      uint    `json:"new_slot_id"`
	ApptId         uint    `json:"appointment_id"`
	SurchargePrice float64 `json:"surcharge_price"`
}

type MepaHandler struct {
	apptSvc *service.AppointmentService
	prodSvc *service.ProductService
	slotSvc *service.SlotService
}

type JWTAppointmentClaim struct {
	ID uint `json:"ID"`
	jwt.StandardClaims
}

var (
	payment_token           = os.Getenv("PAYMENT_TOKEN")
	notification_url string = "https://acc764f836a1.ngrok-free.app" // URL de notificación
)

func NewMepaHandler(prodSvc *service.ProductService, apptSvc *service.AppointmentService, slotSvc *service.SlotService) *MepaHandler {
	return &MepaHandler{apptSvc, prodSvc, slotSvc}
}

func (h *MepaHandler) CreatePreference(c *gin.Context) {

	var (
		prefItem        []preference.ItemRequest
		req             CreatePreferenceRequest
		user_id         uint
		mp_access_token = os.Getenv("MP_ACCESS_TOKEN")
		auth_token      = os.Getenv("AUTH_TOKEN")
	)

	// 1. Validar tokens
	if mp_access_token == "" {
		c.JSON(400, gin.H{"error": "Faltan variables de entorno"})
		return
	}

	// 2. Obtener datos del cliente
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 3. Recuperar id de cliente si es que existe desde la session
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		log.Println("error", err)
	}

	// 4. Validar la cookie
	user, err := custom_jwt.VerfiySessionToken(cookie)
	if err != nil {
		log.Println("token invalido o expirado")
		user_id = 0
	} else {
		user_id = user.ID
	}

	// 5. Evitar citas duplicadas
	existingAppt, err := h.slotSvc.GetByID(c.Request.Context(), req.SlotID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if existingAppt.Is_booked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ya existe una cita para esta fecha y horario"})
		return
	}

	// 6. Iniciar configuracion de Mercado Pago
	cfg, err := config.New(mp_access_token)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al configurar Mercado Pago"})
		return
	}

	// 7. Crear cliente de Mercado Pago
	client := preference.NewClient(cfg)

	// 8. Recuperar productos mediante sus IDs
	for _, prodID := range req.Products {

		existingProd, err := h.prodSvc.GetByID(c.Request.Context(), prodID)
		if err != nil {
			if err == domain.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		price := existingProd.Price

		if req.PaymentPercentage < 100 {
			price = existingProd.Price * float64(req.PaymentPercentage) / 100
		}

		prefItem = append(prefItem, preference.ItemRequest{
			ID:          strconv.Itoa(int(existingProd.ID)),
			Title:       existingProd.Name,
			UnitPrice:   price,
			Quantity:    1,
			Description: existingProd.Description,
		})
	}

	// 9. Crear token temporal para redirigir al usuario si todo va bien
	claim := JWTAppointmentClaim{
		ID: req.SlotID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(payment_token))

	if err != nil {
		log.Println("Algo no fue bien creando el token", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "algo no fue bien creando el token"})
	}

	log.Println("TOKEN", tokenString)

	// 10. Crear preferencia
	request := preference.Request{
		Items: prefItem,
		Payer: &preference.PayerRequest{
			Name:    req.CustomerName,
			Surname: req.CustomerSurname,
		},
		NotificationURL: fmt.Sprintf("%s/api/v1/appointments/", notification_url),
		Metadata: map[string]any{
			"date":               req.Date,
			"slot_id":            req.SlotID,
			"payment_percentage": req.PaymentPercentage,
			"time":               req.Time,
			"user_id":            user_id,
		},
		BackURLs: &preference.BackURLsRequest{
			Success: fmt.Sprintf("http://localhost:5173/payment/success/%s", tokenString),
		},
	}

	// 11. Crear preferencia en Mercado Pago
	preferenceRes, err := client.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al crear preferencia en Mercado Pago"})
		return
	}

	c.JSON(200, gin.H{
		"message":    "Preferencia creada exitosamente",
		"init_point": preferenceRes.InitPoint,
	})
}

func (h *MepaHandler) CreateSurchargePreference(c *gin.Context) {
	var (
		req             CreateSurchargePrefReq
		mp_access_token = os.Getenv("MP_ACCESS_TOKEN")
		auth_token      = os.Getenv("AUTH_TOKEN")
	)

	// 1. Validar tokens
	if mp_access_token == "" {
		c.JSON(400, gin.H{"error": "Faltan variables de entorno"})
		return
	}

	// 2. Obtener datos del cliente
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 3. Verificar sesion del cliente
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error verificando la sesion"})
		return
	}

	// 4. Validar la cookie
	_, err = custom_jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tu sesion expiro"})
		return
	}

	// 5. Iniciar configuracion de Mercado Pago
	cfg, err := config.New(mp_access_token)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al configurar Mercado Pago"})
		return
	}

	// 7. Crear cliente de Mercado Pago
	client := preference.NewClient(cfg)

	// 8. Crear preferencia
	request := preference.Request{
		Items: []preference.ItemRequest{
			{
				Title:     "Reprogramacion de cita",
				UnitPrice: req.SurchargePrice,
				Quantity:  1,
			},
		},
		NotificationURL: fmt.Sprintf("%s/api/v1/appointments/surcharge", notification_url),
		Metadata: map[string]any{
			"old_slot_id": req.OldSlotId,
			"new_slot_id": req.NewSlotId,
			"appt_id":     req.ApptId,
		},
		BackURLs: &preference.BackURLsRequest{
			Success: "http://localhost:5173/appointment",
		},
	}

	// 9. Crear preferencia en Mercado Pago
	preferenceRes, err := client.Create(c.Request.Context(), request)
	if err != nil {
		log.Println("err", err)
		c.JSON(500, gin.H{"error": "Error al crear preferencia en Mercado Pago"})
		return
	}

	c.JSON(200, gin.H{
		"message":    "Preferencia creada exitosamente",
		"init_point": preferenceRes.InitPoint,
	})
}

func (h *MepaHandler) GetPayment(c *gin.Context) {
	tokenStr := c.Param("token")

	if tokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token vacío"})
		return
	}

	tokenData, err := VerfiyToken(tokenStr)
	if err != nil || tokenData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Token Data", tokenData.ID)

	existing, err := h.slotSvc.GetByID(c.Request.Context(), uint(tokenData.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no es posible recuperar la información del slot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "información recuperada correctamente",
		"slot":    existing,
	})
}

type VerifyTokenRes struct {
	ID float64 `json:"ID"`
}

func VerfiyToken(tokenString string) (*VerifyTokenRes, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Asegura que la firma sea la esperada
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(payment_token), nil
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

		slotID := &VerifyTokenRes{
			ID: claims["ID"].(float64),
		}

		return slotID, nil
	}

	return nil, errors.New("invalid token")
}
