package preferences

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// FileRepository guarda as preferências num arquivo JSON (config.json).
type FileRepository struct {
	path string
}

// NewFileRepository cria o repositório para o arquivo path. A pasta do arquivo precisa existir
// (config.EnsureLayout a cria).
func NewFileRepository(path string) *FileRepository {
	return &FileRepository{path: path}
}

// Load lê o arquivo. Arquivo ausente devolve ErrNotFound; leitura que falha ou JSON inválido
// devolve outro erro, com o caminho. Não depende de mais nada, então o main.go o usa sozinho.
func (r *FileRepository) Load() (Preferences, error) {
	data, err := os.ReadFile(r.path)
	if errors.Is(err, fs.ErrNotExist) {
		return Preferences{}, fmt.Errorf("%w: %s", ErrNotFound, r.path)
	}
	if err != nil {
		return Preferences{}, fmt.Errorf("read %s: %w", r.path, err)
	}
	var p Preferences
	if err := json.Unmarshal(data, &p); err != nil {
		return Preferences{}, fmt.Errorf("parse %s: %w", r.path, err)
	}
	return p, nil
}

// Save grava p como JSON UTF-8 indentado com 2 espaços, num arquivo temporário na mesma pasta
// que depois substitui o arquivo final: uma falha no meio nunca deixa o config.json pela metade.
func (r *FileRepository) Save(p Preferences) (err error) {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("encode preferences: %w", err)
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(filepath.Dir(r.path), filepath.Base(r.path)+".*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tmp.Close() // pode já estar fechado; o erro de um segundo Close não importa
			_ = os.Remove(tmp.Name())
		}
	}()

	if _, err = tmp.Write(data); err != nil {
		return fmt.Errorf("write %s: %w", tmp.Name(), err)
	}
	if err = tmp.Sync(); err != nil {
		return fmt.Errorf("sync %s: %w", tmp.Name(), err)
	}
	if err = tmp.Close(); err != nil {
		return fmt.Errorf("close %s: %w", tmp.Name(), err)
	}
	if err = os.Rename(tmp.Name(), r.path); err != nil {
		return fmt.Errorf("replace %s: %w", r.path, err)
	}
	return nil
}
