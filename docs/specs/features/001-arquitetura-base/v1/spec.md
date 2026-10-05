---
status: aprovado # rascunho | aprovado
---

# Tech Spec — Arquitetura base do Dev Tools

- **Feature:** `001-arquitetura-base` · **Versão:** v1 · **Data:** 2026-10-05
- **Origem:** [prd.md](prd.md) · [tech-alignment.md](tech-alignment.md)
- **ADRs aplicáveis:**
  - [`docs/adr/001-stack-base-do-dev-tools.md`](../../../../adr/001-stack-base-do-dev-tools.md)
  - [`docs/adr/002-arquitetura-em-camadas-e-padrao-de-feature.md`](../../../../adr/002-arquitetura-em-camadas-e-padrao-de-feature.md)
  - [`docs/adr/003-aplicacao-portatil.md`](../../../../adr/003-aplicacao-portatil.md)
  - [`docs/adr/004-acesso-a-api-local-token-por-execucao-e-bootstrap-do-front.md`](../../../../adr/004-acesso-a-api-local-token-por-execucao-e-bootstrap-do-front.md)
  - [`docs/adr/005-padrao-de-erros-e-idiomas.md`](../../../../adr/005-padrao-de-erros-e-idiomas.md)
  - [`docs/adr/006-padrao-de-logging.md`](../../../../adr/006-padrao-de-logging.md)
  - [`docs/adr/007-migracoes-do-banco-de-dados.md`](../../../../adr/007-migracoes-do-banco-de-dados.md)
  - [`docs/adr/008-organizacao-do-front.md`](../../../../adr/008-organizacao-do-front.md)

---

## 1. Resumo técnico

Feature de fundação: ainda não existe código. Esta spec é o primeiro retrato concreto do
Dev Tools, derivado das ADRs 001 a 008 e do Tech Alignment. O código fica em `dev-tools/`, na
raiz deste repositório. O motor em Go (Wails v2) monta, em `cmd/dev-tools/main.go`, a
infraestrutura comum (`internal/platform/`) e três features: **projects** (a feature-modelo, um
CRUD com SQLite), **preferences** (tema e idioma, num arquivo JSON) e **system** (saúde e atalho
para os logs). A API HTTP escuta em `127.0.0.1:41790` e exige token nas rotas que não são
públicas. O front React recebe token, endereço e estado por uma única chamada de bootstrap e,
depois disso, fala só com a API. O instalador é o NSIS da integração do Wails, ajustado para
rodar como usuário comum.

**Decisões do Tech Alignment aplicadas:**

| Decisão                                    | Onde está nesta spec         |
| ------------------------------------------ | ---------------------------- |
| D-01 — Token por execução e bootstrap      | §3.1, §3.3.1, §5, §7.2, §8   |
| D-02 — Erros e idiomas                     | §3.2, §8 (Idiomas)           |
| D-03 — Preferências em arquivo             | §3.1, §4.3, §6.2             |
| D-04 — Logs com `slog` + `lumberjack`      | §4.3, §8 (Logs)              |
| D-05 — Migrações embutidas                 | §4.1, §4.2                   |
| D-06 — Banco indisponível (modo degradado) | §5, §6.2, §10 (CT-59, CT-65) |
| D-07 — Caminhos guardados                  | §2.3, §4.1, §5               |
| D-08 — Porta fixa 41790                    | §3.1, §8 (Configuração)      |
| D-09 — Instalador NSIS                     | §3.3.4                       |
| D-10 — Front React puro                    | §2.1, §7, §11                |
| D-11 — Código numa subpasta do repositório | §11 (`dev-tools/`)           |

## 2. Arquitetura

### 2.1 Componentes

| Componente                               | Camada                        | Responsabilidade                                                                                                                           | Situação |
| ---------------------------------------- | ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ | -------- |
| `cmd/dev-tools` (`main.go`, `routes.go`) | montagem                      | Sequência de inicialização (§6.1), opções do Wails, lista de rotas da API                                                                  | novo     |
| `internal/apperr`                        | núcleo comum                  | Tipo de erro de domínio (`Kind`, `Code`, `Field`, causa) e erros comuns; sem I/O nem dependências                                          | novo     |
| `internal/platform/config`               | infraestrutura                | Resolve a pasta base e todos os caminhos (`Paths`), garante a estrutura de pastas e lê as variáveis de ambiente                            | novo     |
| `internal/platform/db`                   | infraestrutura                | Abre o SQLite, aplica as migrações embutidas, `Pinger`                                                                                     | novo     |
| `internal/platform/log`                  | infraestrutura                | Cria o `*slog.Logger` (texto) com `lumberjack` e a rotação ao iniciar                                                                      | novo     |
| `internal/platform/i18n`                 | infraestrutura                | Catálogo do motor (pt-br, en), tradução por código e detecção do idioma do Windows                                                         | novo     |
| `internal/platform/httpserver`           | infraestrutura                | Registro de rotas, middlewares (CORS, token, recuperação, log de acesso), resposta de erro traduzida, `Listen` em loopback, página `/docs` | novo     |
| `internal/platform/desktop`              | infraestrutura                | `Bridge` (binding único `Bootstrap`), `Window` (foco e tema da janela), mensagens do WebView2, diálogo nativo, ID de instância             | novo     |
| `internal/features/projects`             | núcleo + infra + apresentação | Entidade, serviço, interface de repositório, repositório SQLite e handler HTTP de Projetos                                                 | novo     |
| `internal/features/preferences`          | núcleo + infra + apresentação | Entidade, serviço, interface de repositório, repositório em arquivo JSON e handler HTTP de Preferências                                    | novo     |
| `internal/features/system`               | núcleo + infra + apresentação | Serviço de saúde e de abrir a pasta de logs, abridor via Explorer, handler HTTP                                                            | novo     |
| `frontend/` (React)                      | front                         | Shell, telas Projetos, Preferências e Status, cliente da API, catálogo das telas, tema                                                     | novo     |
| `build/windows/installer`                | instalador                    | Script NSIS ajustado (§3.3.4)                                                                                                              | novo     |

### 2.2 Dependências entre camadas

```
cmd/dev-tools/main.go ──monta──► platform/* + features/* + frontend (assets embutidos)

features/<f>/handler.go ──► features/<f>/service.go ──► features/<f>/repository.go (interface)
        │                          │                              ▲
        │                          └──► internal/apperr           │ implementa
        ▼                                                         │
platform/httpserver          features/<f>/sqlite_repository.go ───┘ (ou file_repository.go)
                                        └──► database/sql (o *sql.DB vem de platform/db)

front ──(1x binding Wails)──► platform/desktop.Bridge.Bootstrap
front ──(HTTP + token)──────► platform/httpserver ──► handlers das features
```

O que cada camada **não pode** importar:

- **Núcleo** (`project.go`, `errors.go`, `repository.go`, `service.go` de cada feature): não
  importa `database/sql`, `net/http`, `github.com/wailsapp/*`, `modernc.org/sqlite` nem
  nenhum pacote de `internal/platform/`. Pode importar a biblioteca padrão (`os`,
  `path/filepath`, `log/slog`, `time`, `strings`, `context`) e `internal/apperr`. O CT-51 confere
  isso.
- **Handler** de uma feature: importa o próprio serviço, `internal/apperr` e
  `platform/httpserver`; nunca o repositório nem `database/sql`.
- **Features** não importam umas às outras. Quando uma precisa de algo da outra, `main.go` liga
  as duas com uma função (ex.: o idioma atual chega ao `httpserver` como `func() string`).
- **`internal/platform/*`** não importa `internal/features/*`.
- **`internal/apperr`** só importa a biblioteca padrão. Não conhece HTTP: quem converte `Kind`
  em status HTTP é o `httpserver`.

### 2.3 Layout em disco

A instalação segue o esquema do prompt original. A **pasta base** (ADR 003) é a pasta do
executável, ou `DEV_TOOLS_HOME`. A **raiz da instalação** (D-07) é a pasta acima dela. Na v1 só
existem as pastas que o app usa:

```
C:\dev-tools\                    raiz da instalação (Paths.Root)
├── projects\                    Paths.Projects — pasta de projetos (nunca apagada pelo app nem pelo desinstalador)
│   ├── web\
│   ├── android\
│   └── windows\
└── programs\                    pasta base (Paths.Base)
    ├── dev-tools.exe
    ├── uninstall.exe            (só quando instalado)
    ├── data\                    Paths.Data
    │   ├── dev-tools.db         Paths.DBFile
    │   └── config.json          Paths.ConfigFile
    ├── logs\                    Paths.Logs
    │   └── dev-tools.log        Paths.LogFile (+ backups do lumberjack)
    └── webview\                 Paths.WebView — dados do WebView2 (cache, cookies)
```

## 3. Contratos

### 3.1 API

Base: `http://127.0.0.1:41790`. Corpo e resposta em JSON (`Content-Type: application/json;
charset=utf-8`). **Proteção "token"** = cabeçalho `Authorization: Bearer <token>`, em que
`<token>` é o valor entregue pelo bootstrap (§3.3.1). Toda rota é protegida por padrão: só as
marcadas como `Public` no registro (§3.3.2) dispensam o token.

| Operação            | Método e rota                                                            | Entrada              | Saída (sucesso)                                                            | Erros                                                                                                                                                                                                                                                                                       | Proteção                          |
| ------------------- | ------------------------------------------------------------------------ | -------------------- | -------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------- |
| Saúde               | `GET /api/health`                                                        | —                    | `200` + `{version, database}`                                              | — (responde `200` mesmo com o banco fora do ar)                                                                                                                                                                                                                                             | pública                           |
| Listar projetos     | `GET /api/projects`                                                      | —                    | `200` + `{items: Project[]}`, ordenado por nome sem diferenciar maiúsculas | `401 auth.token_invalid` · `503 database.unavailable` · `500 internal.unexpected`                                                                                                                                                                                                           | token                             |
| Criar projeto       | `POST /api/projects`                                                     | `{name, path, type}` | `201` + `Project`                                                          | `400 request.invalid_body` · `422 project.name_required` · `422 project.path_required` · `422 project.type_required` · `422 project.type_invalid` · `422 project.path_invalid` · `422 project.path_not_found` · `409 project.name_taken` · `409 project.path_taken` · `401` · `503` · `500` | token                             |
| Editar projeto      | `PUT /api/projects/{id}`                                                 | `{name, path, type}` | `200` + `Project`                                                          | os mesmos do criar, mais `404 project.not_found`                                                                                                                                                                                                                                            | token                             |
| Remover projeto     | `DELETE /api/projects/{id}`                                              | —                    | `204` sem corpo                                                            | `404 project.not_found` · `401` · `503` · `500`                                                                                                                                                                                                                                             | token                             |
| Salvar preferências | `PUT /api/preferences`                                                   | `{theme, language}`  | `200` + `{theme, language}`                                                | `400 request.invalid_body` · `422 preferences.theme_invalid` · `422 preferences.language_invalid` · `500 preferences.save_failed` · `401`                                                                                                                                                   | token                             |
| Abrir pasta de logs | `POST /api/system/open-logs-folder`                                      | —                    | `204` sem corpo                                                            | `500 system.open_logs_failed` · `401`                                                                                                                                                                                                                                                       | token                             |
| Documentação        | `GET /docs` · `GET /docs/openapi.yaml` · `GET /docs/redoc.standalone.js` | —                    | `200` (HTML, YAML, JS)                                                     | —                                                                                                                                                                                                                                                                                           | pública (fora do contrato `/api`) |

**Tipos:**

- `Project` = `{"id": 1, "name": "loja-api", "path": "C:\\dev-tools\\projects\\web\\loja-api",
"type": "web", "folderExists": true, "createdAt": "2026-10-05T13:00:00Z", "updatedAt":
"2026-10-05T13:00:00Z"}`. `path` é sempre o **caminho absoluto** resolvido (§5, D-07);
  `folderExists` é calculado a cada leitura (RN-18); datas em RFC 3339, UTC.
- Saúde = `{"version": "0.1.0", "database": "available"}`. `database` é `"available"` ou
  `"unavailable"`.
- Preferências = `{"theme": "dark", "language": "en"}`. `theme` é `"light"`, `"dark"` ou
  `"system"`; `language` é `"pt-br"` ou `"en"`.

**Regras da entrada:**

- `PUT /api/projects/{id}` **substitui** o projeto: `name`, `path` e `type` são obrigatórios,
  como no `POST`. Não é atualização parcial. O formulário do front sempre envia os três.
- `PUT /api/preferences` também substitui: `theme` e `language` são obrigatórios.
- `{id}` que não é um inteiro positivo responde `404 project.not_found`.
- Corpo com JSON malformado, com campo desconhecido ou maior que 1 MiB responde
  `400 request.invalid_body`.
- Campo ausente no JSON vale como string vazia e cai na validação de obrigatório.
- Exemplo de criação: `{"name": "loja-api", "path": "C:\\dev-tools\\projects\\web\\loja-api",
"type": "web"}`.

**Formato de todo erro** (status ≥ 400):

```json
{ "error": { "code": "project.name_taken", "message": "Já existe um projeto com esse nome.", "field": "name" } }
```

`field` só aparece quando o erro é de um campo da entrada (`name`, `path`, `type`, `theme`,
`language`). `message` sai no idioma salvo nas preferências **no momento da requisição**.

**CORS** (todas as respostas): `Access-Control-Allow-Origin: *`. Uma requisição `OPTIONS`
(preflight) responde `204`, sem exigir token, com `Access-Control-Allow-Methods: GET, POST, PUT,
DELETE, OPTIONS`, `Access-Control-Allow-Headers: Authorization, Content-Type`,
`Access-Control-Allow-Private-Network: true` e `Access-Control-Max-Age: 600`. O CORS aberto é
seguro aqui porque o token viaja num cabeçalho, nunca em cookie: uma página de outra origem não
tem como enviá-lo, e quem barra a requisição é o token (D-01).

Rotas fora do registro (ex.: `GET /api/inexistente`) ficam com a resposta padrão do `ServeMux`
(`404`/`405` em texto). O front não chama essas rotas.

### 3.2 Erros

