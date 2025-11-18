package captura

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"time"

	"emisor-pantalla/internal/hub"

	"github.com/kbinani/screenshot"
)

// IniciarCaptura comienza el bucle de captura y transmisión.
// Ahora es agnóstico a la plataforma y no importa 'gohook'.
func IniciarCaptura(h *hub.Hub, calidad int, fps int) {
	if calidad < 1 { calidad = 1 }
	if calidad > 100 { calidad = 100 }
	if fps < 1 { fps = 1 }

	frameDuration := time.Second / time.Duration(fps)
	opcionesJPEG := &jpeg.Options{Quality: calidad}

	// Esta llamada usará la versión real o la falsa dependiendo de la plataforma.
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

		x, y := GetMousePos()
		
		// Esta llamada usará la versión real o la falsa.
		isClicked := CheckForRecentClick()

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

// drawCursor no necesita cambios, ya recibe el booleano 'isClicked'.
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
