package desktop

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Window age sobre a janela do Wails. O contexto do Wails só existe depois do Startup; antes
// disso os métodos não fazem nada. O contexto é protegido porque ApplyTheme é chamado pelas
// goroutines da API.
type Window struct {
	mu  sync.RWMutex
	ctx context.Context
}

// Startup guarda o contexto do Wails (usar como OnStartup).
func (w *Window) Startup(ctx context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.ctx = ctx
}

func (w *Window) context() context.Context {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.ctx
}

// Focus traz a janela para a frente, restaurando-a se estiver minimizada (usar na segunda
// instância, RN-02).
func (w *Window) Focus() {
	ctx := w.context()
	if ctx == nil {
		return
	}
	runtime.WindowUnminimise(ctx)
	runtime.WindowShow(ctx)
}

// ApplyTheme aplica o tema ("light", "dark" ou qualquer outro valor para o do sistema) à janela.
func (w *Window) ApplyTheme(theme string) {
	ctx := w.context()
	if ctx == nil {
		return
	}
	switch theme {
	case "light":
		runtime.WindowSetLightTheme(ctx)
	case "dark":
		runtime.WindowSetDarkTheme(ctx)
	default:
		runtime.WindowSetSystemDefaultTheme(ctx)
	}
}

// WindowsTheme converte o tema das preferências no tema do Wails para Windows.
func WindowsTheme(theme string) windows.Theme {
	switch theme {
	case "light":
		return windows.Light
	case "dark":
		return windows.Dark
	default:
		return windows.SystemDefault
	}
}
