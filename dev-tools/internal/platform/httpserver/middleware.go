package httpserver

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"dev-tools/internal/apperr"
)

// recorder guarda o status escrito e se a resposta já começou.
type recorder struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (rec *recorder) WriteHeader(status int) {
	if !rec.wrote {
		rec.status = status
		rec.wrote = true
	}
	rec.ResponseWriter.WriteHeader(status)
}

func (rec *recorder) Write(b []byte) (int, error) {
	if !rec.wrote {
		rec.WriteHeader(http.StatusOK)
	}
	return rec.ResponseWriter.Write(b)
}

// Unwrap deixa o http.ResponseController alcançar o ResponseWriter original.
func (rec *recorder) Unwrap() http.ResponseWriter { return rec.ResponseWriter }

// asRecorder reaproveita o recorder que um middleware anterior já criou, ou cria um.
func asRecorder(w http.ResponseWriter) *recorder {
	if rec, ok := w.(*recorder); ok {
		return rec
	}
	return &recorder{ResponseWriter: w}
}

// recoverPanic transforma o pânico de um handler em 500 internal.unexpected, registrando o valor
// e a pilha no log.
func recoverPanic(cfg Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := asRecorder(w)
		defer func() {
			p := recover()
			if p == nil {
				return
			}
			if p == http.ErrAbortHandler {
				panic(p) // é como o handler pede para abortar a conexão; não é falha
			}
			err := apperr.ErrInternal.Wrap(fmt.Errorf("panic: %v", p))
			if rec.wrote {
				// A resposta já começou: não dá para trocar por um 500, só registrar.
				cfg.Logger.Error("http.error", "method", r.Method, "path", r.URL.Path,
					"status", rec.status, "code", apperr.ErrInternal.Code, "err", err,
					"stack", string(debug.Stack()))
				return
			}
			writeError(rec, r, cfg, err, "stack", string(debug.Stack()))
		}()
		next.ServeHTTP(rec, r)
	})
}

// logAccess registra http.request (DEBUG) ao fim de cada requisição. Nunca loga o token, o
// cabeçalho Authorization nem corpos.
func logAccess(cfg Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := asRecorder(w)
		start := time.Now()
		next.ServeHTTP(rec, r)
		status := rec.status
		if !rec.wrote {
			status = http.StatusOK // o net/http responde 200 quando o handler não escreve nada
		}
		cfg.Logger.Debug("http.request", "method", r.Method, "path", r.URL.Path,
			"status", status, "duration_ms", time.Since(start).Milliseconds())
	})
}

// cors acrescenta os cabeçalhos de CORS (spec §3.1) a toda resposta e atende o preflight
// (OPTIONS) com 204, sem token e sem chamar o handler. O CORS aberto é seguro porque o token
// viaja em cabeçalho, nunca em cookie.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		h.Set("Access-Control-Allow-Private-Network", "true")
		h.Set("Access-Control-Max-Age", "600")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// requireToken exige "Authorization: Bearer <token>", comparado em tempo constante. Falha
// responde 401 auth.token_invalid sem chamar o handler. Um token configurado vazio recusa tudo.
func requireToken(cfg Config, next http.Handler) http.Handler {
	want := []byte(cfg.Token)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || len(want) == 0 || subtle.ConstantTimeCompare([]byte(got), want) != 1 {
			writeError(w, r, cfg, apperr.ErrTokenInvalid)
			return
		}
		next.ServeHTTP(w, r)
	})
}
