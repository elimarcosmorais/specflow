package httpserver

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"dev-tools/internal/apperr"
	"dev-tools/internal/platform/i18n"
)

const testToken = "test-token"

type errorResponse struct {
	Error struct {
		Code    string  `json:"code"`
		Message string  `json:"message"`
		Field   *string `json:"field"`
	} `json:"error"`
}

// newTestHandler monta o handler com o token "test-token", o idioma lido de *language e um
// logger de texto em nível DEBUG que grava no buffer devolvido.
func newTestHandler(t *testing.T, language *string, routes []Route) (http.Handler, *bytes.Buffer) {
	t.Helper()
	logs := &bytes.Buffer{}
	cfg := Config{
		Token:    testToken,
		Language: func() string { return *language },
		Logger:   slog.New(slog.NewTextHandler(logs, &slog.HandlerOptions{Level: slog.LevelDebug})),
	}
	return NewHandler(cfg, routes), logs
}

func do(h http.Handler, method, target, authorization string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, nil)
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) errorResponse {
	t.Helper()
	var got errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("corpo não é JSON de erro: %v\n%s", err, rec.Body)
	}
	return got
}

func okHandler(calls *int) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		*calls++
		WriteJSON(w, http.StatusOK, map[string]string{"ok": "yes"})
		return nil
	}
}

// CT-14
func TestHandler_Token(t *testing.T) {
	for _, language := range []string{"pt-br", "en"} {
		t.Run(language, func(t *testing.T) {
			calls := 0
			h, _ := newTestHandler(t, &language, []Route{
				{Method: http.MethodGet, Pattern: "/api/ping", Handler: okHandler(&calls)},
			})

			for name, authorization := range map[string]string{
				"sem Authorization": "",
				"token errado":      "Bearer errado",
				"sem o prefixo":     testToken,
				"esquema errado":    "Basic " + testToken,
			} {
				rec := do(h, http.MethodGet, "/api/ping", authorization)
				if rec.Code != http.StatusUnauthorized {
					t.Fatalf("%s: status = %d, quer 401", name, rec.Code)
				}
				got := decodeError(t, rec)
				if got.Error.Code != "auth.token_invalid" {
					t.Errorf("%s: code = %q, quer auth.token_invalid", name, got.Error.Code)
				}
				if want := i18n.T(language, "auth.token_invalid", nil); got.Error.Message != want || want == "auth.token_invalid" {
					t.Errorf("%s: message = %q, quer %q", name, got.Error.Message, want)
				}
			}
			if calls != 0 {
				t.Fatalf("handler chamado %d vez(es) sem token válido, quer 0", calls)
			}

			rec := do(h, http.MethodGet, "/api/ping", "Bearer "+testToken)
			if rec.Code != http.StatusOK {
				t.Errorf("com token: status = %d, quer 200", rec.Code)
			}
			if calls != 1 {
				t.Errorf("com token: handler chamado %d vez(es), quer 1", calls)
			}
		})
	}
}

func TestHandler_TokenMessagesDifferByLanguage(t *testing.T) {
	language := "pt-br"
	h, _ := newTestHandler(t, &language, []Route{
		{Method: http.MethodGet, Pattern: "/api/ping", Handler: okHandler(new(int))},
	})

	pt := decodeError(t, do(h, http.MethodGet, "/api/ping", "")).Error.Message
	language = "en"
	en := decodeError(t, do(h, http.MethodGet, "/api/ping", "")).Error.Message

	if pt == en {
		t.Errorf("a mensagem não mudou com o idioma: %q", pt)
	}
}

func TestHandler_EmptyConfiguredTokenRejectsEverything(t *testing.T) {
	language := "en"
	calls := 0
	logs := &bytes.Buffer{}
	h := NewHandler(Config{
		Token:    "",
		Language: func() string { return language },
		Logger:   slog.New(slog.NewTextHandler(logs, nil)),
	}, []Route{{Method: http.MethodGet, Pattern: "/api/ping", Handler: okHandler(&calls)}})

	for _, authorization := range []string{"", "Bearer ", "Bearer x"} {
		if rec := do(h, http.MethodGet, "/api/ping", authorization); rec.Code != http.StatusUnauthorized {
			t.Errorf("Authorization %q: status = %d, quer 401", authorization, rec.Code)
		}
	}
	if calls != 0 {
		t.Errorf("handler chamado %d vez(es), quer 0", calls)
	}
}

