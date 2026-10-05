package projects

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"dev-tools/internal/apperr"
)

var (
	createdAt = time.Date(2026, 10, 5, 13, 0, 0, 0, time.UTC)
	updatedAt = time.Date(2026, 10, 6, 9, 30, 0, 0, time.UTC)
)

// env é o cenário de um teste: a raiz da instalação (uma pasta temporária que contém
// projects\web\loja-api), o fake do repositório, o serviço e o buffer do logger. now é o relógio
// injetado; o teste o muda para simular a passagem do tempo.
type env struct {
	root string
	repo *fakeRepository
	svc  *Service
	logs *bytes.Buffer
	now  time.Time
}

func newEnv(t *testing.T) *env {
	t.Helper()
	e := &env{root: t.TempDir(), repo: newFakeRepository(), logs: &bytes.Buffer{}, now: createdAt}
	logger := slog.New(slog.NewTextHandler(e.logs, nil))
	e.svc = NewService(e.repo, e.root, logger, func() time.Time { return e.now })
	e.mkdir(t, "projects", "web", "loja-api")
	return e
}

// mkdir cria uma pasta dentro da raiz da instalação e devolve o caminho absoluto.
func (e *env) mkdir(t *testing.T, elem ...string) string {
	t.Helper()
	dir := filepath.Join(append([]string{e.root}, elem...)...)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", dir, err)
	}
	return dir
}

// create cadastra um projeto que tem de ser válido.
func (e *env) create(t *testing.T, name, path string, typ Type) View {
	t.Helper()
	v, err := e.svc.Create(context.Background(), Input{Name: name, Path: path, Type: typ})
	if err != nil {
		t.Fatalf("Create(%q): %v", name, err)
	}
	return v
}

// wantError confere o código, o tipo e o campo do erro, e não só que houve erro.
func wantError(t *testing.T, got error, want *apperr.Error) {
	t.Helper()
	var ae *apperr.Error
	if !errors.As(got, &ae) {
		t.Fatalf("erro = %v, quer %s", got, want.Code)
	}
	if ae.Code != want.Code || ae.Kind != want.Kind || ae.Field != want.Field {
		t.Fatalf("erro = {%s %s field=%q}, quer {%s %s field=%q}",
			ae.Kind, ae.Code, ae.Field, want.Kind, want.Code, want.Field)
	}
}

func wantStored(t *testing.T, repo *fakeRepository, n int) {
	t.Helper()
	if got := len(repo.all()); got != n {
		t.Fatalf("o fake guardou %d projetos, quer %d: %+v", got, n, repo.all())
	}
}

// CT-35
func TestService_Create_InsideInstall(t *testing.T) {
	e := newEnv(t)
	abs := filepath.Join(e.root, "projects", "web", "loja-api")

	v, err := e.svc.Create(context.Background(), Input{Name: "loja-api", Path: abs, Type: TypeWeb})
	if err != nil {
		t.Fatalf("Create() = %v", err)
	}

	want := View{ID: 1, Name: "loja-api", Path: abs, Type: TypeWeb, FolderExists: true,
		CreatedAt: createdAt, UpdatedAt: createdAt}
	if v != want {
		t.Errorf("View = %+v, quer %+v", v, want)
	}
	stored := filepath.Join("projects", "web", "loja-api")
	wantProject := Project{ID: 1, Name: "loja-api", NameKey: "loja-api", Path: stored,
		PathKey: strings.ToLower(stored), Type: TypeWeb, CreatedAt: createdAt, UpdatedAt: createdAt}
	if got := e.repo.all(); !reflect.DeepEqual(got, []Project{wantProject}) {
		t.Errorf("guardado = %+v, quer [%+v]", got, wantProject)
	}
}

