package sse

import (
	"log"
	"sync"
)

type ClientChan chan string

type Hub struct {
	clients       map[ClientChan]bool
	newClients    chan ClientChan
	closedClients chan ClientChan
	broadcast     chan string
	mu            sync.Mutex
}

func NewHub() *Hub {
	h := &Hub{
		clients:       make(map[ClientChan]bool),
		newClients:    make(chan ClientChan),
		closedClients: make(chan ClientChan),
		broadcast:     make(chan string),
	}
	go h.listen()
	return h
}

func (h *Hub) listen() {
	for {
		select {
		case client := <-h.newClients:
			h.clients[client] = true
			log.Printf("[SSE] Nuevo cliente conectado (%d total)", len(h.clients))
		case client := <-h.closedClients:
			delete(h.clients, client)
			close(client)
			log.Printf("[SSE] Cliente desconectado (%d total)", len(h.clients))
		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client <- msg:
				default:
					// Si el cliente no puede recibir, lo ignoramos
				}
			}
		}
	}
}

func (h *Hub) Broadcast(message string) {
	h.broadcast <- message
}

func (h *Hub) Register(client ClientChan) {
	h.newClients <- client
}

func (h *Hub) Unregister(client ClientChan) {
	h.closedClients <- client
}
