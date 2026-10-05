package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// CT-01
func TestResolveBaseDir(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{
			name: "DEV_TOOLS_HOME relativo vira absoluto",
			env:  map[string]string{"DEV_TOOLS_HOME": filepath.Join("rel", "programs")},
			want: filepath.Join(cwd, "rel", "programs"),
		},
		{
			name: "DEV_TOOLS_HOME absoluto fica como está",
			env:  map[string]string{"DEV_TOOLS_HOME": filepath.Join(cwd, "abs", "programs")},
			want: filepath.Join(cwd, "abs", "programs"),
		},
		{
			name: "sem a variável usa a pasta do executável",
			env:  map[string]string{},
			want: filepath.Dir(exe),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveBaseDir(func(k string) string { return tt.env[k] })
			if err != nil {
				t.Fatalf("ResolveBaseDir: erro inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("ResolveBaseDir = %q, quer %q", got, tt.want)
			}
		})
	}
}

// CT-02
func TestNewPaths(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, "programs")

	got := NewPaths(base)

	want := Paths{
		Base:       base,
		Root:       tmp,
		Data:       filepath.Join(base, "data"),
		Logs:       filepath.Join(base, "logs"),
		WebView:    filepath.Join(base, "webview"),
		DBFile:     filepath.Join(base, "data", "dev-tools.db"),
		ConfigFile: filepath.Join(base, "data", "config.json"),
		LogFile:    filepath.Join(base, "logs", "dev-tools.log"),
		Projects:   filepath.Join(tmp, "projects"),
	}
	if got != want {
		t.Errorf("NewPaths(%q) =\n%+v\nquer\n%+v", base, got, want)
	}
}

func TestEnsureLayout_CriaTudoNaPrimeiraExecucao(t *testing.T) {
	p := NewPaths(filepath.Join(t.TempDir(), "programs"))

	created, err := EnsureLayout(p)
	if err != nil {
		t.Fatalf("EnsureLayout: erro inesperado: %v", err)
	}

	wantCreated := []string{
		p.Base, p.Data, p.Logs, p.WebView, p.Projects,
		filepath.Join(p.Projects, "web"),
		filepath.Join(p.Projects, "android"),
		filepath.Join(p.Projects, "windows"),
	}
	if !reflect.DeepEqual(created, wantCreated) {
		t.Errorf("created = %v, quer %v", created, wantCreated)
	}
	for _, dir := range wantCreated {
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			t.Errorf("pasta %q não existe como diretório (err=%v)", dir, err)
		}
	}
	if _, err := os.Stat(filepath.Join(p.Base, ".write-test")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf(".write-test deveria ter sido apagado, Stat devolveu err=%v", err)
	}
}

// CT-03
func TestEnsureLayout_RecriaSoOQueFalta(t *testing.T) {
	p := NewPaths(filepath.Join(t.TempDir(), "programs"))
	if _, err := EnsureLayout(p); err != nil {
		t.Fatalf("EnsureLayout inicial: %v", err)
	}
	dataFile := filepath.Join(p.Data, "a.txt")
	if err := os.WriteFile(dataFile, []byte("conteúdo"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(p.Logs); err != nil {
		t.Fatal(err)
	}

	created, err := EnsureLayout(p)
	if err != nil {
		t.Fatalf("EnsureLayout: erro inesperado: %v", err)
	}

	if want := []string{p.Logs}; !reflect.DeepEqual(created, want) {
		t.Errorf("created = %v, quer %v", created, want)
	}
	if info, err := os.Stat(p.Logs); err != nil || !info.IsDir() {
		t.Errorf("logs não foi recriada (err=%v)", err)
	}
	got, err := os.ReadFile(dataFile)
	if err != nil {
		t.Fatalf("data\\a.txt sumiu: %v", err)
	}
	if string(got) != "conteúdo" {
		t.Errorf("data\\a.txt = %q, quer %q", got, "conteúdo")
	}

	// Estrutura completa: nada a criar.
	created, err = EnsureLayout(p)
	if err != nil {
		t.Fatalf("EnsureLayout (3ª): erro inesperado: %v", err)
	}
	if len(created) != 0 {
		t.Errorf("created = %v, quer vazio", created)
	}
}

// CT-04
func TestEnsureLayout_PastaSemEscrita(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "arquivo.txt")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := NewPaths(filepath.Join(blocker, "programs"))

	created, err := EnsureLayout(p)

	var nwe *NotWritableError
	if !errors.As(err, &nwe) {
		t.Fatalf("err = %v (%T), quer *NotWritableError", err, err)
	}
	if nwe.Path != p.Base {
		t.Errorf("NotWritableError.Path = %q, quer %q", nwe.Path, p.Base)
	}
	if nwe.Unwrap() == nil {
		t.Error("NotWritableError.Unwrap() = nil, quer a causa do sistema")
	}
	if len(created) != 0 {
		t.Errorf("created = %v, quer vazio", created)
	}
}

func TestEnsureLayout_FalhaAoCriarSubpasta(t *testing.T) {
	p := NewPaths(filepath.Join(t.TempDir(), "programs"))
	if err := os.MkdirAll(p.Base, 0o755); err != nil {
		t.Fatal(err)
	}
	// "data" é um arquivo: MkdirAll(Data) não consegue criar a pasta.
	if err := os.WriteFile(p.Data, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := EnsureLayout(p)

	var nwe *NotWritableError
	if !errors.As(err, &nwe) {
		t.Fatalf("err = %v (%T), quer *NotWritableError", err, err)
	}
	if nwe.Path != p.Data {
		t.Errorf("NotWritableError.Path = %q, quer %q", nwe.Path, p.Data)
	}
}
