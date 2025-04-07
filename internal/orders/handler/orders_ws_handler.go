package handler

import (
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

// WEBSOCKET
// Peer estructura para manejar una conexión peer-to-peer
type Peer struct {
	connection []*websocket.Conn
	mu         sync.Mutex
}

var peer Peer

// Configuración del upgrader de WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true }, // Permitir todas las conexiones (ajusta según sea necesario)
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleConnection gestiona conexión WebSocket
func HandleConnection(rw http.ResponseWriter, r *http.Request) {
	// Actualizar a WebSocket
	ws, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println("Error al actualizar la conexión:", err.Error())
		return
	}
	defer ws.Close()

	// Registrar la conexión como la conexión activa
	peer.mu.Lock()
	peer.connection = append(peer.connection, ws)
	peer.mu.Unlock()

	log.Println("[ORDERS] Nueva conexión establecida")

	// Leer mensajes del cliente y reenviar directamente al peer
	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error leyendo mensaje:", err.Error())
			break
		}

		// Reenviar el mensaje al peer
		err = sendMessageToPeer(messageType, msg)
		if err != nil {
			log.Println("Error enviando datos actualizados:", err.Error())
			break
		}
	}

	// Al cerrar, eliminar la conexión activa
	removeConnection(ws)
	log.Println("[ORDERS] Conexion cerrada")
}

// removeConnection elimina conexiones cerradas
func removeConnection(conn *websocket.Conn) {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	for i, c := range peer.connection {
		if c == conn {
			peer.connection = slices.Delete(peer.connection, i, i+1)
			break
		}
	}
}

// sendUpdatedData envía mediante ws
func sendMessageToPeer(messageType int, msg []byte) error {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	if len(peer.connection) == 0 {
		log.Println("No hay peers conectados para recibir el mensaje schedules")
		return nil
	}

	var activeConnections []*websocket.Conn
	for _, conn := range peer.connection {
		err := conn.WriteMessage(messageType, msg)
		if err != nil {
			conn.Close() // Cerrar la conexion fallida
			continue     // continuar con el proceso
		}
		activeConnections = append(activeConnections, conn)
	}

	// Actualizar la lista de conexiones activas
	peer.connection = activeConnections
	return nil
}
