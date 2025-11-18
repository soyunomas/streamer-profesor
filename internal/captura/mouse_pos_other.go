//go:build !linux || (linux && !amd64)

package captura

// GetMousePos es un fallback para sistemas no-Linux.
// Devuelve una posición fuera de la pantalla para que el cursor no se dibuje.
func GetMousePos() (int, int) {
    return -1, -1
}
