package http

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/service"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
	}
}

type UpdateUserRequest struct {
	user *domain.User
}

func (h *UserHandler) Update(c *gin.Context) {
	var (
		idStr      = c.Param("id")
		auth_token = os.Getenv("AUTH_TOKEN")
		reqBody    UpdateUserRequest
	)

	// 1. Vlidar el id
	if idStr == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}

	// 2. Validar la sesion del usuario
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	if _, err := jwt.VerfiySessionToken(cookie); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	// 3. Parsear el id a uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "El id debe ser valido"})
		return
	}

	// 4. Bindear el cuerpo de la peticion
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// 5. Actualizar el usuario
	if err := h.userSvc.Update(c.Request.Context(), &domain.User{
		ID:      uint(id),
		Name:    reqBody.user.Name,
		Surname: reqBody.user.Surname,
		Email:   reqBody.user.Email,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando los datos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario actualizado correctamente",
		"user":    reqBody.user,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	// Implement the logic to get user by ID
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	// Implement the logic to get user by ID
}
