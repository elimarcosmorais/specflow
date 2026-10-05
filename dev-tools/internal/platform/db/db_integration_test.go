//go:build integration

package db

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func newLogger() (*slog.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return slog.New(slog.NewTextHandler(buf, nil)), buf
}

func userVersion(t *testing.T, conn *sql.DB) int {
	t.Helper()
	var v int
	if err := conn.QueryRow("PRAGMA user_version").Scan(&v); err != nil {
		t.Fatal(err)
	}
	return v
}

func tableExists(t *testing.T, conn *sql.DB, name string) bool {
	t.Helper()
	var n int
	err := conn.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?", name).Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	return n == 1
}

// CT-08
func TestOpen_AppliesEmbeddedMigrationsOnce(t *testing.T) {
	ctx := context.Background()
	// a pasta tem espaço no nome de propósito: o caminho vai numa URI do SQLite
	dir := filepath.Join(t.TempDir(), "dev tools")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "dev-tools.db")

	logger, buf := newLogger()
	conn, err := Open(ctx, path, logger)
	if err != nil {
		t.Fatalf("1st Open: %v", err)
	}
	if got := userVersion(t, conn); got != 1 {
		t.Errorf("user_version after 1st Open = %d, want 1", got)
	}
	if !tableExists(t, conn, "projects") {
		t.Error("table projects does not exist after 1st Open")
	}
	if got := strings.Count(buf.String(), "db.migration_applied"); got != 1 {
		t.Errorf("db.migration_applied logged %d times on 1st Open, want 1\n%s", got, buf)
	}
	if !strings.Contains(buf.String(), "version=1 file=0001_create_projects.sql") {
		t.Errorf("db.migration_applied without version/file\n%s", buf)
	}
	if !strings.Contains(buf.String(), "db.opened") || !strings.Contains(buf.String(), "schema_version=1") {
		t.Errorf("db.opened with schema_version=1 not logged\n%s", buf)
	}
	if _, err := conn.Exec(`INSERT INTO projects (name, name_key, path, path_key, type, created_at, updated_at)
		VALUES ('A', 'a', 'p', 'p', 'web', 't', 't')`); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}

	logger, buf = newLogger()
	conn, err = Open(ctx, path, logger)
	if err != nil {
		t.Fatalf("2nd Open: %v", err)
	}
	defer conn.Close()
	if got := strings.Count(buf.String(), "db.migration_applied"); got != 0 {
		t.Errorf("db.migration_applied logged %d times on 2nd Open, want 0\n%s", got, buf)
	}
	if got := userVersion(t, conn); got != 1 {
		t.Errorf("user_version after 2nd Open = %d, want 1", got)
	}
	var n int
	if err := conn.QueryRow("SELECT COUNT(*) FROM projects").Scan(&n); err != nil || n != 1 {
		t.Errorf("projects rows after 2nd Open = %d (err %v), want 1", n, err)
	}
}

// §4.1: restrições da tabela projects.
func TestOpen_ProjectsTableConstraints(t *testing.T) {
	conn, err := Open(context.Background(), filepath.Join(t.TempDir(), "x.db"), slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	insert := func(nameKey, pathKey, typ string) error {
		_, err := conn.Exec(`INSERT INTO projects (name, name_key, path, path_key, type, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, 't', 't')`, nameKey, nameKey, pathKey, pathKey, typ)
		return err
	}
	if err := insert("a", "pa", "web"); err != nil {
		t.Fatalf("valid insert: %v", err)
	}
	cases := map[string]struct{ nameKey, pathKey, typ, want string }{
		"duplicate name_key": {"a", "pb", "web", "name_key"},
		"duplicate path_key": {"b", "pa", "web", "path_key"},
		"invalid type":       {"c", "pc", "linux", "CHECK"},
	}
	for name, tc := range cases {
		err := insert(tc.nameKey, tc.pathKey, tc.typ)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: err = %v, want it to mention %q", name, err, tc.want)
		}
	}
}