// CT-36
func TestService_Create_OutsideInstall(t *testing.T) {
	e := newEnv(t)
	outside := t.TempDir()

	v, err := e.svc.Create(context.Background(), Input{Name: "site", Path: outside, Type: TypeAndroid})
	if err != nil {
		t.Fatalf("Create() = %v", err)
	}

	if v.Path != outside {
		t.Errorf("View.Path = %q, quer %q", v.Path, outside)
	}
	got := e.repo.all()
	if len(got) != 1 || got[0].Path != outside || got[0].PathKey != strings.ToLower(outside) {
		t.Errorf("guardado = %+v, quer Path = %q", got, outside)
	}
}

// Uma pasta cujo nome apenas começa com o da raiz não está dentro dela.
func TestService_Create_SiblingOfInstallRootStaysAbsolute(t *testing.T) {
	e := newEnv(t)
	sibling := e.root + "-extra"
	if err := os.MkdirAll(sibling, 0o755); err != nil {
		t.Fatal(err)
	}

	e.create(t, "extra", sibling, TypeWeb)

	if got := e.repo.all(); len(got) != 1 || got[0].Path != sibling {
		t.Errorf("guardado = %+v, quer Path = %q", got, sibling)
	}
}

// Os espaços das pontas saem do nome e do caminho antes de guardar.
func TestService_Create_TrimsFields(t *testing.T) {
	e := newEnv(t)
	abs := filepath.Join(e.root, "projects", "web", "loja-api")

	v := e.create(t, "  loja-api  ", " "+abs+" ", " web ")

	if v.Name != "loja-api" || v.Path != abs || v.Type != TypeWeb {
		t.Errorf("View = %+v, quer nome, caminho e tipo sem espaços", v)
	}
}

// CT-37
func TestService_Create_Required(t *testing.T) {
	tests := []struct {
		name string
		in   func(dir string) Input
		want *apperr.Error
	}{
		{"nome vazio", func(d string) Input { return Input{"", d, TypeWeb} }, ErrNameRequired},
		{"nome só espaços", func(d string) Input { return Input{"   ", d, TypeWeb} }, ErrNameRequired},
		{"caminho vazio", func(string) Input { return Input{"loja-api", "", TypeWeb} }, ErrPathRequired},
		{"caminho só espaços", func(string) Input { return Input{"loja-api", "  ", TypeWeb} }, ErrPathRequired},
		{"tipo vazio", func(d string) Input { return Input{"loja-api", d, ""} }, ErrTypeRequired},
		{"os três vazios", func(string) Input { return Input{} }, ErrNameRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)
			dir := filepath.Join(e.root, "projects", "web", "loja-api")

			_, err := e.svc.Create(context.Background(), tt.in(dir))

			wantError(t, err, tt.want)
			wantStored(t, e.repo, 0)
		})
	}
}

// CT-38
func TestService_Create_TypeInvalid(t *testing.T) {
	e := newEnv(t)

	_, err := e.svc.Create(context.Background(),
		Input{Name: "loja-api", Path: filepath.Join(e.root, "projects", "web", "loja-api"), Type: "ios"})

	wantError(t, err, ErrTypeInvalid)
	wantStored(t, e.repo, 0)
}

// CT-39
func TestService_Create_RelativePath(t *testing.T) {
	e := newEnv(t)

	_, err := e.svc.Create(context.Background(),
		Input{Name: "loja-api", Path: filepath.Join("projects", "web", "loja-api"), Type: TypeWeb})

	wantError(t, err, ErrPathInvalid)
	wantStored(t, e.repo, 0)
}

// CT-40
func TestService_Create_PathNotFound(t *testing.T) {
	e := newEnv(t)
	file := filepath.Join(e.root, "arquivo.txt")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	tests := map[string]string{
		"pasta inexistente":     filepath.Join(e.root, "nao-existe"),
		"caminho de um arquivo": file,
	}
	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := e.svc.Create(context.Background(), Input{Name: "loja-api", Path: path, Type: TypeWeb})

			wantError(t, err, ErrPathNotFound)
			wantStored(t, e.repo, 0)
		})
	}
}

