package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"dev-tools/frontend"
	"dev-tools/internal/features/preferences"
	"dev-tools/internal/features/projects"
	"dev-tools/internal/features/system"
	"dev-tools/internal/platform/config"
	"dev-tools/internal/platform/db"
	"dev-tools/internal/platform/desktop"
	"dev-tools/internal/platform/httpserver"
	"dev-tools/internal/platform/i18n"
	applog "dev-tools/internal/platform/log"
)

var version = "0.1.0"

// shutdownTimeout é o prazo para as requisições em andamento terminarem ao encerrar.
const shutdownTimeout = 3 * time.Second

func main() {
	os.Exit(run())
}

// run monta o app na sequência de spec §6.1 e devolve o código de saída do processo.
func run() int {
	// 1. Pasta base e caminhos.
	base, err := config.ResolveBaseDir(os.Getenv)
	if err != nil {
		showStartupError(i18n.DetectWindows(), err.Error())
		return 1
	}
	paths := config.NewPaths(base)

	// 2. Estrutura de pastas. Sem escrita na pasta base não há onde logar (RN-04): só o diálogo.
	created, err := config.EnsureLayout(paths)
	if err != nil {
		message := err.Error()
		language := startupLanguage(paths)
		var notWritable *config.NotWritableError
		if errors.As(err, &notWritable) {
			message = i18n.T(language, "startup.folder_not_writable", map[string]string{"path": notWritable.Path})
		}
		showStartupError(language, message)
		return 1
	}

	// 3. e 4. Variáveis de ambiente e log.
	settings, warnings := config.LoadSettings(os.Getenv)
	logger, logCloser, err := applog.New(paths.LogFile, settings.LogLevel, os.Stdout)
	if err != nil {
		showStartupError(startupLanguage(paths), err.Error())
		return 1
	}
	logger.Info("app.started", "version", version, "base_dir", paths.Base)
	for _, dir := range created {
		logger.Info("layout.created", "path", dir)
	}
	for _, warning := range warnings {
		name, value, _ := strings.Cut(warning, "=") // formato "VAR=valor"
		logger.Warn("config.invalid_env", "var", name, "value", value)
	}

	// 5. Preferências.
	window := &desktop.Window{}
	prefs := preferences.NewService(
		preferences.NewFileRepository(paths.ConfigFile),
		paths.ConfigFile,
		i18n.DetectWindows,
		logger,
		func(p preferences.Preferences) { window.ApplyTheme(p.Theme) },
	)
	if err := prefs.Load(); err != nil {
		// Os padrões já estão em memória: o app segue sem conseguir gravar o config.json.
		logger.Error("config.save_failed", "path", paths.ConfigFile, "err", err)
	}

	// 6. Banco. Sem banco o app abre em modo degradado (D-06).
	databaseStatus := "available"
	conn, err := db.Open(context.Background(), paths.DBFile, logger)
	if err != nil {
		logger.Error("db.open_failed", "path", paths.DBFile, "err", err)
		databaseStatus = "unavailable"
	}

	// 7. Serviços, token e handler da API.
	token, err := httpserver.NewToken()
	if err != nil {
		logger.Error("app.token_failed", "err", err)
		showStartupError(prefs.Get().Language, err.Error())
		closeAll(logger, nil, conn, logCloser)
		return 1
	}
	handler := httpserver.NewHandler(
		httpserver.Config{
			Token:    token,
			Language: func() string { return prefs.Get().Language },
			Logger:   logger,
		},
		apiRoutes(
			projects.NewService(projects.NewSQLiteRepository(conn), paths.Root, logger, time.Now),
			prefs,
			system.NewService(version, db.Pinger{DB: conn}, system.ExplorerOpener{}, paths.Logs, logger),
		),
	)

	// 8. API na porta local.
	srv := &http.Server{Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	apiStatus := "ok"
	ln, err := httpserver.Listen(settings.APIPort)
	switch {
	case err == nil:
		go func() {
			if err := srv.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
				logger.Error("api.serve_failed", "err", err)
			}
		}()
		logger.Info("api.listening", "addr", ln.Addr().String())
	case errors.Is(err, httpserver.ErrPortInUse):
		logger.Error("api.port_in_use", "port", settings.APIPort)
		apiStatus = "port_in_use"
	default:
		logger.Error("api.listen_failed", "port", settings.APIPort, "err", err)
		apiStatus = "failed"
	}

	// 9. Janela. O encerramento roda no OnShutdown e, se o Wails falhar antes de abrir, aqui.
	var once sync.Once
	shutdown := func() { once.Do(func() { closeAll(logger, srv, conn, logCloser) }) }

	bridge := desktop.NewBridge(
		desktop.BootstrapInfo{
			Version:        version,
			APIBaseURL:     fmt.Sprintf("http://127.0.0.1:%d", settings.APIPort),
			Token:          token,
			APIStatus:      apiStatus,
			APIPort:        settings.APIPort,
			DatabaseStatus: databaseStatus,
			LogsDir:        paths.Logs,
		},
		func() (theme, language string) {
			p := prefs.Get()
			return p.Theme, p.Language
		},
	)
	initial := prefs.Get()
	err = wails.Run(&options.App{
		Title:       "Dev Tools",
		Width:       1100,
		Height:      720,
		MinWidth:    900,
		MinHeight:   600,
		AssetServer: &assetserver.Options{Assets: frontend.Assets()},
		Bind:        []any{bridge},
		OnStartup:   window.Startup,
		OnShutdown:  func(context.Context) { shutdown() },
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               desktop.UniqueID(paths.Base),
			OnSecondInstanceLaunch: func(options.SecondInstanceData) { window.Focus() },
		},
		Windows: &windows.Options{
			WebviewUserDataPath: paths.WebView,
			Theme:               desktop.WindowsTheme(initial.Theme),
			Messages:            desktop.WebView2Messages(initial.Language),
		},
	})
	code := 0
	if err != nil {
		logger.Error("app.run_failed", "err", err)
		code = 1
	}
	shutdown()
	return code
}

// closeAll encerra o app: para a API (com prazo), fecha o banco, registra app.stopped e fecha o
// log. srv e conn podem ser nulos, quando a etapa correspondente não aconteceu.
func closeAll(logger *slog.Logger, srv *http.Server, conn *sql.DB, logCloser io.Closer) {
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Warn("api.shutdown_failed", "err", err)
		}
		cancel()
	}
	if conn != nil {
		if err := conn.Close(); err != nil {
			logger.Warn("db.close_failed", "err", err)
		}
	}
	logger.Info("app.stopped")
	_ = logCloser.Close() // último passo: não há onde registrar uma falha
}

// startupLanguage devolve o idioma para mensagens antes de existir o serviço de preferências:
// o do config.json ou, se não der para lê-lo, o do Windows.
func startupLanguage(paths config.Paths) string {
	p, err := preferences.NewFileRepository(paths.ConfigFile).Load()
	if err != nil || p.Language == "" {
		return i18n.DetectWindows()
	}
	return p.Language
}

// showStartupError mostra o diálogo nativo de erro fatal da inicialização.
func showStartupError(language, message string) {
	desktop.ShowFatal(i18n.T(language, "startup.title", nil), message)
}
