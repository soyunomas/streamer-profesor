//go:build !linux || (linux && !amd64)

package captura

// iniciarListenerClics es una función vacía para compatibilidad de compilación.
func iniciarListenerClics() {
	// No hacemos nada en sistemas no soportados por el hook.
}

// CheckForRecentClick siempre devuelve false en sistemas no soportados.
func CheckForRecentClick() bool {
	return false
}
