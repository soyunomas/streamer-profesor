package servidor

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"emisor-pantalla/internal/hub"
	"emisor-pantalla/web" // Importamos el paquete web que acabamos de crear

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Iniciar configura las rutas y arranca el servidor HTTP.
func Iniciar(h *hub.Hub, puerto string) {
	// Servir el archivo estático index.html desde la memoria embebida
	http.HandleFunc("/", serveIndex)

	// Manejador para la conexión WebSocket
	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(h, w, r)
	})

	localIP, err := getOutboundIP()
	if err != nil {
		localIP = "localhost"
	}

	listenAddr := fmt.Sprintf(":%s", puerto)
	log.Printf("Servidor iniciado. Abre http://%s%s en el navegador del alumno.", localIP, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

// handleWebSocket gestiona la conexión WebSocket de un cliente.
func handleWebSocket(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al actualizar a WebSocket: %v", err)
		return
	}
	defer conn.Close()

	h.AddClient(conn)
	defer h.RemoveClient(conn)

	// Mantener la conexión abierta hasta que el cliente se desconecte
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

// serveIndex sirve el archivo index.html incrustado en el binario.
func serveIndex(w http.ResponseWriter, r *http.Request) {
	// Leemos el archivo directamente de la variable Assets del paquete web
	content, err := web.Assets.ReadFile("static/index.html")
	if err != nil {
		log.Printf("Error leyendo index.html: %v", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No se encontró el archivo index.html"))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
}

// getOutboundIP obtiene la IP local no-loopback.
func getOutboundIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && ip.IsGlobalUnicast() && ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}
	return "", errors.New("no se encontró una IP de red local adecuada")
}
