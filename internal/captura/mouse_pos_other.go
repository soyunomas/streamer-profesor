//go:build !linux || (linux && !amd64)

package captura

// GetMousePos es un fallback para todos los sistemas que no son linux/amd64.
// Esto incluye Windows, macOS y arquitecturas Linux como ARM (Raspberry Pi).
// Devuelve una posición fuera de la pantalla para que el cursor no se dibuje.
func GetMousePos() (int, int) {
	return -1, -1
}
