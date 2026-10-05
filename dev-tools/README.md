# Dev Tools

App de desktop para Windows (Wails v2: motor em Go + front em React) que cadastra e organiza
projetos de desenvolvimento. É portátil: banco, configurações e logs ficam dentro da pasta do
app, nunca no registro, em `%APPDATA%` nem no `PATH` (ADR 003).

Todos os comandos abaixo partem da pasta `dev-tools/` deste repositório, salvo indicação.

## Pré-requisitos

| Ferramenta | Versão | Para quê |
|------------|--------|----------|
| Go         | 1.26+  | motor, testes e build |
| Wails CLI  | v2.16  | `wails dev`, `wails build` (`go install github.com/wailsapp/wails/v2/cmd/wails@v2.16.0`) |
| Node.js    | LTS atual (testado com 24) | front (React + Vite) |
| NSIS       | 3.x    | só para gerar o instalador (`makensis` no `PATH`) |

O WebView2 já vem no Windows 10 e 11 atualizados.

## Rodar em desenvolvimento

Em `dev-tools/cmd/dev-tools`, aponte a pasta base para uma pasta de desenvolvimento e use outra
porta da API, para não se confundir com o app instalado (porta `41790`):

```powershell
cd cmd\dev-tools
$env:DEV_TOOLS_HOME = "<repo>\dev-tools\.dev-home\programs"
$env:DEV_TOOLS_API_PORT = "41791"
wails dev
```

Troque `<repo>` pelo caminho deste repositório. A pasta `.dev-home\` é criada pelo app e está no
`.gitignore`.

Variáveis de ambiente do app:

| Variável              | Padrão              | Uso |
|-----------------------|---------------------|-----|
| `DEV_TOOLS_HOME`      | pasta do executável | pasta base (banco, `config.json`, logs) |
| `DEV_TOOLS_API_PORT`  | `41790`             | porta da API local (`127.0.0.1`) |
| `DEV_TOOLS_LOG_LEVEL` | `info`              | `debug`, `info`, `warn` ou `error` |

## Testes

```powershell
go test ./...                    # unitários
go test -tags integration ./...  # todos (unitários + integração)
```

## Tipos do front

Os tipos da API são gerados a partir do `openapi.yaml` do motor. Em `dev-tools/frontend`:

```powershell
npm run gen:api     # regenera src/shared/api/schema.d.ts
npm run typecheck   # tsc --noEmit
```

## Build

Em `dev-tools/cmd/dev-tools`:

```powershell
wails build
```

Saída: `dev-tools/build/bin/dev-tools.exe`.

O `wails build` e o `wails dev` regeneram `frontend/package.json.md5`, que está no `.gitignore`.

## Instalador

Precisa do `makensis` no `PATH`. Em `dev-tools/cmd/dev-tools`:

```powershell
wails build -nsis
```

Saída: `dev-tools/build/bin/dev-tools-amd64-installer.exe`.

O `wails build -nsis` baixa `build/windows/installer/tmp/MicrosoftEdgeWebview2Setup.exe`; a pasta
`tmp/` está no `.gitignore`.

O instalador (`build/windows/installer/project.nsi`):

- roda como usuário comum, sem pedir admin (`RequestExecutionLevel user`);
- sugere `C:\dev-tools`, e o usuário pode escolher outra pasta;
- instala o app em `<pasta>\programs\` e cria `programs\{data,logs}` e
  `projects\{web,android,windows}`;
- cria atalhos na Área de Trabalho e no Menu Iniciar (inclusive "Desinstalar Dev Tools");
- não escreve no registro nem associa tipos de arquivo.

O desinstalador (`<pasta>\programs\uninstall.exe`) pergunta se deve remover também os dados e os
logs. **Sim** apaga `programs\` inteira; **Não** mantém `programs\data\` e `programs\logs\`.
Nos dois casos, `projects\` nunca é tocada. Feche o app antes de desinstalar.

O `wails build` regenera `build/windows/installer/wails_tools.nsh` a cada execução. O
`project.nsi` não depende das partes dele que escrevem no registro, então o instalador continua
sem escrever nada lá mesmo com o arquivo regenerado.
