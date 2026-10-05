package desktop

import (
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"dev-tools/internal/platform/i18n"
)

// WebView2Messages devolve os textos que o Wails mostra quando o WebView2 está ausente ou
// desatualizado, no idioma pedido (ADR 005).
func WebView2Messages(language string) *windows.Messages {
	t := func(key string) string { return i18n.T(language, key, nil) }
	return &windows.Messages{
		InstallationRequired: t("webview2.installation_required"),
		UpdateRequired:       t("webview2.update_required"),
		MissingRequirements:  t("webview2.missing_requirements"),
		Webview2NotInstalled: t("webview2.webview2_not_installed"),
		Error:                t("webview2.error"),
		FailedToInstall:      t("webview2.failed_to_install"),
		DownloadPage:         t("webview2.download_page"),
		PressOKToInstall:     t("webview2.press_ok_to_install"),
		ContactAdmin:         t("webview2.contact_admin"),
		InvalidFixedWebview2: t("webview2.invalid_fixed_webview2"),
		WebView2ProcessCrash: t("webview2.webview2_process_crash"),
	}
}
