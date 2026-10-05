package projects

import (
	"go/parser"
	"go/token"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// forbiddenImports devolve os imports de um arquivo que o núcleo não pode ter (ADR 002, spec §2.2).
// src nil lê o arquivo de path.
func forbiddenImports(t *testing.T, path string, src any) []string {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), path, src, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("ParseFile(%s): %v", path, err)
	}
	forbidden := []string{"database/sql", "net/http", "github.com/wailsapp/", "modernc.org/", "dev-tools/internal/platform/"}

	var found []string
	for _, spec := range file.Imports {
		imp, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			t.Fatalf("import inválido %s em %s", spec.Path.Value, path)
		}
		for _, f := range forbidden {
			// "database/sql" e "net/http" valem também para os subpacotes.
			if imp == strings.TrimSuffix(f, "/") || strings.HasPrefix(imp, strings.TrimSuffix(f, "/")+"/") {
				found = append(found, imp)
				break
			}
		}
	}
	return found
}

// CT-51
func TestCore_DoesNotImportInfrastructure(t *testing.T) {
	for _, file := range []string{"project.go", "errors.go", "repository.go", "service.go"} {
		t.Run(file, func(t *testing.T) {
			if found := forbiddenImports(t, file, nil); len(found) > 0 {
				t.Errorf("%s importa infraestrutura: %v", file, found)
			}
		})
	}
}

// O verificador do CT-51 enxerga cada import proibido, e só eles.
func TestForbiddenImports_Detection(t *testing.T) {
	src := `package projects

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"github.com/wailsapp/wails/v2"
	"modernc.org/sqlite"
	"dev-tools/internal/platform/db"
	"dev-tools/internal/apperr"
)
`
	want := []string{"database/sql", "net/http", "net/http/httptest", "github.com/wailsapp/wails/v2",
		"modernc.org/sqlite", "dev-tools/internal/platform/db"}

	if got := forbiddenImports(t, "synthetic.go", src); !slices.Equal(got, want) {
		t.Errorf("forbiddenImports = %v, quer %v", got, want)
	}
}
