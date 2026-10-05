package config

import (
	"log/slog"
	"reflect"
	"testing"
)

// CT-05
func TestLoadSettings(t *testing.T) {
	tests := []struct {
		name         string
		env          map[string]string
		want         Settings
		wantWarnings []string
	}{
		{
			name: "vazio usa os padrões",
			env:  map[string]string{},
			want: Settings{APIPort: 41790, LogLevel: slog.LevelInfo},
		},
		{
			name: "valores válidos",
			env:  map[string]string{"DEV_TOOLS_API_PORT": "41800", "DEV_TOOLS_LOG_LEVEL": "debug"},
			want: Settings{APIPort: 41800, LogLevel: slog.LevelDebug},
		},
		{
			name:         "valores inválidos voltam ao padrão com dois avisos",
			env:          map[string]string{"DEV_TOOLS_API_PORT": "abc", "DEV_TOOLS_LOG_LEVEL": "xyz"},
			want:         Settings{APIPort: 41790, LogLevel: slog.LevelInfo},
			wantWarnings: []string{"DEV_TOOLS_API_PORT=abc", "DEV_TOOLS_LOG_LEVEL=xyz"},
		},
		{
			name: "limite superior da porta é aceito",
			env:  map[string]string{"DEV_TOOLS_API_PORT": "65535"},
			want: Settings{APIPort: 65535, LogLevel: slog.LevelInfo},
		},
		{
			name: "limite inferior da porta é aceito",
			env:  map[string]string{"DEV_TOOLS_API_PORT": "1"},
			want: Settings{APIPort: 1, LogLevel: slog.LevelInfo},
		},
		{
			name:         "porta zero é recusada",
			env:          map[string]string{"DEV_TOOLS_API_PORT": "0"},
			want:         Settings{APIPort: 41790, LogLevel: slog.LevelInfo},
			wantWarnings: []string{"DEV_TOOLS_API_PORT=0"},
		},
		{
			name:         "porta acima de 65535 é recusada",
			env:          map[string]string{"DEV_TOOLS_API_PORT": "65536"},
			want:         Settings{APIPort: 41790, LogLevel: slog.LevelInfo},
			wantWarnings: []string{"DEV_TOOLS_API_PORT=65536"},
		},
		{
			name:         "porta negativa é recusada",
			env:          map[string]string{"DEV_TOOLS_API_PORT": "-1"},
			want:         Settings{APIPort: 41790, LogLevel: slog.LevelInfo},
			wantWarnings: []string{"DEV_TOOLS_API_PORT=-1"},
		},
		{
			name: "nível info",
			env:  map[string]string{"DEV_TOOLS_LOG_LEVEL": "info"},
			want: Settings{APIPort: 41790, LogLevel: slog.LevelInfo},
		},
		{
			name: "nível warn",
			env:  map[string]string{"DEV_TOOLS_LOG_LEVEL": "warn"},
			want: Settings{APIPort: 41790, LogLevel: slog.LevelWarn},
		},
		{
			name: "nível error",
			env:  map[string]string{"DEV_TOOLS_LOG_LEVEL": "error"},
			want: Settings{APIPort: 41790, LogLevel: slog.LevelError},
		},
		{
			name:         "só a porta inválida avisa; o nível válido vale",
			env:          map[string]string{"DEV_TOOLS_API_PORT": "99999", "DEV_TOOLS_LOG_LEVEL": "error"},
			want:         Settings{APIPort: 41790, LogLevel: slog.LevelError},
			wantWarnings: []string{"DEV_TOOLS_API_PORT=99999"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, warnings := LoadSettings(func(k string) string { return tt.env[k] })
			if got != tt.want {
				t.Errorf("Settings = %+v, quer %+v", got, tt.want)
			}
			if len(warnings) == 0 && len(tt.wantWarnings) == 0 {
				return
			}
			if !reflect.DeepEqual(warnings, tt.wantWarnings) {
				t.Errorf("avisos = %q, quer %q", warnings, tt.wantWarnings)
			}
		})
	}
}
