//go:build linux && amd64

package captura

import "github.com/go-vgo/robotgo"

// GetMousePos usa robotgo para obtener la posición real.
// Esta versión SOLO se compila en sistemas Linux de 64 bits (amd64).
func GetMousePos() (int, int) {
	return robotgo.GetMousePos()
}
