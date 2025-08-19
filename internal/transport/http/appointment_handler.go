package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/service"
	"github.com/ezep02/rodeo/internal/transport/sse"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/ezep02/rodeo/utils"

	"github.com/gin-gonic/gin"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"gopkg.in/gomail.v2"
)

type AppointmentHandler struct {
	svc       *service.AppointmentService
	couponSvc *service.CouponService
	slotSvc   *service.SlotService
	sseServer *sse.SSEHandler
}

func NewAppointmentHandler(
	apptService *service.AppointmentService,
	couponSvc *service.CouponService,
	sseServer *sse.SSEHandler,
	slotSvc *service.SlotService) *AppointmentHandler {
	return &AppointmentHandler{apptService, couponSvc, slotSvc, sseServer}
}

type UpdateAppointmentRequest struct {
	OldSlotId uint `json:"old_slot_id"`
	NewSlotId uint `json:"new_slot_id"`
}

type CreateAppointmentReq struct {
	Additional_info payment.AdditionalInfoResponse
	Metadata_info   map[string]any `json:"metadata"`
}

func (h *AppointmentHandler) Create(c *gin.Context) {

	// 1. Leer ACCESS_TOKEN
	mpAccessToken := os.Getenv("MP_ACCESS_TOKEN")
	if mpAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan variables de entorno"})
		return
	}

	// 2. Decodificar payload enviado por mp
	var payload map[string]any
	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// 3. Recuperar del payload el campo id almacenado dentro de data
	data, ok := payload["data"].(map[string]any)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'data' inválido"})
		return
	}

	paymentStr := fmt.Sprintf("%v", data["id"])

	paymentID, err := strconv.ParseInt(paymentStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de pago inválido"})
		return
	}

	// 4. Inicializar el cliente de Mercado Pago
	cfg, err := config.New(mpAccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo configurar el cliente"})
		return
	}

	// 5. Consultar pago utilizanodo el ID
	paymentClient := payment.NewClient(cfg)

	paymentInfo, err := paymentClient.Get(context.Background(), int(paymentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pago no encontrado"})
		return
	}

	// 6. Parsear la metadata
	metadata, err := utils.MetadataParser(paymentInfo.Metadata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 7. Mapear IDs a domain.Product
	products := make([]domain.Product, len(paymentInfo.AdditionalInfo.Items))

	for i, prodID := range paymentInfo.AdditionalInfo.Items {
		prodIDUint, err := strconv.ParseUint(prodID.ID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de producto inválido"})
			return
		}
		products[i] = domain.Product{ID: uint(prodIDUint)}
	}

	if metadata.CouponCode != "" {

		if err := h.couponSvc.UpdateStatus(c.Request.Context(), strings.ToUpper(metadata.CouponCode)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando el estado del cupón"})
			return
		}
		log.Println("[CUPON] Cupón utilizado:", metadata.CouponCode)
	}

	// 8. Crear la cita
	newAppt := domain.Appointment{
		ClientName:        paymentInfo.AdditionalInfo.Payer.FirstName,
		ClientSurname:     paymentInfo.AdditionalInfo.Payer.LastName,
		SlotID:            metadata.SlotID,
		PaymentPercentage: metadata.PaymentPercentage,
		UserID:            metadata.UserID,
		Status:            "active",
		Products:          products,
	}

	if err := h.svc.Schedule(c.Request.Context(), &newAppt); err != nil {
		log.Println("Error creando cita:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 9. Recuperar la orden
	slot, err := h.slotSvc.GetByID(c.Request.Context(), newAppt.SlotID)
	if err != nil || slot == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve slot"})
		return
	}
	newAppt.Slot = *slot

	log.Println("[CREANDO CITA] Nueva cita creada con ID:", newAppt.ID)

	// 10. Enviar stream de datos
	ssePayload := sse.SSEMessage{
		Type: "appointment_created",
		Data: newAppt,
	}

	jsonMsg, _ := json.Marshal(ssePayload)
	h.sseServer.Hub.Broadcast(string(jsonMsg))

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cita creada exitosamente",
	})
}

