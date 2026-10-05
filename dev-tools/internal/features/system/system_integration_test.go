//go:build integration

package system_test

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"dev-tools/internal/features/system"
	"dev-tools/internal/platform/config"
	"dev-tools/internal/platform/db"
	"dev-tools/internal/platform/httpserver"
)

const testToken = "test-token"

type recordingOpener struct {
	err   error
	paths []string
}

func (r *recordingOpener) Open(path string) error {
	r.paths = append(r.paths, path)
	return r.err
}

type okPinger struct{}

func (okPinger) Ping(context.Context) error { return nil }

// CT-65
func TestSystem_HealthViaAPI(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	paths := config.NewPaths(filepath.Join(t.TempDir(), "programs"))
	conn, err := db.Open(context.Background(), filepath.Join(t.TempDir(), "dev-tools.db"), logger)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() }) // Close é idempotente; o teste fecha antes de propósito

	svc := system.NewService("0.1.0", db.Pinger{DB: conn}, &recordingOpener{}, paths.Logs, logger)
	handler := httpserver.NewHandler(httpserver.Config{Token: testToken}, system.Routes(svc))

	get := func() *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/health", nil))
		return rec
	}

	rec := get()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, quer 200\n%s", rec.Code, rec.Body)
	}
	if want := `{"version":"0.1.0","database":"available"}`; rec.Body.String() != want {
		t.Errorf("corpo = %s, quer %s", rec.Body, want)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q", ct)
	}

	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}

	rec = get()
	if rec.Code != http.StatusOK {
		t.Fatalf("status com o banco fechado = %d, quer 200\n%s", rec.Code, rec.Body)
	}
	if want := `{"version":"0.1.0","database":"unavailable"}`; rec.Body.String() != want {
		t.Errorf("corpo com o banco fechado = %s, quer %s", rec.Body, want)
	}
}

func TestSystem_OpenLogsFolderRequiresToken(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	opener := &recordingOpener{}
	svc := system.NewService("0.1.0", okPinger{}, opener, `C:\logs`, logger)
	handler := httpserver.NewHandler(httpserver.Config{Token: testToken}, system.Routes(svc))

	tests := []struct {
		name          string
		authorization string
		wantStatus    int
		wantOpened    []string
	}{
		{"sem token", "", http.StatusUnauthorized, nil},
		{"com token", "Bearer " + testToken, http.StatusNoContent, []string{`C:\logs`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opener.paths = nil
			req := httptest.NewRequest(http.MethodPost, "/api/system/open-logs-folder", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, quer %d\n%s", rec.Code, tt.wantStatus, rec.Body)
			}
			if tt.wantStatus == http.StatusNoContent && rec.Body.Len() != 0 {
				t.Errorf("corpo = %q, quer vazio", rec.Body)
			}
			if len(opener.paths) != len(tt.wantOpened) || (len(tt.wantOpened) > 0 && opener.paths[0] != tt.wantOpened[0]) {
				t.Errorf("Open recebeu %q, quer %q", opener.paths, tt.wantOpened)
			}
		})
	}
}

func TestSystem_OpenLogsFolderFailure(t *testing.T) {
	opener := &recordingOpener{err: errors.New("explorer.exe not found")}
	svc := system.NewService("0.1.0", okPinger{}, opener, `C:\logs`, slog.New(slog.DiscardHandler))
	cfg := httpserver.Config{Token: testToken, Language: func() string { return "en" }}
	handler := httpserver.NewHandler(cfg, system.Routes(svc))
	req := httptest.NewRequest(http.MethodPost, "/api/system/open-logs-folder", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, quer 500: %s", rec.Code, rec.Body)
	}
	want := `{"error":{"code":"system.open_logs_failed","message":"The logs folder could not be opened."}}`
	if rec.Body.String() != want {
		t.Errorf("corpo = %s, quer %s", rec.Body, want)
	}
}