Toda mensagem existe em pt-br e em en no catálogo do motor (`internal/platform/i18n/catalog.go`).
A tradução para en é feita pelo executor, com o mesmo sentido do texto abaixo.

| Código                         | Quando acontece                                                                                         | Status HTTP | Mensagem (pt-br)                                                                                 |
| ------------------------------ | ------------------------------------------------------------------------------------------------------- | ----------- | ------------------------------------------------------------------------------------------------ |
| `auth.token_invalid`           | `Authorization` ausente ou com token diferente do da execução                                           | `401`       | Operação recusada: a requisição não veio do Dev Tools.                                           |
| `request.invalid_body`         | JSON malformado, campo desconhecido ou corpo > 1 MiB                                                    | `400`       | Os dados enviados estão em um formato inválido.                                                  |
| `project.name_required`        | `name` vazio após remover espaços                                                                       | `422`       | Informe o nome do projeto.                                                                       |
| `project.path_required`        | `path` vazio após remover espaços                                                                       | `422`       | Informe a pasta do projeto.                                                                      |
| `project.type_required`        | `type` vazio                                                                                            | `422`       | Escolha o tipo do projeto.                                                                       |
| `project.type_invalid`         | `type` diferente de `web`, `android` e `windows`                                                        | `422`       | Tipo inválido. Escolha web, android ou windows.                                                  |
| `project.path_invalid`         | `path` não é um caminho absoluto                                                                        | `422`       | Informe o caminho completo da pasta (ex.: C:\dev-tools\projects\web\loja-api).                   |
| `project.path_not_found`       | `path` não existe ou não é uma pasta                                                                    | `422`       | A pasta informada não existe.                                                                    |
| `project.name_taken`           | Outro projeto tem o mesmo nome, sem diferenciar maiúsculas                                              | `409`       | Já existe um projeto com esse nome.                                                              |
| `project.path_taken`           | Outro projeto aponta para a mesma pasta, sem diferenciar maiúsculas                                     | `409`       | Essa pasta já está cadastrada em outro projeto.                                                  |
| `project.not_found`            | `{id}` inválido ou sem projeto                                                                          | `404`       | Projeto não encontrado. Atualize a lista e tente de novo.                                        |
| `preferences.theme_invalid`    | `theme` fora de `light`, `dark` e `system`                                                              | `422`       | Tema inválido. Escolha Claro, Escuro ou Sistema.                                                 |
| `preferences.language_invalid` | `language` fora de `pt-br` e `en`                                                                       | `422`       | Idioma inválido. Escolha pt-br ou en.                                                            |
| `preferences.save_failed`      | Falha ao gravar `config.json`                                                                           | `500`       | Não foi possível salvar as preferências. Veja os logs para mais detalhes.                        |
| `system.open_logs_failed`      | Falha ao iniciar o Explorer                                                                             | `500`       | Não foi possível abrir a pasta de logs.                                                          |
| `database.unavailable`         | Banco não abriu na inicialização, conexão fechada ou qualquer falha do SQLite que não seja de unicidade | `503`       | O banco de dados está indisponível. Reinicie o Dev Tools; se o problema continuar, veja os logs. |
| `internal.unexpected`          | Erro que não é `apperr`, ou pânico no handler                                                           | `500`       | Ocorreu um erro inesperado. Veja os logs para mais detalhes.                                     |

Mapeamento `apperr.Kind` → status, feito só no `httpserver`: `BadRequest` → 400,
`Unauthorized` → 401, `NotFound` → 404, `Conflict` → 409, `Invalid` → 422, `Internal` → 500,
`Unavailable` → 503.

Mensagens do motor que não são erros da API (mesmo catálogo):

| Chave                         | Uso                                                                                                                                                                                                                                                                          | Texto (pt-br)                                                                                                                                                  |
| ----------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `startup.title`               | Título do diálogo nativo                                                                                                                                                                                                                                                     | Dev Tools                                                                                                                                                      |
| `startup.folder_not_writable` | Diálogo nativo da RN-04; parâmetro `{path}`                                                                                                                                                                                                                                  | O Dev Tools não pode abrir: a pasta {path} não permite escrita. Instale o app numa pasta do seu usuário (ex.: C:\dev-tools\) ou ajuste as permissões da pasta. |
| `webview2.*` (11 chaves)      | `options.Windows.Messages` do Wails: `installation_required`, `update_required`, `missing_requirements`, `webview2_not_installed`, `error`, `failed_to_install`, `download_page`, `press_ok_to_install`, `contact_admin`, `invalid_fixed_webview2`, `webview2_process_crash` | Tradução do texto padrão do Wails (`windows.DefaultMessages()`)                                                                                                |

### 3.3 Outras interfaces

#### 3.3.1 Bootstrap (binding do Wails)

A única chamada pelos bindings (ADR 004): `Bridge.Bootstrap() BootstrapInfo`, em
`internal/platform/desktop`. O `Bridge` não tem nenhum outro método exportado, porque o Wails
expõe ao front todo método exportado da struct registrada em `Bind` (o CT-25 confere isso).
O front chama `Bootstrap()` uma vez, ao montar o `App`, pelo módulo gerado em
`frontend/wailsjs/go/desktop/Bridge`.

| Campo (JSON)     | Tipo                                    | Exemplo                           | Observação                                                              |
| ---------------- | --------------------------------------- | --------------------------------- | ----------------------------------------------------------------------- |
| `version`        | string                                  | `"0.1.0"`                         | `version` de `main.go`                                                  |
| `apiBaseUrl`     | string                                  | `"http://127.0.0.1:41790"`        | Com a porta efetiva (pode vir de `DEV_TOOLS_API_PORT`)                  |
| `token`          | string                                  | 64 caracteres hexadecimais        | 32 bytes de `crypto/rand`, gerado a cada execução, só em memória        |
| `apiStatus`      | `"ok"` \| `"port_in_use"` \| `"failed"` | `"port_in_use"`                   | Resultado do `Listen`                                                   |
| `apiPort`        | number                                  | `41790`                           |                                                                         |
| `databaseStatus` | `"available"` \| `"unavailable"`        | `"available"`                     | Estado na abertura; depois disso a tela de Status usa `GET /api/health` |
| `logsDir`        | string                                  | `"C:\\dev-tools\\programs\\logs"` | Mostrado no aviso de porta ocupada                                      |
| `theme`          | `"light"` \| `"dark"` \| `"system"`     | `"system"`                        | Valor atual em memória                                                  |
| `language`       | `"pt-br"` \| `"en"`                     | `"pt-br"`                         | Valor atual em memória                                                  |

Assinaturas (Go):

```go
// internal/platform/desktop
type BootstrapInfo struct { /* campos da tabela acima, com as tags json */ }
func NewBridge(static BootstrapInfo, prefs func() (theme, language string)) *Bridge
func (b *Bridge) Bootstrap() BootstrapInfo          // único método exportado

type Window struct{ /* ctx do Wails */ }
func (w *Window) Startup(ctx context.Context)        // OnStartup
func (w *Window) Focus()                             // OnSecondInstanceLaunch: WindowUnminimise + Show
func (w *Window) ApplyTheme(theme string)            // runtime.WindowSetLightTheme/DarkTheme/SystemDefaultTheme; nada antes do Startup

func WindowsTheme(theme string) windows.Theme        // "light"→Light, "dark"→Dark, outro→SystemDefault
func WebView2Messages(language string) *windows.Messages
func UniqueID(baseDir string) string                 // "dev-tools-" + 12 primeiros hex de sha256(minúsculas(Clean(baseDir)))
func ShowFatal(title, message string)                // MessageBoxW com MB_OK|MB_ICONERROR (golang.org/x/sys/windows)
```

#### 3.3.2 Contratos internos do motor

```go
// internal/apperr
type Kind string // BadRequest, Unauthorized, NotFound, Conflict, Invalid, Internal, Unavailable
type Error struct { Kind Kind; Code string; Field string; Err error }
func New(kind Kind, code string) *Error
func (e *Error) WithField(field string) *Error      // cópia com Field
func (e *Error) Wrap(cause error) *Error            // cópia com a causa (só para log; nunca vai ao cliente)
func (e *Error) Error() string; func (e *Error) Unwrap() error; func (e *Error) Is(target error) bool // compara Code
var ErrInvalidBody, ErrTokenInvalid, ErrDatabaseUnavailable, ErrInternal *Error

// internal/platform/config
type Paths struct { Base, Root, Data, Logs, WebView, DBFile, ConfigFile, LogFile, Projects string }
func ResolveBaseDir(getenv func(string) string) (string, error) // DEV_TOOLS_HOME (filepath.Abs) ou pasta de os.Executable()
func NewPaths(base string) Paths                                // Root = filepath.Dir(base); demais conforme §2.3
func EnsureLayout(p Paths) (created []string, err error)        // erro: *NotWritableError{Path}
type Settings struct { APIPort int; LogLevel slog.Level }
func LoadSettings(getenv func(string) string) (Settings, []string) // segundo retorno: avisos de valores inválidos

// internal/platform/db
func Open(ctx context.Context, path string, logger *slog.Logger) (*sql.DB, error) // Open + Ping + Migrate
func Migrate(ctx context.Context, db *sql.DB, migrations fs.FS, logger *slog.Logger) (version int, err error)
var ErrSchemaNewer error
type Pinger struct{ DB *sql.DB }
func (p Pinger) Ping(ctx context.Context) error     // DB nil → erro

// internal/platform/log
func New(logFile string, level slog.Level, console io.Writer) (*slog.Logger, io.Closer, error)

// internal/platform/i18n
func T(language, key string, params map[string]string) string // troca {chave}; sem tradução → en; sem en → a própria chave
func FromLocale(locale string) string                          // "pt-BR"/"pt_br" (sem diferenciar maiúsculas) → "pt-br"; resto → "en"
func DetectWindows() string                                    // GetUserPreferredUILanguages(MUI_LANGUAGE_NAME)[0] → FromLocale; erro → "en"
func Keys(language string) []string                            // para o teste de paridade do catálogo

// internal/platform/httpserver
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error
type Route struct { Method, Pattern string; Public bool; Handler HandlerFunc } // Pattern: "/api/projects/{id}"
type Config struct { Token string; Language func() string; Logger *slog.Logger }
func NewHandler(cfg Config, routes []Route) http.Handler // mux + /docs + middlewares (recuperação → log de acesso → CORS → token por rota)
func NewToken() (string, error)
func Listen(port int) (net.Listener, error)              // "127.0.0.1:<port>"; WSAEADDRINUSE → ErrPortInUse
var ErrPortInUse error
func WriteJSON(w http.ResponseWriter, status int, v any)
func DecodeJSON(r *http.Request, dst any) error          // MaxBytesReader 1 MiB + DisallowUnknownFields → apperr.ErrInvalidBody
```

As features expõem `func Routes(s *Service) []httpserver.Route`. `cmd/dev-tools/routes.go` junta
as rotas de todas elas em `apiRoutes(...)`, que o `main.go` usa e o CT-23 compara com o
`openapi.yaml`.

#### 3.3.3 Variáveis de ambiente

| Variável              | Padrão              | Uso                                                                                                                                                        |
| --------------------- | ------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `DEV_TOOLS_HOME`      | pasta do executável | Pasta base (ADR 003). Usada em desenvolvimento                                                                                                             |
| `DEV_TOOLS_API_PORT`  | `41790`             | Porta da API; inteiro de 1 a 65535. Valor inválido → padrão e aviso no log (`config.invalid_env`). Usada para rodar o `wails dev` ao lado do app instalado |
| `DEV_TOOLS_LOG_LEVEL` | `info`              | `debug`, `info`, `warn` ou `error`. Valor inválido → `info` e aviso no log                                                                                 |

#### 3.3.4 Instalador e desinstalador (NSIS)