// CT-15
func TestHandler_PublicRoute(t *testing.T) {
	language := "en"
	calls := 0
	h, _ := newTestHandler(t, &language, []Route{
		{Method: http.MethodGet, Pattern: "/api/health", Public: true, Handler: okHandler(&calls)},
	})

	rec := do(h, http.MethodGet, "/api/health", "")

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, quer 200", rec.Code)
	}
	if calls != 1 {
		t.Errorf("handler chamado %d vez(es), quer 1", calls)
	}
}

// CT-16
func TestHandler_Preflight(t *testing.T) {
	language := "en"
	calls := 0
	h, _ := newTestHandler(t, &language, []Route{
		{Method: http.MethodGet, Pattern: "/api/projects", Handler: okHandler(&calls)},
		{Method: http.MethodPost, Pattern: "/api/projects", Handler: okHandler(&calls)},
	})

	req := httptest.NewRequest(http.MethodOptions, "/api/projects", nil)
	req.Header.Set("Origin", "http://wails.localhost")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, quer 204", rec.Code)
	}
	want := map[string]string{
		"Access-Control-Allow-Origin":          "*",
		"Access-Control-Allow-Methods":         "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers":         "Authorization, Content-Type",
		"Access-Control-Allow-Private-Network": "true",
		"Access-Control-Max-Age":               "600",
	}
	for name, value := range want {
		if got := rec.Header().Get(name); got != value {
			t.Errorf("%s = %q, quer %q", name, got, value)
		}
	}
	if calls != 0 {
		t.Errorf("handler chamado %d vez(es) no preflight, quer 0", calls)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("preflight com corpo: %q", rec.Body)
	}
}

func TestHandler_CORSOriginOnEveryResponse(t *testing.T) {
	language := "en"
	h, _ := newTestHandler(t, &language, []Route{
		{Method: http.MethodGet, Pattern: "/api/ping", Handler: okHandler(new(int))},
	})

	tests := []struct {
		name          string
		authorization string
		target        string
		wantStatus    int
	}{
		{"sucesso", "Bearer " + testToken, "/api/ping", http.StatusOK},
		{"401 sem token", "", "/api/ping", http.StatusUnauthorized},
		{"rota fora do registro", "Bearer " + testToken, "/api/inexistente", http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(h, http.MethodGet, tt.target, tt.authorization)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, quer %d", rec.Code, tt.wantStatus)
			}
			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
				t.Errorf("Access-Control-Allow-Origin = %q, quer *", got)
			}
		})
	}
}

// CT-17
func TestHandler_ErrorResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) error {
		return apperr.New(apperr.Invalid, "project.name_required").WithField("name")
	}
	messages := map[string]string{}

	for _, language := range []string{"pt-br", "en"} {
		t.Run(language, func(t *testing.T) {
			h, _ := newTestHandler(t, &language, []Route{
				{Method: http.MethodPost, Pattern: "/api/projects", Handler: handler},
			})

			rec := do(h, http.MethodPost, "/api/projects", "Bearer "+testToken)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, quer 422", rec.Code)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
				t.Errorf("Content-Type = %q", got)
			}
			want := i18n.T(language, "project.name_required", nil)
			if want == "project.name_required" {
				t.Fatalf("catálogo sem tradução em %s", language)
			}
			wantBody := `{"error":{"code":"project.name_required","message":` + jsonString(t, want) + `,"field":"name"}}`
			if got := rec.Body.String(); got != wantBody {
				t.Errorf("corpo = %s\nquer    %s", got, wantBody)
			}
			messages[language] = want
		})
	}
	if messages["pt-br"] == messages["en"] {
		t.Errorf("mensagens iguais nos dois idiomas: %q", messages["pt-br"])
	}
}

