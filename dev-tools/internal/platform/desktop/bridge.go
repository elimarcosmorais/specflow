// Package desktop reúne o que é específico da janela nativa: o binding único do Wails
// (Bootstrap), o foco e o tema da janela, as mensagens do WebView2, o diálogo nativo de erro
// fatal e o ID da instância única.
package desktop

// BootstrapInfo é o que o front recebe na única chamada pelos bindings do Wails (ADR 004,
// spec §3.3.1).
type BootstrapInfo struct {
	Version        string `json:"version"`
	APIBaseURL     string `json:"apiBaseUrl"`
	Token          string `json:"token"`
	APIStatus      string `json:"apiStatus"` // "ok", "port_in_use" ou "failed"
	APIPort        int    `json:"apiPort"`
	DatabaseStatus string `json:"databaseStatus"` // "available" ou "unavailable"
	LogsDir        string `json:"logsDir"`
	Theme          string `json:"theme"`
	Language       string `json:"language"`
}

// Bridge é a struct registrada em Bind no Wails, que expõe ao front todo método exportado dela.
// Por isso Bootstrap deve ser o único (ADR 004).
type Bridge struct {
	static BootstrapInfo
	prefs  func() (theme, language string)
}

// NewBridge cria o Bridge. static traz os campos fixos da execução; prefs devolve o tema e o
// idioma atuais e é chamada a cada Bootstrap.
func NewBridge(static BootstrapInfo, prefs func() (theme, language string)) *Bridge {
	return &Bridge{static: static, prefs: prefs}
}

// Bootstrap devolve os campos fixos e as preferências atuais.
func (b *Bridge) Bootstrap() BootstrapInfo {
	info := b.static
	info.Theme, info.Language = b.prefs()
	return info
}