Gerado por `wails build -nsis` a partir de `dev-tools/cmd/dev-tools`. Saída:
`dev-tools/build/bin/dev-tools-amd64-installer.exe`. O NSIS (`makensis`) fica instalado em
`C:\dev\programs\nsis\nsis-v3.x.x\`, com o atalho `current` no `PATH`. Ajustes no
`build/windows/installer/project.nsi` e no `wails_tools.nsh` gerados pelo Wails:

| Ajuste                 | Valor                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| ---------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Nível de execução      | `RequestExecutionLevel user` (nunca `admin`)                                                                                                                                                                                                                                                                                                                                                                                                             |
| Pasta sugerida         | `InstallDir "C:\dev-tools"`, com a página de escolha de pasta (`MUI_PAGE_DIRECTORY`)                                                                                                                                                                                                                                                                                                                                                                     |
| Arquivos               | `$INSTDIR\programs\dev-tools.exe`; cria `$INSTDIR\programs\{data,logs}` e `$INSTDIR\projects\{web,android,windows}`                                                                                                                                                                                                                                                                                                                                      |
| Atalhos                | `SetShellVarContext current`; `$SMPROGRAMS\Dev Tools.lnk`, `$DESKTOP\Dev Tools.lnk` e `$SMPROGRAMS\Desinstalar Dev Tools.lnk`                                                                                                                                                                                                                                                                                                                            |
| Desinstalador          | `WriteUninstaller "$INSTDIR\programs\uninstall.exe"`. Removidas as escritas de registro do modelo (`wails.writeUninstaller` grava em `HKLM\...\Uninstall`) e as associações de arquivo (RN-03)                                                                                                                                                                                                                                                           |
| WebView2               | Mantida a macro `wails.webview2runtime` do modelo (baixa o bootstrapper se faltar)                                                                                                                                                                                                                                                                                                                                                                       |
| Idiomas                | `MUI_LANGUAGE "English"` e `"PortugueseBR"`, com `LangString MSG_REMOVE_DATA` nos dois                                                                                                                                                                                                                                                                                                                                                                   |
| Seção de desinstalação | `MessageBox MB_YESNO\|MB_ICONQUESTION "$(MSG_REMOVE_DATA)"` ("Remover também os dados e os logs do Dev Tools?"). **Sim** → `RMDir /r "$INSTDIR\programs"`. **Não** → apaga `dev-tools.exe`, `uninstall.exe` e `RMDir /r "$INSTDIR\programs\webview"`, mantendo `data\` e `logs\`. Nos dois casos: apaga os atalhos, `RMDir "$INSTDIR\programs"` e `RMDir "$INSTDIR"` **sem `/r`** (só removem pastas vazias). Nenhum comando toca em `$INSTDIR\projects` |

#### 3.3.5 Abrir a pasta de logs

`system.ExplorerOpener.Open(path)` executa `explorer.exe <path>` com `exec.Command(...).Start()`.
Usa `Start`, e não `Run`, porque o `explorer.exe` devolve código de saída 1 mesmo quando dá certo.

## 4. Dados

### 4.1 Modelo

| Tabela     | Coluna       | Tipo      | Restrições                                             | Observação                                                                                                                                                                            |
| ---------- | ------------ | --------- | ------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `projects` | `id`         | `INTEGER` | `PRIMARY KEY AUTOINCREMENT`                            |                                                                                                                                                                                       |
| `projects` | `name`       | `TEXT`    | `NOT NULL`                                             | Nome como digitado, sem espaços nas pontas                                                                                                                                            |
| `projects` | `name_key`   | `TEXT`    | `NOT NULL UNIQUE`                                      | `strings.ToLower(name)`, calculado no núcleo (`projects.NameKey`); rede de segurança da RN-13                                                                                         |
| `projects` | `path`       | `TEXT`    | `NOT NULL`                                             | Forma guardada (D-07): relativa à raiz da instalação se a pasta estiver dentro dela (ex.: `projects\web\loja-api`), senão absoluta (ex.: `D:\trabalho\site`). Sempre `filepath.Clean` |
| `projects` | `path_key`   | `TEXT`    | `NOT NULL UNIQUE`                                      | `strings.ToLower(path)` da forma guardada (`projects.PathKey`); rede de segurança da RN-14                                                                                            |
| `projects` | `type`       | `TEXT`    | `NOT NULL CHECK (type IN ('web','android','windows'))` | RN-12                                                                                                                                                                                 |
| `projects` | `created_at` | `TEXT`    | `NOT NULL`                                             | RFC 3339, UTC                                                                                                                                                                         |
| `projects` | `updated_at` | `TEXT`    | `NOT NULL`                                             | RFC 3339, UTC                                                                                                                                                                         |

Conexão: `sql.Open("sqlite", "file:<DBFile>?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")`
e `SetMaxOpenConns(1)` (um processo só, ADR 003; evita disputa de escrita).

### 4.2 Migrações

| Arquivo                                                    | O que faz                        |
| ---------------------------------------------------------- | -------------------------------- |
| `internal/platform/db/migrations/0001_create_projects.sql` | Cria a tabela `projects` de §4.1 |

Regras do mecanismo (ADR 007):

- Nome `NNNN_<descricao>.sql`; a versão é o `NNNN`. As migrações moram todas nesta pasta, numa
  única ordem global, e cada feature nova acrescenta o próximo número.
- No `Open`: lê `PRAGMA user_version`; aplica, em ordem, cada migração com versão maior, uma
  transação por arquivo, que termina com `PRAGMA user_version = NNNN`. Loga
  `db.migration_applied`.
- Numeração com buraco ou repetida → erro (bug de programação, pego pelo CT-08).
- `user_version` maior que a última migração conhecida (banco de uma versão mais nova do app) →
  `ErrSchemaNewer`, e o app sobe em modo degradado (§6.2).

### 4.3 Outros dados guardados

| Dado              | Local                                                          | Formato                                                                                                                                       | Ausente ou inválido                                                                                                                                                                                                                                                                 |
| ----------------- | -------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Preferências      | `programs\data\config.json`                                    | `{"theme": "system", "language": "pt-br"}` (UTF-8, indentado com 2 espaços). Gravado num arquivo temporário na mesma pasta e depois renomeado | **Ausente** → tema `system` e idioma detectado do Windows (RN-09, RN-10); grava o arquivo e loga `config.created`. **JSON inválido** → mesmo padrão, regrava e loga `config.invalid` (WARN). **Um campo inválido** → só esse campo volta ao padrão; regrava e loga `config.invalid` |
| Logs              | `programs\logs\dev-tools.log` + backups `dev-tools-<data>.log` | Texto do `slog` (§8)                                                                                                                          | Pasta recriada ao iniciar (RN-05)                                                                                                                                                                                                                                                   |
| Dados do WebView2 | `programs\webview\`                                            | Interno do WebView2 (`options.Windows.WebviewUserDataPath`)                                                                                   | Recriada ao iniciar                                                                                                                                                                                                                                                                 |

Nenhum dado vai para registro, `%APPDATA%`, `%LOCALAPPDATA%` ou `PATH` (RN-03).

## 5. Regras e validações

| Regra                                                                                                                                                                                                                                                                                                  | Origem   | Onde é aplicada                                                        | Em caso de falha                                                                                                                                                                     |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- | ---------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| O instalador roda como usuário comum; o exe declara `asInvoker` no manifesto (`build/windows/wails.exe.manifest`)                                                                                                                                                                                      | RN-01    | NSIS (§3.3.4), manifesto                                               | — (nunca há pedido de elevação)                                                                                                                                                      |
| `SingleInstanceLock` com `UniqueId = desktop.UniqueID(Paths.Base)`; a segunda instância chama `Window.Focus` na primeira e encerra                                                                                                                                                                     | RN-02    | `main.go`, `desktop`                                                   | —                                                                                                                                                                                    |
| Todo caminho do app sai de `config.Paths`; o WebView2 usa `Paths.WebView`                                                                                                                                                                                                                              | RN-03    | `config`, `main.go`                                                    | —                                                                                                                                                                                    |
| Antes de tudo, `EnsureLayout` cria e apaga `Paths.Base\.write-test`                                                                                                                                                                                                                                    | RN-04    | `config.EnsureLayout`, `main.go`                                       | `*NotWritableError` → `desktop.ShowFatal` com `startup.folder_not_writable` no idioma salvo (ou detectado, se `config.json` não puder ser lido) → sai com código 1, sem abrir janela |
| `EnsureLayout` faz `MkdirAll` em `Data`, `Logs`, `WebView` e `Projects\{web,android,windows}`; nunca apaga nem altera conteúdo                                                                                                                                                                         | RN-05    | `config.EnsureLayout`                                                  | Falha de criação → `*NotWritableError` (como na RN-04). Cada pasta criada → `layout.created` (INFO)                                                                                  |
| Front, `openapi.yaml`, Redoc e migrações vão embutidos (`embed`); o front usa a pilha de fontes do sistema (`system-ui`)                                                                                                                                                                               | RN-06    | `frontend/embed.go`, `httpserver/docs.go`, `db`                        | —                                                                                                                                                                                    |
| `Listen(port)` em `127.0.0.1`; porta ocupada → `ErrPortInUse`                                                                                                                                                                                                                                          | RN-07    | `httpserver.Listen`, `main.go`                                         | Loga `api.port_in_use port=<n>` (ERROR); `apiStatus = "port_in_use"`; a janela abre assim mesmo; o front mostra o aviso (§7.2)                                                       |
| A API escuta só em `127.0.0.1`                                                                                                                                                                                                                                                                         | RN-08    | `httpserver.Listen`                                                    | — (outros computadores não recebem resposta)                                                                                                                                         |
| Toda rota exige `Authorization: Bearer <token>`, comparado com `subtle.ConstantTimeCompare`, exceto as `Public` (`GET /api/health` e `/docs*`)                                                                                                                                                         | RN-08    | middleware de token do `httpserver`                                    | `401 auth.token_invalid`; o handler não é chamado                                                                                                                                    |
| Primeira execução: idioma = `i18n.FromLocale(primeiro idioma de exibição do Windows)`                                                                                                                                                                                                                  | RN-09    | `preferences.Service.Load` (detector injetado)                         | Falha na detecção → `en`                                                                                                                                                             |
| Idioma ∈ {`pt-br`, `en`}                                                                                                                                                                                                                                                                               | RN-09    | `preferences.Service.Update`                                           | `422 preferences.language_invalid`                                                                                                                                                   |
| Primeira execução: tema `system`. Tema ∈ {`light`, `dark`, `system`}                                                                                                                                                                                                                                   | RN-10    | `preferences.Service`                                                  | `422 preferences.theme_invalid`                                                                                                                                                      |
| O motor traduz o que diz (erros da API, diálogo nativo, mensagens do WebView2); o front traduz o que desenha                                                                                                                                                                                           | RN-11    | `httpserver` (resposta de erro), `desktop`, `frontend/src/shared/i18n` | Chave sem tradução → texto em en (motor); no front, chave faltando é erro de compilação (o tipo de `en.ts` é derivado de `pt-br.ts`)                                                 |
| `name`, `path` e `type` obrigatórios (após `strings.TrimSpace`); `type` ∈ {`web`, `android`, `windows`}. Ordem da validação: obrigatórios (name → path → type), formato (`type_invalid`, `path_invalid`), existência da pasta, unicidade (nome → pasta). Responde o **primeiro** erro encontrado       | RN-12    | `projects.Service`                                                     | `422 project.name_required` · `path_required` · `type_required` · `type_invalid`                                                                                                     |
| Nome único: `NameKey(name)` não pode existir em outro projeto (na edição, o próprio projeto é ignorado: trocar só maiúsculas do próprio nome é permitido)                                                                                                                                              | RN-13    | `projects.Service` (via `Repository.ExistsName`) + `UNIQUE(name_key)`  | `409 project.name_taken`. Violação do `UNIQUE` (corrida) é mapeada pelo repositório para o mesmo erro                                                                                |
| Pasta única: `PathKey(forma guardada)` não pode existir em outro projeto                                                                                                                                                                                                                               | RN-14    | `projects.Service` (via `Repository.ExistsPath`) + `UNIQUE(path_key)`  | `409 project.path_taken` (mesmo tratamento de corrida)                                                                                                                               |
| `path` precisa ser absoluto (`filepath.IsAbs`) e, após `filepath.Clean`, existir como pasta (`os.Stat` + `IsDir`), em qualquer unidade                                                                                                                                                                 | RN-15    | `projects.Service`                                                     | `422 project.path_invalid` · `422 project.path_not_found`                                                                                                                            |
| Forma guardada: se o caminho limpo está dentro de `Paths.Root` (comparação sem diferenciar maiúsculas), guarda `filepath.Rel(Root, path)`; senão guarda o absoluto. Na leitura, relativo → `filepath.Join(Root, path)`                                                                                 | — (D-07) | `projects.Service` (`installRoot` no construtor)                       | —                                                                                                                                                                                    |
| O serviço de Projetos só **lê** o disco (`os.Stat`); nenhum código de Projetos cria, move, copia ou apaga arquivo ou pasta                                                                                                                                                                             | RN-16    | `projects.Service`                                                     | — (CT-49 confere)                                                                                                                                                                    |
| Remover exige confirmação no front; a API em si não pede confirmação                                                                                                                                                                                                                                   | RN-17    | `ConfirmDialog` em `ProjectsPage`                                      | Cancelar → nenhuma chamada à API                                                                                                                                                     |
| `folderExists` é calculado a cada `List`; projeto sem pasta continua cadastrado, editável e removível                                                                                                                                                                                                  | RN-18    | `projects.Service.List`                                                | —                                                                                                                                                                                    |
| `lumberjack`: `MaxSize 10`, `MaxBackups 5`, `MaxAge 14`, `LocalTime true`, `Compress false`. Ao iniciar, se `dev-tools.log` existe e a última modificação é de um dia anterior, chama `Rotate()` (que também dispara o descarte por idade)                                                             | RN-19    | `platform/log.New`                                                     | Falha ao rotacionar → segue gravando no arquivo atual e loga `log.rotate_failed` (WARN)                                                                                              |
| O desinstalador sempre remove o app, pergunta sobre dados e logs e nunca toca em `projects\`                                                                                                                                                                                                           | RN-20    | NSIS (§3.3.4)                                                          | —                                                                                                                                                                                    |
| Banco que não abre (erro do `Open`, `ErrSchemaNewer`) → o app segue sem banco (`*sql.DB` nulo). Repositório com `DB` nulo, ou qualquer erro do SQLite que não seja de unicidade nem `sql.ErrNoRows` → `apperr.ErrDatabaseUnavailable` (com a causa, para o log). Reabrir o banco só ao reiniciar o app | — (D-06) | `main.go`, `projects.SQLiteRepository`, `db.Pinger`                    | `503 database.unavailable`; saúde `"unavailable"`                                                                                                                                    |

## 6. Fluxos

### 6.1 Fluxo principal

**Inicialização** (`cmd/dev-tools/main.go`, nesta ordem):

1. `config.ResolveBaseDir(os.Getenv)` → `config.NewPaths(base)`.
2. `config.EnsureLayout(paths)`. Se falhar: lê o idioma (`preferences.FileRepository.Load`,
   senão `i18n.DetectWindows`), chama `desktop.ShowFatal` e encerra com `os.Exit(1)` (RN-04).
3. `config.LoadSettings(os.Getenv)`.
4. `log.New(paths.LogFile, settings.LogLevel, os.Stdout)`. Loga `app.started version=… base_dir=…`,
   cada pasta de `created` (`layout.created`) e cada aviso de `LoadSettings`.
5. `preferences.NewService(preferences.NewFileRepository(paths.ConfigFile), detectLanguage,
logger, onChange)` e `Load()`. `onChange` chama `window.ApplyTheme(p.Theme)`.
6. `db.Open(ctx, paths.DBFile, logger)`. Se falhar: loga `db.open_failed` (ERROR) e segue com
   `*sql.DB` nulo (modo degradado).
7. Monta os serviços (`projects.NewService(projects.NewSQLiteRepository(conn), paths.Root,
logger, time.Now)`, `system.NewService(version, db.Pinger{DB: conn}, system.ExplorerOpener{},
paths.Logs, logger)`), gera o token (`httpserver.NewToken`) e o handler
   (`httpserver.NewHandler(cfg, apiRoutes(...))`; `cfg.Language` lê `prefs.Get().Language`).
8. `httpserver.Listen(settings.APIPort)`. Sucesso → `go srv.Serve(ln)` e loga `api.listening
addr=127.0.0.1:41790`. `ErrPortInUse` → loga `api.port_in_use port=41790`;
   outro erro → `api.listen_failed` (ERROR).
9. `wails.Run` com: título `Dev Tools`, 1100×720 (mínimo 900×600), `AssetServer.Assets =
frontend.Assets()`, `Bind = []any{bridge}`, `OnStartup = window.Startup`,
   `OnShutdown = shutdown`, `SingleInstanceLock{UniqueId: desktop.UniqueID(paths.Base),
OnSecondInstanceLaunch: window.Focus}`, `Windows{WebviewUserDataPath: paths.WebView,
Theme: desktop.WindowsTheme(prefs.Theme), Messages: desktop.WebView2Messages(prefs.Language)}`.
10. **Front:** `App` chama `Bootstrap()` → `configureApi(apiBaseUrl, token)`, aplica idioma e
    tema e abre a tela inicial: **Projetos** se `apiStatus = "ok"`, senão **Status** com o aviso.

**Cadastrar projeto** (PRD, passos 7 e 8):

1. Front (`ProjectForm`) → `POST /api/projects` com `Authorization: Bearer <token>`.
2. `httpserver`: CORS → token válido → handler.
3. `projects.handler` → `httpserver.DecodeJSON` → `Service.Create(ctx, Input)`.
4. `Service`: valida (§5) → calcula a forma guardada, `NameKey` e `PathKey` →
   `Repository.ExistsName/ExistsPath` → `Repository.Create` → loga `project.created id name path`
   → devolve a `View` (caminho absoluto, `folderExists = true`).
5. Handler → `201` + `Project`. Front fecha o formulário e recarrega a lista (`GET /api/projects`).

Editar (`PUT`) e remover (`DELETE`, depois do `ConfirmDialog`) seguem o mesmo caminho. Remover
só apaga a linha do banco (RN-16).

**Encerramento** (`OnShutdown`): `srv.Shutdown` com prazo de 3 s → `conn.Close()` → loga
`app.stopped` → fecha o logger.

### 6.2 Fluxos alternativos e de erro

- **App já aberto** → o `SingleInstanceLock` passa o controle para a primeira instância, que
  chama `Window.Focus` (`WindowUnminimise` + `Show`); a segunda encerra sem criar janela.
- **Pasta base sem permissão de escrita** → passo 2 de §6.1: diálogo nativo traduzido, saída com
  código 1; nada é logado (a pasta de logs pode não ser gravável).
- **Porta ocupada** → passo 8: log ERROR; a janela abre; `Bootstrap` informa `port_in_use`; o
  front mostra o aviso na tela de Status e desabilita Projetos e Preferências. Para tentar de
  novo, o usuário fecha e abre o app.
- **Pasta da estrutura apagada** → recriada no passo 2 e logada no passo 4 (`layout.created`).
- **Banco indisponível na abertura** → passo 6: modo degradado. `GET /api/health` responde
  `"unavailable"`; `/api/projects*` responde `503 database.unavailable` no idioma salvo; a tela
  de Status mostra "Indisponível"; a tela de Projetos mostra o erro e o botão "Tentar de novo".
- **Banco cai durante o uso** (no teste: `conn.Close()`) → o `Ping` da saúde falha →
  `"unavailable"`; o repositório recebe `sql: database is closed` → `503 database.unavailable`.
- **`config.json` ausente ou inválido** → §4.3; o app abre com os padrões.
- **Validação ou conflito no cadastro ou edição** → `422`/`409` com `code`, `message` e `field`;
  nada é gravado; o front mostra a mensagem sob o campo indicado em `field`, ou no topo do
  formulário quando não há `field`.
- **Pasta do projeto apagada depois do cadastro** → `List` devolve `folderExists: false`; o front
  mostra o selo "pasta não encontrada"; editar e remover continuam disponíveis.
- **Projeto removido por outro caminho entre a lista e a ação** → `404 project.not_found`; o
  front mostra a mensagem e recarrega a lista.
- **Usuário cancela a remoção** → o `ConfirmDialog` fecha sem chamar a API.
- **Chamada sem token** (outro programa ou página do navegador) → `401 auth.token_invalid`;
  nada muda.
- **Erro inesperado ou pânico no handler** → middleware de recuperação → `500
internal.unexpected`; loga `http.error` (ERROR) com a causa; o cliente não recebe detalhes.
- **Motor deixa de responder durante o uso** (falha de rede no `fetch`) → o front mostra
  `errors.network` do próprio catálogo.
- **WebView2 ausente** → o Wails mostra as mensagens de `options.Windows.Messages`, traduzidas.
- **Desinstalação** → §3.3.4.

## 7. Interface

### 7.1 Telas e navegação

Sem biblioteca de rotas: o `App` guarda a tela atual num estado
(`"projects" | "preferences" | "status"`) e mostra um menu lateral fixo com os três itens.

| Tela                                     | Como se chega                                     | O que mostra                                                                                                                                                                                                                              | Ações                                                                 |
| ---------------------------------------- | ------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------- |
| Diálogo nativo de erro                   | Abrir o app numa pasta sem permissão de escrita   | `startup.folder_not_writable` com o caminho                                                                                                                                                                                               | OK (fecha)                                                            |
| Shell                                    | Abrir o app                                       | Menu lateral (Projetos, Preferências, Status) e a área da tela atual                                                                                                                                                                      | Trocar de tela                                                        |
| Projetos                                 | Menu "Projetos" (tela inicial com a API no ar)    | Tabela: nome, tipo (Web, Android, Windows, traduzido), pasta (caminho absoluto) e o selo "pasta não encontrada" quando `folderExists = false`                                                                                             | "Novo projeto"; "Editar" e "Remover" em cada linha                    |
| Formulário de projeto (painel modal)     | "Novo projeto" ou "Editar"                        | Campos **Nome** (texto), **Pasta** (texto, com a dica "Cole o caminho completo da pasta, ex.: C:\dev-tools\projects\web\loja-api") e **Tipo** (lista com "Selecione…" vazio, Web, Android, Windows); na edição, os campos vêm preenchidos | Salvar (`POST`/`PUT`), Cancelar                                       |
| Confirmação de remoção (`ConfirmDialog`) | "Remover" numa linha                              | "Remover o projeto "{nome}" do Dev Tools? A pasta e os arquivos continuam no disco."                                                                                                                                                      | Remover (`DELETE`), Cancelar                                          |
| Preferências                             | Menu "Preferências"                               | Tema: Claro, Escuro, Sistema (opções de rádio); Idioma: Português (Brasil), English (lista)                                                                                                                                               | Trocar → `PUT /api/preferences` na hora; a tela muda ao receber `200` |
| Status                                   | Menu "Status" (tela inicial com a API fora do ar) | Versão do app; banco de dados (Disponível / Indisponível); endereço da API e da documentação (`http://127.0.0.1:41790/docs`, texto selecionável); pasta de logs                                                                           | "Abrir pasta de logs" (`POST /api/system/open-logs-folder`)           |