func jsonString(t *testing.T, s string) string {
	t.Helper()
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestHandler_ErrorFieldOmittedWhenEmpty(t *testing.T) {
	language := "en"
	h, _ := newTestHandler(t, &language, []Route{
		{Method: http.MethodGet, Pattern: "/api/x", Handler: func(w http.ResponseWriter, r *http.Request) error {
			return apperr.New(apperr.NotFound, "project.not_found")
		}},
	})

	rec := do(h, http.MethodGet, "/api/x", "Bearer "+testToken)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, quer 404", rec.Code)
	}
	if strings.Contains(rec.Body.String(), `"field"`) {
		t.Errorf("corpo tem field: %s", rec.Body)
	}
}

func TestStatusFor(t *testing.T) {
	tests := []struct {
		kind apperr.Kind
		want int
	}{
		{apperr.BadRequest, 400},
		{apperr.Unauthorized, 401},
		{apperr.NotFound, 404},
		{apperr.Conflict, 409},
		{apperr.Invalid, 422},
		{apperr.Internal, 500},
		{apperr.Unavailable, 503},
		{apperr.Kind("desconhecido"), 500},
	}
	for _, tt := range tests {
		t.Run(string(tt.kind), func(t *testing.T) {
			if got := statusFor(tt.kind); got != tt.want {
				t.Errorf("statusFor(%q) = %d, quer %d", tt.kind, got, tt.want)
			}
		})
	}
}

// CT-18
func TestHandler_UnexpectedError(t *testing.T) {
	tests := []struct {
		name    string
		handler HandlerFunc
		cause   string
	}{
		{"erro que não é apperr", func(w http.ResponseWriter, r *http.Request) error { return errors.New("boom") }, "boom"},
		{"pânico", func(w http.ResponseWriter, r *http.Request) error { panic("kaboom") }, "kaboom"},
		{"erro composto que não é apperr", func(w http.ResponseWriter, r *http.Request) error {
			return errors.Join(errors.New("boom"), errors.New("outro"))
		}, "boom"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			language := "pt-br"
			h, logs := newTestHandler(t, &language, []Route{
				{Method: http.MethodGet, Pattern: "/api/x", Handler: tt.handler},
			})

			rec := do(h, http.MethodGet, "/api/x", "Bearer "+testToken)

			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, quer 500", rec.Code)
			}
			got := decodeError(t, rec)
			if got.Error.Code != "internal.unexpected" {
				t.Errorf("code = %q, quer internal.unexpected", got.Error.Code)
			}
			if want := i18n.T("pt-br", "internal.unexpected", nil); got.Error.Message != want {
				t.Errorf("message = %q, quer %q", got.Error.Message, want)
			}
			if strings.Contains(rec.Body.String(), tt.cause) {
				t.Errorf("a causa %q vazou no corpo: %s", tt.cause, rec.Body)
			}
			line := logLine(t, logs, "http.error")
			for _, want := range []string{"level=ERROR", "method=GET", "path=/api/x", "status=500", "code=internal.unexpected", tt.cause} {
				if !strings.Contains(line, want) {
					t.Errorf("log sem %q:\n%s", want, line)
				}
			}
		})
	}
}

func TestHandler_PanicAfterResponseStartedKeepsStatus(t *testing.T) {
	language := "en"
	h, logs := newTestHandler(t, &language, []Route{
		{Method: http.MethodGet, Pattern: "/api/x", Handler: func(w http.ResponseWriter, r *http.Request) error {
			w.WriteHeader(http.StatusAccepted)
			panic("late")
		}},
	})

	rec := do(h, http.MethodGet, "/api/x", "Bearer "+testToken)

	if rec.Code != http.StatusAccepted {
		t.Errorf("status = %d, quer 202 (já escrito)", rec.Code)
	}
	line := logLine(t, logs, "http.error")
	for _, want := range []string{"level=ERROR", "status=202", "late"} {
		if !strings.Contains(line, want) {
			t.Errorf("log sem %q:\n%s", want, line)
		}
	}
}

