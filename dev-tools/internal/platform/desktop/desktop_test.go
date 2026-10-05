package desktop

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"dev-tools/internal/platform/i18n"
)

// CT-25: o Wails expõe ao front todo método exportado do Bridge; só Bootstrap pode existir.
func TestBridge_OnlyBootstrapIsExported(t *testing.T) {
	typ := reflect.TypeOf(&Bridge{})
	if typ.NumMethod() != 1 {
		names := make([]string, 0, typ.NumMethod())
		for i := 0; i < typ.NumMethod(); i++ {
			names = append(names, typ.Method(i).Name)
		}
		t.Fatalf("exported methods = %v, want only [Bootstrap]", names)
	}
	if got := typ.Method(0).Name; got != "Bootstrap" {
		t.Errorf("exported method = %q, want %q", got, "Bootstrap")
	}
}

// CT-26
func TestBridge_Bootstrap_ReadsPreferencesOnEveryCall(t *testing.T) {
	static := BootstrapInfo{
		Version:        "0.1.0",
		APIBaseURL:     "http://127.0.0.1:41790",
		Token:          "tok",
		APIStatus:      "port_in_use",
		APIPort:        41790,
		DatabaseStatus: "unavailable",
		LogsDir:        `C:\dev-tools\programs\logs`,
	}
	theme, language := "dark", "en"
	calls := 0
	bridge := NewBridge(static, func() (string, string) {
		calls++
		return theme, language
	})

	want := static
	want.Theme, want.Language = "dark", "en"
	if got := bridge.Bootstrap(); got != want {
		t.Errorf("first Bootstrap() = %+v, want %+v", got, want)
	}

	theme, language = "light", "pt-br"
	want.Theme, want.Language = "light", "pt-br"
	if got := bridge.Bootstrap(); got != want {
		t.Errorf("second Bootstrap() = %+v, want %+v", got, want)
	}
	if calls != 2 {
		t.Errorf("prefs called %d times, want 2 (once per Bootstrap)", calls)
	}
}

func TestBootstrapInfo_JSONTags(t *testing.T) {
	want := map[string]string{
		"Version":        "version",
		"APIBaseURL":     "apiBaseUrl",
		"Token":          "token",
		"APIStatus":      "apiStatus",
		"APIPort":        "apiPort",
		"DatabaseStatus": "databaseStatus",
		"LogsDir":        "logsDir",
		"Theme":          "theme",
		"Language":       "language",
	}
	typ := reflect.TypeOf(BootstrapInfo{})
	if typ.NumField() != len(want) {
		t.Fatalf("BootstrapInfo has %d fields, want %d", typ.NumField(), len(want))
	}
	for field, tag := range want {
		f, ok := typ.FieldByName(field)
		if !ok {
			t.Errorf("BootstrapInfo has no field %s", field)
			continue
		}
		if got := f.Tag.Get("json"); got != tag {
			t.Errorf("%s json tag = %q, want %q", field, got, tag)
		}
	}
}

// CT-27
func TestUniqueID(t *testing.T) {
	a := UniqueID(`C:\dev-tools\programs`)
	b := UniqueID(`c:\DEV-TOOLS\programs\`)
	c := UniqueID(`D:\outro\programs`)

	if a != b {
		t.Errorf("UniqueID differs for the same folder written differently: %q vs %q", a, b)
	}
	if a == c {
		t.Errorf("UniqueID is equal for different folders: %q", a)
	}
	for _, id := range []string{a, b, c} {
		if !strings.HasPrefix(id, "dev-tools-") {
			t.Errorf("UniqueID %q has no %q prefix", id, "dev-tools-")
		}
	}

	sum := sha256.Sum256([]byte(`c:\dev-tools\programs`))
	if want := "dev-tools-" + hex.EncodeToString(sum[:])[:12]; a != want {
		t.Errorf("UniqueID = %q, want %q", a, want)
	}
}

// CT-28
func TestWebView2Messages(t *testing.T) {
	pt := WebView2Messages("pt-br")
	en := WebView2Messages("en")

	fields := func(m *windows.Messages) map[string]string {
		v := reflect.ValueOf(*m)
		out := make(map[string]string, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			out[v.Type().Field(i).Name] = v.Field(i).String()
		}
		return out
	}
	ptFields, enFields := fields(pt), fields(en)
	if len(ptFields) != 11 {
		t.Fatalf("windows.Messages has %d fields, want 11", len(ptFields))
	}
	for name, ptText := range ptFields {
		if ptText == "" || enFields[name] == "" {
			t.Errorf("%s is empty (pt-br=%q, en=%q)", name, ptText, enFields[name])
		}
		if ptText == enFields[name] {
			t.Errorf("%s is the same in pt-br and en: %q", name, ptText)
		}
	}

	// Cada campo do Wails recebe a chave certa do catálogo.
	keys := map[string]string{
		"InstallationRequired": "webview2.installation_required",
		"UpdateRequired":       "webview2.update_required",
		"MissingRequirements":  "webview2.missing_requirements",
		"Webview2NotInstalled": "webview2.webview2_not_installed",
		"Error":                "webview2.error",
		"FailedToInstall":      "webview2.failed_to_install",
		"DownloadPage":         "webview2.download_page",
		"PressOKToInstall":     "webview2.press_ok_to_install",
		"ContactAdmin":         "webview2.contact_admin",
		"InvalidFixedWebview2": "webview2.invalid_fixed_webview2",
		"WebView2ProcessCrash": "webview2.webview2_process_crash",
	}
	for field, key := range keys {
		if want := i18n.T("pt-br", key, nil); ptFields[field] != want {
			t.Errorf("pt-br %s = %q, want catalog %q = %q", field, ptFields[field], key, want)
		}
		if want := i18n.T("en", key, nil); enFields[field] != want {
			t.Errorf("en %s = %q, want catalog %q = %q", field, enFields[field], key, want)
		}
	}
}

func TestWindowsTheme(t *testing.T) {
	tests := []struct {
		theme string
		want  windows.Theme
	}{
		{"light", windows.Light},
		{"dark", windows.Dark},
		{"system", windows.SystemDefault},
		{"", windows.SystemDefault},
		{"unknown", windows.SystemDefault},
	}
	for _, tt := range tests {
		if got := WindowsTheme(tt.theme); got != tt.want {
			t.Errorf("WindowsTheme(%q) = %v, want %v", tt.theme, got, tt.want)
		}
	}
}

// Antes do Startup não há contexto do Wails: chamar o runtime com ctx nulo daria pânico.
func TestWindow_BeforeStartup_DoesNothing(t *testing.T) {
	w := &Window{}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic before Startup: %v", r)
		}
	}()
	w.ApplyTheme("dark")
	w.Focus()
}
