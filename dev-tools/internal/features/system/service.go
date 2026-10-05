// Package system reúne a saúde do motor e a abertura da pasta de logs. Não guarda dados (sem
// repositório) e não importa internal/platform: a infraestrutura satisfaz as interfaces, e o
// main.go liga as peças.
package system

import (
	"context"
	"log/slog"
	"time"

	"dev-tools/internal/apperr"
)

const (
	DatabaseAvailable   = "available"
	DatabaseUnavailable = "unavailable"

	// pingTimeout limita o Ping: um banco travado não pode travar a saúde.
	pingTimeout = 2 * time.Second
)

// ErrOpenLogsFailed é devolvido quando o Explorer não inicia (spec §3.2). A causa é acrescentada
// com Wrap.
var ErrOpenLogsFailed = apperr.New(apperr.Internal, "system.open_logs_failed")

// Pinger verifica se o banco responde. db.Pinger o satisfaz.
type Pinger interface {
	Ping(ctx context.Context) error
}

// FolderOpener abre uma pasta no gerenciador de arquivos do sistema.
type FolderOpener interface {
	Open(path string) error
}

// Health é o corpo de GET /api/health.
type Health struct {
	Version  string `json:"version"`
	Database string `json:"database"`
}

// Service responde a saúde e abre a pasta de logs.
type Service struct {
	version string
	pinger  Pinger
	opener  FolderOpener
	logsDir string
	logger  *slog.Logger
}

// NewService cria o serviço. version é a versão do app; logsDir é a pasta aberta por
// OpenLogsFolder.
func NewService(version string, pinger Pinger, opener FolderOpener, logsDir string, logger *slog.Logger) *Service {
	return &Service{version: version, pinger: pinger, opener: opener, logsDir: logsDir, logger: logger}
}

// Health devolve a versão e o estado do banco. Falha no Ping é "unavailable", nunca um erro:
// a saúde precisa responder justamente quando o banco está fora do ar.
func (s *Service) Health(ctx context.Context) Health {
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	database := DatabaseAvailable
	if err := s.pinger.Ping(ctx); err != nil {
		database = DatabaseUnavailable
	}
	return Health{Version: s.version, Database: database}
}

// OpenLogsFolder abre a pasta de logs no Explorer. Falha vira ErrOpenLogsFailed com a causa.
func (s *Service) OpenLogsFolder() error {
	if err := s.opener.Open(s.logsDir); err != nil {
		return ErrOpenLogsFailed.Wrap(err)
	}
	s.logger.Info("system.logs_folder_opened", "path", s.logsDir)
	return nil
}
