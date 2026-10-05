package system

import "os/exec"

// ExplorerOpener abre pastas no Windows Explorer.
type ExplorerOpener struct{}

// Open executa explorer.exe com path. Usa Start, e não Run, porque o explorer.exe devolve código
// de saída 1 mesmo quando dá certo (spec §3.3.5). O processo não é aguardado.
func (ExplorerOpener) Open(path string) error {
	cmd := exec.Command("explorer.exe", path)
	if err := cmd.Start(); err != nil {
		return err
	}
	// Libera os recursos do processo filho sem esperar o fim dele.
	go func() { _ = cmd.Wait() }()
	return nil
}
