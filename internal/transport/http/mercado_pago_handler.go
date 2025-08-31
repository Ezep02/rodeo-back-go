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
	CouponCode        string `json:"coupon_code,omitempty"` // Optional coupon code
}

type CreateSurchargePrefReq struct {
	OldSlotId      uint    `json:"old_slot_id"`
	NewSlotId      uint    `json:"new_slot_id"`
	ApptId         uint    `json:"appointment_id"`
	SurchargePrice float64 `json:"surcharge_price"`
}

type MepaHandler struct {
	apptSvc   *service.AppointmentService
	prodSvc   *service.ProductService
	slotSvc   *service.SlotService
	couponSvc *service.CouponService
}

type JWTAppointmentClaim struct {
	ID uint `json:"ID"`
	jwt.StandardClaims
}

var (
	payment_token           = os.Getenv("PAYMENT_TOKEN")
	notification_url string = "https://c38c518e6523.ngrok-free.app" // URL de notificación
)

func NewMepaHandler(
	prodSvc *service.ProductService,
	apptSvc *service.AppointmentService,
	slotSvc *service.SlotService,
	couponSvc *service.CouponService) *MepaHandler {
	return &MepaHandler{apptSvc, prodSvc, slotSvc, couponSvc}
}

func (h *MepaHandler) CreatePreference(c *gin.Context) {
	var (
		req           CreatePreferenceRequest
		prefItems     []preference.ItemRequest
		userID        uint
		mpAccessToken = os.Getenv("MP_ACCESS_TOKEN")
		authToken     = os.Getenv("AUTH_TOKEN")
		couponToApply *domain.Coupon
	)

	// 1. Validar tokens
	if mpAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan variables de entorno"})
		return
	}

	// 2. Parsear request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 3. Recuperar id de cliente desde cookie
	cookie, _ := c.Cookie(authToken)
	if user, err := custom_jwt.VerfiySessionToken(cookie); err == nil {
		userID = user.ID
	}

	// 4. Evitar citas duplicadas
	existingAppt, err := h.slotSvc.GetByID(c.Request.Context(), req.SlotID)
	if err != nil || existingAppt.IsBooked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ya existe una cita para esta fecha y horario"})
		return
	}

	// 5. Validar cupón
	if req.CouponCode != "" {
		coupon, err := h.couponSvc.GetByCode(c.Request.Context(), req.CouponCode)
		if err != nil || time.Now().After(coupon.ExpireAt) || !coupon.IsAvailable {
			c.JSON(http.StatusBadRequest, gin.H{"error": "el codigo de descuento ingresado no es valido o expiro"})
			return
		}
		couponToApply = coupon
	}

	// 6. Configurar Mercado Pago
	cfg, err := config.New(mpAccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al configurar Mercado Pago"})
		return
	}
	client := preference.NewClient(cfg)

	// 7. Recuperar productos y calcular precio final
	for _, prodID := range req.Products {
		prod, err := h.prodSvc.GetByID(c.Request.Context(), prodID)
		if err != nil {
			status := http.StatusInternalServerError
			if err == domain.ErrNotFound {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		unitPrice := prod.Price

		// Aplicar descuento de promoción
		if prod.PromotionDiscount > 0 && prod.PromotionEndDate.After(time.Now()) {
			unitPrice *= (1 - float64(prod.PromotionDiscount)/100)
		}

		// Aplicar descuento de cupón
		if couponToApply != nil {
			unitPrice *= (1 - float64(couponToApply.DiscountPercentage)/100)
		}

		// Aplicar porcentaje de pago
		if req.PaymentPercentage < 100 {
			unitPrice *= float64(req.PaymentPercentage) / 100
		}

		prefItems = append(prefItems, preference.ItemRequest{
			ID:          strconv.Itoa(int(prod.ID)),
			Title:       prod.Name,
			UnitPrice:   unitPrice,
			Quantity:    1,
			Description: prod.Description,
		})
	}

	// 8. Crear token temporal
	claim := JWTAppointmentClaim{
		ID: req.SlotID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(payment_token))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando token"})
		return
	}

	// 9. Crear preferencia
	request := preference.Request{
		Items: prefItems,
		Payer: &preference.PayerRequest{
			Name:    req.CustomerName,
			Surname: req.CustomerSurname,
		},
		NotificationURL: fmt.Sprintf("%s/api/v1/appointments/", notification_url),
		Metadata: map[string]any{
			"date":               req.Date,
			"slot_id":            req.SlotID,
			"payment_percentage": req.PaymentPercentage,
			"user_id":            userID,
			"coupon_code":        req.CouponCode,
		},
		BackURLs: &preference.BackURLsRequest{
			Success: fmt.Sprintf("http://localhost:5173/payment/success/%s", tokenString),
		},
	}

	// 10. Crear preferencia en Mercado Pago
	preferenceRes, err := client.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear preferencia en Mercado Pago"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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
