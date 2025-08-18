package http

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ezep02/rodeo/internal/domain"
	"github.com/ezep02/rodeo/internal/service"
	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postSvc *service.PostService
}

func NewPostHandler(postSvc *service.PostService) *PostHandler {
	return &PostHandler{postSvc}
}

type CreatePostRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	PreviewImage string `json:"preview_url"`
}

func (h *PostHandler) Create(c *gin.Context) {
	var (
		req        CreatePostRequest
		auth_token = os.Getenv("AUTH_TOKEN")
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	log.Println("Creating post with title:", req)

	// 2. Verificar sesion
	sessionToken, err := c.Cookie(auth_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sesion expirada o token invalido",
		})
		return
	}

	exitingUser, err := jwt.VerfiySessionToken(sessionToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 3. Crear el post
	post := &domain.Post{
		UserID:      exitingUser.ID,
		Title:       req.Title,
		PreviewUrl:  req.PreviewImage,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err = h.postSvc.Create(c.Request.Context(), post); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}

func (h *PostHandler) List(c *gin.Context) {

	// 1. Recuperar el offset de la query
	offset := c.Query("offset")
	if offset == "" {
		offset = "0" // Valor por defecto si no se proporciona
	}

	// 2. Convertir offset a int
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
		return
	}

	posts, err := h.postSvc.List(c.Request.Context(), offsetInt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func (h *PostHandler) GetByID(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(400, gin.H{"error": "Post ID is required"})
		return
	}

	// Convertir postID a int
	paymentID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de pago inválido"})
		return
	}

	post, err := h.postSvc.GetByID(c.Request.Context(), uint(paymentID))
	if err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func (h *PostHandler) Update(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(400, gin.H{"error": "Post ID is required"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Convertir postID a int
	paymentID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post := &domain.Post{
		ID:          uint(paymentID),
		Title:       req.Title,
		Description: req.Description,
		PreviewUrl:  req.PreviewImage,
	}

	if err = h.postSvc.Update(c.Request.Context(), post, uint(paymentID)); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
	})
}

func (h *PostHandler) Delete(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(400, gin.H{"error": "Post ID es requerido"})
		return
	}

	// Convertir postID a int
	paymentID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post ID invalid"})
		return
	}

	if err := h.postSvc.Delete(c.Request.Context(), uint(paymentID)); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post eliminado correctamente",
	})
}

func (h *PostHandler) Count(c *gin.Context) {
	count, err := h.postSvc.Count(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to count posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_post": count,
	})
}
