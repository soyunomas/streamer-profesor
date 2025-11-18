package main

import (
	"flag"
	"fmt"

	"emisor-pantalla/internal/captura"
	"emisor-pantalla/internal/hub"
	"emisor-pantalla/internal/servidor"
)

func main() {
	// 1. Configuración de parámetros de línea de comandos
	port := flag.Int("port", 8080, "Puerto en el que el servidor escuchará")
	quality := flag.Int("quality", 50, "Calidad de imagen JPEG (1-100). Menor valor = mayor velocidad.")
	fps := flag.Int("fps", 15, "Tasa de cuadros por segundo (FPS) objetivo.")
	
	flag.Parse()

	// 2. Creación del Hub de WebSockets
	h := hub.NewHub()

	// 3. Inicio de la captura de pantalla en una goroutine
	// Pasamos la calidad y los FPS configurados
	go captura.IniciarCaptura(h, *quality, *fps)

	// 4. Inicio del servidor web (esta llamada es bloqueante)
	servidor.Iniciar(h, fmt.Sprintf("%d", *port))
}
