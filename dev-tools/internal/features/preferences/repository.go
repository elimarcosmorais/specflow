package preferences

import "errors"

// ErrNotFound é o erro de Repository.Load quando ainda não há preferências guardadas (primeira
// execução). Qualquer outro erro de Load significa conteúdo ilegível ou inválido.
var ErrNotFound = errors.New("preferences not found")

// Repository guarda as preferências. Load devolve os valores como estão guardados, sem
// validá-los: quem valida campo a campo é o Service.
type Repository interface {
	Load() (Preferences, error)
	Save(Preferences) error
}
