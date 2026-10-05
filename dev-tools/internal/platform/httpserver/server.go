// Package httpserver é o servidor HTTP comum do motor: registro de rotas com token por padrão,
// CORS, recuperação de pânico, log de acesso, resposta de erro traduzida (ADR 005), Listen só
// em loopback e a página /docs. As features só declaram as próprias rotas.
package httpserver

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strconv"

	"golang.org/x/sys/windows"
)

// HandlerFunc é um handler que devolve o erro em vez de escrevê-lo: o httpserver o converte
// em resposta traduzida. Devolver nil significa que o handler já escreveu a resposta.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Route é uma rota da API. Pattern não leva o método (ex.: "/api/projects/{id}"). Toda rota
// exige o token, exceto as marcadas como Public.
type Route struct {
	Method, Pattern string
	Public          bool
	Handler         HandlerFunc
}

// Config reúne o que o servidor precisa de fora. Language devolve o idioma atual a cada
// chamada, para a mensagem de erro seguir as preferências no momento da requisição.
type Config struct {
	Token    string
	Language func() string
	Logger   *slog.Logger
}

// ErrPortInUse indica que a porta pedida a Listen já está em uso.
var ErrPortInUse = errors.New("port already in use")

// NewHandler monta o ServeMux com as rotas (padrões "MÉTODO /caminho") e as rotas de /docs, e o
// envolve nos middlewares: recuperação → log de acesso → CORS → token por rota. O CORS fica fora
// do mux para atender o preflight antes dele (senão um OPTIONS numa rota só de GET/POST viraria
// 405).
func NewHandler(cfg Config, routes []Route) http.Handler {
	if cfg.Language == nil {
		cfg.Language = func() string { return "en" }
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.New(slog.DiscardHandler)
	}

	mux := http.NewServeMux()
	for _, route := range append(slices.Clone(routes), docsRoutes()...) {
		h := adapt(cfg, route.Handler)
		if !route.Public {
			h = requireToken(cfg, h)
		}
		mux.Handle(route.Method+" "+route.Pattern, h)
	}
	return recoverPanic(cfg, logAccess(cfg, cors(mux)))
}

// adapt converte um HandlerFunc em http.Handler, escrevendo o erro devolvido.
func adapt(cfg Config, h HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			writeError(w, r, cfg, err)
		}
	})
}

// NewToken gera o token da execução: 32 bytes de crypto/rand em hexadecimal (64 caracteres).
func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// Listen abre a porta em 127.0.0.1 (a API nunca escuta em outro endereço). Porta ocupada
// devolve um erro em que errors.Is(err, ErrPortInUse) é verdadeiro.
func Listen(port int) (net.Listener, error) {
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		if errors.Is(err, windows.WSAEADDRINUSE) {
			return nil, fmt.Errorf("listen on %s: %w: %w", addr, ErrPortInUse, err)
		}
		return nil, fmt.Errorf("listen on %s: %w", addr, err)
	}
	return ln, nil
}
