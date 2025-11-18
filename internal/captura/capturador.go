package captura

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"sync"
	"time"

	"emisor-pantalla/internal/hub"

	"github.com/kbinani/screenshot"
	hook "github.com/robotn/gohook"
)

// Variables para gestionar el estado del clic
var (
	lastClickTime time.Time
	clickMutex    sync.Mutex
)

// IniciarCaptura comienza el bucle de captura y transmisión.
func IniciarCaptura(h *hub.Hub, calidad int, fps int) {
	// 1. Validaciones
	if calidad < 1 { calidad = 1 }
	if calidad > 100 { calidad = 100 }
	if fps < 1 { fps = 1 }

	frameDuration := time.Second / time.Duration(fps)
	opcionesJPEG := &jpeg.Options{Quality: calidad}

	// 2. Iniciamos el "oyente" de clics en segundo plano
	go iniciarListenerClics()

	for {
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
		
		// Intentamos obtener posición del mouse
		x, y := GetMousePos()
		
		// Verificamos si hubo un clic hace poco (menos de 300ms)
		clickMutex.Lock()
		isClicked := time.Since(lastClickTime) < 300*time.Millisecond
		clickMutex.Unlock()

		var imgFinal image.Image = img
		if x >= 0 && y >= 0 {
			imgFinal = drawCursor(img, x, y, isClicked)
		}
		
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, imgFinal, opcionesJPEG); err != nil {
			continue
		}

		if buf.Len() > 0 {
			h.Broadcast(buf.Bytes())
		}

		time.Sleep(frameDuration)
	}
}

// iniciarListenerClics escucha los eventos globales del sistema.
func iniciarListenerClics() {
	// Iniciamos el hook
	evChan := hook.Start()
	// CORRECCIÓN: Usamos hook.End() en lugar de hook.Stop()
	defer hook.End()

	// Escuchamos el canal de eventos
	for ev := range evChan {
		// Hook.MouseDown significa botón presionado. Button 1 es el izquierdo.
		if ev.Kind == hook.MouseDown && ev.Button == 1 {
			clickMutex.Lock()
			lastClickTime = time.Now()
			clickMutex.Unlock()
		}
	}
}

// drawCursor dibuja el cursor y un efecto si hay clic.
func drawCursor(img image.Image, x, y int, isClicked bool) image.Image {
	bounds := img.Bounds()
	drawableImg := image.NewRGBA(bounds)
	draw.Draw(drawableImg, bounds, img, image.Point{}, draw.Src)

	// Colores
	colCursor := color.RGBA{255, 0, 0, 255}       // Rojo sólido
	colClick := color.RGBA{255, 120, 0, 200}      // Amarillo semi-transparente
	
	radioCursor := 8

	// 1. Dibujar el efecto de clic (Anillo exterior)
	if isClicked {
		radioEfecto := 25 // Radio grande para la onda
		grosor := 4       // Grosor del anillo

		for r := -radioEfecto; r <= radioEfecto; r++ {
			for c := -radioEfecto; c <= radioEfecto; c++ {
				dist := r*r + c*c
				// Fórmula para hacer un anillo hueco
				if dist <= radioEfecto*radioEfecto && dist >= (radioEfecto-grosor)*(radioEfecto-grosor) {
					if x+c >= bounds.Min.X && x+c < bounds.Max.X && y+r >= bounds.Min.Y && y+r < bounds.Max.Y {
						drawableImg.Set(x+c, y+r, colClick)
					}
				}
			}
		}
	}

	// 2. Dibujar el cursor normal (Círculo central)
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