// Os obrigatórios vêm antes do formato, e o formato antes da existência (RN-12).
func TestService_Create_ValidationOrder(t *testing.T) {
	e := newEnv(t)
	missing := filepath.Join(e.root, "nao-existe")

	tests := []struct {
		name string
		in   Input
		want *apperr.Error
	}{
		{"tipo inválido antes do caminho relativo", Input{"loja-api", "relativo", "ios"}, ErrTypeInvalid},
		{"tipo vazio antes do caminho relativo", Input{"loja-api", "relativo", ""}, ErrTypeRequired},
		{"caminho relativo antes da existência", Input{"loja-api", "relativo", TypeWeb}, ErrPathInvalid},
		{"tipo inválido antes da existência", Input{"loja-api", missing, "ios"}, ErrTypeInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.svc.Create(context.Background(), tt.in)

			wantError(t, err, tt.want)
		})
	}
}

// CT-41
func TestService_Create_NameTaken(t *testing.T) {
	e := newEnv(t)
	e.create(t, "loja-api", filepath.Join(e.root, "projects", "web", "loja-api"), TypeWeb)
	other := e.mkdir(t, "projects", "web", "outra")

	_, err := e.svc.Create(context.Background(), Input{Name: "Loja-API", Path: other, Type: TypeWeb})

	wantError(t, err, ErrNameTaken)
	wantStored(t, e.repo, 1)
}

// CT-42
func TestService_Create_PathTaken(t *testing.T) {
	sep := string(os.PathSeparator)
	tests := []struct {
		name         string
		path         func(root string) string
		needsWindows bool
	}{
		{"maiúsculas e barra final", func(root string) string {
			return strings.ToUpper(root) + sep + "Projects" + sep + "Web" + sep + "Loja-API" + sep
		}, true},
		{"segmento ponto", func(root string) string {
			return root + sep + "projects" + sep + "web" + sep + "." + sep + "loja-api"
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.needsWindows && runtime.GOOS != "windows" {
				t.Skip("o sistema de arquivos só ignora maiúsculas no Windows")
			}
			e := newEnv(t)
			e.create(t, "loja-api", filepath.Join(e.root, "projects", "web", "loja-api"), TypeWeb)

			_, err := e.svc.Create(context.Background(),
				Input{Name: "outro-nome", Path: tt.path(e.root), Type: TypeWeb})

			wantError(t, err, ErrPathTaken)
			wantStored(t, e.repo, 1)
		})
	}
}

// CT-43
func TestService_List(t *testing.T) {
	e := newEnv(t)
	for _, name := range []string{"zeta", "Alfa", "beta"} {
		e.create(t, name, e.mkdir(t, "projects", "web", name), TypeWindows)
	}
	if err := os.RemoveAll(filepath.Join(e.root, "projects", "web", "beta")); err != nil {
		t.Fatal(err)
	}

	got, err := e.svc.List(context.Background())
	if err != nil {
		t.Fatalf("List() = %v", err)
	}

	want := []struct {
		name   string
		exists bool
	}{{"Alfa", true}, {"beta", false}, {"zeta", true}}
	if len(got) != len(want) {
		t.Fatalf("List() devolveu %d itens, quer %d: %+v", len(got), len(want), got)
	}
	for i, w := range want {
		v := got[i]
		if v.Name != w.name || v.FolderExists != w.exists || v.Type != TypeWindows ||
			v.Path != filepath.Join(e.root, "projects", "web", w.name) {
			t.Errorf("item %d = %+v, quer nome %q, FolderExists %v, tipo windows, caminho absoluto",
				i, v, w.name, w.exists)
		}
	}
}

