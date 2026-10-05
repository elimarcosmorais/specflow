// Package projects cadastra, lista, edita e remove os projetos do usuário: valida a entrada
// (RN-12 a RN-15), calcula a forma guardada do caminho (D-07) e expõe /api/projects. O núcleo
// (project.go, errors.go, repository.go e service.go) só depende da biblioteca padrão e de
// internal/apperr (ADR 002).
package projects

import (
	"strings"
	"time"
)

// Type é o tipo do projeto.
type Type string

// Valores válidos de Type (RN-12).
const (
	TypeWeb     Type = "web"
	TypeAndroid Type = "android"
	TypeWindows Type = "windows"
)

// Project é o registro guardado. Path está na forma guardada (D-07): relativo à raiz da instalação
// quando a pasta está dentro dela, senão absoluto. NameKey e PathKey são as chaves de unicidade
// (RN-13, RN-14).
type Project struct {
	ID        int64
	Name      string
	NameKey   string
	Path      string
	PathKey   string
	Type      Type
	CreatedAt time.Time
	UpdatedAt time.Time
}

// View é o projeto como sai para fora do núcleo: Path é sempre absoluto e FolderExists é calculado
// a cada leitura (RN-18).
type View struct {
	ID           int64
	Name         string
	Path         string
	Type         Type
	FolderExists bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Input são os dados de criação e de edição de um projeto.
type Input struct {
	Name string
	Path string
	Type Type
}

// NameKey é a chave de unicidade do nome (RN-13): o nome sem diferenciar maiúsculas.
func NameKey(name string) string { return strings.ToLower(name) }

// PathKey é a chave de unicidade da pasta (RN-14), calculada sobre a forma guardada do caminho.
func PathKey(stored string) string { return strings.ToLower(stored) }

func (t Type) valid() bool {
	return t == TypeWeb || t == TypeAndroid || t == TypeWindows
}
