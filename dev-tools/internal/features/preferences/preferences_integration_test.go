//go:build integration

package preferences_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dev-tools/internal/features/preferences"
	"dev-tools/internal/platform/config"
	"dev-tools/internal/platform/httpserver"
)

const testToken = "test-token"

type env struct {
	paths   config.Paths
	svc     *preferences.Service
	handler http.Handler
}

// newEnv monta a pasta do app em t.TempDir(), o FileRepository real, o serviço já carregado
// (primeira execução, idioma detectado "pt-br") e o handler HTTP com o token "test-token".
func newEnv(t *testing.T) env {
	t.Helper()
	paths := config.NewPaths(filepath.Join(t.TempDir(), "programs"))
	if _, err := config.EnsureLayout(paths); err != nil {
		t.Fatal(err)
	}
	svc := newService(paths)
	if err := svc.Load(); err != nil {
		t.Fatal(err)
	}
	cfg := httpserver.Config{
		Token:    testToken,
		Language: func() string { return svc.Get().Language },
	}
	return env{paths: paths, svc: svc, handler: httpserver.NewHandler(cfg, preferences.Routes(svc))}
}

func newService(paths config.Paths) *preferences.Service {
	return preferences.NewService(
		preferences.NewFileRepository(paths.ConfigFile),
		paths.ConfigFile,
		func() string { return "pt-br" },
		slog.New(slog.DiscardHandler),
		nil,
	)
}

func (e env) put(body, authorization string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, "/api/preferences", strings.NewReader(body))
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	rec := httptest.NewRecorder()
	e.handler.ServeHTTP(rec, req)
	return rec
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// CT-33
func TestPreferences_PersistAcrossRuns(t *testing.T) {
	e := newEnv(t)

	rec := e.put(`{"theme":"dark","language":"en"}`, "Bearer "+testToken)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, quer 200\n%s", rec.Code, rec.Body)
	}
	var got preferences.Preferences
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("corpo não é JSON: %v\n%s", err, rec.Body)
	}
	if want := (preferences.Preferences{Theme: "dark", Language: "en"}); got != want {
		t.Errorf("corpo = %+v, quer %+v", got, want)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q", ct)
	}

	if want := "{\n  \"theme\": \"dark\",\n  \"language\": \"en\"\n}\n"; readFile(t, e.paths.ConfigFile) != want {
		t.Errorf("config.json = %q, quer %q", readFile(t, e.paths.ConfigFile), want)
	}
	entries, err := os.ReadDir(e.paths.Data)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "config.json" {
		t.Errorf("pasta data = %v, quer só config.json (sem temporários)", entries)
	}

	next := newService(e.paths)
	if err := next.Load(); err != nil {
		t.Fatal(err)
	}
	if want := (preferences.Preferences{Theme: "dark", Language: "en"}); next.Get() != want {
		t.Errorf("novo serviço Get() = %+v, quer %+v", next.Get(), want)
	}
}

func TestPreferences_FirstRunCreatesConfigFile(t *testing.T) {
	e := newEnv(t)

	if want := "{\n  \"theme\": \"system\",\n  \"language\": \"pt-br\"\n}\n"; readFile(t, e.paths.ConfigFile) != want {
		t.Errorf("config.json = %q, quer %q", readFile(t, e.paths.ConfigFile), want)
	}
}

