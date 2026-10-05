//go:build integration

package log

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"
)

const backupTimeFormat = "2006-01-02T15-04-05.000"

func backupName(daysAgo int) string {
	return "dev-tools-" + time.Now().AddDate(0, 0, -daysAgo).Format(backupTimeFormat) + ".log"
}

// listBackups devolve os nomes dos backups (dev-tools-*.log) em dir, em ordem alfabética.
func listBackups(t *testing.T, dir string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(dir, "dev-tools-*.log"))
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, len(matches))
	for i, m := range matches {
		names[i] = filepath.Base(m)
	}
	sort.Strings(names)
	return names
}

// waitFor repete cond até ela ser verdadeira ou o prazo de 5 s acabar (o descarte do
// lumberjack é assíncrono).
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("prazo de 5 s esgotado esperando: %s", what)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// CT-11
func TestNew_FormatoTexto(t *testing.T) {
	logFile := filepath.Join(tempDir(t), "logs", "dev-tools.log")

	logger, closer, err := New(logFile, slog.LevelInfo, nil)
	if err != nil {
		t.Fatalf("New: erro inesperado: %v", err)
	}
	logger.Info("project.created", "id", 1)
	if err := closer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	lines := strings.Split(strings.TrimRight(readFile(t, logFile), "\r\n"), "\n")
	if len(lines) != 1 {
		t.Fatalf("o arquivo tem %d linhas, quer 1: %q", len(lines), lines)
	}
	want := regexp.MustCompile(`^time=\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\S* level=INFO msg=project\.created id=1\r?$`)
	if !want.MatchString(lines[0]) {
		t.Errorf("linha %q não bate com %s", lines[0], want)
	}
}

func TestNew_NivelFiltraEEscreveNoConsole(t *testing.T) {
	logFile := filepath.Join(tempDir(t), "dev-tools.log")
	var console bytes.Buffer

	logger, closer, err := New(logFile, slog.LevelWarn, &console)
	if err != nil {
		t.Fatalf("New: erro inesperado: %v", err)
	}
	logger.Info("ignorada")
	logger.Warn("mantida", "k", "v")
	if err := closer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	for name, got := range map[string]string{"arquivo": readFile(t, logFile), "console": console.String()} {
		if strings.Contains(got, "ignorada") {
			t.Errorf("%s contém a mensagem abaixo do nível: %q", name, got)
		}
		if !strings.Contains(got, "level=WARN msg=mantida k=v") {
			t.Errorf("%s = %q, quer conter %q", name, got, "level=WARN msg=mantida k=v")
		}
	}
}

func TestNew_ArquivoVazioErro(t *testing.T) {
	logger, closer, err := New("", slog.LevelInfo, nil)
	if err == nil {
		t.Fatal("New(\"\") devolveu erro nil")
	}
	if logger != nil || closer != nil {
		t.Errorf("com erro, logger=%v e closer=%v deveriam ser nil", logger, closer)
	}
}

// CT-12
func TestNew_DescartePorIdadeEQuantidade(t *testing.T) {
	dir := tempDir(t)
	logFile := filepath.Join(dir, "dev-tools.log")
	seeded := map[int]string{}
	for _, days := range []int{1, 2, 3, 4, 5, 6, 20} {
		seeded[days] = backupName(days)
		writeFile(t, filepath.Join(dir, seeded[days]), "antigo\n")
	}

	logger, closer, err := New(logFile, slog.LevelInfo, nil)
	if err != nil {
		t.Fatalf("New: erro inesperado: %v", err)
	}
	// Passa de 10 MB: o lumberjack rotaciona por tamanho e dispara o descarte.
	chunk := strings.Repeat("x", 1024)
	for i := 0; i < 11*1024; i++ {
		logger.Info("enchendo", "chunk", chunk)
	}
	if err := closer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	waitFor(t, "descartar o backup de 20 dias e deixar 5 backups", func() bool {
		return len(listBackups(t, dir)) == 5
	})

	// Restam os 5 mais recentes: o criado pela rotação por tamanho e os de 1 a 4 dias atrás.
	got := listBackups(t, dir)
	isSeeded := map[string]bool{}
	for _, name := range seeded {
		isSeeded[name] = true
	}
	var fresh []string
	for _, name := range got {
		if !isSeeded[name] {
			fresh = append(fresh, name)
		}
	}
	if len(fresh) != 1 {
		t.Errorf("backups novos = %v, quer exatamente 1 (rotação por tamanho); restaram %v", fresh, got)
	}
	for _, days := range []int{1, 2, 3, 4} {
		if _, err := os.Stat(filepath.Join(dir, seeded[days])); err != nil {
			t.Errorf("backup de %d dia(s) atrás deveria ter ficado: %v", days, err)
		}
	}
	for _, days := range []int{5, 6, 20} {
		if _, err := os.Stat(filepath.Join(dir, seeded[days])); !os.IsNotExist(err) {
			t.Errorf("backup de %d dias atrás deveria ter sido descartado (Stat err=%v)", days, err)
		}
	}
}

