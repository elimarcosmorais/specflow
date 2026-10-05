// Package apperr define o erro de domínio do motor: um tipo, um código estável e, opcionalmente,
// o campo da entrada e a causa. Só depende da biblioteca padrão e não conhece HTTP: quem converte
// Kind em status é o platform/httpserver.
package apperr

// Kind classifica o erro. O httpserver o traduz para o status HTTP.
type Kind string

const (
	BadRequest   Kind = "BadRequest"
	Unauthorized Kind = "Unauthorized"
	NotFound     Kind = "NotFound"
	Conflict     Kind = "Conflict"
	Invalid      Kind = "Invalid"
	Internal     Kind = "Internal"
	Unavailable  Kind = "Unavailable"
)

// Error é o erro de domínio. Code é o código estável da API (ex.: "project.name_taken"), Field é o
// campo da entrada a que o erro se refere (vazio quando não há) e Err é a causa interna, que só
// vai para o log e nunca para o cliente.
type Error struct {
	Kind  Kind
	Code  string
	Field string
	Err   error
}

// Erros comuns a toda feature.
var (
	ErrInvalidBody         = New(BadRequest, "request.invalid_body")
	ErrTokenInvalid        = New(Unauthorized, "auth.token_invalid")
	ErrDatabaseUnavailable = New(Unavailable, "database.unavailable")
	ErrInternal            = New(Internal, "internal.unexpected")
)

// New cria um erro com o tipo e o código dados.
func New(kind Kind, code string) *Error {
	return &Error{Kind: kind, Code: code}
}

// WithField devolve uma cópia do erro com Field preenchido. O original não muda.
func (e *Error) WithField(field string) *Error {
	c := *e
	c.Field = field
	return &c
}

// Wrap devolve uma cópia do erro com a causa. O original não muda.
func (e *Error) Wrap(cause error) *Error {
	c := *e
	c.Err = cause
	return &c
}

// Error devolve o código e, se houver, a causa.
func (e *Error) Error() string {
	if e.Err == nil {
		return e.Code
	}
	return e.Code + ": " + e.Err.Error()
}

// Unwrap devolve a causa, para errors.Is e errors.As.
func (e *Error) Unwrap() error {
	return e.Err
}

// Is compara só o Code: dois *Error com o mesmo código são o mesmo erro, qualquer que seja o
// campo ou a causa.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t != nil && e.Code == t.Code
}
