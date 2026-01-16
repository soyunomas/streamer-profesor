package main

import (
	"flag"
	"fmt"

	"emisor-pantalla/internal/captura"
	"emisor-pantalla/internal/hub"
	"emisor-pantalla/internal/servidor"
)

func main() {
	port := flag.Int("port", 8080, "Puerto del servidor")
	quality := flag.Int("quality", 50, "Calidad JPEG (1-100)")
	fps := flag.Int("fps", 15, "FPS objetivo")
	pin := flag.String("pin", "", "PIN de administración (si se omite, cualquiera puede controlar)") // NUEVO
	
	flag.Parse()

	if *pin == "" {
		fmt.Println("⚠️ ADVERTENCIA: No se ha establecido un PIN (-pin 1234).")
		fmt.Println("Cualquier usuario podrá pausar la transmisión.")
	}

	h := hub.NewHub()

	// Pasamos el PIN al capturador para que gestione los permisos
	go captura.IniciarCaptura(h, *quality, *fps, *pin)

	servidor.Iniciar(h, fmt.Sprintf("%d", *port))
}
