// Package frontend embute o front compilado (Vite grava em dist/app) no executável.
package frontend

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

// Assets devolve os arquivos do front, com a raiz em dist/app.
func Assets() fs.FS {
	assets, err := fs.Sub(dist, "dist/app")
	if err != nil {
		// Só falha com caminho inválido; o caminho é constante.
		panic(err)
	}
	return assets
}