// CT-13
func TestNew_RotacaoAoIniciar(t *testing.T) {
	dir := tempDir(t)
	logFile := filepath.Join(dir, "dev-tools.log")
	oldBackup := filepath.Join(dir, backupName(20))
	writeFile(t, oldBackup, "muito antigo\n")
	writeFile(t, logFile, "de ontem\n")
	yesterday := time.Now().AddDate(0, 0, -1)
	if err := os.Chtimes(logFile, yesterday, yesterday); err != nil {
		t.Fatal(err)
	}

	logger, closer, err := New(logFile, slog.LevelInfo, nil)
	if err != nil {
		t.Fatalf("New: erro inesperado: %v", err)
	}
	defer closer.Close()

	waitFor(t, "descartar o backup de 20 dias", func() bool {
		_, err := os.Stat(oldBackup)
		return os.IsNotExist(err)
	})
	backups := listBackups(t, dir)
	if len(backups) != 1 {
		t.Fatalf("backups = %v, quer exatamente 1", backups)
	}
	if got := readFile(t, filepath.Join(dir, backups[0])); got != "de ontem\n" {
		t.Errorf("conteúdo do backup = %q, quer %q", got, "de ontem\n")
	}
	if got := readFile(t, logFile); got != "" {
		t.Errorf("dev-tools.log novo = %q, quer vazio", got)
	}

	logger.Info("depois")
	if got := readFile(t, logFile); !strings.Contains(got, "msg=depois") {
		t.Errorf("dev-tools.log = %q, quer conter msg=depois", got)
	}
}

func TestNew_SemRotacaoQuandoOArquivoEDeHoje(t *testing.T) {
	dir := tempDir(t)
	logFile := filepath.Join(dir, "dev-tools.log")
	writeFile(t, logFile, "de hoje\n")

	logger, closer, err := New(logFile, slog.LevelInfo, nil)
	if err != nil {
		t.Fatalf("New: erro inesperado: %v", err)
	}
	logger.Info("depois")
	if err := closer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if backups := listBackups(t, dir); len(backups) != 0 {
		t.Errorf("backups = %v, quer nenhum", backups)
	}
	if got := readFile(t, logFile); !strings.HasPrefix(got, "de hoje\n") || !strings.Contains(got, "msg=depois") {
		t.Errorf("dev-tools.log = %q, quer o conteúdo de hoje preservado e msg=depois acrescentada", got)
	}
}

// RN-19: falha ao rotacionar não impede o log. No Windows, um arquivo aberto por outro
// processo sem compartilhar exclusão não pode ser renomeado.
func TestNew_FalhaAoRotacionarSegueGravando(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("depende de o Windows recusar o rename de arquivo aberto")
	}
	dir := tempDir(t)
	logFile := filepath.Join(dir, "dev-tools.log")
	writeFile(t, logFile, "de ontem\n")
	yesterday := time.Now().AddDate(0, 0, -1)
	if err := os.Chtimes(logFile, yesterday, yesterday); err != nil {
		t.Fatal(err)
	}
	held, err := os.Open(logFile)
	if err != nil {
		t.Fatal(err)
	}
	defer held.Close()

	logger, closer, err := New(logFile, slog.LevelInfo, nil)
	if err != nil {
		t.Fatalf("New: erro inesperado: %v", err)
	}
	logger.Info("depois")
	if err := closer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	held.Close()

	got := readFile(t, logFile)
	if !strings.HasPrefix(got, "de ontem\n") {
		t.Errorf("o conteúdo anterior deveria continuar no arquivo: %q", got)
	}
	if !regexp.MustCompile(`level=WARN msg=log\.rotate_failed err=`).MatchString(got) {
		t.Errorf("dev-tools.log = %q, quer conter level=WARN msg=log.rotate_failed err=", got)
	}
	if !strings.Contains(got, "msg=depois") {
		t.Errorf("dev-tools.log = %q, quer conter msg=depois", got)
	}
	if backups := listBackups(t, dir); len(backups) != 0 {
		t.Errorf("backups = %v, quer nenhum", backups)
	}
}

// tempDir é um t.TempDir() que tenta de novo ao limpar: no Windows, antivírus e o descarte
// assíncrono do lumberjack podem segurar um arquivo por instantes.
func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "dev-tools-log-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		deadline := time.Now().Add(5 * time.Second)
		for {
			err := os.RemoveAll(dir)
			if err == nil || time.Now().After(deadline) {
				if err != nil {
					t.Errorf("limpar %s: %v", dir, err)
				}
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	})
	return dir
}
