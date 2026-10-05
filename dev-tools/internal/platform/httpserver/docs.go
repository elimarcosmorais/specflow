package httpserver

import (
	"embed"
	"net/http"
)

// docsFS embute a página /docs: o Redoc e o contrato, para a documentação funcionar offline (RN-06).
//
//go:embed docs
var docsFS embed.FS

// docsRoutes são as rotas públicas da documentação, fora do contrato /api.
func docsRoutes() []Route {
	return []Route{
		docsRoute("/docs", "docs/index.html", "text/html; charset=utf-8"),
		docsRoute("/docs/openapi.yaml", "docs/openapi.yaml", "application/yaml; charset=utf-8"),
		docsRoute("/docs/redoc.standalone.js", "docs/redoc.standalone.js", "text/javascript; charset=utf-8"),
	}
}

func docsRoute(pattern, file, contentType string) Route {
	return Route{
		Method:  http.MethodGet,
		Pattern: pattern,
		Public:  true,
		Handler: func(w http.ResponseWriter, r *http.Request) error {
			body, err := docsFS.ReadFile(file)
			if err != nil {
				return err // não ocorre: o embed garante o arquivo; se ocorresse, vira 500
			}
			w.Header().Set("Content-Type", contentType)
			_, err = w.Write(body)
			return err
		},
	}
}
