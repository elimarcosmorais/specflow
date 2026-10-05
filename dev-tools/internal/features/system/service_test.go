package system

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"dev-tools/internal/apperr"
)

type fakePinger struct {
	err   error
	calls int
}

func (f *fakePinger) Ping(context.Context) error {
	f.calls++
	return f.err
}

type fakeOpener struct {
	err   error
	paths []string
}

func (f *fakeOpener) Open(path string) error {
	f.paths = append(f.paths, path)
	return f.err
}

const testLogsDir = `C:\dev-tools\programs\logs`

func newTestService(p Pinger, o FolderOpener) (*Service, *bytes.Buffer) {
	logs := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(logs, nil))
	return NewService("0.1.0", p, o, testLogsDir, logger), logs
}

// CT-63
func TestService_Health(t *testing.T) {
	tests := []struct {
		name     string
		pingErr  error
		wantBase string
	}{
		{"banco disponível", nil, "available"},
		{"banco indisponível", errors.New("sql: database is closed"), "unavailable"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pinger := &fakePinger{err: tt.pingErr}
			svc, _ := newTestService(pinger, &fakeOpener{})

			got := svc.Health(context.Background())

			if want := (Health{Version: "0.1.0", Database: tt.wantBase}); got != want {
				t.Errorf("Health() = %+v, quer %+v", got, want)
			}
			if pinger.calls != 1 {
				t.Errorf("Ping chamado %d vezes, quer 1", pinger.calls)
			}
		})
	}
}

// CT-64
func TestService_OpenLogsFolder_Success(t *testing.T) {
	opener := &fakeOpener{}
	svc, _ := newTestService(&fakePinger{}, opener)

	if err := svc.OpenLogsFolder(); err != nil {
		t.Fatalf("OpenLogsFolder() = %v, quer nil", err)
	}

	if len(opener.paths) != 1 || opener.paths[0] != testLogsDir {
		t.Errorf("Open recebeu %q, quer [%q]", opener.paths, testLogsDir)
	}
}

// CT-64
func TestService_OpenLogsFolder_Failure(t *testing.T) {
	cause := errors.New("explorer.exe not found")
	opener := &fakeOpener{err: cause}
	svc, _ := newTestService(&fakePinger{}, opener)

	err := svc.OpenLogsFolder()

	if !errors.Is(err, ErrOpenLogsFailed) {
		t.Fatalf("OpenLogsFolder() = %v, quer system.open_logs_failed", err)
	}
	var ae *apperr.Error
	if !errors.As(err, &ae) || ae.Kind != apperr.Internal || ae.Code != "system.open_logs_failed" {
		t.Errorf("erro = %+v, quer Kind Internal e código system.open_logs_failed", ae)
	}
	if !errors.Is(err, cause) {
		t.Errorf("erro %v não carrega a causa %v", err, cause)
	}
	if len(opener.paths) != 1 || opener.paths[0] != testLogsDir {
		t.Errorf("Open recebeu %q, quer [%q]", opener.paths, testLogsDir)
	}
}

func TestService_OpenLogsFolder_LogsOnlyOnSuccess(t *testing.T) {
	svc, logs := newTestService(&fakePinger{}, &fakeOpener{})
	if err := svc.OpenLogsFolder(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(logs.String(), "level=INFO msg=system.logs_folder_opened") {
		t.Errorf("log sem system.logs_folder_opened:\n%s", logs)
	}

	svc, logs = newTestService(&fakePinger{}, &fakeOpener{err: errors.New("boom")})
	if err := svc.OpenLogsFolder(); err == nil {
		t.Fatal("OpenLogsFolder() = nil, quer erro")
	}
	if strings.Contains(logs.String(), "system.logs_folder_opened") {
		t.Errorf("log de sucesso numa falha:\n%s", logs)
	}
}