**Tema no front:** `applyTheme(theme)` grava `data-theme="light|dark|system"` no `<html>`.
`theme.css` define as variáveis de cor para `[data-theme=light]` e `[data-theme=dark]`, e o
`system` usa `@media (prefers-color-scheme: dark)`. A barra de título acompanha pelo
`Window.ApplyTheme` (motor).

**Idioma no front:** `I18nProvider` (contexto React) com `useT()`. A chave muda ao receber `200`
do `PUT /api/preferences`, e as mensagens de erro que vêm da API já chegam traduzidas pelo motor.

### 7.2 Estados

| Tela         | Estado                              | O que o usuário vê                                                                                                                                                                                                                                          |
| ------------ | ----------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Shell        | API fora do ar (`apiStatus ≠ "ok"`) | Itens Projetos e Preferências desabilitados, com a dica "Indisponível: a API do Dev Tools não está no ar"; abre direto na tela de Status                                                                                                                    |
| Projetos     | carregando                          | "Carregando projetos…"                                                                                                                                                                                                                                      |
| Projetos     | vazio                               | "Nenhum projeto cadastrado." e o botão "Novo projeto"                                                                                                                                                                                                       |
| Projetos     | lista                               | Tabela (§7.1)                                                                                                                                                                                                                                               |
| Projetos     | erro                                | A `message` da API (ex.: `database.unavailable`) ou `errors.network`, e o botão "Tentar de novo" (refaz o `GET`)                                                                                                                                            |
| Formulário   | salvando                            | Botão "Salvando…" desabilitado; campos bloqueados                                                                                                                                                                                                           |
| Formulário   | erro de campo                       | Borda de erro e a `message` sob o campo de `field`                                                                                                                                                                                                          |
| Formulário   | erro geral                          | A `message` no topo do formulário                                                                                                                                                                                                                           |
| Confirmação  | removendo                           | Botão "Removendo…" desabilitado; em erro, a `message` no diálogo                                                                                                                                                                                            |
| Preferências | salvando / erro                     | Opções desabilitadas enquanto salva; em erro, a `message` e a volta à opção anterior                                                                                                                                                                        |
| Status       | API no ar                           | Dados de `GET /api/health` (carrega ao abrir a tela)                                                                                                                                                                                                        |
| Status       | banco indisponível                  | "Banco de dados: Indisponível", em destaque de alerta                                                                                                                                                                                                       |
| Status       | porta ocupada                       | Aviso: "A porta {apiPort} está em uso por outro programa. Feche esse programa e abra o Dev Tools de novo. O ocorrido foi registrado nos logs: {logsDir}". Versão e banco vêm do bootstrap. O botão "Abrir pasta de logs" fica desabilitado (depende da API) |
| Status       | `apiStatus = "failed"`              | Aviso: "Não foi possível iniciar a API local do Dev Tools. Veja os logs: {logsDir}"                                                                                                                                                                         |
| Status       | erro ao abrir a pasta               | A `message` de `system.open_logs_failed`                                                                                                                                                                                                                    |

## 8. Aspectos transversais

- **Segurança:**
  - Rede: só `127.0.0.1` (RN-08); nenhum outro endereço.
  - Token: 32 bytes de `crypto/rand` em hexadecimal, gerado a cada execução, só em memória
    (nunca em disco nem em log), comparado em tempo constante. Exigido em toda rota não `Public`.
  - Entrada: corpo de até 1 MiB, campos desconhecidos recusados, caminho validado com
    `filepath.IsAbs` + `filepath.Clean` + `os.Stat`. O SQL usa sempre parâmetros (`?`).
  - Erros: a causa interna (texto do SQLite, pilha do pânico) vai só para o log, nunca para o
    cliente.
- **Logs** (`log/slog`, `TextHandler`, em `logs\dev-tools.log` e em `os.Stdout`; linha no
  formato `time=2026-10-05T10:00:00.000-03:00 level=INFO msg=project.created id=1
name=loja-api path=…`):

  | Evento (`msg`)                                            | Nível                    | Atributos                                          |
  | --------------------------------------------------------- | ------------------------ | -------------------------------------------------- |
  | `app.started`                                             | INFO                     | `version`, `base_dir`                              |
  | `app.stopped`                                             | INFO                     | —                                                  |
  | `layout.created`                                          | INFO                     | `path`                                             |
  | `config.created`                                          | INFO                     | `theme`, `language`                                |
  | `config.invalid`                                          | WARN                     | `path`, `err`                                      |
  | `config.invalid_env`                                      | WARN                     | `var`, `value`                                     |
  | `db.opened`                                               | INFO                     | `path`, `schema_version`                           |
  | `db.migration_applied`                                    | INFO                     | `version`, `file`                                  |
  | `db.open_failed`                                          | ERROR                    | `path`, `err`                                      |
  | `api.listening`                                           | INFO                     | `addr`                                             |
  | `api.port_in_use`                                         | ERROR                    | `port`                                             |
  | `api.listen_failed`                                       | ERROR                    | `port`, `err`                                      |
  | `http.request`                                            | DEBUG                    | `method`, `path`, `status`, `duration_ms`          |
  | `http.error`                                              | WARN (4xx) / ERROR (5xx) | `method`, `path`, `status`, `code`, `err` (só 5xx) |
  | `project.created` / `project.updated` / `project.deleted` | INFO                     | `id`, `name`, `path`                               |
  | `preferences.updated`                                     | INFO                     | `theme`, `language`                                |
  | `log.rotate_failed`                                       | WARN                     | `err`                                              |

  Nunca aparecem: o token, o cabeçalho `Authorization` e corpos de requisição. O `os.Stdout`
  só é visível no `wails dev` (o exe final não tem console).

- **Idiomas:**
  - Motor: `internal/platform/i18n/catalog.go` com dois mapas (`"pt-br"`, `"en"`) de chave →
    texto: os 17 códigos de §3.2, `startup.*` e `webview2.*`. O idioma de cada resposta vem de
    `cfg.Language()`, que lê o valor atual do `preferences.Service`.
  - Front: `src/shared/i18n/pt-br.ts` (objeto `as const`; o tipo `Messages` é derivado dele) e
    `en.ts` (tipado como `Messages`). Contém rótulos, títulos, estados de §7.2, nomes dos tipos
    e temas, e `errors.network`.
  - Instalador: textos do NSIS em English e PortugueseBR (§3.3.4).
- **Configuração:** variáveis de §3.3.3; preferências em `config.json` (§4.3); constantes:
  porta `41790`, `version = "0.1.0"` em `main.go` (sobrescrevível com
  `-ldflags "-X main.version=…"`), política de logs da RN-19.
- **Inicialização e encerramento:** §6.1. A barra de título já abre no tema salvo
  (`Windows.Theme`) e muda junto com o front (`Window.ApplyTheme`).