func TestService_List_Empty(t *testing.T) {
	e := newEnv(t)

	got, err := e.svc.List(context.Background())

	if err != nil {
		t.Fatalf("List() = %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("List() = %#v, quer lista vazia não nula", got)
	}
}

// CT-44
func TestService_Update(t *testing.T) {
	e := newEnv(t)
	v := e.create(t, "loja-api", filepath.Join(e.root, "projects", "web", "loja-api"), TypeWeb)
	outside := t.TempDir()
	e.now = updatedAt

	got, err := e.svc.Update(context.Background(), v.ID, Input{Name: "novo", Path: outside, Type: TypeAndroid})
	if err != nil {
		t.Fatalf("Update() = %v", err)
	}

	want := View{ID: v.ID, Name: "novo", Path: outside, Type: TypeAndroid, FolderExists: true,
		CreatedAt: createdAt, UpdatedAt: updatedAt}
	if got != want {
		t.Errorf("View = %+v, quer %+v", got, want)
	}
	saved := e.repo.items[v.ID]
	if saved.Name != "novo" || saved.NameKey != "novo" || saved.Path != outside ||
		saved.PathKey != strings.ToLower(outside) || saved.Type != TypeAndroid ||
		saved.CreatedAt != createdAt || saved.UpdatedAt != updatedAt {
		t.Errorf("guardado = %+v", saved)
	}
	wantStored(t, e.repo, 1)
}

// CT-45
func TestService_Update_DataOfAnother(t *testing.T) {
	e := newEnv(t)
	dirA := e.mkdir(t, "projects", "web", "a")
	dirB := e.mkdir(t, "projects", "web", "b")
	a := e.create(t, "a-api", dirA, TypeWeb)
	e.create(t, "b-api", dirB, TypeWeb)
	before := e.repo.items[a.ID]
	e.now = updatedAt

	t.Run("nome de outro", func(t *testing.T) {
		_, err := e.svc.Update(context.Background(), a.ID, Input{Name: "b-api", Path: dirA, Type: TypeWeb})

		wantError(t, err, ErrNameTaken)
		if e.repo.items[a.ID] != before {
			t.Errorf("A mudou: %+v, quer %+v", e.repo.items[a.ID], before)
		}
	})
	t.Run("pasta de outro", func(t *testing.T) {
		_, err := e.svc.Update(context.Background(), a.ID, Input{Name: "a-api", Path: dirB, Type: TypeWeb})

		wantError(t, err, ErrPathTaken)
		if e.repo.items[a.ID] != before {
			t.Errorf("A mudou: %+v, quer %+v", e.repo.items[a.ID], before)
		}
	})
	t.Run("próprio nome em maiúsculas", func(t *testing.T) {
		v, err := e.svc.Update(context.Background(), a.ID, Input{Name: "A-API", Path: dirA, Type: TypeWeb})

		if err != nil {
			t.Fatalf("Update() = %v", err)
		}
		if v.Name != "A-API" || e.repo.items[a.ID].Name != "A-API" || e.repo.items[a.ID].NameKey != "a-api" {
			t.Errorf("View = %+v, guardado = %+v, quer nome A-API e chave a-api", v, e.repo.items[a.ID])
		}
	})
}

// CT-46
func TestService_Update_PathNotFound(t *testing.T) {
	e := newEnv(t)
	a := e.create(t, "loja-api", filepath.Join(e.root, "projects", "web", "loja-api"), TypeWeb)
	before := e.repo.items[a.ID]

	_, err := e.svc.Update(context.Background(), a.ID,
		Input{Name: "loja-api", Path: filepath.Join(e.root, "nao-existe"), Type: TypeWeb})

	wantError(t, err, ErrPathNotFound)
	if e.repo.items[a.ID] != before {
		t.Errorf("A mudou: %+v, quer %+v", e.repo.items[a.ID], before)
	}
}

// CT-47
func TestService_NotFound(t *testing.T) {
	e := newEnv(t)
	dir := filepath.Join(e.root, "projects", "web", "loja-api")

	t.Run("Update", func(t *testing.T) {
		_, err := e.svc.Update(context.Background(), 999, Input{Name: "x", Path: dir, Type: TypeWeb})

		wantError(t, err, ErrNotFound)
	})
	t.Run("Delete", func(t *testing.T) {
		wantError(t, e.svc.Delete(context.Background(), 999), ErrNotFound)
	})
	if e.logs.Len() != 0 {
		t.Errorf("nada deveria ser logado, mas o log tem:\n%s", e.logs)
	}
}

// CT-48
func TestService_Delete(t *testing.T) {
	e := newEnv(t)
	dir := filepath.Join(e.root, "projects", "web", "loja-api")
	v := e.create(t, "loja-api", dir, TypeWeb)

	if err := e.svc.Delete(context.Background(), v.ID); err != nil {
		t.Fatalf("Delete() = %v", err)
	}

	wantStored(t, e.repo, 0)
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("a pasta do projeto deveria continuar no disco: %v", err)
	}
}

