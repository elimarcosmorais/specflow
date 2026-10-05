package i18n

import "golang.org/x/sys/windows"

// DetectWindows devolve o idioma do motor para o idioma de exibição preferido do usuário no
// Windows (o primeiro de GetUserPreferredUILanguages). Em erro ou lista vazia devolve "en".
func DetectWindows() string {
	languages, err := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil || len(languages) == 0 {
		return langEn
	}
	return FromLocale(languages[0])
}
