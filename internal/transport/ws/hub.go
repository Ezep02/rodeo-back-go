package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients map[*websocket.Conn]bool
	Mux     sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.Mux.Lock()
	defer h.Mux.Unlock()
	// si todo va bien, marcamos conexion como activa
	h.Clients[conn] = true
}

func (h *Hub) Unregister(conn *websocket.Conn) {

	h.Mux.Lock()
	defer h.Mux.Unlock()

	if _, ok := h.Clients[conn]; ok {
		delete(h.Clients, conn)
		conn.Close()
	}
}

func (h *Hub) Broadcast(message []byte) {

	h.Mux.Lock()
	defer h.Mux.Unlock()

	// Envia el mensaje a todas las conexiones activas
	for conn := range h.Clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			conn.Close()
			delete(h.Clients, conn)
		}
	}
}