func (h *AppointmentHandler) ListByDateRange(c *gin.Context) {

	startStr := c.Param("start")
	endStr := c.Param("end")

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start date invalid"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end date invalid"})
		return
	}

	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end date must be after start date"})
		return
	}

	appts, err := h.svc.ListByDateRange(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch appointments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"appointments": appts,
		"total":        len(appts),
	})
}

func (h *AppointmentHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	appt, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cita no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching appointment"})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (h *AppointmentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var req UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 3. Actualizar la cita
	updatedAppt := &domain.Appointment{
		ID:     uint(id),
		SlotID: req.NewSlotId,
	}

	if err := h.svc.Update(c.Request.Context(), uint(id), req.OldSlotId, updatedAppt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4 Recuperar slot actualizado
	updatedSlot, err := h.slotSvc.GetByID(c.Request.Context(), req.NewSlotId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error recuperando slot"})
		return
	}

	// 5. Crear evento
	ssePayload := sse.SSEMessage{
		Type: "appointment_updated",
		Data: domain.Appointment{
			ID:   uint(id),
			Slot: *updatedSlot,
		},
	}
	log.Println("[ENVIANDO EVENTO ACTUALIZAR]")

	// 6. Despachar evento
	jsonMsg, _ := json.Marshal(ssePayload)
	h.sseServer.Hub.Broadcast(string(jsonMsg))

	c.JSON(http.StatusOK, updatedAppt)
}

type CancelReq struct {
	Recharge float64   `json:"recharge"`
	ExpireAt time.Time `json:"expire_at"`
}

func (h *AppointmentHandler) Cancel(c *gin.Context) {

	var (
		req        CancelReq
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	// 3. Recuperar id de cliente si es que existe desde la session
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	// 4. Validar la cookie
	user, err := jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		log.Println("Cancel error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "solicitud invalida",
		})
		return
	}

	// Verificar que no este cancelada
	exist, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cita no encontrada"})
		return
	}

	if exist.Status == "cancelled" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No podes cancelar dos veces la misma cita"})
		return
	}

	// Si pago completo, se crea un cupon
	if req.Recharge > 0 {

		coupon, err := h.couponSvc.GenerateCoupon(12)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando cupon"})
			return
		}

		log.Println("[CUPON] CREANDO CUPON DE DESCUENTO", req.Recharge)

		go func(userID uint, recharge float64, code string) {
			ctx := context.Background()
			log.Println("TIEMPO ASIGNADO", req.ExpireAt)

			err := h.couponSvc.Create(ctx, &domain.Coupon{
				Code:               strings.ToUpper(code),
				UserID:             userID,
				DiscountPercentage: recharge,
				IsAvailable:        true,
				CreatedAt:          time.Now(),
				ExpireAt:           req.ExpireAt,
			})

			if err != nil {
				log.Printf("Error creando cupon: %v\n", err)
				return
			}

			log.Println("Cupon creado correctamente")
		}(user.ID, req.Recharge, coupon)
	}

	if err := h.svc.Cancel(c.Request.Context(), uint(id)); err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cita no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Comunicar al dashboard de la cancelacion
	ssePayload := sse.SSEMessage{
		Type: "appointment_cancelled",
		Data: exist,
	}

	jsonMsg, _ := json.Marshal(ssePayload)
	h.sseServer.Hub.Broadcast(string(jsonMsg))
	// Si hay recargo, generar cupon
	c.JSON(http.StatusOK, "cita cancelada exitosamente")

}

func (h *AppointmentHandler) GetTotal(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	total, err := h.svc.GetTotalPrice(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"appointment_id": id,
		"total_price":    total,
	})
}

func (h *AppointmentHandler) GetByUserID(c *gin.Context) {

	auth_token := os.Getenv("AUTH_TOKEN")

	// TODO: reemplazarlo por un middleware de session
	sessionToken, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sesion expirada o token invalido",
		})
		return
	}

	_, err = jwt.VerfiySessionToken(sessionToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "posible id invalido",
		})
	}

	appt, err := h.svc.GetByUserID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fue posible recuperar las citas",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"appointments": appt,
	})
}

