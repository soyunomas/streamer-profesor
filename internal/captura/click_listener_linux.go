//go:build linux && amd64

package captura

import (
	"sync"
	"time"

	hook "github.com/robotn/gohook"
)

// Estas variables solo existirán en la compilación para linux/amd64
var (
	lastClickTime time.Time
	clickMutex    sync.Mutex
)

// iniciarListenerClics inicia el hook de eventos real en linux/amd64.
func iniciarListenerClics() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		if ev.Kind == hook.MouseDown && ev.Button == 1 {
			clickMutex.Lock()
			lastClickTime = time.Now()
			clickMutex.Unlock()
		}
	}
}

// CheckForRecentClick verifica si hubo un clic reciente.
func CheckForRecentClick() bool {
	clickMutex.Lock()
	defer clickMutex.Unlock()
	return time.Since(lastClickTime) < 300*time.Millisecond
}
