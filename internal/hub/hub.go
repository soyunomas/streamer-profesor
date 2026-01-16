package hub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client envuelve la conexión WebSocket para poder identificarla
type Client struct {
	Conn *websocket.Conn
	Hub  *Hub
}

// SendJSON envía un objeto JSON al cliente (helper)
func (c *Client) SendJSON(v interface{}) {
	c.Conn.WriteJSON(v)
}

// ClientMessage estructura para pasar mensajes del socket al motor de captura
type ClientMessage struct {
	Client *Client
	Data   []byte
}

// Hub mantiene el conjunto de clientes activos.
type Hub struct {
	clients map[*Client]bool
	cmdChan chan ClientMessage // Canal para comandos (Pause, Auth)
	mutex   sync.Mutex
}

// NewHub crea una nueva instancia de Hub.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		cmdChan: make(chan ClientMessage, 10), // Buffer pequeño
	}
}

// GetCommandChannel devuelve el canal para que el capturador escuche
func (h *Hub) GetCommandChannel() chan ClientMessage {
	return h.cmdChan
}

// AddClient registra un nuevo cliente en el hub.
func (h *Hub) AddClient(conn *websocket.Conn) *Client {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	client := &Client{Conn: conn, Hub: h}
	h.clients[client] = true
	log.Println("Nuevo cliente conectado. Total:", len(h.clients))
	return client
}

// RemoveClient elimina un cliente del hub.
func (h *Hub) RemoveClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		client.Conn.Close()
		log.Println("Cliente desconectado. Total:", len(h.clients))
	}
}

// ListenClientMessages escucha mensajes del cliente (JS) y los manda al canal.
// Esta función es bloqueante y se ejecuta en una goroutine por cliente.
func (h *Hub) ListenClientMessages(client *Client) {
	defer h.RemoveClient(client)
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}
		// Enviamos al canal principal para que captura.go lo procese de forma segura
		h.cmdChan <- ClientMessage{Client: client, Data: message}
	}
}

// Broadcast envía un mensaje binario (imagen JPEG) a todos.
func (h *Hub) Broadcast(message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.clients {
		if err := client.Conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
			log.Printf("Error de escritura, cerrando: %v", err)
			client.Conn.Close()
			delete(h.clients, client)
		}
	}
}

// BroadcastJSON envía un mensaje de texto (JSON estado) a todos.
func (h *Hub) BroadcastJSON(message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.clients {
		client.Conn.WriteMessage(websocket.TextMessage, message)
	}
}
