// Package config resolve onde o app guarda seus arquivos (ADR 003) e lê as variáveis de ambiente.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Paths reúne todo caminho usado pelo app (RN-03). Layout em disco: spec §2.3.
type Paths struct {
	Base       string // pasta do executável, ou DEV_TOOLS_HOME
	Root       string // raiz da instalação: a pasta acima de Base
	Data       string
	Logs       string
	WebView    string
	DBFile     string
	ConfigFile string
	LogFile    string
	Projects   string
}

// NotWritableError indica que o app não conseguiu criar ou gravar em Path.
type NotWritableError struct {
	Path string
	Err  error
}

func (e *NotWritableError) Error() string {
	return fmt.Sprintf("folder %q is not writable: %v", e.Path, e.Err)
}

func (e *NotWritableError) Unwrap() error { return e.Err }

// ResolveBaseDir devolve a pasta base: DEV_TOOLS_HOME (absoluta) ou a pasta do executável.
func ResolveBaseDir(getenv func(string) string) (string, error) {
	if home := getenv("DEV_TOOLS_HOME"); home != "" {
		return filepath.Abs(home)
	}
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve executable path: %w", err)
	}
	return filepath.Dir(exe), nil
}

// NewPaths monta os caminhos do app a partir da pasta base.
func NewPaths(base string) Paths {
	root := filepath.Dir(base)
	data := filepath.Join(base, "data")
	logs := filepath.Join(base, "logs")
	return Paths{
		Base:       base,
		Root:       root,
		Data:       data,
		Logs:       logs,
		WebView:    filepath.Join(base, "webview"),
		DBFile:     filepath.Join(data, "dev-tools.db"),
		ConfigFile: filepath.Join(data, "config.json"),
		LogFile:    filepath.Join(logs, "dev-tools.log"),
		Projects:   filepath.Join(root, "projects"),
	}
}

// EnsureLayout confirma que a pasta base aceita escrita (RN-04) e cria as pastas que faltam
// (RN-05). Devolve só as pastas que não existiam; nunca apaga nem altera conteúdo. Qualquer
// falha devolve *NotWritableError.
func EnsureLayout(p Paths) (created []string, err error) {
	// Base primeiro: no desenvolvimento (DEV_TOOLS_HOME) ela pode não existir ainda.
	baseMade, err := ensureDir(p.Base)
	if err != nil {
		return nil, &NotWritableError{Path: p.Base, Err: err}
	}
	if err := checkWritable(p.Base); err != nil {
		return nil, err
	}
	if baseMade {
		created = append(created, p.Base)
	}

	for _, dir := range subDirs(p) {
		made, err := ensureDir(dir)
		if err != nil {
			return created, &NotWritableError{Path: dir, Err: err}
		}
		if made {
			created = append(created, dir)
		}
	}
	return created, nil
}

// subDirs lista as pastas do layout abaixo de Base e de Root, pais antes dos filhos.
func subDirs(p Paths) []string {
	return []string{
		p.Data, p.Logs, p.WebView, p.Projects,
		filepath.Join(p.Projects, "web"),
		filepath.Join(p.Projects, "android"),
		filepath.Join(p.Projects, "windows"),
	}
}

// ensureDir cria dir (e os pais que faltarem) e informa se ele ainda não existia.
func ensureDir(dir string) (made bool, err error) {
	info, err := os.Stat(dir)
	switch {
	case err == nil && info.IsDir():
		return false, nil
	case err == nil:
		return false, fmt.Errorf("%q exists and is not a folder", dir)
	case !errors.Is(err, fs.ErrNotExist):
		return false, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false, err
	}
	return true, nil
}

// checkWritable cria e apaga Base\.write-test (RN-04).
func checkWritable(base string) error {
	probe := filepath.Join(base, ".write-test")
	f, err := os.Create(probe)
	if err != nil {
		return &NotWritableError{Path: base, Err: err}
	}
	if err := f.Close(); err != nil {
		return &NotWritableError{Path: base, Err: err}
	}
	if err := os.Remove(probe); err != nil {
		return &NotWritableError{Path: base, Err: err}
	}
	return nil
}
