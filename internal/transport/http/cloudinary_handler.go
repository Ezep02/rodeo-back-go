package http

import (
	"net/http"

	"github.com/ezep02/rodeo/internal/service"
	"github.com/gin-gonic/gin"
)

type CloudinaryHandler struct {
	svc *service.CloudinaryService
}

func NewCloudinaryHandler(svc *service.CloudinaryService) *CloudinaryHandler {
	return &CloudinaryHandler{svc}
}

func (h *CloudinaryHandler) Images(c *gin.Context) {

	// Llama al Admin API para obtener imágenes
	res, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching images from Cloudinary",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assets": res,
	})
}

func (h *CloudinaryHandler) Video(c *gin.Context) {

	// Llama al Admin API para obtener imágenes
	res, err := h.svc.Video(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching images from Cloudinary",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assets": res,
	})
}
