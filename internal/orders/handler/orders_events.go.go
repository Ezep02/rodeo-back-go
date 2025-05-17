package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	clients = make(map[chan []byte]bool)
	mu      sync.Mutex
)

func PushOrderUpdate(msgBytes []byte) {

	for ch := range clients {
		select {
		case ch <- msgBytes:
		default:
			log.Println("[WARNING] cliente no responde")
		}
	}

}

// Send orders real time data to the client using SSE protocol
func SendOrderEvents(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming no soportado", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan []byte)

	mu.Lock()
	clients[messageChan] = true
	mu.Unlock()

	log.Println("[NEW CLIENT CONNECTED] Total clients:", len(clients))

	defer func() {
		mu.Lock()
		delete(clients, messageChan)
		mu.Unlock()
		close(messageChan)
		log.Println("[CLIENT DISCONNECTED] Total clients:", len(clients))
	}()

	notify := r.Context().Done()

	log.Println("[SENDING ORDER DATA]")
	for {
		select {
		case <-notify:
			log.Println("Client disconnected")
			return
		case msg, ok := <-messageChan:
			if !ok {
				return
			}
			fmt.Fprintf(rw, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}
