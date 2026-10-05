package httpserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"dev-tools/internal/apperr"
	"dev-tools/internal/platform/i18n"
)

// maxBodyBytes é o tamanho máximo do corpo de uma requisição (1 MiB).
const maxBodyBytes = 1 << 20

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// WriteJSON escreve v como JSON com o status dado. Um v que não serializa é erro de programação:
// vira pânico, e a recuperação responde 500.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body) // falha aqui é o cliente que foi embora; nada a fazer
}

// DecodeJSON lê o corpo em dst, com limite de 1 MiB, recusando campo desconhecido e conteúdo
// depois do valor. Qualquer falha devolve apperr.ErrInvalidBody com a causa.
func DecodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(http.MaxBytesReader(nil, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return apperr.ErrInvalidBody.Wrap(err)
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return apperr.ErrInvalidBody.Wrap(errors.New("unexpected data after the JSON value"))
	}
	return nil
}

// writeError responde err no formato {"error":{"code","message","field"}}, com a mensagem
// traduzida para o idioma atual e o status do Kind, e registra http.error: WARN em 4xx, ERROR em
// 5xx (só aí com a causa). Erro que não é apperr vira 500 internal.unexpected. A causa nunca vai
// ao cliente. extra são atributos acrescentados ao log.
func writeError(w http.ResponseWriter, r *http.Request, cfg Config, err error, extra ...any) {
	var ae *apperr.Error
	if !errors.As(err, &ae) {
		ae = apperr.ErrInternal.Wrap(err)
		err = ae
	}
	status := statusFor(ae.Kind)

	attrs := []any{"method", r.Method, "path", r.URL.Path, "status", status, "code", ae.Code}
	if status >= http.StatusInternalServerError {
		cfg.Logger.Error("http.error", append(append(attrs, "err", err), extra...)...)
	} else {
		cfg.Logger.Warn("http.error", append(attrs, extra...)...)
	}

	WriteJSON(w, status, errorBody{Error: errorDetail{
		Code:    ae.Code,
		Message: i18n.T(cfg.Language(), ae.Code, nil),
		Field:   ae.Field,
	}})
}

// statusFor converte o Kind em status HTTP (spec §3.2). Kind desconhecido é 500.
func statusFor(kind apperr.Kind) int {
	switch kind {
	case apperr.BadRequest:
		return http.StatusBadRequest
	case apperr.Unauthorized:
		return http.StatusUnauthorized
	case apperr.NotFound:
		return http.StatusNotFound
	case apperr.Conflict:
		return http.StatusConflict
	case apperr.Invalid:
		return http.StatusUnprocessableEntity
	case apperr.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