- **Desenvolvimento:**
  - Rodar: em `dev-tools/cmd/dev-tools`, com `DEV_TOOLS_HOME=<repo>\dev-tools\.dev-home\programs`
    e `DEV_TOOLS_API_PORT=41791`, executar `wails dev`. O `UniqueID` vem da pasta base, então o
    app de desenvolvimento não se confunde com o instalado.
  - Testes: em `dev-tools/`, `go test ./...` (unitários) e `go test -tags integration ./...`
    (todos).
  - Tipos do front: em `dev-tools/frontend`, `npm run gen:api` gera
    `src/shared/api/schema.d.ts` a partir de `../internal/platform/httpserver/docs/openapi.yaml`
    (`openapi-typescript`, dependência de desenvolvimento). `npm run typecheck` roda o
    `tsc --noEmit`.
  - Build do front: o Vite grava em `frontend/dist/app/` (`build.outDir`). O
    `frontend/embed.go` embute `all:dist` e expõe `Assets()` como `fs.Sub(..., "dist/app")`. A
    pasta `frontend/dist/` existe no repositório com um `.gitkeep`, que o Vite não apaga por
    estar fora do `outDir`. Assim, `go test ./...` compila mesmo antes do primeiro build do front.
  - `wails.json` (em `cmd/dev-tools/`): `"name": "dev-tools"`, `"outputfilename": "dev-tools"`,
    `"frontend:dir": "../../frontend"`, `"build:dir": "../../build"`,
    `"frontend:install": "npm install"`, `"frontend:build": "npm run build"`,
    `"frontend:dev:watcher": "npm run dev"`, `"frontend:dev:serverUrl": "auto"`.

## 9. Cobertura do PRD

| História | Atendida por (componentes, operações, telas)                                                                                       | Seções                      |
| -------- | ---------------------------------------------------------------------------------------------------------------------------------- | --------------------------- |
| US-01    | Instalador NSIS como usuário comum, em `C:\dev-tools`                                                                              | §3.3.4, §5 (RN-01)          |
| US-02    | Seção de desinstalação com a pergunta sobre dados e logs                                                                           | §3.3.4, §5 (RN-20)          |
| US-03    | Manifesto `asInvoker`, `SingleInstanceLock` + `Window.Focus`, diálogo nativo, modo degradado com a porta ocupada, assets embutidos | §5, §6.1, §6.2, §7.2        |
| US-04    | `config.EnsureLayout`                                                                                                              | §2.3, §5 (RN-05)            |
| US-05    | `config.Paths` relativo ao executável, forma guardada relativa (D-07), `WebviewUserDataPath`                                       | §2.3, §4.1, §4.3, §5        |
| US-06    | `GET /api/health`, `system.Service`, tela de Status                                                                                | §3.1, §7                    |
| US-07    | `preferences` (tema), `PUT /api/preferences`, `applyTheme`, `Window.ApplyTheme`, tela de Preferências                              | §3.1, §4.3, §7              |
| US-08    | `preferences` (idioma), `i18n.DetectWindows`, catálogos do motor e do front                                                        | §3.2, §4.3, §8              |
| US-09    | `POST /api/projects`, `projects.Service.Create`, formulário                                                                        | §3.1, §5, §6.1, §7          |
| US-10    | `GET /api/projects`, `folderExists`, tela de Projetos                                                                              | §3.1, §5 (RN-18), §7        |
| US-11    | `PUT /api/projects/{id}`, `projects.Service.Update`, formulário em modo edição                                                     | §3.1, §5, §7                |
| US-12    | `DELETE /api/projects/{id}`, `ConfirmDialog`                                                                                       | §3.1, §5 (RN-16, RN-17), §7 |
| US-13    | `openapi.yaml` + Redoc em `/docs`, teste de contrato, token e loopback                                                             | §3.1, §3.3.2, §10           |
| US-14    | `platform/log` (`slog` + `lumberjack`), eventos, `POST /api/system/open-logs-folder`                                               | §4.3, §5 (RN-19), §8        |
| US-15    | Estrutura de `internal/features/projects` (núcleo, infra e testes separados), teste de imports                                     | §2, §11, §10 (CT-51)        |
| US-16    | Testes unitários com fake, testes de integração com SQLite real, banco fechado no teste                                            | §10                         |

## 10. Estratégia de testes