func (h *AppointmentHandler) Surcharge(c *gin.Context) {

	// 1. Leer ACCESS_TOKEN
	mpAccessToken := os.Getenv("MP_ACCESS_TOKEN")
	if mpAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan variables de entorno"})
		return
	}

	// 2. Decodificar payload enviado por mp
	var payload map[string]any
	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// 3. Recuperar del payload el campo id almacenado dentro de data
	data, ok := payload["data"].(map[string]any)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'data' inválido"})
		return
	}

	paymentStr := fmt.Sprintf("%v", data["id"])

	paymentID, err := strconv.ParseInt(paymentStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de pago inválido"})
		return
	}

	// 4. Inicializar el cliente de Mercado Pago
	cfg, err := config.New(mpAccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo configurar el cliente"})
		return
	}

	// 5. Consultar pago utilizanodo el ID
	paymentClient := payment.NewClient(cfg)

	paymentInfo, err := paymentClient.Get(context.Background(), int(paymentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pago no encontrado"})
		return
	}

	metadata, err := utils.SurchargeMetadataParcer(paymentInfo.Metadata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Actualizar la cita
	updatedAppt := &domain.Appointment{
		ID:     metadata.ApptId,
		SlotID: metadata.NewSlotId,
		Status: "updated",
	}

	if err := h.svc.Update(c.Request.Context(), metadata.ApptId, metadata.OldSlotId, updatedAppt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4 Recuperar slot actualizado
	updatedSlot, err := h.slotSvc.GetByID(c.Request.Context(), metadata.NewSlotId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error recuperando slot"})
		return
	}

	// 5. Crear evento
	ssePayload := sse.SSEMessage{
		Type: "appointment_updated",
		Data: domain.Appointment{
			ID:   uint(metadata.ApptId),
			Slot: *updatedSlot,
		},
	}
	log.Println("[ENVIANDO EVENTO ACTUALIZAR SURCHARGE]")

	// 6. Despachar evento
	jsonMsg, _ := json.Marshal(ssePayload)
	h.sseServer.Hub.Broadcast(string(jsonMsg))

	c.JSON(http.StatusOK, updatedAppt)
}

func (h *AppointmentHandler) Reminder(c *gin.Context) {
	key := os.Getenv("EMAIL_APP_PASSWORD")

	m := gomail.NewMessage()
	m.SetAddressHeader("From", "91b38a002@smtp-brevo.com", "El Rodeo Barbería") // o tu correo verificado
	m.SetHeader("Reply-To", "reservas@tubarberia.com")                          // opcional
	m.SetHeader("To", "pereyraezequiel15617866@outlook.es")
	m.SetHeader("Subject", "¡Recordatorio de tu cita!")
	m.SetBody("text/plain", "Hola, recordá que tenés una cita programada.")
	m.AddAlternative("text/html", `
		<!DOCTYPE html>
		<html>
		<body style="font-family: Arial; padding: 10px;">
			<h3>Hola,</h3>
			<p>Este es un recordatorio de tu cita con <strong>El Rodeo Barbería</strong>.</p>
			<p>📅 <strong>Hora:</strong> 15:00 hs</p>
			<p>¡Te esperamos!</p>
		</body>
		</html>`)

	d := gomail.NewDialer("smtp-relay.brevo.com", 587, "91b38a002@smtp-brevo.com", key)
	d.SSL = false
	d.TLSConfig = &tls.Config{ServerName: "smtp-relay.brevo.com"}

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error al enviar correo:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo enviar el correo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Correo enviado correctamente"})
}

type CreateWithCouponReq struct {
	CouponCode string `json:"coupon_code"`
	Items      []uint `json:"items"`
}

func (h *AppointmentHandler) CreateWithCoupon(c *gin.Context) {
	// TODO: Implementar la creación de cita con cupón
}
