// Package preferences guarda o tema e o idioma do app em config.json: valida os valores, aplica
// os padrões da primeira execução (RN-09, RN-10) e expõe PUT /api/preferences.
package preferences

import "slices"

// Valores válidos de Preferences.Theme.
const (
	ThemeLight  = "light"
	ThemeDark   = "dark"
	ThemeSystem = "system"
)

// Valores válidos de Preferences.Language.
const (
	LanguagePtBR = "pt-br"
	LanguageEn   = "en"
)

// DefaultTheme é o tema da primeira execução (RN-10).
const DefaultTheme = ThemeSystem

var (
	themes    = []string{ThemeLight, ThemeDark, ThemeSystem}
	languages = []string{LanguagePtBR, LanguageEn}
)

// Preferences são as escolhas do usuário. O JSON é o mesmo da API e do config.json.
type Preferences struct {
	Theme    string `json:"theme"`
	Language string `json:"language"`
}

func validTheme(theme string) bool { return slices.Contains(themes, theme) }

func validLanguage(language string) bool { return slices.Contains(languages, language) }