// logLine devolve a única linha de log com msg=<msg>; falha se não houver exatamente uma.
func logLine(t *testing.T, logs *bytes.Buffer, msg string) string {
	t.Helper()
	var found []string
	for _, line := range strings.Split(logs.String(), "\n") {
		if strings.Contains(line, "msg="+msg+" ") || strings.HasSuffix(line, "msg="+msg) {
			found = append(found, line)
		}
	}
	if len(found) != 1 {
		t.Fatalf("esperava 1 linha com msg=%s, achei %d:\n%s", msg, len(found), logs)
	}
	return found[0]
}

func TestHandler_AccessLog(t *testing.T) {
	language := "en"
	h, logs := newTestHandler(t, &language, []Route{
		{Method: http.MethodGet, Pattern: "/api/ping", Handler: okHandler(new(int))},
		{Method: http.MethodGet, Pattern: "/api/missing", Handler: func(w http.ResponseWriter, r *http.Request) error {
			return apperr.New(apperr.NotFound, "project.not_found").Wrap(errors.New("segredo-interno"))
		}},
	})

	do(h, http.MethodGet, "/api/ping", "Bearer "+testToken)
	line := logLine(t, logs, "http.request")
	if ok, _ := regexp.MatchString(`level=DEBUG msg=http\.request method=GET path=/api/ping status=200 duration_ms=\d+$`, line); !ok {
		t.Errorf("linha de acesso inesperada:\n%s", line)
	}

	logs.Reset()
	do(h, http.MethodGet, "/api/missing", "Bearer "+testToken)
	line = logLine(t, logs, "http.error")
	for _, want := range []string{"level=WARN", "method=GET", "path=/api/missing", "status=404", "code=project.not_found"} {
		if !strings.Contains(line, want) {
			t.Errorf("log de 4xx sem %q:\n%s", want, line)
		}
	}
	if strings.Contains(line, "err=") || strings.Contains(logs.String(), "segredo-interno") {
		t.Errorf("log de 4xx não deve ter a causa:\n%s", logs)
	}
}

func TestHandler_LogsNeverContainTokenOrAuthorization(t *testing.T) {
	language := "en"
	h, logs := newTestHandler(t, &language, []Route{
		{Method: http.MethodGet, Pattern: "/api/ping", Handler: okHandler(new(int))},
	})

	do(h, http.MethodGet, "/api/ping", "Bearer "+testToken)
	do(h, http.MethodGet, "/api/ping", "Bearer errado-secreto")

	for _, secret := range []string{testToken, "errado-secreto", "Authorization", "Bearer"} {
		if strings.Contains(logs.String(), secret) {
			t.Errorf("o log contém %q:\n%s", secret, logs)
		}
	}
}

func TestHandler_UnknownRouteKeepsMuxResponse(t *testing.T) {
	language := "en"
	h, _ := newTestHandler(t, &language, []Route{
		{Method: http.MethodGet, Pattern: "/api/ping", Handler: okHandler(new(int))},
	})

	if rec := do(h, http.MethodGet, "/api/inexistente", "Bearer "+testToken); rec.Code != http.StatusNotFound {
		t.Errorf("rota inexistente: status = %d, quer 404", rec.Code)
	}
	if rec := do(h, http.MethodPost, "/api/ping", "Bearer "+testToken); rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("método não registrado: status = %d, quer 405", rec.Code)
	}
}

