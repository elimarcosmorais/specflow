// Package db abre o SQLite do app e aplica as migrações embutidas (ADR 001, ADR 007).
package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	_ "modernc.org/sqlite" // driver "sqlite", sem CGO
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

// ErrSchemaNewer indica que o banco foi criado por uma versão mais nova do app: o
// PRAGMA user_version é maior que a última migração conhecida.
var ErrSchemaNewer = errors.New("database schema is newer than this app version")

// migrationName reconhece NNNN_<descricao>.sql; o grupo 1 é a versão.
var migrationName = regexp.MustCompile(`^(\d{4})_[A-Za-z0-9_-]+\.sql$`)

// Open abre (ou cria) o banco em path, confere a conexão e aplica as migrações embutidas.
// Em erro, a conexão é fechada.
func Open(ctx context.Context, path string, logger *slog.Logger) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", dsn(path))
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	// Um processo só (ADR 003): uma conexão evita disputa de escrita.
	conn.SetMaxOpenConns(1)

	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	migrations, err := fs.Sub(embeddedMigrations, "migrations")
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("embedded migrations: %w", err)
	}
	version, err := Migrate(ctx, conn, migrations, logger)
	if err != nil {
		conn.Close()
		return nil, err
	}
	logger.Info("db.opened", "path", path, "schema_version", version)
	return conn, nil
}

// dsn monta a string de conexão. Os caracteres que a sintaxe de URI do SQLite trata como
// especiais no caminho (%, ? e #) são escapados.
func dsn(path string) string {
	escaped := strings.NewReplacer("%", "%25", "?", "%3f", "#", "%23").Replace(filepath.ToSlash(path))
	return "file:" + escaped + "?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"
}

type migration struct {
	version int
	file    string
}

// Migrate aplica, em ordem, as migrações de migrations (arquivos NNNN_<descricao>.sql na raiz)
// com versão maior que o PRAGMA user_version, uma transação por arquivo. Devolve a versão
// do esquema ao final.
func Migrate(ctx context.Context, db *sql.DB, migrations fs.FS, logger *slog.Logger) (int, error) {
	known, err := listMigrations(migrations)
	if err != nil {
		return 0, err
	}

	var current int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&current); err != nil {
		return 0, fmt.Errorf("read user_version: %w", err)
	}
	if current > len(known) {
		return current, fmt.Errorf("%w: database is at version %d, latest known is %d", ErrSchemaNewer, current, len(known))
	}

	for _, m := range known[current:] {
		if err := apply(ctx, db, migrations, m); err != nil {
			return current, err
		}
		current = m.version
		logger.Info("db.migration_applied", "version", m.version, "file", m.file)
	}
	return current, nil
}

// listMigrations lê e valida a numeração: 1, 2, 3… sem buraco nem repetição.
func listMigrations(migrations fs.FS) ([]migration, error) {
	entries, err := fs.ReadDir(migrations, ".")
	if err != nil {
		return nil, fmt.Errorf("read migrations: %w", err)
	}
	var known []migration
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		match := migrationName.FindStringSubmatch(e.Name())
		if match == nil {
			return nil, fmt.Errorf("migration %q: name must be NNNN_<description>.sql", e.Name())
		}
		version, _ := strconv.Atoi(match[1])
		known = append(known, migration{version: version, file: e.Name()})
	}
	sort.Slice(known, func(i, j int) bool { return known[i].version < known[j].version })
	for i, m := range known {
		if want := i + 1; m.version != want {
			return nil, fmt.Errorf("migration %q: expected version %d (numbering must be sequential, without gaps or repeats)", m.file, want)
		}
	}
	return known, nil
}

// apply executa uma migração e grava a nova versão na mesma transação.
func apply(ctx context.Context, db *sql.DB, migrations fs.FS, m migration) (err error) {
	script, err := fs.ReadFile(migrations, m.file)
	if err != nil {
		return fmt.Errorf("read migration %q: %w", m.file, err)
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("migration %q: begin: %w", m.file, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if _, err = tx.ExecContext(ctx, string(script)); err != nil {
		return fmt.Errorf("migration %q: %w", m.file, err)
	}
	// PRAGMA não aceita parâmetro; a versão vem de um inteiro validado.
	if _, err = tx.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", m.version)); err != nil {
		return fmt.Errorf("migration %q: set user_version: %w", m.file, err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("migration %q: commit: %w", m.file, err)
	}
	return nil
}

// Pinger verifica se o banco responde (usado pela saúde).
type Pinger struct{ DB *sql.DB }

// Ping devolve erro se o banco é nulo ou não responde.
func (p Pinger) Ping(ctx context.Context) error {
	if p.DB == nil {
		return errors.New("database is not available")
	}
	return p.DB.PingContext(ctx)
}
