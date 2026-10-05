package desktop

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
)

// UniqueID devolve o ID da instância única (RN-02) para a pasta base: o mesmo para a mesma
// pasta escrita de formas diferentes (maiúsculas, barra final) e diferente para outra pasta.
// Assim o app de desenvolvimento (DEV_TOOLS_HOME) não se confunde com o instalado.
func UniqueID(baseDir string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(filepath.Clean(baseDir))))
	return "dev-tools-" + hex.EncodeToString(sum[:])[:12]
}
