package sse

import (
	"io"

	"github.com/gin-gonic/gin"
)

type SSEHandler struct {
	Hub *Hub
}

func NewSSEHandler(hub *Hub) *SSEHandler {
	return &SSEHandler{Hub: hub}
}

func (h *SSEHandler) Handle(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	clientChan := make(ClientChan)
	h.Hub.Register(clientChan)

	// Manejar desconexión
	notify := c.Writer.CloseNotify()
	go func() {
		<-notify
		h.Hub.Unregister(clientChan)
		for range clientChan {
			// drenamos el canal
		}
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-clientChan; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
