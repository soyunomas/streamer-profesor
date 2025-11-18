package hub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub mantiene el conjunto de clientes activos y transmite mensajes a los clientes.
type Hub struct {
	clients map[*websocket.Conn]bool
	mutex   sync.Mutex
}

// NewHub crea una nueva instancia de Hub.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

// AddClient registra un nuevo cliente en el hub.
func (h *Hub) AddClient(conn *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.clients[conn] = true
	log.Println("Nuevo cliente conectado. Total:", len(h.clients))
}

// RemoveClient elimina un cliente del hub.
func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	// Comprobar si el cliente existe antes de borrar para evitar pánicos si se llama dos veces
	if _, ok := h.clients[conn]; ok {
		delete(h.clients, conn)
		log.Println("Cliente desconectado. Total:", len(h.clients))
	}
}

// Broadcast envía un mensaje a todos los clientes conectados.
func (h *Hub) Broadcast(message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.clients {
		if err := client.WriteMessage(websocket.BinaryMessage, message); err != nil {
			log.Printf("Error de escritura, desconectando cliente: %v", err)
			// Es seguro cerrar y eliminar el cliente aquí mismo mientras se itera
			client.Close()
			delete(h.clients, client)
		}
	}
}