// CT-08 (numeração) e atomicidade de cada migração.
func TestMigrate(t *testing.T) {
	ctx := context.Background()
	open := func(t *testing.T) *sql.DB {
		conn, err := sql.Open("sqlite", dsn(filepath.Join(t.TempDir(), "m.db")))
		if err != nil {
			t.Fatal(err)
		}
		conn.SetMaxOpenConns(1)
		t.Cleanup(func() { conn.Close() })
		return conn
	}
	logger := slog.New(slog.DiscardHandler)

	t.Run("applies in order and returns the final version", func(t *testing.T) {
		conn := open(t)
		fsys := fstest.MapFS{
			"0001_a.sql": {Data: []byte("CREATE TABLE a (id INTEGER);")},
			"0002_b.sql": {Data: []byte("CREATE TABLE b (a_id INTEGER REFERENCES a(id));")},
		}
		version, err := Migrate(ctx, conn, fsys, logger)
		if err != nil || version != 2 {
			t.Fatalf("Migrate = %d, %v; want 2, nil", version, err)
		}
		if got := userVersion(t, conn); got != 2 || !tableExists(t, conn, "a") || !tableExists(t, conn, "b") {
			t.Errorf("user_version = %d, tables a/b not both present", got)
		}
	})

	t.Run("gap in numbering is an error and applies nothing", func(t *testing.T) {
		conn := open(t)
		fsys := fstest.MapFS{
			"0001_a.sql": {Data: []byte("CREATE TABLE a (id INTEGER);")},
			"0003_c.sql": {Data: []byte("CREATE TABLE c (id INTEGER);")},
		}
		_, err := Migrate(ctx, conn, fsys, logger)
		if err == nil || !strings.Contains(err.Error(), `"0003_c.sql": expected version 2`) {
			t.Fatalf("err = %v, want gap error naming 0003_c.sql", err)
		}
		if got := userVersion(t, conn); got != 0 || tableExists(t, conn, "a") {
			t.Errorf("user_version = %d, table a exists = %v; want nothing applied", got, tableExists(t, conn, "a"))
		}
	})

	t.Run("repeated number is an error", func(t *testing.T) {
		conn := open(t)
		fsys := fstest.MapFS{
			"0001_a.sql": {Data: []byte("CREATE TABLE a (id INTEGER);")},
			"0001_b.sql": {Data: []byte("CREATE TABLE b (id INTEGER);")},
		}
		_, err := Migrate(ctx, conn, fsys, logger)
		if err == nil || !strings.Contains(err.Error(), `"0001_b.sql": expected version 2`) {
			t.Fatalf("err = %v, want repeated-number error naming 0001_b.sql", err)
		}
	})

	t.Run("invalid file name is an error", func(t *testing.T) {
		conn := open(t)
		fsys := fstest.MapFS{"create_a.sql": {Data: []byte("CREATE TABLE a (id INTEGER);")}}
		_, err := Migrate(ctx, conn, fsys, logger)
		if err == nil || !strings.Contains(err.Error(), "NNNN_<description>.sql") {
			t.Fatalf("err = %v, want file name error", err)
		}
	})

	t.Run("failing migration rolls back and keeps the previous version", func(t *testing.T) {
		conn := open(t)
		fsys := fstest.MapFS{
			"0001_a.sql": {Data: []byte("CREATE TABLE a (id INTEGER);")},
			"0002_b.sql": {Data: []byte("CREATE TABLE b (id INTEGER); INSERT INTO missing VALUES (1);")},
		}
		version, err := Migrate(ctx, conn, fsys, logger)
		if err == nil || !strings.Contains(err.Error(), `migration "0002_b.sql"`) {
			t.Fatalf("err = %v, want error naming 0002_b.sql", err)
		}
		if version != 1 {
			t.Errorf("version = %d, want 1", version)
		}
		if got := userVersion(t, conn); got != 1 || tableExists(t, conn, "b") {
			t.Errorf("user_version = %d, table b exists = %v; want 1 and no table b", got, tableExists(t, conn, "b"))
		}
	})
}

// CT-09
func TestOpen_SchemaNewer(t *testing.T) {
	path := filepath.Join(t.TempDir(), "newer.db")
	raw, err := sql.Open("sqlite", dsn(path))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := raw.Exec("PRAGMA user_version = 99"); err != nil {
		t.Fatal(err)
	}
	if err := raw.Close(); err != nil {
		t.Fatal(err)
	}

	logger, buf := newLogger()
	conn, err := Open(context.Background(), path, logger)
	if !errors.Is(err, ErrSchemaNewer) {
		t.Fatalf("Open err = %v, want ErrSchemaNewer", err)
	}
	if conn != nil {
		conn.Close()
		t.Error("Open returned a connection together with the error")
	}
	if strings.Contains(buf.String(), "db.opened") {
		t.Errorf("db.opened logged for a rejected database\n%s", buf)
	}

	// o arquivo continua intacto
	raw, err = sql.Open("sqlite", dsn(path))
	if err != nil {
		t.Fatal(err)
	}
	defer raw.Close()
	if got := userVersion(t, raw); got != 99 {
		t.Errorf("user_version = %d, want 99 (untouched)", got)
	}
	if tableExists(t, raw, "projects") {
		t.Error("migrations were applied to a newer database")
	}
}

// CT-10
func TestPinger(t *testing.T) {
	ctx := context.Background()
	conn, err := Open(ctx, filepath.Join(t.TempDir(), "ping.db"), slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	if err := (Pinger{DB: conn}).Ping(ctx); err != nil {
		t.Errorf("Ping on open database = %v, want nil", err)
	}

	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	if err := (Pinger{DB: conn}).Ping(ctx); err == nil || !strings.Contains(err.Error(), "closed") {
		t.Errorf("Ping on closed database = %v, want a \"closed\" error", err)
	}

	if err := (Pinger{DB: nil}).Ping(ctx); err == nil || err.Error() != "database is not available" {
		t.Errorf("Ping with nil DB = %v, want \"database is not available\"", err)
	}
}
