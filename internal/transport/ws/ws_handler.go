package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	Hub *Hub
}

func NewWSHandler(hub *Hub) *WSHandler {
	return &WSHandler{Hub: hub}
}

func (h *WSHandler) HandleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("[ws error]", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "upgrade failed"})
		return
	}

	// Registrar cliente en el hub
	log.Println("[NEW CLIENT CONNECTED]")
	h.Hub.Register(conn)

	// Limpieza de los recursos
	defer func() {
		h.Hub.Unregister(conn)
		log.Println("[CLIENT DISCONNECTED]")
	}()

	// Bucle de lectura para detectar desconexiones
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
