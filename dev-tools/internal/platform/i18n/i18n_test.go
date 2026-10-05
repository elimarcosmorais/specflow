package i18n

import (
	"reflect"
	"slices"
	"testing"
)

// CT-06
func TestFromLocale(t *testing.T) {
	tests := []struct {
		locale string
		want   string
	}{
		{"pt-BR", "pt-br"},
		{"pt-br", "pt-br"},
		{"pt_BR", "pt-br"},
		{"pt-PT", "en"},
		{"es-ES", "en"},
		{"en-US", "en"},
		{"", "en"},
	}
	for _, tt := range tests {
		t.Run(tt.locale, func(t *testing.T) {
			if got := FromLocale(tt.locale); got != tt.want {
				t.Errorf("FromLocale(%q) = %q, quer %q", tt.locale, got, tt.want)
			}
		})
	}
}

var errorCodes = []string{
	"auth.token_invalid",
	"request.invalid_body",
	"project.name_required",
	"project.path_required",
	"project.type_required",
	"project.type_invalid",
	"project.path_invalid",
	"project.path_not_found",
	"project.name_taken",
	"project.path_taken",
	"project.not_found",
	"preferences.theme_invalid",
	"preferences.language_invalid",
	"preferences.save_failed",
	"system.open_logs_failed",
	"database.unavailable",
	"internal.unexpected",
}

var webview2Keys = []string{
	"webview2.installation_required",
	"webview2.update_required",
	"webview2.missing_requirements",
	"webview2.webview2_not_installed",
	"webview2.error",
	"webview2.failed_to_install",
	"webview2.download_page",
	"webview2.press_ok_to_install",
	"webview2.contact_admin",
	"webview2.invalid_fixed_webview2",
	"webview2.webview2_process_crash",
}

// CT-07
func TestCatalogoCompleto(t *testing.T) {
	ptKeys, enKeys := Keys("pt-br"), Keys("en")

	if !reflect.DeepEqual(ptKeys, enKeys) {
		t.Errorf("pt-br e en têm chaves diferentes:\npt-br: %v\nen:    %v", ptKeys, enKeys)
	}

	want := slices.Concat(errorCodes, []string{"startup.title", "startup.folder_not_writable"}, webview2Keys)
	if len(ptKeys) != len(want) {
		t.Errorf("catálogo com %d chaves, quer %d", len(ptKeys), len(want))
	}
	for _, lang := range []string{"pt-br", "en"} {
		for _, key := range want {
			if got := T(lang, key, nil); got == key || got == "" {
				t.Errorf("T(%q, %q) sem tradução: %q", lang, key, got)
			}
		}
	}
}

func TestKeysIdiomaDesconhecido(t *testing.T) {
	if got := Keys("fr"); len(got) != 0 {
		t.Errorf("Keys(\"fr\") = %v, quer vazio", got)
	}
}

func TestT(t *testing.T) {
	tests := []struct {
		name     string
		language string
		key      string
		params   map[string]string
		want     string
	}{
		{
			name:     "chave inexistente devolve a própria chave",
			language: "en",
			key:      "chave.inexistente",
			want:     "chave.inexistente",
		},
		{
			name:     "chave inexistente em pt-br também",
			language: "pt-br",
			key:      "chave.inexistente",
			want:     "chave.inexistente",
		},
		{
			name:     "substitui o parâmetro",
			language: "pt-br",
			key:      "startup.folder_not_writable",
			params:   map[string]string{"path": `C:\x`},
			want:     `O Dev Tools não pode abrir: a pasta C:\x não permite escrita. Instale o app numa pasta do seu usuário (ex.: C:\dev-tools\) ou ajuste as permissões da pasta.`,
		},
		{
			name:     "sem params mantém o marcador",
			language: "pt-br",
			key:      "startup.folder_not_writable",
			want:     `O Dev Tools não pode abrir: a pasta {path} não permite escrita. Instale o app numa pasta do seu usuário (ex.: C:\dev-tools\) ou ajuste as permissões da pasta.`,
		},
		{
			name:     "texto pt-br da tabela",
			language: "pt-br",
			key:      "project.name_required",
			want:     "Informe o nome do projeto.",
		},
		{
			name:     "texto padrão do Wails em en",
			language: "en",
			key:      "webview2.error",
			want:     "Error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := T(tt.language, tt.key, tt.params); got != tt.want {
				t.Errorf("T(%q, %q, %v) = %q, quer %q", tt.language, tt.key, tt.params, got, tt.want)
			}
		})
	}
}

func TestTIdiomaDesconhecidoUsaEn(t *testing.T) {
	en := T("en", "project.name_required", nil)
	if en == "project.name_required" || en == T("pt-br", "project.name_required", nil) {
		t.Fatalf("pré-condição: en deveria ter texto próprio, veio %q", en)
	}
	if got := T("fr", "project.name_required", nil); got != en {
		t.Errorf("T(\"fr\", ...) = %q, quer o texto en %q", got, en)
	}
}
