package middleware

import (
	"net/http"
	"os"

	"github.com/ezep02/rodeo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// Middleware encargado de verificar si el usuario es admin o no
func AuthorizeAdmin() gin.HandlerFunc {

	var auth_token = os.Getenv("AUTH_TOKEN")

	return func(c *gin.Context) {

		sessionToken, err := c.Cookie(auth_token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "sesion expirada o token invalido",
			})
			return
		}

		session, err := jwt.VerfiySessionToken(sessionToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if !session.Is_admin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Acceso solo permitido para administradores"})
			c.Abort()
			return
		}
		c.Next()
	}
}