// CT-19
func TestDecodeJSON_Invalid(t *testing.T) {
	type input struct {
		Name string `json:"name"`
	}
	tests := []struct {
		name string
		body string
	}{
		{"malformado", `{`},
		{"campo desconhecido", `{"x":1}`},
		{"maior que 1 MiB", `{"name":"` + strings.Repeat("a", 2<<20) + `"}`},
		{"corpo vazio", ``},
		{"conteúdo depois do valor", `{"name":"a"}{"name":"b"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))

			err := DecodeJSON(req, &input{})

			if !errors.Is(err, apperr.ErrInvalidBody) {
				t.Fatalf("err = %v, quer apperr.ErrInvalidBody", err)
			}
			var ae *apperr.Error
			if !errors.As(err, &ae) || ae.Err == nil {
				t.Errorf("o erro não guardou a causa: %#v", err)
			}
		})
	}
}

func TestDecodeJSON_TooLargeKeepsCause(t *testing.T) {
	body := `{"name":"` + strings.Repeat("a", 2<<20) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	err := DecodeJSON(req, &struct {
		Name string `json:"name"`
	}{})

	var tooLarge *http.MaxBytesError
	if !errors.As(err, &tooLarge) || tooLarge.Limit != 1<<20 {
		t.Errorf("a causa não é *http.MaxBytesError com limite de 1 MiB: %v", err)
	}
}

func TestDecodeJSON_Valid(t *testing.T) {
	var got struct {
		Name string `json:"name"`
	}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(` {"name":"loja-api"} `))

	if err := DecodeJSON(req, &got); err != nil {
		t.Fatalf("err = %v, quer nil", err)
	}
	if got.Name != "loja-api" {
		t.Errorf("Name = %q, quer loja-api", got.Name)
	}
}

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteJSON(rec, http.StatusCreated, map[string]int{"id": 1})

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, quer 201", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q", got)
	}
	if got := rec.Body.String(); got != `{"id":1}` {
		t.Errorf("corpo = %q", got)
	}
}

// CT-20
func TestHandler_Docs(t *testing.T) {
	language := "en"
	h, _ := newTestHandler(t, &language, nil)

	tests := []struct {
		path        string
		contentType string
		contains    string
	}{
		{"/docs", "text/html; charset=utf-8", "<redoc"},
		{"/docs/openapi.yaml", "application/yaml; charset=utf-8", "openapi: 3.1"},
		{"/docs/redoc.standalone.js", "text/javascript; charset=utf-8", "Redoc"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			rec := do(h, http.MethodGet, tt.path, "")

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, quer 200", rec.Code)
			}
			if got := rec.Header().Get("Content-Type"); got != tt.contentType {
				t.Errorf("Content-Type = %q, quer %q", got, tt.contentType)
			}
			if !strings.Contains(rec.Body.String(), tt.contains) {
				t.Errorf("corpo não contém %q", tt.contains)
			}
			want, err := docsFS.ReadFile("docs/" + strings.TrimPrefix(tt.path, "/docs/"))
			if tt.path == "/docs" {
				want, err = docsFS.ReadFile("docs/index.html")
			}
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(rec.Body.Bytes(), want) {
				t.Error("o corpo difere do arquivo embutido")
			}
		})
	}
}

func TestDocs_IndexUsesOnlyLocalPaths(t *testing.T) {
	index, err := docsFS.ReadFile("docs/index.html")
	if err != nil {
		t.Fatal(err)
	}

	refs := regexp.MustCompile(`(?i)\b(?:src|href|spec-url)\s*=\s*["']([^"']*)["']`).FindAllStringSubmatch(string(index), -1)
	if len(refs) == 0 {
		t.Fatal("o index.html não referencia nenhum recurso")
	}
	for _, ref := range refs {
		if !strings.HasPrefix(ref[1], "/docs/") {
			t.Errorf("referência fora de /docs/: %q", ref[1])
		}
	}
	if !strings.Contains(string(index), "/docs/redoc.standalone.js") || !strings.Contains(string(index), "/docs/openapi.yaml") {
		t.Error("o index.html deve carregar redoc.standalone.js e openapi.yaml")
	}
}

func TestNewToken(t *testing.T) {
	a, err := NewToken()
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewToken()
	if err != nil {
		t.Fatal(err)
	}

	if len(a) != 64 {
		t.Errorf("len(token) = %d, quer 64", len(a))
	}
	if _, err := hex.DecodeString(a); err != nil {
		t.Errorf("token não é hexadecimal: %v", err)
	}
	if a == b {
		t.Error("dois tokens seguidos iguais")
	}
}