> **Resumo:** 82 casos · unitários 43 · integração 22 · e2e 0 · manuais 17
> **Padrão (ADR 002):** `go test`. Unitários em `<arquivo>_test.go` ao lado do código, com fakes
> das interfaces (o fake do repositório de Projetos fica em `fake_repository_test.go`).
> Integração em `*_integration_test.go`, com `//go:build integration`, SQLite real em
> `t.TempDir()` e `httptest`. A integração monta a pilha real (`config.NewPaths` sobre uma pasta
> temporária com `programs\`, `db.Open`, serviço, repositório, `httpserver.NewHandler` com token
> `"test-token"`) e um logger que grava num `bytes.Buffer`. Testes em tabela quando o mesmo
> comportamento varia só na entrada. Nomes: `TestService_Create_NameTaken`. O front não tem
> testes automatizados na v1 (nenhuma ADR prevê ferramenta de teste para ele): o que é só visual
> é coberto por casos manuais. E2e: N/A, pelo mesmo motivo.

### 10.1 Critérios de aceite → testes

| CA    | Resumo                               | Testes                                   |
| ----- | ------------------------------------ | ---------------------------------------- |
| CA-01 | Instalador sem admin, pasta sugerida | CT-66                                    |
| CA-02 | Desinstalar com "sim"                | CT-67                                    |
| CA-03 | Desinstalar com "não"                | CT-68                                    |
| CA-04 | Abrir sem admin                      | CT-69                                    |
| CA-05 | Instância única                      | CT-27, CT-70                             |
| CA-06 | Pasta sem escrita                    | CT-04, CT-71                             |
| CA-07 | Porta ocupada                        | CT-21, CT-25, CT-26, CT-72               |
| CA-08 | Offline                              | CT-20, CT-73                             |
| CA-09 | Recriar a pasta de logs              | CT-03                                    |
| CA-10 | Cópia para outro PC                  | CT-01, CT-02, CT-35, CT-36, CT-60, CT-74 |
| CA-11 | Status com o banco disponível        | CT-63, CT-65, CT-75                      |
| CA-12 | Status com o banco indisponível      | CT-09, CT-26, CT-63, CT-75               |
| CA-13 | Saúde pelo navegador                 | CT-15, CT-65                             |
| CA-14 | Tema Sistema na primeira execução    | CT-28, CT-29, CT-30, CT-76               |
| CA-15 | Tema Escuro guardado                 | CT-31, CT-33, CT-76                      |
| CA-16 | pt-br no Windows pt-BR               | CT-06, CT-29, CT-77                      |
| CA-17 | en no Windows em outro idioma        | CT-06, CT-29, CT-77                      |
| CA-18 | en guardado, inclusive erros         | CT-07, CT-17, CT-28, CT-31, CT-33, CT-77 |
| CA-19 | Cadastrar sem mexer na pasta         | CT-35, CT-49, CT-52                      |
| CA-20 | Pasta inexistente                    | CT-39, CT-40, CT-53                      |
| CA-21 | Nome repetido (maiúsculas)           | CT-41, CT-54                             |
| CA-22 | Pasta repetida                       | CT-42, CT-54                             |
| CA-23 | Campo obrigatório vazio              | CT-17, CT-37, CT-55, CT-78               |
| CA-24 | Cadastrar com o banco indisponível   | CT-59                                    |
| CA-25 | Lista com nome, tipo e pasta         | CT-43, CT-52                             |
| CA-26 | "Pasta não encontrada"               | CT-43, CT-57                             |
| CA-27 | Editar sem mexer no disco            | CT-44, CT-49, CT-52                      |
| CA-28 | Editar para nome ou pasta de outro   | CT-45, CT-56                             |
| CA-29 | Editar para pasta inexistente        | CT-46, CT-56                             |
| CA-30 | Remover com confirmação              | CT-48, CT-49, CT-52, CT-79               |
| CA-31 | Cancelar a remoção                   | CT-79                                    |
| CA-32 | Documentação da API                  | CT-20, CT-23, CT-81                      |
| CA-33 | Outro programa não altera projetos   | CT-14, CT-24, CT-34, CT-58               |
| CA-34 | Outro computador sem resposta        | CT-22, CT-82                             |
| CA-35 | Registro de eventos e erros          | CT-11, CT-18, CT-50, CT-61               |
| CA-36 | Abrir a pasta de logs                | CT-64, CT-80                             |
| CA-37 | Descarte de logs antigos             | CT-12, CT-13                             |
| CA-38 | Camadas separadas                    | CT-51                                    |
| CA-39 | RN-12 a RN-16 sem banco e sem API    | CT-35 a CT-50                            |
| CA-40 | Integração com banco real e API      | CT-08, CT-52 a CT-58, CT-65              |
| CA-41 | Banco indisponível simulado          | CT-10, CT-59, CT-65                      |

### 10.2 Casos de teste

#### Configuração e caminhos (`internal/platform/config/paths_test.go`, `settings_test.go`)

| CT    | Camada   | O que verifica                | Entrada / cenário                                                                                                   | Resultado esperado                                                                                                                                                                      | CA       |
| ----- | -------- | ----------------------------- | ------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| CT-01 | unitário | Resolução da pasta base       | `getenv` com `DEV_TOOLS_HOME=rel\programs`; depois `getenv` vazio                                                   | Caminho absoluto de `rel\programs`; sem a variável, a pasta de `os.Executable()`                                                                                                        | CA-10    |
| CT-02 | unitário | Caminhos derivados            | `NewPaths("<tmp>\\programs")`                                                                                       | `Root=<tmp>`, `Data=<tmp>\programs\data`, `DBFile=…\data\dev-tools.db`, `ConfigFile=…\data\config.json`, `LogFile=…\logs\dev-tools.log`, `WebView=…\webview`, `Projects=<tmp>\projects` | CA-10    |
| CT-03 | unitário | Recria só o que falta (RN-05) | Estrutura completa com `data\a.txt`; apagar `logs\`; `EnsureLayout`                                                 | `logs\` recriada; `created = [logs]`; `data\a.txt` intacto                                                                                                                              | CA-09    |
| CT-04 | unitário | Pasta sem escrita (RN-04)     | `Paths.Base` dentro de um caminho que é arquivo (`<tmp>\arquivo.txt\programs`)                                      | `*NotWritableError` com `Path` igual à pasta base                                                                                                                                       | CA-06    |
| CT-05 | unitário | Variáveis de ambiente         | Vazio; `DEV_TOOLS_API_PORT=41800`, `DEV_TOOLS_LOG_LEVEL=debug`; `DEV_TOOLS_API_PORT=abc`, `DEV_TOOLS_LOG_LEVEL=xyz` | `41790`/`Info`; `41800`/`Debug`; `41790`/`Info` com dois avisos                                                                                                                         | — (D-08) |

#### Idiomas do motor (`internal/platform/i18n/i18n_test.go`)

| CT    | Camada   | O que verifica            | Entrada / cenário                                                                                                                   | Resultado esperado                                                                                                               | CA           |
| ----- | -------- | ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | ------------ |
| CT-06 | unitário | Locale → idioma (RN-09)   | `pt-BR`, `pt-br`, `pt_BR`, `pt-PT`, `es-ES`, `en-US`, `""`                                                                          | `pt-br`, `pt-br`, `pt-br`, `en`, `en`, `en`, `en`                                                                                | CA-16, CA-17 |
| CT-07 | unitário | Catálogo completo (RN-11) | `Keys("pt-br")` × `Keys("en")`; `T("en", "chave.inexistente", nil)`; `T("pt-br", "startup.folder_not_writable", {"path": "C:\\x"})` | Mesmas chaves nos dois idiomas, incluindo os 17 códigos de §3.2; chave inexistente devolve a própria chave; `{path}` substituído | CA-18        |

#### Banco e migrações (`internal/platform/db/db_integration_test.go`)

| CT    | Camada     | O que verifica            | Entrada / cenário                                                                                                     | Resultado esperado                                                                                                                                     | CA    |
| ----- | ---------- | ------------------------- | --------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ | ----- |
| CT-08 | integração | Migrações (ADR 007)       | `Open` num arquivo novo; fechar; `Open` de novo. Em seguida, `Migrate` com um `fstest.MapFS` contendo `0001` e `0003` | 1º `Open`: `user_version = 1`, tabela `projects` existe, `db.migration_applied` logado uma vez. 2º `Open`: nada reaplicado. Buraco na numeração → erro | CA-40 |
| CT-09 | integração | Banco de versão mais nova | Arquivo com `PRAGMA user_version = 99`                                                                                | `Open` devolve `ErrSchemaNewer`                                                                                                                        | CA-12 |
| CT-10 | integração | `Pinger`                  | Banco aberto; depois `conn.Close()`; depois `Pinger{DB: nil}`                                                         | `nil`; erro; erro                                                                                                                                      | CA-41 |

#### Logs (`internal/platform/log/log_integration_test.go`)

| CT    | Camada     | O que verifica                          | Entrada / cenário                                                                                              | Resultado esperado                                                                                                                               | CA    |
| ----- | ---------- | --------------------------------------- | -------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ----- |
| CT-11 | integração | Formato (ADR 006)                       | `New(<tmp>\logs\dev-tools.log, Info, nil)`; `logger.Info("project.created", "id", 1)`; fechar                  | O arquivo contém uma linha com `time=` (data e hora), `level=INFO`, `msg=project.created` e `id=1`                                               | CA-35 |
| CT-12 | integração | Descarte por idade e quantidade (RN-19) | Criar 6 backups `dev-tools-<data>.log` de 1 a 6 dias atrás e 1 de 20 dias atrás; `New`; escrever mais de 10 MB | Após a rotação (aguardar o descarte assíncrono com prazo de 5 s): o backup de 20 dias foi apagado e restam no máximo 5 backups, os mais recentes | CA-37 |
| CT-13 | integração | Rotação ao iniciar                      | `dev-tools.log` com modificação de ontem (`os.Chtimes`) e um backup de 20 dias atrás; `New`                    | Um backup novo com o conteúdo de ontem; `dev-tools.log` novo; o backup de 20 dias foi apagado                                                    | CA-37 |

#### Servidor HTTP (`internal/platform/httpserver/server_test.go`, `server_integration_test.go`)

| CT    | Camada     | O que verifica             | Entrada / cenário                                                                                                            | Resultado esperado                                                                                                                                      | CA           |
| ----- | ---------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------ |
| CT-14 | unitário   | Token (RN-08)              | Rota protegida de teste chamada sem `Authorization`, com `Bearer errado` e com `Bearer test-token`                           | `401` + `auth.token_invalid` e mensagem no idioma de `cfg.Language`, sem chamar o handler; `401`; handler chamado                                       | CA-33        |
| CT-15 | unitário   | Rota pública               | Rota `Public` chamada sem `Authorization`                                                                                    | Handler chamado, `200`                                                                                                                                  | CA-13        |
| CT-16 | unitário   | Preflight CORS             | `OPTIONS /api/projects` com `Origin: http://wails.localhost`, sem token                                                      | `204` com os cinco cabeçalhos de §3.1; handler não chamado                                                                                              | — (D-01)     |
| CT-17 | unitário   | Resposta de erro (ADR 005) | Handler devolve `apperr.New(Invalid, "project.name_required").WithField("name")`, com `cfg.Language` = `pt-br` e depois `en` | `422`; corpo `{"error":{"code":"project.name_required","message":<texto do catálogo no idioma>,"field":"name"}}`; mensagens diferentes nos dois idiomas | CA-18, CA-23 |
| CT-18 | unitário   | Erro inesperado            | Handler devolve `errors.New("boom")`; outro handler entra em pânico                                                          | `500` + `internal.unexpected` nos dois; `"boom"` não aparece no corpo; log com `level=ERROR msg=http.error` e a causa                                   | CA-35        |
| CT-19 | unitário   | Corpo inválido             | `DecodeJSON` com `{`, com `{"x":1}` em struct sem `x` e com 2 MiB                                                            | `apperr.ErrInvalidBody` nos três                                                                                                                        | —            |
| CT-20 | unitário   | `/docs` offline (RN-06)    | `GET /docs`, `GET /docs/openapi.yaml`, `GET /docs/redoc.standalone.js`                                                       | `200` nos três, servidos do `embed`; o `index.html` não referencia `http://` nem `https://` em `src` ou `href`                                          | CA-08, CA-32 |
| CT-21 | integração | Porta ocupada (RN-07)      | Ocupar uma porta com `net.Listen("tcp", "127.0.0.1:0")`; `Listen(<mesma porta>)`                                             | `errors.Is(err, ErrPortInUse)`                                                                                                                          | CA-07        |
| CT-22 | integração | Só loopback (RN-08)        | `Listen(0)`                                                                                                                  | `ln.Addr()` é `127.0.0.1`                                                                                                                               | CA-34        |

#### Rotas e contrato (`cmd/dev-tools/routes_test.go`)

| CT    | Camada   | O que verifica                   | Entrada / cenário                                                                                 | Resultado esperado                                                | CA    |
| ----- | -------- | -------------------------------- | ------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- | ----- |
| CT-23 | unitário | Rotas × `openapi.yaml` (ADR 001) | `apiRoutes(...)` com dependências nulas × `paths` do `openapi.yaml` (lido com `gopkg.in/yaml.v3`) | Todo método + caminho registrado existe no contrato, e vice-versa | CA-32 |
| CT-24 | unitário | Rotas públicas                   | `apiRoutes(...)`                                                                                  | Só `GET /api/health` tem `Public = true`                          | CA-33 |

#### Desktop (`internal/platform/desktop/desktop_test.go`)

| CT    | Camada   | O que verifica                         | Entrada / cenário                                                                                                                                              | Resultado esperado                                                                       | CA           |
| ----- | -------- | -------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ------------ |
| CT-25 | unitário | Um único binding (ADR 004)             | `reflect.TypeOf(&Bridge{})`                                                                                                                                    | Exatamente um método exportado: `Bootstrap`                                              | CA-07        |
| CT-26 | unitário | Conteúdo do bootstrap                  | `NewBridge` com `APIStatus: "port_in_use"`, `APIPort: 41790`, `DatabaseStatus: "unavailable"`, `LogsDir`; `prefs` devolve `dark`/`en` e depois `light`/`pt-br` | `Bootstrap()` traz os campos estáticos e as preferências atuais a cada chamada           | CA-07, CA-12 |
| CT-27 | unitário | ID de instância (RN-02)                | `UniqueID("C:\\dev-tools\\programs")`, `UniqueID("c:\\DEV-TOOLS\\programs\\")`, `UniqueID("D:\\outro\\programs")`                                              | Os dois primeiros iguais; o terceiro diferente; todos com o prefixo `dev-tools-`         | CA-05        |
| CT-28 | unitário | Mensagens do WebView2 e tema da janela | `WebView2Messages("pt-br")`, `("en")`; `WindowsTheme("light"/"dark"/"system")`                                                                                 | Os 11 campos preenchidos e diferentes entre os idiomas; `Light`, `Dark`, `SystemDefault` | CA-14, CA-18 |

#### Preferências (`internal/features/preferences/service_test.go`, `preferences_integration_test.go`)

| CT    | Camada     | O que verifica                   | Entrada / cenário                                                                                                                         | Resultado esperado                                                                                                                                 | CA                  |
| ----- | ---------- | -------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| CT-29 | unitário   | Primeira execução (RN-09, RN-10) | Fake sem arquivo; detector devolve `pt-br`, depois `en` (tabela)                                                                          | `Get()` = `system` + idioma detectado; o fake recebeu `Save` com esses valores                                                                     | CA-14, CA-16, CA-17 |
| CT-30 | unitário   | Arquivo inválido                 | Fake devolve erro de JSON; depois `{"theme":"roxo","language":"en"}`                                                                      | Padrões (`system` + detectado); no segundo, `system` + `en`; `Save` chamado; log `config.invalid` (WARN)                                           | CA-14               |
| CT-31 | unitário   | Salvar                           | `Update({dark, en})`                                                                                                                      | `Save` chamado; `Get()` = `dark`/`en`; `onChange` chamado uma vez; log `preferences.updated`                                                       | CA-15, CA-18        |
| CT-32 | unitário   | Valores inválidos                | `Update({"roxo","en"})`, `Update({"dark","fr"})`, `Update({"",""})`                                                                       | `preferences.theme_invalid` (field `theme`); `preferences.language_invalid` (field `language`); `theme_invalid`; nenhum `Save`; `Get()` inalterado | — (RN-09, RN-10)    |
| CT-33 | integração | Persistência entre execuções     | `PUT /api/preferences {"theme":"dark","language":"en"}` com token; criar novo `Service` + `FileRepository` na mesma pasta e chamar `Load` | `200` + `{dark, en}`; `config.json` com esses valores; o novo serviço carrega `dark`/`en`                                                          | CA-15, CA-18        |
| CT-34 | integração | Token nas preferências           | `PUT /api/preferences` sem `Authorization`                                                                                                | `401`; `config.json` inalterado                                                                                                                    | CA-33               |

#### Projetos — núcleo (`internal/features/projects/service_test.go`, `arch_test.go`)

Fake do repositório em memória; pastas reais em `t.TempDir()`, com `installRoot` = uma pasta
temporária que contém `projects\web\loja-api`.

| CT    | Camada   | O que verifica                            | Entrada / cenário                                                                                                                        | Resultado esperado                                                                                                              | CA                  |
| ----- | -------- | ----------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| CT-35 | unitário | Criar dentro da instalação (D-07)         | `Create({"loja-api", <root>\projects\web\loja-api, "web"})`                                                                              | `View.Path` absoluto, `FolderExists = true`; o fake guardou `Path = projects\web\loja-api`, `NameKey = loja-api`                | CA-10, CA-19        |
| CT-36 | unitário | Criar fora da instalação                  | Pasta em outro `t.TempDir()`                                                                                                             | O fake guardou o caminho absoluto                                                                                               | CA-10               |
| CT-37 | unitário | Obrigatórios (RN-12)                      | Tabela: `name` vazio ou só espaços; `path` vazio; `type` vazio; os três vazios                                                           | `project.name_required`/`name`; `path_required`/`path`; `type_required`/`type`; com os três vazios, `name_required`; fake vazio | CA-23               |
| CT-38 | unitário | Tipo inválido (RN-12)                     | `type = "ios"`                                                                                                                           | `project.type_invalid`, field `type`                                                                                            | CA-39               |
| CT-39 | unitário | Caminho relativo (RN-15)                  | `path = "projects\\web\\loja-api"`                                                                                                       | `project.path_invalid`, field `path`                                                                                            | CA-20               |
| CT-40 | unitário | Pasta inexistente (RN-15)                 | `path` inexistente; `path` de um arquivo                                                                                                 | `project.path_not_found` nos dois; fake vazio                                                                                   | CA-20               |
| CT-41 | unitário | Nome único (RN-13)                        | Existe `loja-api`; criar `Loja-API` com outra pasta                                                                                      | `project.name_taken`, field `name`; fake com 1 projeto                                                                          | CA-21               |
| CT-42 | unitário | Pasta única (RN-14)                       | Existe projeto em `<root>\projects\web\loja-api`; criar outro com `<ROOT>\Projects\Web\Loja-API\` e com `<root>\projects\web\.\loja-api` | `project.path_taken`, field `path`, nos dois                                                                                    | CA-22               |
| CT-43 | unitário | Lista (RN-18)                             | Projetos `zeta`, `Alfa`, `beta`; apagar a pasta de `beta`                                                                                | Ordem `Alfa`, `beta`, `zeta`; `beta.FolderExists = false`; os outros `true`; nome, tipo e caminho absoluto em cada item         | CA-25, CA-26        |
| CT-44 | unitário | Editar                                    | `Update(id, {novo nome, outra pasta existente, "android"})`                                                                              | `View` com os três valores novos; `UpdatedAt` = `now()` injetado; `CreatedAt` mantido                                           | CA-27               |
| CT-45 | unitário | Editar para dados de outro (RN-13, RN-14) | Projetos A e B; editar A com o nome de B; editar A com a pasta de B; editar A com o próprio nome em maiúsculas                           | `name_taken`; `path_taken`; o terceiro aceito; nos dois primeiros, A inalterado no fake                                         | CA-28               |
| CT-46 | unitário | Editar para pasta inexistente             | Editar A com pasta inexistente                                                                                                           | `project.path_not_found`; A inalterado                                                                                          | CA-29               |
| CT-47 | unitário | Id inexistente                            | `Update(999, …)`, `Delete(999)`                                                                                                          | `project.not_found` nos dois                                                                                                    | CA-39               |
| CT-48 | unitário | Remover                                   | `Delete(id)`                                                                                                                             | Fake sem o projeto                                                                                                              | CA-30               |
| CT-49 | unitário | Nada muda no disco (RN-16)                | Retrato recursivo (nomes, tamanhos, datas) da pasta do projeto e da raiz antes e depois de `Create`, `Update` e `Delete`                 | Retratos iguais                                                                                                                 | CA-19, CA-27, CA-30 |
| CT-50 | unitário | Eventos de log                            | `Create`, `Update`, `Delete` com logger em `bytes.Buffer`                                                                                | Linhas `level=INFO` com `msg=project.created`, `project.updated`, `project.deleted`, cada uma com `id` e `name`                 | CA-35               |
| CT-51 | unitário | Regra de dependência (ADR 002)            | `go/parser` lê os imports de `project.go`, `errors.go`, `repository.go` e `service.go`                                                   | Nenhum import de `database/sql`, `net/http`, `github.com/wailsapp/`, `modernc.org/` ou `dev-tools/internal/platform/`           | CA-38               |

#### Projetos — integração (`internal/features/projects/projects_integration_test.go`)

| CT    | Camada     | O que verifica                   | Entrada / cenário                                                                                                                                                      | Resultado esperado                                                                                                      | CA                                |
| ----- | ---------- | -------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- | --------------------------------- |
| CT-52 | integração | CRUD pela API                    | Com token: `POST` (`loja-api`, pasta em `projects\web`, `web`) → `GET` → `PUT` (nome `loja`, tipo `android`) → `DELETE` → `GET`. Retrato da pasta antes e depois       | `201` + `Project`; lista com 1 item (nome, tipo, pasta); `200` com os dados novos; `204`; lista vazia; pasta intacta    | CA-19, CA-25, CA-27, CA-30, CA-40 |
| CT-53 | integração | Pasta inexistente                | `POST` com pasta inexistente                                                                                                                                           | `422` + `project.path_not_found` com a mensagem pt-br do catálogo; `GET` vazio                                          | CA-20, CA-40                      |
| CT-54 | integração | Duplicidade                      | Existe `loja-api`; `POST Loja-API` (outra pasta); `POST outro-nome` (mesma pasta)                                                                                      | `409 project.name_taken`; `409 project.path_taken`; `GET` com 1 item                                                    | CA-21, CA-22, CA-40               |
| CT-55 | integração | Obrigatórios                     | `POST {"name":"","path":"","type":""}` e `POST {}`                                                                                                                     | `422 project.name_required` com `field: "name"` nos dois; `GET` vazio                                                   | CA-23, CA-40                      |
| CT-56 | integração | Edição recusada                  | Projetos A e B; `PUT A` com o nome de B; `PUT A` com pasta inexistente                                                                                                 | `409 project.name_taken`; `422 project.path_not_found`; `GET` mostra A com os dados anteriores                          | CA-28, CA-29, CA-40               |
| CT-57 | integração | Pasta apagada depois do cadastro | Criar o projeto; `os.RemoveAll` na pasta; `GET`; `PUT` com outra pasta existente; `DELETE`                                                                             | `folderExists: false`; `200`; `204`                                                                                     | CA-26, CA-40                      |
| CT-58 | integração | Sem token (RN-08)                | `POST`, `PUT`, `DELETE` e `GET` sem `Authorization` e com `Bearer errado`                                                                                              | `401 auth.token_invalid` em todos; `SELECT COUNT(*)` e linhas inalterados                                               | CA-33, CA-40                      |
| CT-59 | integração | Banco indisponível (D-06)        | Criar 1 projeto; `conn.Close()`; `POST` e `GET` com `cfg.Language` = `pt-br`; repetir o `POST` com `en`                                                                | `503` + `database.unavailable` com a mensagem do catálogo em pt-br e depois em en; `http.error` (ERROR) no log          | CA-24, CA-41                      |
| CT-60 | integração | Portabilidade (D-07)             | Em `A\`: estrutura, `config.json` (`dark`/`en`) e projeto em `A\projects\web\loja-api`. Copiar `A\` inteira para `B\`, apagar `A\` e montar a pilha sobre `B\programs` | `GET` mostra o caminho `B\projects\web\loja-api` com `folderExists: true`; `preferences.Load` devolve `dark`/`en`       | CA-10                             |
| CT-61 | integração | Log de eventos                   | `POST` válido; `POST` com nome repetido                                                                                                                                | Buffer com `msg=project.created` (INFO) e `msg=http.error status=409 code=project.name_taken` (WARN), ambos com `time=` | CA-35                             |
| CT-62 | integração | `UNIQUE` como rede de segurança  | `SQLiteRepository.Create` direto com `NameKey` repetido; depois com `PathKey` repetido                                                                                 | `project.name_taken`; `project.path_taken`                                                                              | — (RN-13, RN-14)                  |

#### Sistema (`internal/features/system/service_test.go`, `system_integration_test.go`)

| CT    | Camada     | O que verifica        | Entrada / cenário                                          | Resultado esperado                                                          | CA                         |
| ----- | ---------- | --------------------- | ---------------------------------------------------------- | --------------------------------------------------------------------------- | -------------------------- |
| CT-63 | unitário   | Saúde                 | Fake `Pinger` sem erro; depois com erro                    | `{version: "0.1.0", database: "available"}`; depois `"unavailable"`         | CA-11, CA-12               |
| CT-64 | unitário   | Abrir a pasta de logs | Fake `FolderOpener` sem erro; depois com erro              | O fake recebeu `Paths.Logs`; depois `system.open_logs_failed`               | CA-36                      |
| CT-65 | integração | Saúde pela API        | `GET /api/health` sem token; `conn.Close()`; `GET` de novo | `200 {"version":"0.1.0","database":"available"}`; `200` com `"unavailable"` | CA-11, CA-13, CA-40, CA-41 |

#### Manuais

Feitos no build de release (`wails build -nsis`) num Windows 10 ou 11, salvo indicação.

| CT    | Camada | O que verifica         | Entrada / cenário                                                                                                                                                                                                        | Resultado esperado                                                                                                                                                                                                                                                                                                   | CA                  |
| ----- | ------ | ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| CT-66 | manual | Instalar sem admin     | 1. Entrar numa conta **padrão** (sem admin). 2. Executar o instalador. 3. Conferir a pasta sugerida; trocar para `C:\outra-pasta` e voltar para `C:\dev-tools`. 4. Concluir                                              | Nenhum pedido do UAC; sugestão `C:\dev-tools`; `C:\dev-tools\programs\dev-tools.exe`, `uninstall.exe` e `C:\dev-tools\projects\{web,android,windows}` criados; atalhos na Área de Trabalho e no Menu Iniciar (incluindo "Desinstalar Dev Tools"); nada novo em `HKLM` nem em `HKCU\Software` (conferir no `regedit`) | CA-01               |
| CT-67 | manual | Desinstalar com "sim"  | 1. Cadastrar um projeto em `C:\dev-tools\projects\web\teste` com um arquivo dentro. 2. Fechar o app. 3. "Desinstalar Dev Tools" → **Sim**                                                                                | `C:\dev-tools\programs\` removida; `C:\dev-tools\projects\web\teste\` e o arquivo continuam no disco; sem UAC                                                                                                                                                                                                        | CA-02               |
| CT-68 | manual | Desinstalar com "não"  | Reinstalar, abrir e fechar o app; desinstalar → **Não**                                                                                                                                                                  | `dev-tools.exe`, `uninstall.exe` e `webview\` removidos; `data\` (com `dev-tools.db` e `config.json`), `logs\` e `projects\` continuam                                                                                                                                                                               | CA-03               |
| CT-69 | manual | Abrir sem admin        | Na conta padrão, abrir pelo atalho                                                                                                                                                                                       | A janela abre sem pedido do UAC                                                                                                                                                                                                                                                                                      | CA-04               |
| CT-70 | manual | Instância única        | Com o app aberto e minimizado, abrir pelo atalho de novo                                                                                                                                                                 | Nenhuma segunda janela nem processo (conferir no Gerenciador de Tarefas); a janela existente volta para a frente                                                                                                                                                                                                     | CA-05               |
| CT-71 | manual | Pasta sem escrita      | 1. Copiar `programs\` para `C:\teste-ro\programs`. 2. `icacls C:\teste-ro\programs /deny "%USERNAME%":(W)`. 3. Abrir `dev-tools.exe` com `config.json` em pt-br e depois em en                                           | Diálogo nativo com `startup.folder_not_writable` e o caminho, no idioma do `config.json`; nenhuma janela do app                                                                                                                                                                                                      | CA-06               |
| CT-72 | manual | Porta ocupada          | 1. Em PowerShell: `$l=[System.Net.Sockets.TcpListener]::new([Net.IPAddress]::Loopback,41790); $l.Start()`. 2. Abrir o app                                                                                                | A janela abre na tela de Status com o aviso da porta 41790 e a pasta de logs; Projetos e Preferências desabilitados; `logs\dev-tools.log` com `msg=api.port_in_use port=41790`                                                                                                                                       | CA-07               |
| CT-73 | manual | Offline (RN-06)        | Desativar todas as redes; abrir o app; cadastrar, listar, editar e remover um projeto; abrir `http://127.0.0.1:41790/docs` no navegador                                                                                  | Tudo funciona; a documentação carrega completa                                                                                                                                                                                                                                                                       | CA-08               |
| CT-74 | manual | Outro computador       | 1. Com projetos, tema Escuro e idioma en, copiar `C:\dev-tools\` para um pendrive. 2. Em outro PC Windows, colar em `D:\dev-tools\` e abrir. 3. Conferir `%APPDATA%` e `%LOCALAPPDATA%`                                  | Mesmos projetos (os de dentro de `projects\` com o caminho em `D:\`; os de fora aparecem como "pasta não encontrada" se não existirem no PC), tema Escuro e idioma en; nenhuma pasta `dev-tools*` criada em `%APPDATA%` ou `%LOCALAPPDATA%`                                                                          | CA-10               |
| CT-75 | manual | Tela de Status         | 1. Abrir o app e ir a Status. 2. Fechar, substituir `data\dev-tools.db` por um arquivo de texto qualquer e abrir de novo, em pt-br e em en                                                                               | 1: versão e "Disponível". 2: "Indisponível" (en: "Unavailable"); a tela de Projetos mostra a mensagem de `database.unavailable` no idioma escolhido                                                                                                                                                                  | CA-11, CA-12        |
| CT-76 | manual | Tema                   | 1. Sem `config.json`, com o Windows em modo escuro, abrir o app. 2. Mudar o Windows para claro. 3. No app, escolher Escuro, fechar e abrir de novo                                                                       | 1: app escuro (Sistema). 2: app claro na hora. 3: abre no Escuro, inclusive a barra de título                                                                                                                                                                                                                        | CA-14, CA-15        |
| CT-77 | manual | Idioma                 | 1. Sem `config.json`, com o idioma de exibição do Windows em Português (Brasil), abrir. 2. Repetir com o Windows em Español. 3. Escolher en, fechar e abrir; percorrer todas as telas e provocar um erro (nome repetido) | 1: pt-br. 2: en. 3: todas as telas, estados e o erro em en                                                                                                                                                                                                                                                           | CA-16, CA-17, CA-18 |
| CT-78 | manual | Formulário             | Salvar sem nome; depois sem pasta; depois sem tipo                                                                                                                                                                       | Em cada caso, borda de erro e a mensagem sob o campo que falta; nenhum projeto novo na lista                                                                                                                                                                                                                         | CA-23               |
| CT-79 | manual | Confirmação de remoção | "Remover" → **Cancelar**; "Remover" → **Remover**                                                                                                                                                                        | Cancelar: o projeto continua na lista. Remover: sai da lista, e a pasta continua no Explorer                                                                                                                                                                                                                         | CA-30, CA-31        |
| CT-80 | manual | Abrir a pasta de logs  | Status → "Abrir pasta de logs"                                                                                                                                                                                           | O Explorer abre em `C:\dev-tools\programs\logs`                                                                                                                                                                                                                                                                      | CA-36               |
| CT-81 | manual | Documentação           | Abrir `http://127.0.0.1:41790/docs` no navegador                                                                                                                                                                         | Redoc lista `GET /api/health`, as 4 operações de Projetos, `PUT /api/preferences` e `POST /api/system/open-logs-folder`, com corpo, respostas e erros de cada uma                                                                                                                                                    | CA-32               |
| CT-82 | manual | Rede                   | De outro computador da mesma rede: `curl http://<ip-do-pc>:41790/api/health`                                                                                                                                             | Conexão recusada ou tempo esgotado; nenhuma resposta da API                                                                                                                                                                                                                                                          | CA-34               |

## 11. Arquivos

Ainda não existe código: tudo é novo, inclusive a pasta `dev-tools/` (D-11).

```
ai-software-engineering/
├── docs/adr/001…008                                  [R] ADRs que esta spec cumpre (§12)
└── dev-tools/
    ├── go.mod                                        [N] module dev-tools; go 1.25; wails/v2, modernc.org/sqlite, gopkg.in/natefinch/lumberjack.v2, golang.org/x/sys, gopkg.in/yaml.v3 (testes)
    ├── go.sum                                        [N] somas das dependências
    ├── .gitignore                                    [N] build/bin/, frontend/node_modules/, frontend/dist/app/, .dev-home/
    ├── README.md                                     [N] como rodar em desenvolvimento, testar, gerar tipos, build e instalador (§8)
    ├── cmd/dev-tools/
    │   ├── main.go                                   [N] sequência de §6.1; var version = "0.1.0"
    │   ├── routes.go                                 [N] apiRoutes(...) junta as rotas das features
    │   ├── routes_test.go                            [N] CT-23, CT-24
    │   └── wails.json                                [N] projeto Wails (§8, Desenvolvimento)
    ├── build/
    │   ├── appicon.png                               [N] ícone do app (modelo do Wails)
    │   └── windows/
    │       ├── icon.ico                              [N] ícone do exe
    │       ├── info.json                             [N] metadados do exe (nome, versão)
    │       ├── wails.exe.manifest                    [N] requestedExecutionLevel asInvoker
    │       └── installer/
    │           ├── project.nsi                       [N] script NSIS ajustado (§3.3.4)
    │           └── wails_tools.nsh                   [N] macros do Wails, sem escrita de registro
    ├── internal/
    │   ├── apperr/
    │   │   └── apperr.go                             [N] Kind, Error, New, WithField, Wrap, erros comuns
    │   ├── platform/
    │   │   ├── config/
    │   │   │   ├── paths.go                          [N] Paths, ResolveBaseDir, NewPaths, EnsureLayout, NotWritableError
    │   │   │   ├── settings.go                       [N] Settings, LoadSettings
    │   │   │   ├── paths_test.go                     [N] CT-01 a CT-04
    │   │   │   └── settings_test.go                  [N] CT-05
    │   │   ├── db/
    │   │   │   ├── db.go                             [N] Open, Migrate, ErrSchemaNewer, Pinger, embed das migrações
    │   │   │   ├── migrations/0001_create_projects.sql [N] tabela projects (§4.1)
    │   │   │   └── db_integration_test.go            [N] CT-08 a CT-10
    │   │   ├── log/
    │   │   │   ├── log.go                            [N] New (slog texto + lumberjack + rotação ao iniciar)
    │   │   │   └── log_integration_test.go           [N] CT-11 a CT-13
    │   │   ├── i18n/
    │   │   │   ├── i18n.go                           [N] T, FromLocale, Keys
    │   │   │   ├── catalog.go                        [N] textos pt-br e en (§3.2)
    │   │   │   ├── detect_windows.go                 [N] DetectWindows (GetUserPreferredUILanguages)
    │   │   │   └── i18n_test.go                      [N] CT-06, CT-07
    │   │   ├── httpserver/
    │   │   │   ├── server.go                         [N] HandlerFunc, Route, Config, NewHandler, NewToken, Listen, ErrPortInUse
    │   │   │   ├── middleware.go                     [N] recuperação, log de acesso, CORS, token
    │   │   │   ├── respond.go                        [N] WriteJSON, DecodeJSON, escrita do erro traduzido (Kind → status)
    │   │   │   ├── docs.go                           [N] embed e rotas /docs
    │   │   │   ├── docs/openapi.yaml                 [N] contrato OpenAPI 3.1 de §3.1 e §3.2
    │   │   │   ├── docs/index.html                   [N] página do Redoc, só com caminhos locais
    │   │   │   ├── docs/redoc.standalone.js          [N] Redoc 2.x, versão fixada
    │   │   │   ├── server_test.go                    [N] CT-14 a CT-20
    │   │   │   └── server_integration_test.go        [N] CT-21, CT-22
    │   │   └── desktop/
    │   │       ├── bridge.go                         [N] BootstrapInfo, Bridge, NewBridge, Bootstrap
    │   │       ├── window.go                         [N] Window (Startup, Focus, ApplyTheme), WindowsTheme
    │   │       ├── messages.go                       [N] WebView2Messages
    │   │       ├── instance.go                       [N] UniqueID
    │   │       ├── dialog_windows.go                 [N] ShowFatal (MessageBoxW)
    │   │       └── desktop_test.go                   [N] CT-25 a CT-28
    │   └── features/
    │       ├── projects/
    │       │   ├── project.go                        [N] Project, View, Input, Type, NameKey, PathKey
    │       │   ├── errors.go                         [N] erros project.* (apperr)
    │       │   ├── repository.go                     [N] interface Repository (List, Get, ExistsName, ExistsPath, Create, Update, Delete)
    │       │   ├── service.go                        [N] Service: List, Create, Update, Delete; forma guardada do caminho
    │       │   ├── sqlite_repository.go              [N] SQLiteRepository; mapeia UNIQUE, ErrNoRows e falhas do SQLite
    │       │   ├── handler.go                        [N] Routes(s) e DTOs JSON
    │       │   ├── fake_repository_test.go           [N] fake em memória
    │       │   ├── service_test.go                   [N] CT-35 a CT-50
    │       │   ├── arch_test.go                      [N] CT-51
    │       │   └── projects_integration_test.go      [N] CT-52 a CT-62
    │       ├── preferences/
    │       │   ├── preferences.go                    [N] Preferences, Theme, Language, valores padrão
    │       │   ├── errors.go                         [N] erros preferences.*
    │       │   ├── repository.go                     [N] interface Repository (Load, Save)
    │       │   ├── service.go                        [N] Service: Load, Get, Update; detector e onChange injetados
    │       │   ├── file_repository.go                [N] FileRepository (config.json, gravação atômica)
    │       │   ├── handler.go                        [N] Routes(s)
    │       │   ├── service_test.go                   [N] CT-29 a CT-32
    │       │   └── preferences_integration_test.go   [N] CT-33, CT-34
    │       └── system/
    │           ├── service.go                        [N] Service: Health, OpenLogsFolder; interfaces Pinger e FolderOpener; errors system.*
    │           ├── opener_windows.go                 [N] ExplorerOpener
    │           ├── handler.go                        [N] Routes(s)
    │           ├── service_test.go                   [N] CT-63, CT-64
    │           └── system_integration_test.go        [N] CT-65
    └── frontend/
        ├── package.json                              [N] react 19, react-dom, vite, typescript, @vitejs/plugin-react, openapi-typescript (dev); scripts dev, build, gen:api, typecheck
        ├── package-lock.json                         [N] versões travadas
        ├── tsconfig.json                             [N] strict
        ├── vite.config.ts                            [N] build.outDir = dist/app
        ├── index.html                                [N] raiz do app
        ├── embed.go                                  [N] package frontend; embed all:dist; Assets()
        ├── dist/.gitkeep                             [N] mantém a pasta do embed
        ├── wailsjs/                                  [N] gerado pelo Wails (binding Bootstrap e runtime)
        └── src/
            ├── main.tsx                              [N] monta o App
            ├── app/App.tsx                           [N] bootstrap, shell, menu, tela atual, estado da API
            ├── app/bootstrap.ts                      [N] chama Bootstrap() e tipa o BootstrapInfo
            ├── app/app.css                           [N] layout do shell
            ├── shared/api/client.ts                  [N] configureApi, apiFetch (token + JSON), ApiError {status, code, message, field}
            ├── shared/api/schema.d.ts                [N] gerado do openapi.yaml
            ├── shared/i18n/I18nProvider.tsx          [N] contexto, useT, setLanguage
            ├── shared/i18n/pt-br.ts                  [N] textos das telas (origem do tipo Messages)
            ├── shared/i18n/en.ts                     [N] textos das telas em en
            ├── shared/theme/theme.css                [N] variáveis de cor claro/escuro/sistema
            ├── shared/theme/applyTheme.ts            [N] data-theme no <html>
            ├── shared/ui/ConfirmDialog.tsx           [N] diálogo de confirmação
            ├── features/projects/api.ts              [N] listProjects, createProject, updateProject, deleteProject
            ├── features/projects/ProjectsPage.tsx    [N] lista, estados, remover
            ├── features/projects/ProjectForm.tsx     [N] criar/editar, erro por campo
            ├── features/preferences/api.ts           [N] savePreferences
            ├── features/preferences/PreferencesPage.tsx [N] tema e idioma
            ├── features/system/api.ts                [N] getHealth, openLogsFolder
            └── features/system/StatusPage.tsx        [N] saúde, aviso de porta, abrir logs
```

Legenda: `[N]` novo · `[M]` modificado · `[R]` referência (só leitura)

## 12. Conformidade com ADRs

| ADR                                                                                          | Aplica? | Como esta spec cumpre                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| -------------------------------------------------------------------------------------------- | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [ADR 001](../../../../adr/001-stack-base-do-dev-tools.md)                                    | sim     | "Usar Go 1.25+ com Wails v2 (WebView2) no app desktop, React 19 + TypeScript + Vite no front, SQLite via `modernc.org/sqlite` (sem CGO), `net/http` da biblioteca padrão na API, contrato OpenAPI escrito à mão (`openapi.yaml`) servido numa página `/docs`, e `go test` nos testes do motor." → `go.mod` com `go 1.25` e as dependências de §11; rotas no `ServeMux` com método e parâmetro; `openapi.yaml` escrito à mão, servido em `/docs` com Redoc embutido (§3.1); testes em `go test` (§10). O teste de contrato pedido nos custos é o CT-23. A exceção do bootstrap à consequência neutra ("o front fala com o motor pela API HTTP") é a da ADR 004                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| [ADR 002](../../../../adr/002-arquitetura-em-camadas-e-padrao-de-feature.md)                 | sim     | "Organizar o motor em pacotes por feature (`internal/features/<feature>/`), cada um com entidade, serviço, interface de repositório, implementação SQLite e handler HTTP, e a infraestrutura comum em `internal/platform/` (db, log, i18n, httpserver, config). O núcleo (entidade e serviço) nunca importa infraestrutura: depende de interfaces que a infra implementa, e `cmd/dev-tools/main.go` monta tudo." → `projects` segue o molde completo (§11); `preferences` tem o mesmo molde com um repositório em arquivo no lugar do SQLite, como manda a D-03; `system` não guarda dados e por isso não tem repositório. A infraestrutura comum fica em `platform/` (mais `desktop`, a camada do Wails). `internal/apperr` fica fora de `platform/` por não ser infraestrutura (tipo puro, sem I/O): assim o núcleo importa só ele e a biblioteca padrão (§2.2, CT-51). `main.go` fica em `cmd/dev-tools/`, com o `wails.json` ao lado (§13). "Testes unitários ficam no `_test.go` ao lado do código, com fakes das interfaces; testes de integração ficam em `*_integration_test.go` na própria feature, com `//go:build integration`, SQLite real em pasta temporária e `httptest`." → §10 (Padrão) |
| [ADR 003](../../../../adr/003-aplicacao-portatil.md)                                         | sim     | "Guardar tudo (banco, configurações, logs, cache…) dentro da pasta do app, que é a pasta do executável (`os.Executable`) ou, se definida, a da variável `DEV_TOOLS_HOME`; nunca usar registro, `%APPDATA%` nem `PATH`, nem pedir admin. O app roda em instância única (`SingleInstanceLock` do Wails v2: abrir de novo foca a janela aberta), recria ao iniciar as pastas que faltarem e, se a pasta não permitir escrita, fecha com uma mensagem clara." → `config.ResolveBaseDir`/`Paths` (§2.3, §3.3.2); dados do WebView2 em `programs\webview`; instalador sem registro e sem admin (§3.3.4); `SingleInstanceLock` + `Window.Focus`; `EnsureLayout` e o diálogo nativo (§5). A pasta `projects\` fica na raiz da instalação, fora da pasta do app, porque guarda projetos do usuário, não dados do app. O `UniqueId` vem da pasta base (§13)                                                                                                                                                                                                                                                                                                                                                        |
| [ADR 004](../../../../adr/004-acesso-a-api-local-token-por-execucao-e-bootstrap-do-front.md) | sim     | "Gerar a cada execução um token aleatório, mantido só em memória e exigido nas operações que alteram dados. O front recebe o token por uma única chamada interna do Wails de bootstrap, que entrega também o endereço da API, o estado da inicialização e as preferências; depois dela, toda comunicação segue pela API HTTP." → `httpserver.NewToken` + middleware de token (§3.1, §8); `Bridge.Bootstrap` como único método exportado (§3.3.1, CT-25). O token é exigido também na listagem de projetos, para cumprir a RN-08 ("operações de Projetos"); isso é mais restritivo que a ADR, e não contrário a ela. O tema da barra de título é aplicado pelo motor (`Window.ApplyTheme`), sem binding novo                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| [ADR 005](../../../../adr/005-padrao-de-erros-e-idiomas.md)                                  | sim     | "Cada camada traduz o que ela mesma produz: o motor devolve todo erro com um código estável e a mensagem já traduzida no idioma salvo, usando o próprio catálogo (erros e diálogos nativos); o front tem o catálogo das telas (rótulos, títulos, textos de status)." → formato `{error: {code, message, field}}` e 17 códigos (§3.1, §3.2); catálogo do motor com erros, `startup.*` e `webview2.*`; catálogo do front em `shared/i18n` (§8)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| [ADR 006](../../../../adr/006-padrao-de-logging.md)                                          | sim     | "Usar `log/slog` da biblioteca padrão em formato texto (`chave=valor`), com rotação e descarte por `lumberjack`: arquivos de até 10 MB, no máximo 5 antigos e descarte após 14 dias (no máximo ~60 MB em disco)." → `platform/log.New` com `TextHandler` e `lumberjack` (`MaxSize 10`, `MaxBackups 5`, `MaxAge 14`) (§5, RN-19); a rotação ao iniciar só antecipa a rotação do próprio `lumberjack`; eventos em §8                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| [ADR 007](../../../../adr/007-migracoes-do-banco-de-dados.md)                                | sim     | "Usar migrações SQL numeradas, embutidas no executável (`embed`) e aplicadas ao iniciar, com a versão aplicada guardada no próprio arquivo do banco (`PRAGMA user_version`); cada feature nova acrescenta uma migração." → `platform/db/migrations/0001_create_projects.sql`, `db.Migrate` com uma transação por arquivo e `PRAGMA user_version` (§4.2)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| [ADR 008](../../../../adr/008-organizacao-do-front.md)                                       | sim     | "Usar React puro com CSS próprio: pastas por feature espelhando o motor, tema por variáveis CSS (o tema "Sistema" segue `prefers-color-scheme`), chamadas com `fetch` e tipos gerados do `openapi.yaml`, sem bibliotecas de estado nem de componentes." → `src/features/{projects,preferences,system}` espelhando o motor; `theme.css` com variáveis e `prefers-color-scheme`; `apiFetch` sobre `fetch`; `schema.d.ts` gerado por `openapi-typescript` (ferramenta de build, não biblioteca de runtime); nenhuma dependência de estado, componentes ou rotas (§7, §11)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |

## 13. Riscos e hipóteses

| Risco ou `[HIPÓTESE]`                                                                                                                                                                                                                                                          | Impacto                                                                    | Como validar ou mitigar                                                                                                                                                                   |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `[HIPÓTESE]` O Wails v2 compila, roda o `wails dev` e gera os bindings com o `wails.json` em `cmd/dev-tools/`, usando `../../frontend` e `../../build` (conferido no fonte da v2.16.0: o build roda `go build` na pasta do projeto, e o `go.mod` é procurado nas pastas acima) | Alto: é o que permite cumprir o caminho `cmd/dev-tools/main.go` da ADR 002 | Validar na primeira task, com `wails build` e `wails dev` de ponta a ponta. Se falhar: parar e registrar uma ADR que substitua a 002 só nesse caminho (`main.go` na raiz de `dev-tools/`) |
| `[HIPÓTESE]` (Tech Alignment, D-04) A rotação do `lumberjack` funciona no Windows com o arquivo aberto, e o descarte por idade roda ao iniciar                                                                                                                                 | Médio: logs crescendo além do limite                                       | CT-12 e CT-13; a rotação ao iniciar de §5 (RN-19) garante o descarte ao menos uma vez por dia de uso                                                                                      |
| `[HIPÓTESE]` (Tech Alignment, D-09) O modelo NSIS do Wails, ajustado (§3.3.4), instala em `C:\dev-tools` como usuário comum                                                                                                                                                    | Médio: CA-01 a CA-03                                                       | CT-66 a CT-68. Se não der, o PRD (seção 10) permite levar o instalador para a v2                                                                                                          |
| `[HIPÓTESE]` O bootstrapper do WebView2 chamado pelo instalador instala sem admin quando falta o runtime                                                                                                                                                                       | Baixo: Windows 10 e 11 atualizados já trazem o WebView2                    | Testar numa VM sem WebView2. As mensagens do Wails saem traduzidas (`webview2.*`)                                                                                                         |
| `[HIPÓTESE]` (Tech Alignment, D-10) O `openapi-typescript` gera tipos corretos a partir do nosso `openapi.yaml` 3.1                                                                                                                                                            | Médio: tipos do front                                                      | Rodar `npm run gen:api` + `npm run typecheck` na primeira task do front. Se falhar, parar e perguntar: escrever os tipos à mão diverge da ADR 008                                         |
| `[HIPÓTESE]` O WebView2 (origem `wails.localhost`) chama `http://127.0.0.1:41790` com `Authorization` sem ser barrado por CORS ou pelo Private Network Access do Chromium                                                                                                      | Alto: sem isso o front não fala com a API                                  | Os cabeçalhos de §3.1 já incluem `Access-Control-Allow-Private-Network: true`. Validar no `wails dev` e no build de release na primeira task que liga o front à API                       |
| `[HIPÓTESE]` `runtime.WindowSetDarkTheme/LightTheme/SystemDefaultTheme` muda a barra de título na hora                                                                                                                                                                         | Baixo: só a barra de título                                                | CT-76. Se não mudar na hora, a barra acompanha no próximo início (`Windows.Theme`)                                                                                                        |
| `[HIPÓTESE]` (PRD) Criar `C:\dev-tools` não exige admin numa conta padrão                                                                                                                                                                                                      | Médio: CA-01                                                               | CT-66                                                                                                                                                                                     |
| O `UniqueId` por pasta base permite abrir **duas instalações diferentes** ao mesmo tempo; a segunda fica com a porta ocupada                                                                                                                                                   | Baixo: caso raro, e cada instalação tem o próprio banco                    | Aceito: é o que permite rodar o `wails dev` ao lado do app instalado. O comportamento é o da porta ocupada (RN-07)                                                                        |
| Banco aberto por uma versão mais nova do app (`ErrSchemaNewer`)                                                                                                                                                                                                                | Baixo                                                                      | Modo degradado + `db.open_failed` no log; nada é alterado no arquivo                                                                                                                      |
| Desinstalar com o app aberto deixa o exe para trás                                                                                                                                                                                                                             | Baixo                                                                      | O NSIS informa que não conseguiu apagar o arquivo. O roteiro dos CT-67 e CT-68 fecha o app antes                                                                                          |
| O Redoc embutido aumenta o exe (cerca de 1 MB)                                                                                                                                                                                                                                 | Baixo                                                                      | Aceito: é o preço de a documentação funcionar offline (RN-06)                                                                                                                             |

## 14. Próximo passo

- **Próxima fase:** Task Plan → `/agent-spec-task-plan docs/specs/features/001-arquitetura-base/v1/spec.md`
- **ADR a registrar antes:** nenhuma
