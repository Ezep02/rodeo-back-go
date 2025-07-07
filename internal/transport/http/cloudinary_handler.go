package http

import (
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/gin-gonic/gin"
)

func GetCloudinaryImages(c *gin.Context) {

	// Carga variables de entorno
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error initializing Cloudinary",
		})
		return
	}

	// Inicializa Cloudinary
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error initializing Cloudinary",
		})
		return
	}

	// Llama al Admin API para obtener imágenes
	res, err := cld.Admin.Assets(c.Request.Context(), admin.AssetsParams{
		AssetType:  "image",
		MaxResults: 10,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error fetching images from Cloudinary",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assets": res.Assets,
	})
}
