# Débitos — Arquitetura base do Dev Tools · v1

Problemas médios ou baixos encontrados nesta versão: não bloquearam a entrega, mas devem ser
corrigidos depois. São pagos pela fase Debt (`/agent-spec-debt`).

---

## DT-01 — Arquivo gerado pelo Wails fora do .gitignore

- **Status:** pago
- **Severidade:** baixa
- **Origem:** revisão da T1
- **Arquivo:** `dev-tools/.gitignore:4`
- **Problema:** todo `wails build` ou `wails dev` recria `frontend/package.json.md5`. O arquivo não está no `.gitignore`, então volta como não rastreado e pode ser commitado por engano.
- **Correção sugerida:** acrescentar `frontend/package.json.md5` ao `dev-tools/.gitignore`.
- **Pago em:** v2 · T1

## DT-02 — Log config.invalid sem o atributo path exigido pela spec §8

- **Status:** pago
- **Severidade:** média
- **Origem:** revisão da T6
- **Arquivo:** `dev-tools/internal/features/preferences/service.go:43`
- **Problema:** a spec §8 define `config.invalid` (WARN) com os atributos `path` e `err`, mas o serviço loga só `err`. Quando um campo é inválido, o erro não traz o caminho do arquivo, e o operador não vê qual `config.json` foi regravado.
- **Correção sugerida:** expor o caminho ao serviço sem importar infraestrutura (parâmetro opcional `path string` em `NewService`, ou interface opcional `Path() string` no `Repository`); logar `"path", path` junto de `"err"` nos dois ramos de `config.invalid` e cobrir isso no CT-30.
- **Pago em:** v2 · T2

## DT-03 — wails_tools.nsh contém macros de escrita de registro, contra a T13 e a spec §3.3.4

- **Status:** aberto
- **Severidade:** média
- **Origem:** revisão da T13
- **Arquivo:** `dev-tools/build/windows/installer/wails_tools.nsh:107`
- **Problema:** a seção 3 da T13 e a spec §3.3.4 pedem o `.nsh` "sem escrita de registro", mas o `wails build` regenera o arquivo a cada execução ("DO NOT EDIT"), com `wails.writeUninstaller`, `wails.deleteUninstaller`, `APP_ASSOCIATE` e `CUSTOM_PROTOCOL_ASSOCIATE`. O `project.nsi` não chama nenhuma delas, então o instalador cumpre a ADR 003; o desvio é entre o artefato e o texto da spec.
- **Correção sugerida:** corrigir a spec §3.3.4 e a seção 3 da T13: o "sem registro" vale para o `project.nsi`, que não invoca as macros de registro; o `wails_tools.nsh` fica como o Wails gera. Opcional: acrescentar ao README uma checagem como `grep -n "writeUninstaller\|associate" project.nsi`.
- **Pago em:** —

## DT-04 — Pasta tmp do instalador fora do .gitignore

- **Status:** pago
- **Severidade:** baixa
- **Origem:** revisão da T13
- **Arquivo:** `dev-tools/.gitignore:0`
- **Problema:** cada `wails build -nsis` recria `build/windows/installer/tmp/MicrosoftEdgeWebview2Setup.exe` (cerca de 1,7 MB, binário baixado). Nada o ignora, então ele pode entrar num commit por engano.
- **Correção sugerida:** adicionar `build/windows/installer/tmp/` ao `dev-tools/.gitignore` e citar isso no README, que já avisa sobre `package.json.md5`.
- **Pago em:** v2 · T1
