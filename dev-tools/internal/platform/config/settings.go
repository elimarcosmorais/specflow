package config

import (
	"log/slog"
	"strconv"
	"strings"
)

const (
	envAPIPort  = "DEV_TOOLS_API_PORT"
	envLogLevel = "DEV_TOOLS_LOG_LEVEL"

	defaultAPIPort  = 41790
	defaultLogLevel = slog.LevelInfo
)

// Settings são as opções lidas das variáveis de ambiente (spec §3.3.3).
type Settings struct {
	APIPort  int
	LogLevel slog.Level
}

var logLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// LoadSettings lê DEV_TOOLS_API_PORT e DEV_TOOLS_LOG_LEVEL. Variável ausente usa o padrão em
// silêncio; valor inválido usa o padrão e gera um aviso no formato "VARIAVEL=valor", com o
// valor recusado como veio, para o chamador logar.
func LoadSettings(getenv func(string) string) (Settings, []string) {
	s := Settings{APIPort: defaultAPIPort, LogLevel: defaultLogLevel}
	var warnings []string

	if raw := getenv(envAPIPort); raw != "" {
		port, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil || port < 1 || port > 65535 {
			warnings = append(warnings, envAPIPort+"="+raw)
		} else {
			s.APIPort = port
		}
	}

	if raw := getenv(envLogLevel); raw != "" {
		level, ok := logLevels[strings.ToLower(strings.TrimSpace(raw))]
		if !ok {
			warnings = append(warnings, envLogLevel+"="+raw)
		} else {
			s.LogLevel = level
		}
	}

	return s, warnings
}
