// Package log monta o logger do app: slog em texto, gravado com rotação e descarte (ADR 006).
package log

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// New devolve um logger que grava em logFile (texto do slog, com rotação do lumberjack) e,
// quando console não é nil, também em console. O io.Closer fecha o arquivo.
//
// Se logFile já existe e foi modificado num dia anterior, ele é rotacionado antes do primeiro
// registro (RN-19). Falha ao rotacionar não impede o log: segue no arquivo atual e registra
// log.rotate_failed.
func New(logFile string, level slog.Level, console io.Writer) (*slog.Logger, io.Closer, error) {
	if logFile == "" {
		return nil, nil, errors.New("log file path is empty")
	}
	info, err := os.Stat(logFile)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, nil, fmt.Errorf("inspect log file: %w", err)
	}

	file := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     14, // dias
		LocalTime:  true,
		Compress:   false,
	}
	var out io.Writer = file
	if console != nil {
		out = io.MultiWriter(file, console)
	}
	logger := slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: level}))

	if err == nil && modifiedBeforeToday(info.ModTime(), time.Now()) {
		if err := file.Rotate(); err != nil {
			logger.Warn("log.rotate_failed", "err", err)
		}
	}
	return logger, file, nil
}

// modifiedBeforeToday informa se mod é de um dia (no fuso local) anterior ao de now.
func modifiedBeforeToday(mod, now time.Time) bool {
	y, m, d := now.Date()
	startOfToday := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	return mod.Before(startOfToday)
}
