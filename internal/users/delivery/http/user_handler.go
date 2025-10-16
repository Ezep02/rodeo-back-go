package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ezep02/rodeo/internal/auth/domain"
	"github.com/ezep02/rodeo/internal/users/domain/user"
	"github.com/ezep02/rodeo/internal/users/usecase"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc       *usecase.UserService
	cloudinarySvc *usecase.CloudinaryService
}

func NewUserHandler(userSvc *usecase.UserService, cloudinarySvc *usecase.CloudinaryService) *UserHandler {
	return &UserHandler{userSvc, cloudinarySvc}
}

func (h *UserHandler) Update(c *gin.Context) {
	var (
		idStr      = c.Param("id")
		auth_token = os.Getenv("AUTH_TOKEN")
		reqBody    *user.User
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
		log.Println("Error binding JSON:", err)
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// 5. Actualizar el usuario
	if err := h.userSvc.Update(c.Request.Context(), &user.User{
		ID:           uint(id),
		Name:         reqBody.Name,
		Surname:      reqBody.Surname,
		Email:        reqBody.Email,
		Phone_number: reqBody.Phone_number,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario actualizado correctamente",
		"user":    reqBody,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {

	var (
		idStr      = c.Param("id")
		auth_token = os.Getenv("AUTH_TOKEN")
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

	// 4. Obtener el usuario por ID
	u, err := h.userSvc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.User{
			ID:           u.ID,
			Name:         u.Name,
			Surname:      u.Surname,
			Email:        u.Email,
			Phone_number: u.Phone_number,
			CreatedAt:    u.CreatedAt,
			UpdatedAt:    u.UpdatedAt,
			Avatar:       u.Avatar,
		},
	})
}

func (h *UserHandler) UserInfo(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	existing, err := jwt.VerifyUserSession(c, auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 4. Obtener el usuario por ID
	user, err := h.userSvc.GetByID(c.Request.Context(), existing.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo el usuario"})
		return
	}

	// TODO: agregar el avatar a la respuesta

	c.JSON(http.StatusOK, domain.User{
		ID:             user.ID,
		Name:           user.Name,
		Surname:        user.Surname,
		Email:          user.Email,
		Phone_number:   user.Phone_number,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		LastNameChange: user.LastNameChange,
		Username:       user.Username,
		Avatar:         user.Avatar,
		Is_admin:       user.Is_admin,
		Is_barber:      user.Is_barber,
	})
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	// Implement the logic to get user by ID
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {

	var (
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	// 1. Validar la sesion del usuario
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	existingUser, err := jwt.VerfiySessionToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	// Obtener el archivo desde la request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se recibió ningún archivo o el campo es incorrecto",
		})
		return
	}

	// Abrimos el archivo para leerlo
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo abrir el archivo",
		})
		return
	}
	defer openedFile.Close()

	// Crear un nombre unico para el archivo
	filename := fmt.Sprintf("avatar_%s-:user_id:%d", file.Filename, existingUser.ID)

	secure_url, err := h.cloudinarySvc.UploadAvatar(
		c.Request.Context(),
		openedFile,
		filename,
	)

	if err != nil {
		log.Println("ERROR CLOUDINARY UPLOAD:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Error al subir el archivo %s", file.Filename),
		})
		return
	}

	go func(secure_url string) {

		ctx := context.Background()
		// Actualizar el avatar del usuario en la base de datos
		if err := h.userSvc.UpdateAvatar(ctx, secure_url, existingUser.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error actualizando el avatar del usuario en la base de datos",
			})
			return
		}

	}(secure_url)

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Avatar subido correctamente",
		"avatar":  secure_url,
	})
}

type UpdateUsernameRequest struct {
	NewUsername string `json:"new_username" binding:"required"`
}

func (h *UserHandler) UpdateUsername(c *gin.Context) {

	var (
		idStr      = c.Param("id")
		auth_token = os.Getenv("AUTH_TOKEN")
		reqBody    UpdateUsernameRequest
	)

	// 1. Validar la sesion del usuario
	cookie, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autorizado"})
		return
	}

	if _, err := jwt.VerfiySessionToken(cookie); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalido o expirado"})
		return
	}

	// 2. Vlidar el id
	if idStr == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
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

	// 5. Actualizar el username del usuario
	if err := h.userSvc.UpdateUsername(c.Request.Context(), reqBody.NewUsername, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message":  "Avatar subido correctamente",
		"username": reqBody.NewUsername,
	})
}
