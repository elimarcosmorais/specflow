// Package i18n traduz por chave as mensagens que o motor produz (erros da API e diálogos
// nativos), em pt-br e en (ADR 005).
package i18n

import (
	"slices"
	"strings"
)

const (
	langPtBR = "pt-br"
	langEn   = "en"
)

// T devolve o texto de key no idioma language, trocando cada {chave} pelo valor de params.
// Sem tradução no idioma usa o texto em en; sem en devolve a própria chave.
func T(language, key string, params map[string]string) string {
	text, ok := catalog[language][key]
	if !ok {
		text, ok = catalog[langEn][key]
	}
	if !ok {
		return key
	}
	if len(params) == 0 {
		return text
	}
	pairs := make([]string, 0, len(params)*2)
	for name, value := range params {
		pairs = append(pairs, "{"+name+"}", value)
	}
	return strings.NewReplacer(pairs...).Replace(text)
}

// FromLocale converte um locale do Windows em idioma do motor: pt-br (sem diferenciar
// maiúsculas, com "-" ou "_") vira "pt-br"; qualquer outro valor vira "en".
func FromLocale(locale string) string {
	if strings.EqualFold(strings.ReplaceAll(locale, "_", "-"), "pt-br") {
		return langPtBR
	}
	return langEn
}

// Keys devolve, em ordem alfabética, as chaves do catálogo do idioma. Serve ao teste de
// paridade entre os idiomas.
func Keys(language string) []string {
	keys := make([]string, 0, len(catalog[language]))
	for key := range catalog[language] {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}