// entry é o que o retrato de uma pasta guarda de cada item.
type entry struct {
	dir     bool
	size    int64
	modTime int64
}

// snapshot percorre dir recursivamente e anota nome, tamanho e data de cada item.
func snapshot(t *testing.T, dir string) map[string]entry {
	t.Helper()
	snap := map[string]entry{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// os.Stat e não d.Info(): no Windows a data de uma pasta que vem da listagem do pai fica
		// defasada, e o retrato mudaria sem ninguém ter mexido no disco.
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		snap[path] = entry{dir: d.IsDir(), size: info.Size(), modTime: info.ModTime().UnixNano()}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(%s): %v", dir, err)
	}
	return snap
}

// CT-49
func TestService_NeverChangesDisk(t *testing.T) {
	e := newEnv(t)
	inside := filepath.Join(e.root, "projects", "web", "loja-api")
	if err := os.WriteFile(filepath.Join(inside, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "index.html"), []byte("<html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	rootBefore, outsideBefore := snapshot(t, e.root), snapshot(t, outside)

	steps := []struct {
		name string
		run  func() error
	}{
		{"Create", func() error {
			_, err := e.svc.Create(context.Background(), Input{Name: "loja-api", Path: inside, Type: TypeWeb})
			return err
		}},
		{"List", func() error {
			_, err := e.svc.List(context.Background())
			return err
		}},
		{"Update", func() error {
			_, err := e.svc.Update(context.Background(), 1, Input{Name: "site", Path: outside, Type: TypeWindows})
			return err
		}},
		{"Delete", func() error { return e.svc.Delete(context.Background(), 1) }},
	}
	for _, step := range steps {
		if err := step.run(); err != nil {
			t.Fatalf("%s() = %v", step.name, err)
		}
		if got := snapshot(t, e.root); !reflect.DeepEqual(got, rootBefore) {
			t.Errorf("depois de %s a raiz mudou:\nantes:  %v\ndepois: %v", step.name, rootBefore, got)
		}
		if got := snapshot(t, outside); !reflect.DeepEqual(got, outsideBefore) {
			t.Errorf("depois de %s a pasta externa mudou:\nantes:  %v\ndepois: %v", step.name, outsideBefore, got)
		}
	}
}

// CT-50
func TestService_LogsEvents(t *testing.T) {
	e := newEnv(t)
	dir := filepath.Join(e.root, "projects", "web", "loja-api")
	v := e.create(t, "loja-api", dir, TypeWeb)
	if _, err := e.svc.Update(context.Background(), v.ID, Input{Name: "loja", Path: dir, Type: TypeWeb}); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Delete(context.Background(), v.ID); err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(e.logs.String()), "\n")
	want := []string{
		"level=INFO msg=project.created id=1 name=loja-api path=",
		"level=INFO msg=project.updated id=1 name=loja path=",
		"level=INFO msg=project.deleted id=1 name=loja path=",
	}
	if len(lines) != len(want) {
		t.Fatalf("log com %d linhas, quer %d:\n%s", len(lines), len(want), e.logs)
	}
	for i, w := range want {
		if !strings.Contains(lines[i], w) {
			t.Errorf("linha %d = %q, quer conter %q", i, lines[i], w)
		}
	}
}

// Nenhum evento é logado quando a operação falha.
func TestService_Create_Failure_LogsNothing(t *testing.T) {
	e := newEnv(t)

	_, err := e.svc.Create(context.Background(), Input{})

	wantError(t, err, ErrNameRequired)
	if e.logs.Len() != 0 {
		t.Errorf("log = %q, quer vazio", e.logs)
	}
}
