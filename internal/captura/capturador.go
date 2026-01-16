package captura

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"sync"
	"time"

	"emisor-pantalla/internal/hub"

	"github.com/kbinani/screenshot"
)

// WSMessage define la estructura de los comandos que vienen del WebSocket
type WSMessage struct {
	Type    string `json:"type"`    // "AUTH", "CMD"
	Payload string `json:"payload"` // PIN o Comando (PAUSE/RESUME)
}

// ServerState gestiona el estado global de la transmisión
type ServerState struct {
	IsPaused             bool
	LastFrame            []byte
	AdminPIN             string
	AuthenticatedClients map[*hub.Client]bool // Controlamos quién ha metido el PIN
	mu                   sync.Mutex
}

var state ServerState

// IniciarCaptura comienza el bucle de captura y transmisión.
func IniciarCaptura(h *hub.Hub, calidad int, fps int, pin string) {
	if calidad < 1 {
		calidad = 1
	}
	if calidad > 100 {
		calidad = 100
	}
	if fps < 1 {
		fps = 1
	}

	// Inicializar estado
	state = ServerState{
		AdminPIN:             pin,
		AuthenticatedClients: make(map[*hub.Client]bool),
	}

	fmt.Printf("--> Sistema de Captura iniciado. PIN configurado: '%s'\n", pin)

	// Canal para recibir mensajes desde el Hub (WebSockets de clientes)
	cmdChan := h.GetCommandChannel()

	frameDuration := time.Second / time.Duration(fps)
	opcionesJPEG := &jpeg.Options{Quality: calidad}

	// Iniciar listener de clics (platform-specific)
	go iniciarListenerClics()

	// --- BUCLE PRINCIPAL ---
	for {
		// 1. Procesar Comandos de Clientes (Non-blocking)
		select {
		case msg := <-cmdChan:
			handleCommand(msg, h)
		default:
			// No hay comandos nuevos, seguimos
		}

		// 2. Comprobar Estado de Pausa
		state.mu.Lock()
		paused := state.IsPaused
		lastFrame := state.LastFrame
		state.mu.Unlock()

		if paused {
			// Si estamos pausados, enviamos el último frame periódicamente
			if len(lastFrame) > 0 {
				h.Broadcast(lastFrame)
			}
			// Dormimos para ahorrar CPU drásticamente
			time.Sleep(1 * time.Second)
			continue
		}

		// 3. Captura de Pantalla Normal
		n := screenshot.NumActiveDisplays()
		if n <= 0 {
			time.Sleep(2 * time.Second)
			continue
		}

		bounds := screenshot.GetDisplayBounds(0)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		x, y := GetMousePos()
		isClicked := CheckForRecentClick()

		var imgFinal image.Image = img
		if x >= 0 && y >= 0 {
			imgFinal = drawCursor(img, x, y, isClicked)
		}

		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, imgFinal, opcionesJPEG); err != nil {
			continue
		}

		b := buf.Bytes()

		// Guardar frame en memoria por si pausamos después
		state.mu.Lock()
		state.LastFrame = b
		state.mu.Unlock()

		if buf.Len() > 0 {
			h.Broadcast(b)
		}

		time.Sleep(frameDuration)
	}
}

// handleCommand procesa mensajes JSON (Auth/Pausa)
func handleCommand(msg hub.ClientMessage, h *hub.Hub) {
	var req WSMessage
	// Imprimir el mensaje crudo para depuración
	// fmt.Printf("DEBUG RAW: %s\n", string(msg.Data))

	if err := json.Unmarshal(msg.Data, &req); err != nil {
		// Si falla el unmarshal, probablemente es basura o un intento de hackeo, lo ignoramos
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	// Autenticación
	if req.Type == "AUTH" {
		fmt.Printf("--> Intento de Login. PIN Recibido: '%s' | PIN Correcto: '%s'\n", req.Payload, state.AdminPIN)
		
		// Si el PIN coincide O no hay PIN configurado en el servidor
		if req.Payload == state.AdminPIN || state.AdminPIN == "" {
			state.AuthenticatedClients[msg.Client] = true
			fmt.Println("--> Login EXITOSO. Enviando permisos de Admin.")
			
			// --- CORRECCIÓN AQUÍ: Enviamos el mapa DIRECTAMENTE, sin json.Marshal ---
			// El método SendJSON del Hub ya se encarga de convertirlo a JSON.
			msg.Client.SendJSON(map[string]interface{}{
				"type":   "AUTH_OK",
				"paused": state.IsPaused,
			})
		} else {
			fmt.Println("--> Login FALLIDO. PIN Incorrecto.")
			// --- CORRECCIÓN AQUÍ TAMBIÉN ---
			msg.Client.SendJSON(map[string]string{"type": "AUTH_FAIL"})
		}
		return
	}

	// Verificar permisos para comandos CMD
	if !state.AuthenticatedClients[msg.Client] && state.AdminPIN != "" {
		fmt.Println("--> Comando RECHAZADO: Cliente no autenticado.")
		return 
	}

	// Ejecutar Comandos
	if req.Type == "CMD" {
		if req.Payload == "PAUSE" {
			fmt.Println("--> COMANDO: PAUSAR SERVIDOR")
			state.IsPaused = true
			broadcastStatus(h, true)
		} else if req.Payload == "RESUME" {
			fmt.Println("--> COMANDO: REANUDAR SERVIDOR")
			state.IsPaused = false
			broadcastStatus(h, false)
		}
	}
}

func broadcastStatus(h *hub.Hub, paused bool) {
	// Aquí SÍ usamos json.Marshal porque h.BroadcastJSON espera []byte (texto raw)
	// a diferencia de msg.Client.SendJSON que espera interface{}
	msg, _ := json.Marshal(map[string]interface{}{"type": "STATUS", "paused": paused})
	h.BroadcastJSON(msg)
}

// drawCursor dibuja el cursor y el efecto de clic sobre la imagen.
func drawCursor(img image.Image, x, y int, isClicked bool) image.Image {
	bounds := img.Bounds()
	drawableImg := image.NewRGBA(bounds)
	draw.Draw(drawableImg, bounds, img, image.Point{}, draw.Src)

	colCursor := color.RGBA{255, 0, 0, 255}
	colClick := color.RGBA{255, 120, 0, 200}
	radioCursor := 8

	if isClicked {
		radioEfecto := 25
		grosor := 4
		for r := -radioEfecto; r <= radioEfecto; r++ {
			for c := -radioEfecto; c <= radioEfecto; c++ {
				dist := r*r + c*c
				if dist <= radioEfecto*radioEfecto && dist >= (radioEfecto-grosor)*(radioEfecto-grosor) {
					if x+c >= bounds.Min.X && x+c < bounds.Max.X && y+r >= bounds.Min.Y && y+r < bounds.Max.Y {
						drawableImg.Set(x+c, y+r, colClick)
					}
				}
			}
		}
	}

	for r := -radioCursor; r <= radioCursor; r++ {
		for c := -radioCursor; c <= radioCursor; c++ {
			if r*r+c*c <= radioCursor*radioCursor {
				if x+c >= bounds.Min.X && x+c < bounds.Max.X && y+r >= bounds.Min.Y && y+r < bounds.Max.Y {
					drawableImg.Set(x+c, y+r, colCursor)
				}
			}
		}
	}
	return drawableImg
}
