//go:build linux && amd64

package captura

import "github.com/go-vgo/robotgo"

// GetMousePos usa robotgo para obtener la posición real en Linux.
func GetMousePos() (int, int) {
    return robotgo.GetMousePos()
}
