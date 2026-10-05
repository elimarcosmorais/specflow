//go:build integration

package projects_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"dev-tools/internal/apperr"
	"dev-tools/internal/features/preferences"
	"dev-tools/internal/features/projects"
	"dev-tools/internal/platform/config"
	"dev-tools/internal/platform/db"
	"dev-tools/internal/platform/httpserver"
)

const (
	testToken = "test-token"
	bearer    = "Bearer " + testToken
)

var (
	t1300 = time.Date(2026, 10, 5, 13, 0, 0, 0, time.UTC)
	t1400 = time.Date(2026, 10, 5, 14, 0, 0, 0, time.UTC)
)

// stack é a pilha real da spec §10: pastas temporárias, SQLite real, serviço, repositório e
// handler HTTP com o token "test-token". language e now podem ser mudados pelo teste.
type stack struct {
	paths    config.Paths
	conn     *sql.DB
	repo     *projects.SQLiteRepository
	handler  http.Handler
	logs     *bytes.Buffer
	language string
	now      time.Time
}

func newStack(t *testing.T) *stack {
	t.Helper()
	return openStack(t, filepath.Join(t.TempDir(), "programs"))
}

// openStack monta a pilha sobre a pasta programs (a raiz da instalação é a pasta acima dela).
func openStack(t *testing.T, programs string) *stack {
	t.Helper()
	paths := config.NewPaths(programs)
	if _, err := config.EnsureLayout(paths); err != nil {
		t.Fatal(err)
	}
	logs := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(logs, nil))
	conn, err := db.Open(context.Background(), paths.DBFile, logger)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() }) // Close é idempotente; alguns testes fecham antes

	s := &stack{paths: paths, conn: conn, logs: logs, language: "pt-br", now: t1300}
	s.repo = projects.NewSQLiteRepository(conn)
	svc := projects.NewService(s.repo, paths.Root, logger, func() time.Time { return s.now })
	cfg := httpserver.Config{
		Token:    testToken,
		Language: func() string { return s.language },
		Logger:   logger,
	}
	s.handler = httpserver.NewHandler(cfg, projects.Routes(svc))
	return s
}

func (s *stack) do(method, target, body, authorization string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)
	return rec
}

// api faz a requisição com o token certo.
func (s *stack) api(method, target, body string) *httptest.ResponseRecorder {
	return s.do(method, target, body, bearer)
}

