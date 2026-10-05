package desktop

import (
	"strings"

	"golang.org/x/sys/windows"
)

// ShowFatal mostra um diálogo nativo de erro com o botão OK e espera o usuário fechá-lo. É para
// falhas antes de existir janela do app (RN-04).
func ShowFatal(title, message string) {
	// Um NUL no texto impediria a conversão para UTF-16; o texto vem do catálogo e de caminhos,
	// então basta removê-lo.
	text := windows.StringToUTF16Ptr(strings.ReplaceAll(message, "\x00", ""))
	caption := windows.StringToUTF16Ptr(strings.ReplaceAll(title, "\x00", ""))
	// Sem retorno útil: com o diálogo já fechado ou falho, o app encerra do mesmo jeito.
	_, _ = windows.MessageBox(0, text, caption, windows.MB_OK|windows.MB_ICONERROR)
}