func TestPreferences_CorruptFileIsRewrittenWithDefaults(t *testing.T) {
	e := newEnv(t)
	if err := os.WriteFile(e.paths.ConfigFile, []byte(`{"theme": "dark", `), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := newService(e.paths)
	if err := svc.Load(); err != nil {
		t.Fatal(err)
	}

	if want := (preferences.Preferences{Theme: "system", Language: "pt-br"}); svc.Get() != want {
		t.Errorf("Get() = %+v, quer %+v", svc.Get(), want)
	}
	if want := "{\n  \"theme\": \"system\",\n  \"language\": \"pt-br\"\n}\n"; readFile(t, e.paths.ConfigFile) != want {
		t.Errorf("config.json = %q, quer %q", readFile(t, e.paths.ConfigFile), want)
	}
}

func TestPreferences_PartiallyInvalidFileKeepsValidField(t *testing.T) {
	e := newEnv(t)
	if err := os.WriteFile(e.paths.ConfigFile, []byte(`{"theme":"roxo","language":"en"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := newService(e.paths)
	if err := svc.Load(); err != nil {
		t.Fatal(err)
	}

	if want := (preferences.Preferences{Theme: "system", Language: "en"}); svc.Get() != want {
		t.Errorf("Get() = %+v, quer %+v", svc.Get(), want)
	}
	if want := "{\n  \"theme\": \"system\",\n  \"language\": \"en\"\n}\n"; readFile(t, e.paths.ConfigFile) != want {
		t.Errorf("config.json = %q, quer %q", readFile(t, e.paths.ConfigFile), want)
	}
}

// CT-34
func TestPreferences_RequiresToken(t *testing.T) {
	tests := []struct {
		name          string
		authorization string
	}{
		{"sem Authorization", ""},
		{"token errado", "Bearer other-token"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)
			before := readFile(t, e.paths.ConfigFile)

			rec := e.put(`{"theme":"dark","language":"en"}`, tt.authorization)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, quer 401", rec.Code)
			}
			if !strings.Contains(rec.Body.String(), `"code":"auth.token_invalid"`) {
				t.Errorf("corpo sem auth.token_invalid: %s", rec.Body)
			}
			if after := readFile(t, e.paths.ConfigFile); after != before {
				t.Errorf("config.json mudou: %q → %q", before, after)
			}
			if want := (preferences.Preferences{Theme: "system", Language: "pt-br"}); e.svc.Get() != want {
				t.Errorf("Get() = %+v, quer %+v (inalterado)", e.svc.Get(), want)
			}
		})
	}
}

func TestPreferences_HTTPErrors(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantCode   string
		wantField  string
	}{
		{"tema inválido", `{"theme":"roxo","language":"en"}`, 422, "preferences.theme_invalid", "theme"},
		{"idioma inválido", `{"theme":"dark","language":"fr"}`, 422, "preferences.language_invalid", "language"},
		{"campos ausentes valem como vazio", `{}`, 422, "preferences.theme_invalid", "theme"},
		{"JSON malformado", `{"theme":`, 400, "request.invalid_body", ""},
		{"campo desconhecido", `{"theme":"dark","language":"en","extra":1}`, 400, "request.invalid_body", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)
			before := readFile(t, e.paths.ConfigFile)

			rec := e.put(tt.body, "Bearer "+testToken)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, quer %d\n%s", rec.Code, tt.wantStatus, rec.Body)
			}
			var resp struct {
				Error struct {
					Code  string `json:"code"`
					Field string `json:"field"`
				} `json:"error"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("corpo não é JSON de erro: %v\n%s", err, rec.Body)
			}
			if resp.Error.Code != tt.wantCode || resp.Error.Field != tt.wantField {
				t.Errorf("erro = {%s, %s}, quer {%s, %s}", resp.Error.Code, resp.Error.Field, tt.wantCode, tt.wantField)
			}
			if after := readFile(t, e.paths.ConfigFile); after != before {
				t.Errorf("config.json mudou: %q → %q", before, after)
			}
		})
	}
}

func TestPreferences_SaveFailureReturns500(t *testing.T) {
	e := newEnv(t)
	// Um diretório no lugar do config.json faz o rename do arquivo temporário falhar.
	if err := os.Remove(e.paths.ConfigFile); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(e.paths.ConfigFile, 0o755); err != nil {
		t.Fatal(err)
	}

	rec := e.put(`{"theme":"dark","language":"en"}`, "Bearer "+testToken)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, quer 500\n%s", rec.Code, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), `"code":"preferences.save_failed"`) {
		t.Errorf("corpo sem preferences.save_failed: %s", rec.Body)
	}
	if want := (preferences.Preferences{Theme: "system", Language: "pt-br"}); e.svc.Get() != want {
		t.Errorf("Get() = %+v, quer %+v (inalterado)", e.svc.Get(), want)
	}
	entries, err := os.ReadDir(e.paths.Data)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("pasta data = %v, quer só o config.json (temporário apagado)", entries)
	}
}