// folder cria uma pasta de projeto em Projects\<parts> e devolve o caminho.
func (s *stack) folder(t *testing.T, parts ...string) string {
	t.Helper()
	dir := filepath.Join(append([]string{s.paths.Projects}, parts...)...)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// create faz POST /api/projects e confere o 201.
func (s *stack) create(t *testing.T, name, path, typ string) projectJSON {
	t.Helper()
	rec := s.api(http.MethodPost, "/api/projects", payload(name, path, typ))
	requireStatus(t, rec, http.StatusCreated)
	return decodeProject(t, rec)
}

func (s *stack) list(t *testing.T) []projectJSON {
	t.Helper()
	rec := s.api(http.MethodGet, "/api/projects", "")
	requireStatus(t, rec, http.StatusOK)
	var resp struct {
		Items []projectJSON `json:"items"`
	}
	decodeStrict(t, rec, &resp)
	return resp.Items
}

func (s *stack) requireEmptyList(t *testing.T) {
	t.Helper()
	rec := s.api(http.MethodGet, "/api/projects", "")
	requireStatus(t, rec, http.StatusOK)
	if got, want := rec.Body.String(), `{"items":[]}`; got != want {
		t.Errorf("lista = %s, quer %s", got, want)
	}
}

// rows devolve o conteúdo da tabela projects, uma string por linha, ordenado por id.
func (s *stack) rows(t *testing.T) []string {
	t.Helper()
	rs, err := s.conn.Query("SELECT id, name, name_key, path, path_key, type, created_at, updated_at FROM projects ORDER BY id")
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Close()
	var out []string
	for rs.Next() {
		var id int64
		var name, nameKey, path, pathKey, typ, createdAt, updatedAt string
		if err := rs.Scan(&id, &name, &nameKey, &path, &pathKey, &typ, &createdAt, &updatedAt); err != nil {
			t.Fatal(err)
		}
		out = append(out, fmt.Sprint(id, "|", name, "|", nameKey, "|", path, "|", pathKey, "|", typ, "|", createdAt, "|", updatedAt))
	}
	if err := rs.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

// logLines devolve as linhas do log que contêm todos os trechos.
func (s *stack) logLines(substrs ...string) []string {
	var found []string
	for _, line := range strings.Split(s.logs.String(), "\n") {
		if line != "" && !slices.ContainsFunc(substrs, func(sub string) bool { return !strings.Contains(line, sub) }) {
			found = append(found, line)
		}
	}
	return found
}

// projectJSON é o Project da API (spec §3.1). Os testes o leem sem aceitar campo extra.
type projectJSON struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	Type         string `json:"type"`
	FolderExists bool   `json:"folderExists"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
}

func payload(name, path, typ string) string {
	b, err := json.Marshal(map[string]string{"name": name, "path": path, "type": typ})
	if err != nil {
		panic(err)
	}
	return string(b)
}

func decodeStrict(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(rec.Body.Bytes()))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		t.Fatalf("corpo inválido: %v\n%s", err, rec.Body)
	}
}

func decodeProject(t *testing.T, rec *httptest.ResponseRecorder) projectJSON {
	t.Helper()
	var p projectJSON
	decodeStrict(t, rec, &p)
	return p
}

func requireStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status = %d, quer %d\n%s", rec.Code, want, rec.Body)
	}
}

// requireError confere status, código e campo do erro, e devolve o erro para outras conferências.
func requireError(t *testing.T, rec *httptest.ResponseRecorder, status int, code, field string) apiError {
	t.Helper()
	requireStatus(t, rec, status)
	var resp struct {
		Error apiError `json:"error"`
	}
	decodeStrict(t, rec, &resp)
	if resp.Error.Code != code || resp.Error.Field != field {
		t.Errorf("erro = {%s, %s}, quer {%s, %s}", resp.Error.Code, resp.Error.Field, code, field)
	}
	if resp.Error.Message == "" {
		t.Errorf("erro %s sem mensagem", code)
	}
	return resp.Error
}

// snapshot é o retrato de uma árvore de pastas: caminho relativo, tipo e conteúdo de cada entrada.
func snapshot(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if d.IsDir() {
			out = append(out, rel+"|dir")
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out = append(out, rel+"|file|"+string(data))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// CT-52
func TestProjects_CRUDViaAPI(t *testing.T) {
	s := newStack(t)
	folder := s.folder(t, "web", "loja-api")
	writeFile(t, filepath.Join(folder, "index.php"), "<?php echo 1;")
	before := snapshot(t, s.paths.Projects)

	rec := s.api(http.MethodPost, "/api/projects", payload("loja-api", folder, "web"))
	requireStatus(t, rec, http.StatusCreated)
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q", ct)
	}
	wantCreated := projectJSON{ID: 1, Name: "loja-api", Path: folder, Type: "web", FolderExists: true,
		CreatedAt: "2026-10-05T13:00:00Z", UpdatedAt: "2026-10-05T13:00:00Z"}
	if got := decodeProject(t, rec); got != wantCreated {
		t.Errorf("POST = %+v, quer %+v", got, wantCreated)
	}
	// D-07: a pasta dentro da raiz da instalação é guardada em forma relativa.
	if want := []string{"1|loja-api|loja-api|projects\\web\\loja-api|projects\\web\\loja-api|web|2026-10-05T13:00:00Z|2026-10-05T13:00:00Z"}; !slices.Equal(s.rows(t), want) {
		t.Errorf("linhas = %q, quer %q", s.rows(t), want)
	}

	if got := s.list(t); !slices.Equal(got, []projectJSON{wantCreated}) {
		t.Errorf("GET = %+v, quer [%+v]", got, wantCreated)
	}

	s.now = t1400
	rec = s.api(http.MethodPut, "/api/projects/1", payload("loja", folder, "android"))
	requireStatus(t, rec, http.StatusOK)
	wantUpdated := projectJSON{ID: 1, Name: "loja", Path: folder, Type: "android", FolderExists: true,
		CreatedAt: "2026-10-05T13:00:00Z", UpdatedAt: "2026-10-05T14:00:00Z"}
	if got := decodeProject(t, rec); got != wantUpdated {
		t.Errorf("PUT = %+v, quer %+v", got, wantUpdated)
	}
	if got := s.list(t); !slices.Equal(got, []projectJSON{wantUpdated}) {
		t.Errorf("GET depois do PUT = %+v, quer [%+v]", got, wantUpdated)
	}

	rec = s.api(http.MethodDelete, "/api/projects/1", "")
	requireStatus(t, rec, http.StatusNoContent)
	if rec.Body.Len() != 0 {
		t.Errorf("DELETE devolveu corpo: %q", rec.Body)
	}
	s.requireEmptyList(t)

	if after := snapshot(t, s.paths.Projects); !slices.Equal(after, before) {
		t.Errorf("a pasta mudou:\nantes:  %q\ndepois: %q", before, after)
	}
}

func TestProjects_ListIsSortedByNameIgnoringCase(t *testing.T) {
	s := newStack(t)
	s.create(t, "beta", s.folder(t, "web", "beta"), "web")
	s.create(t, "Alfa", s.folder(t, "web", "alfa"), "web")
	s.create(t, "carro", s.folder(t, "web", "carro"), "web")

	var names []string
	for _, p := range s.list(t) {
		names = append(names, p.Name)
	}
	if want := []string{"Alfa", "beta", "carro"}; !slices.Equal(names, want) {
		t.Errorf("nomes = %q, quer %q", names, want)
	}
}

// CT-53
func TestProjects_Create_PathNotFound(t *testing.T) {
	s := newStack(t)
	missing := filepath.Join(s.paths.Projects, "web", "nao-existe")

	rec := s.api(http.MethodPost, "/api/projects", payload("loja-api", missing, "web"))

	e := requireError(t, rec, http.StatusUnprocessableEntity, "project.path_not_found", "path")
	if want := "A pasta informada não existe."; e.Message != want {
		t.Errorf("mensagem = %q, quer %q", e.Message, want)
	}
	s.requireEmptyList(t)
}

// CT-54
func TestProjects_Create_Duplicates(t *testing.T) {
	s := newStack(t)
	folderA := s.folder(t, "web", "loja-api")
	folderB := s.folder(t, "web", "outra")
	s.create(t, "loja-api", folderA, "web")

	rec := s.api(http.MethodPost, "/api/projects", payload("Loja-API", folderB, "web"))
	requireError(t, rec, http.StatusConflict, "project.name_taken", "name")

	rec = s.api(http.MethodPost, "/api/projects", payload("outro-nome", folderA, "web"))
	requireError(t, rec, http.StatusConflict, "project.path_taken", "path")

	got := s.list(t)
	if len(got) != 1 || got[0].Name != "loja-api" || got[0].Path != folderA {
		t.Errorf("lista = %+v, quer só loja-api em %s", got, folderA)
	}
}

// CT-55
func TestProjects_Create_RequiredFields(t *testing.T) {
	for name, body := range map[string]string{
		"campos vazios":   `{"name":"","path":"","type":""}`,
		"objeto sem nada": `{}`,
	} {
		t.Run(name, func(t *testing.T) {
			s := newStack(t)

			rec := s.api(http.MethodPost, "/api/projects", body)

			requireError(t, rec, http.StatusUnprocessableEntity, "project.name_required", "name")
			s.requireEmptyList(t)
		})
	}
}

func TestProjects_InvalidBodyAndID(t *testing.T) {
	s := newStack(t)
	folder := s.folder(t, "web", "loja-api")
	s.create(t, "loja-api", folder, "web")
	before := s.rows(t)
	valid := payload("outra", folder, "web")

	tests := []struct {
		name, method, target, body string
		status                     int
		code                       string
	}{
		{"POST com campo desconhecido", http.MethodPost, "/api/projects", `{"name":"x","path":"y","type":"web","extra":1}`, 400, "request.invalid_body"},
		{"POST com JSON malformado", http.MethodPost, "/api/projects", `{"name":`, 400, "request.invalid_body"},
		{"PUT com campo desconhecido", http.MethodPut, "/api/projects/1", `{"name":"x","path":"y","type":"web","extra":1}`, 400, "request.invalid_body"},
		{"PUT de id que não existe", http.MethodPut, "/api/projects/99", valid, 404, "project.not_found"},
		{"DELETE de id que não existe", http.MethodDelete, "/api/projects/99", "", 404, "project.not_found"},
		{"PUT com id não numérico", http.MethodPut, "/api/projects/abc", valid, 404, "project.not_found"},
		{"PUT com id zero", http.MethodPut, "/api/projects/0", valid, 404, "project.not_found"},
		{"PUT com id negativo", http.MethodPut, "/api/projects/-1", valid, 404, "project.not_found"},
		{"DELETE com id com sinal", http.MethodDelete, "/api/projects/+1", "", 404, "project.not_found"},
		{"DELETE com id que estoura int64", http.MethodDelete, "/api/projects/99999999999999999999", "", 404, "project.not_found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := s.api(tt.method, tt.target, tt.body)

			requireError(t, rec, tt.status, tt.code, "")
			if after := s.rows(t); !slices.Equal(after, before) {
				t.Errorf("linhas mudaram: %q → %q", before, after)
			}
		})
	}
}

// CT-56
func TestProjects_Update_Refused(t *testing.T) {
	s := newStack(t)
	folderA := s.folder(t, "web", "a")
	folderB := s.folder(t, "web", "b")
	a := s.create(t, "proj-a", folderA, "web")
	s.create(t, "proj-b", folderB, "android")
	before := s.list(t)
	target := fmt.Sprintf("/api/projects/%d", a.ID)

	rec := s.api(http.MethodPut, target, payload("PROJ-B", folderA, "web"))
	requireError(t, rec, http.StatusConflict, "project.name_taken", "name")

	rec = s.api(http.MethodPut, target, payload("proj-a", filepath.Join(s.paths.Projects, "web", "nao-existe"), "web"))
	requireError(t, rec, http.StatusUnprocessableEntity, "project.path_not_found", "path")

	if after := s.list(t); !slices.Equal(after, before) {
		t.Errorf("lista mudou:\nantes:  %+v\ndepois: %+v", before, after)
	}
	if after := s.list(t); after[0].Name != "proj-a" || after[0].Path != folderA || after[0].Type != "web" {
		t.Errorf("projeto A = %+v, quer proj-a em %s, web", after[0], folderA)
	}
}

// CT-57
func TestProjects_FolderDeletedAfterRegistration(t *testing.T) {
	s := newStack(t)
	folder := s.folder(t, "web", "loja-api")
	other := s.folder(t, "web", "nova")
	p := s.create(t, "loja-api", folder, "web")
	if err := os.RemoveAll(folder); err != nil {
		t.Fatal(err)
	}

	got := s.list(t)
	if len(got) != 1 || got[0].FolderExists || got[0].Path != folder {
		t.Errorf("GET = %+v, quer 1 item com folderExists false em %s", got, folder)
	}

	target := fmt.Sprintf("/api/projects/%d", p.ID)
	rec := s.api(http.MethodPut, target, payload("loja-api", other, "web"))
	requireStatus(t, rec, http.StatusOK)
	if updated := decodeProject(t, rec); updated.Path != other || !updated.FolderExists {
		t.Errorf("PUT = %+v, quer path %s e folderExists true", updated, other)
	}

	requireStatus(t, s.api(http.MethodDelete, target, ""), http.StatusNoContent)
	s.requireEmptyList(t)
}

// CT-58
func TestProjects_RequireToken(t *testing.T) {
	s := newStack(t)
	folder := s.folder(t, "web", "loja-api")
	other := s.folder(t, "web", "outra")
	s.create(t, "loja-api", folder, "web")
	before := s.rows(t)

	requests := []struct{ method, target, body string }{
		{http.MethodGet, "/api/projects", ""},
		{http.MethodPost, "/api/projects", payload("nova", other, "web")},
		{http.MethodPut, "/api/projects/1", payload("renomeada", other, "android")},
		{http.MethodDelete, "/api/projects/1", ""},
	}
	for _, auth := range []string{"", "Bearer errado"} {
		for _, r := range requests {
			t.Run(fmt.Sprintf("%s %s auth=%q", r.method, r.target, auth), func(t *testing.T) {
				rec := s.do(r.method, r.target, r.body, auth)

				requireError(t, rec, http.StatusUnauthorized, "auth.token_invalid", "")
				if after := s.rows(t); !slices.Equal(after, before) {
					t.Errorf("linhas mudaram: %q → %q", before, after)
				}
			})
		}
	}

	var count int
	if err := s.conn.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("COUNT(*) = %d, quer 1", count)
	}
}

// CT-59
func TestProjects_DatabaseUnavailable(t *testing.T) {
	s := newStack(t)
	s.create(t, "loja-api", s.folder(t, "web", "loja-api"), "web")
	next := payload("nova", s.folder(t, "web", "nova"), "web")
	if err := s.conn.Close(); err != nil {
		t.Fatal(err)
	}

	const ptMessage = "O banco de dados está indisponível. Reinicie o Dev Tools; se o problema continuar, veja os logs."
	const enMessage = "The database is unavailable. Restart Dev Tools; if the problem persists, see the logs."

	s.language = "pt-br"
	for name, rec := range map[string]*httptest.ResponseRecorder{
		"POST": s.api(http.MethodPost, "/api/projects", next),
		"GET":  s.api(http.MethodGet, "/api/projects", ""),
	} {
		e := requireError(t, rec, http.StatusServiceUnavailable, "database.unavailable", "")
		if e.Message != ptMessage {
			t.Errorf("%s: mensagem = %q, quer %q", name, e.Message, ptMessage)
		}
		if strings.Contains(rec.Body.String(), "closed") {
			t.Errorf("%s: a causa vazou para o cliente: %s", name, rec.Body)
		}
	}

	s.language = "en"
	e := requireError(t, s.api(http.MethodPost, "/api/projects", next), http.StatusServiceUnavailable, "database.unavailable", "")
	if e.Message != enMessage {
		t.Errorf("mensagem em en = %q, quer %q", e.Message, enMessage)
	}

	lines := s.logLines("level=ERROR", "msg=http.error", "status=503", "code=database.unavailable", "database is closed")
	if len(lines) != 3 {
		t.Errorf("%d registros http.error (ERROR) com a causa, quer 3\n%s", len(lines), s.logs)
	}
}

func TestSQLiteRepository_NilDatabase(t *testing.T) {
	repo := projects.NewSQLiteRepository(nil)
	ctx := context.Background()
	p := projects.Project{ID: 1, Name: "a", NameKey: "a", Path: "p", PathKey: "p", Type: projects.TypeWeb}

	calls := map[string]func() error{
		"List":       func() error { _, err := repo.List(ctx); return err },
		"Get":        func() error { _, err := repo.Get(ctx, 1); return err },
		"ExistsName": func() error { _, err := repo.ExistsName(ctx, "a", 0); return err },
		"ExistsPath": func() error { _, err := repo.ExistsPath(ctx, "p", 0); return err },
		"Create":     func() error { _, err := repo.Create(ctx, p); return err },
		"Update":     func() error { return repo.Update(ctx, p) },
		"Delete":     func() error { return repo.Delete(ctx, 1) },
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			err := call()

			var ae *apperr.Error
			if !errors.As(err, &ae) || ae.Code != "database.unavailable" || ae.Kind != apperr.Unavailable {
				t.Fatalf("erro = %v, quer database.unavailable", err)
			}
			if ae.Err == nil {
				t.Error("erro sem a causa")
			}
		})
	}
}

func TestSQLiteRepository_NilDatabaseViaAPI(t *testing.T) {
	logs := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(logs, nil))
	svc := projects.NewService(projects.NewSQLiteRepository(nil), t.TempDir(), logger, time.Now)
	handler := httpserver.NewHandler(httpserver.Config{Token: testToken, Logger: logger}, projects.Routes(svc))

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.Header.Set("Authorization", bearer)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	requireError(t, rec, http.StatusServiceUnavailable, "database.unavailable", "")
}

// CT-60
func TestProjects_PortableInstallation(t *testing.T) {
	tmp := t.TempDir()
	rootA := filepath.Join(tmp, "A")
	rootB := filepath.Join(tmp, "B")
	external := filepath.Join(tmp, "externo", "site")
	if err := os.MkdirAll(external, 0o755); err != nil {
		t.Fatal(err)
	}

	a := openStack(t, filepath.Join(rootA, "programs"))
	prefs := preferences.NewService(preferences.NewFileRepository(a.paths.ConfigFile), a.paths.ConfigFile,
		func() string { return "pt-br" }, slog.New(slog.DiscardHandler), nil)
	if err := prefs.Update(preferences.Preferences{Theme: "dark", Language: "en"}); err != nil {
		t.Fatal(err)
	}
	a.create(t, "loja-api", a.folder(t, "web", "loja-api"), "web")
	a.create(t, "site", external, "web")
	if err := a.conn.Close(); err != nil {
		t.Fatal(err)
	}

	if err := copyTree(rootA, rootB); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(rootA); err != nil {
		t.Fatal(err)
	}

	b := openStack(t, filepath.Join(rootB, "programs"))
	got := b.list(t)

	want := []projectJSON{
		{ID: 1, Name: "loja-api", Path: filepath.Join(rootB, "projects", "web", "loja-api"), Type: "web",
			FolderExists: true, CreatedAt: "2026-10-05T13:00:00Z", UpdatedAt: "2026-10-05T13:00:00Z"},
		// Fora da raiz da instalação, o caminho é guardado absoluto e não acompanha a cópia.
		{ID: 2, Name: "site", Path: external, Type: "web",
			FolderExists: true, CreatedAt: "2026-10-05T13:00:00Z", UpdatedAt: "2026-10-05T13:00:00Z"},
	}
	if !slices.Equal(got, want) {
		t.Errorf("GET em B = %+v\nquer     %+v", got, want)
	}

	prefsB := preferences.NewService(preferences.NewFileRepository(b.paths.ConfigFile), b.paths.ConfigFile,
		func() string { return "pt-br" }, slog.New(slog.DiscardHandler), nil)
	if err := prefsB.Load(); err != nil {
		t.Fatal(err)
	}
	if want := (preferences.Preferences{Theme: "dark", Language: "en"}); prefsB.Get() != want {
		t.Errorf("preferências em B = %+v, quer %+v", prefsB.Get(), want)
	}
}

// copyTree copia a árvore src em dst, arquivo por arquivo.
func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, in); err != nil {
			out.Close()
			return err
		}
		return out.Close()
	})
}

// CT-61
func TestProjects_LogsEvents(t *testing.T) {
	s := newStack(t)
	s.create(t, "loja-api", s.folder(t, "web", "loja-api"), "web")
	rec := s.api(http.MethodPost, "/api/projects", payload("loja-api", s.folder(t, "web", "outra"), "web"))
	requireStatus(t, rec, http.StatusConflict)

	created := s.logLines("level=INFO", "msg=project.created")
	if len(created) != 1 || !strings.HasPrefix(created[0], "time=") {
		t.Errorf("registros project.created = %q, quer 1 INFO com time=", created)
	}
	conflict := s.logLines("level=WARN", "msg=http.error", "status=409", "code=project.name_taken")
	if len(conflict) != 1 || !strings.HasPrefix(conflict[0], "time=") {
		t.Errorf("registros http.error 409 = %q, quer 1 WARN com time=", conflict)
	}
}

// CT-62
func TestSQLiteRepository_UniqueIsSafetyNet(t *testing.T) {
	s := newStack(t)
	ctx := context.Background()
	newProject := func(name, path string) projects.Project {
		return projects.Project{Name: name, NameKey: projects.NameKey(name), Path: path, PathKey: projects.PathKey(path),
			Type: projects.TypeWeb, CreatedAt: t1300, UpdatedAt: t1300}
	}
	if _, err := s.repo.Create(ctx, newProject("Loja", `projects\web\loja`)); err != nil {
		t.Fatal(err)
	}

	_, err := s.repo.Create(ctx, newProject("LOJA", `projects\web\outra`))
	if !errors.Is(err, projects.ErrNameTaken) {
		t.Errorf("NameKey repetida: erro = %v, quer project.name_taken", err)
	}
	_, err = s.repo.Create(ctx, newProject("outra", `PROJECTS\web\loja`))
	if !errors.Is(err, projects.ErrPathTaken) {
		t.Errorf("PathKey repetida: erro = %v, quer project.path_taken", err)
	}

	if got := s.rows(t); len(got) != 1 {
		t.Errorf("linhas = %q, quer só a primeira", got)
	}
}

func TestSQLiteRepository_UpdateUniqueViolation(t *testing.T) {
	s := newStack(t)
	ctx := context.Background()
	newProject := func(name, path string) projects.Project {
		return projects.Project{Name: name, NameKey: projects.NameKey(name), Path: path, PathKey: projects.PathKey(path),
			Type: projects.TypeWeb, CreatedAt: t1300, UpdatedAt: t1300}
	}
	if _, err := s.repo.Create(ctx, newProject("a", "pa")); err != nil {
		t.Fatal(err)
	}
	id, err := s.repo.Create(ctx, newProject("b", "pb"))
	if err != nil {
		t.Fatal(err)
	}

	clash := newProject("a", "pb")
	clash.ID = id
	if err := s.repo.Update(ctx, clash); !errors.Is(err, projects.ErrNameTaken) {
		t.Errorf("Update com nome repetido: erro = %v, quer project.name_taken", err)
	}
	clash = newProject("b", "pa")
	clash.ID = id
	if err := s.repo.Update(ctx, clash); !errors.Is(err, projects.ErrPathTaken) {
		t.Errorf("Update com pasta repetida: erro = %v, quer project.path_taken", err)
	}
}

func TestSQLiteRepository_RoundTripAndNotFound(t *testing.T) {
	s := newStack(t)
	ctx := context.Background()
	// Data em fuso diferente e com frações de segundo: volta como o mesmo instante, em UTC.
	zone := time.FixedZone("BRT", -3*60*60)
	created := time.Date(2026, 10, 5, 10, 0, 0, 123456789, zone)
	p := projects.Project{Name: "Loja", NameKey: "loja", Path: `projects\web\loja`, PathKey: `projects\web\loja`,
		Type: projects.TypeAndroid, CreatedAt: created, UpdatedAt: created}

	id, err := s.repo.Create(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Errorf("id = %d, quer 1", id)
	}

	got, err := s.repo.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	p.ID = id
	if got.ID != p.ID || got.Name != p.Name || got.NameKey != p.NameKey || got.Path != p.Path ||
		got.PathKey != p.PathKey || got.Type != p.Type {
		t.Errorf("Get = %+v, quer %+v", got, p)
	}
	if !got.CreatedAt.Equal(created) || !got.UpdatedAt.Equal(created) || got.CreatedAt.Location() != time.UTC {
		t.Errorf("datas = %v e %v, quer o instante %v em UTC", got.CreatedAt, got.UpdatedAt, created)
	}

	all, err := s.repo.List(ctx)
	if err != nil || len(all) != 1 || all[0] != got {
		t.Errorf("List = %+v, %v, quer [%+v]", all, err, got)
	}

	// ExistsName e ExistsPath ignoram o próprio projeto.
	for name, check := range map[string]func(ignoreID int64) (bool, error){
		"ExistsName": func(ignoreID int64) (bool, error) { return s.repo.ExistsName(ctx, "loja", ignoreID) },
		"ExistsPath": func(ignoreID int64) (bool, error) { return s.repo.ExistsPath(ctx, `projects\web\loja`, ignoreID) },
	} {
		if found, err := check(0); err != nil || !found {
			t.Errorf("%s(ignoreID 0) = %v, %v, quer true", name, found, err)
		}
		if found, err := check(id); err != nil || found {
			t.Errorf("%s(ignoreID %d) = %v, %v, quer false", name, id, found, err)
		}
	}
	if found, err := s.repo.ExistsName(ctx, "outro", 0); err != nil || found {
		t.Errorf("ExistsName(outro) = %v, %v, quer false", found, err)
	}

	// Update muda os dados e a data de alteração, e preserva a data de criação.
	p.Name, p.NameKey, p.UpdatedAt = "Loja 2", "loja 2", t1400
	if err := s.repo.Update(ctx, p); err != nil {
		t.Fatal(err)
	}
	if got, err = s.repo.Get(ctx, id); err != nil || got.Name != "Loja 2" || !got.UpdatedAt.Equal(t1400) || !got.CreatedAt.Equal(created) {
		t.Errorf("Get depois do Update = %+v, %v", got, err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		t.Fatal(err)
	}
	if _, err := s.repo.Get(ctx, id); !errors.Is(err, projects.ErrNotFound) {
		t.Errorf("Get depois do Delete: erro = %v, quer project.not_found", err)
	}
	if err := s.repo.Delete(ctx, id); !errors.Is(err, projects.ErrNotFound) {
		t.Errorf("Delete repetido: erro = %v, quer project.not_found", err)
	}
	p.ID = 99
	if err := s.repo.Update(ctx, p); !errors.Is(err, projects.ErrNotFound) {
		t.Errorf("Update de id inexistente: erro = %v, quer project.not_found", err)
	}
}
