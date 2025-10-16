package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ezep02/rodeo/internal/users/usecase"
	"github.com/gin-gonic/gin"
)

type CloudinaryHandler struct {
	svc *usecase.CloudinaryService
}

func NewCloudinaryHandler(svc *usecase.CloudinaryService) *CloudinaryHandler {
	return &CloudinaryHandler{svc}
}

type ImagesReq struct {
	NextCursor string `json:"next_cursor"`
}

func (h *CloudinaryHandler) Images(c *gin.Context) {
	cursor := c.Query("next_cursor")

	res, nextCursor, err := h.svc.List(c.Request.Context(), cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching images from Cloudinary",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assets":      res,
		"next_cursor": nextCursor,
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

func (h *CloudinaryHandler) Upload(c *gin.Context) {

	// Parsear multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error parsing multipart form",
		})
		return
	}

	// Obteener todos los archivos bajo la clave "file"
	files := form.File["file"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No se recibio ningun archivo",
		})
		return
	}

	// Subir archivos
	for _, file := range files {
		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "No se pudo abrir el archivo",
			})
			return
		}
		defer openedFile.Close()

		// Subir a Cloudinary pasando io.Reader directamente
		if err := h.svc.Upload(c.Request.Context(), openedFile, file.Filename); err != nil {
			log.Println("ERROR CLOUD UPLOAD", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Error al subir el archivo: %s a Cloudinary", file.Filename),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
	})
}
