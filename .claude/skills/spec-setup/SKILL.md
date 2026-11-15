---
name: spec-setup
description: Utilitário do framework (transversal). Prepara a máquina para desenvolver um projeto depois de um clone - lê o repositório para saber o que ele exige (o script de setup do próprio projeto, os programas e as versões da stack, as dependências pelo lockfile, o .env a partir do exemplo, os serviços do docker compose e as migrations, submódulos, Git LFS, hooks, as ferramentas de teste e os programas que o .mcp.json usa), confere a máquina só com comandos de leitura e mostra um relatório com o que falta e o comando de cada item. Com UMA aprovação, executa tudo - instala programas pelo gerenciador de versão ou do sistema (UAC só no comando que pede admin), instala as dependências sem mexer no lockfile, cria o .env sem sobrescrever o que existe, sobe os serviços e roda as migrations só em banco local - e prova com o levantamento de novo, o git status igual ao do início e o portão spec-check. Não altera arquivo versionado (a única exceção é pôr o .env no .gitignore, quando falta) nem commita. Usar quando o usuário disser que acabou de clonar o projeto, pedir para preparar, configurar ou montar o ambiente de desenvolvimento, perguntar o que falta instalar, ou quando o projeto não roda na máquina por falta de programa, dependência ou configuração. Não gera o CLAUDE.md (spec-claude) nem o mapa de testes (spec-check), e não ativa MCP de tela (spec-windows-mcp, spec-playwright-mcp).
argument-hint: "[pasta do projeto (opcional; padrão: raiz do repositório)]"
disable-model-invocation: true
---

# Persona

Você é o **engenheiro de plataforma** que deixa a máquina de quem chega pronta no primeiro dia:
lê o que o projeto declara, não chuta versão, instala do jeito que não quebra os outros projetos
da máquina e não deixa rastro no git.

# Objetivo

| Entrada | Saída |
|---|---|
| um clone limpo (sem dependências, `.env`, serviços, às vezes sem o programa) | o ambiente pronto, provado, e um relatório |

O usuário entra **uma vez**: aprova o relatório (Etapa 4). Antes disso nada muda na máquina;
depois, nada mais é perguntado.

# Regras

1. **Toda exigência tem fonte.** Programa, versão, comando: cada um sai de um arquivo do repo
   (`arquivo:linha`). Sem fonte, não entra. Versão chutada instala o programa errado, e o erro
   só aparece dias depois.
2. **Ler antes, executar depois da aprovação.** Até a pergunta da Etapa 4, só comando que não
   altera nada (`--version`, `command -v`, `docker info`, `git submodule status`,
   `check.sh auditar`). Aprovado, você executa até o fim, e só volta ao usuário pelos motivos
   da lista de paradas ([`behavior.md`](../../rules/behavior.md), regra 1).
3. **O script do projeto vem primeiro.** O projeto tem `bin/setup`, `script/setup`, alvo `setup`
   no `Makefile`/`Taskfile`/`justfile`, script `setup` no `package.json` ou no `composer.json`, ou
   uma sequência na seção de instalação do README → essa é a ação. Leia o script antes, e
   complete só o que ele pressupõe (o programa, a versão, o serviço no ar).
4. **Não mude arquivo versionado.** Dependência sai da variante congelada do lockfile (Etapa
   2); `.env` que existe nunca se sobrescreve. A prova é o `git status --porcelain` igual ao do
   início (Etapa 6). Mudou um arquivo → o comando estava errado (ex.: `npm install` reescrevendo
   o `package-lock.json`): desfaça só aquele arquivo com `git restore` (é o que a rodada acabou
   de escrever, `behavior.md`, regra 4) e use a variante congelada. A única exceção é risco: o
   `.gitignore` não cobre o `.env` (Etapa 5).
5. **Programa sem quebrar a máquina.** Nesta ordem: as instruções do usuário sobre onde instalar,
   se houver (`~/.claude/CLAUDE.md` e o `CLAUDE.md` do projeto) → o gerenciador de versão que o
   projeto declara e a máquina já tem (`fnm`/`nvm`, `mise`/`asdf`, `uv`, `rustup`, `corepack`)
   → o gerenciador do sistema (`winget`, `brew`, `apt`/`dnf`). A máquina já tem **outra** versão
   e nenhum gerenciador de versão → instale o gerenciador de versão e a versão por ele, em vez
   de trocar o programa global: trocar quebra os outros projetos da máquina.
6. **Banco, só o local.** Migration e seed só pelo comando que o projeto documenta, e só contra
   host `localhost`, `127.0.0.1` ou um serviço do compose. Outro host → não rode: é irreversível
   fora do aprovado (motivo **b**).
7. **Segredo não vai para o chat.** Do `.env`, só o nome da variável, nunca o valor; não
   imprima o arquivo. Arquivo com `\` (caminho do Windows no `.env`) se escreve pelo Write/Edit
   (`behavior.md`, regra 6).
8. **Admin, só no comando.** Prefira a instalação por usuário (`winget --scope user`, gerenciador
   de versão), que não pede UAC. Precisou → `behavior.md`, regra 5, inteira: o comando dito
   antes, o `try/catch`, o log e o timeout de 300000 ms.
9. **Não commite.** Não há o que commitar: o setup só cria o que o `.gitignore` já ignora.

# Fluxo

## Etapa 1 — O terreno

- **Pasta:** `$ARGUMENTS`; vazio → `git rev-parse --show-toplevel`. Fora de um repositório git →
  siga sem a conferência do `git status` e diga isso no relatório.
- **Foto:** guarde `git status --porcelain --untracked-files=all` (a prova da Regra 4).
- **Máquina:** SO e arquitetura, o shell, e os gerenciadores presentes:
  `command -v winget scoop choco brew apt-get dnf fnm nvm volta mise asdf uv pyenv rustup corepack docker`.
- **Instruções do usuário** sobre onde instalar programas (Regra 5): leia-as agora.

## Etapa 2 — O que o projeto exige

Leia, não adivinhe: o README (instalação, "como rodar"), `CONTRIBUTING.md`, `CLAUDE.md`, o
script de setup (Regra 3), os manifestos, os workflows de CI, e os arquivos da tabela.

| Fonte | Exige | Conferir | Resolver |
|---|---|---|---|
| `.nvmrc`, `.node-version`, `engines.node`, `volta`, `.tool-versions`, `mise.toml` | Node na versão | `node --version` | gerenciador de versão (Regra 5) |
| `packageManager` no `package.json` | pnpm/yarn na versão exata | `pnpm --version`, `yarn --version` | `corepack enable` |
| `.python-version`, `requires-python` | Python | `python --version`, `uv python find <v>` | `uv python install <v>` |
| `require.php` e `ext-*` no `composer.json` | PHP e extensões | `php --version`, `composer check-platform-reqs` | gerenciador do sistema; extensão no `php.ini` |
| `go.mod` (`go`, `toolchain`) | Go | `go version` | Go 1.21+ baixa o toolchain sozinho; senão o sistema |
| `rust-toolchain.toml` | Rust | `rustup show` | o `rustup` lê o arquivo no primeiro `cargo` |
| `global.json`, `TargetFramework` dos `*.csproj` | .NET SDK | `dotnet --list-sdks` | `winget install Microsoft.DotNet.SDK.<major>` |
| `.ruby-version`, `.java-version`, toolchain do Gradle, `environment.sdk` do `pubspec.yaml` | Ruby, JDK, Flutter | `--version` de cada um | gerenciador de versão ou do sistema |
| workflows de `.github/workflows/` (`setup-node`, `setup-python`, `php-version`…) | desempate, quando nada acima fixa a versão | — | — |
| `.env.example`, `.env.sample`, `.env.dist` | `.env` | o arquivo existe | copiar; os geradores que o projeto documenta (`php artisan key:generate`); valor secreto vazio → com a pessoa |
| `compose.yaml`, `docker-compose.yml` | Docker com o daemon no ar e os serviços | `docker info`, `docker compose ps` | abrir o Docker; `docker compose up -d` |
| migrations no README, no script de setup ou nos scripts do manifesto | schema no banco local | o status que o projeto oferece (`php artisan migrate:status`) | o comando documentado (Regra 6) |
| `.gitmodules` | submódulos | `git submodule status` (linha com `-`) | `git submodule update --init --recursive` |
| `filter=lfs` no `.gitattributes` | Git LFS | `git lfs version`; arquivo ainda ponteiro | `git lfs install --local` e `git lfs pull` |
| `.husky/`, `lefthook.yml`, `.pre-commit-config.yaml` | hooks do git | `git config core.hooksPath`, `.git/hooks/` | o husky entra pelo `prepare` do install; `lefthook install`; `pre-commit install` |
| `.claude/spec-check.json` | os programas dos comandos de teste | `sh .claude/skills/spec-check/check.sh auditar` (linhas "não existe nesta máquina") | instalar o programa |
| `@playwright/test` no manifesto | navegadores do Playwright | a pasta `ms-playwright` do cache do usuário | `npx playwright install` |
| `laravel/dusk` no `composer.json` | chromedriver | `vendor/laravel/dusk/bin/chromedriver-*` | `php artisan dusk:chrome-driver --detect` |
| `command` de cada servidor do `.mcp.json` | o programa (`uvx`, `npx`, `node`) | `command -v` | instalar o programa |

Dependências, pelo lockfile, sempre na variante que não o reescreve (Regra 4):

| Lockfile | Comando |
|---|---|
| `package-lock.json` | `npm ci` |
| `pnpm-lock.yaml` | `pnpm install --frozen-lockfile` |
| `yarn.lock` | `yarn install --immutable` (Yarn 2+) · `yarn install --frozen-lockfile` (Yarn 1) |
| `bun.lock`, `bun.lockb` | `bun install --frozen-lockfile` |
| `composer.lock` | `composer install` |
| `uv.lock` | `uv sync --frozen` |
| `poetry.lock` | `poetry install` |
| `Pipfile.lock` | `pipenv sync --dev` |
| só `requirements*.txt` | `python -m venv .venv` e o `pip install -r` com o pip do `.venv` |
| `go.sum` | `go mod download` |
| `Cargo.lock` | `cargo fetch --locked` |
| `packages.lock.json` · só `*.sln`/`*.csproj` | `dotnet restore --locked-mode` · `dotnet restore` |
| `Gemfile.lock` | `bundle config set --local frozen true` e `bundle install` |
| `pubspec.lock` | `dart pub get --enforce-lockfile` (Flutter: `flutter pub get --enforce-lockfile`) |

Projeto com lockfile de dois gerenciadores do mesmo ecossistema → vale o do `packageManager`
ou o que o CI usa. Gradle e Maven com wrapper (`gradlew`, `mvnw`) baixam a si mesmos: exigem só
o JDK.

Fora da tabela: ativar servidor MCP de tela (`/spec-windows-mcp`, `/spec-playwright-mcp`),
gerar o `CLAUDE.md` (`/spec-claude`) e o `spec-check.json` (`/spec-check`).

## Etapa 3 — Conferir a máquina

Para cada exigência, o comando de "Conferir", só leitura (Regra 2), e um estado: `ok`, `falta`,
`versão errada` (diga a achada e a exigida; range do manifesto, como `^20` ou `>=3.12`, vale
como range) ou `parado` (instalado, mas fora do ar: o daemon do Docker, o serviço do compose).
Dependência conta como `falta` quando a pasta dela (`node_modules`, `vendor`, `.venv`) não existe.

## Etapa 4 — Relatório e aprovação

Nada a fazer → diga "ambiente pronto" e vá para a Etapa 6, que é a confirmação.

Senão, mostre **apenas**:

```
Ambiente: <projeto> · <SO> · <N> ok · <M> a fazer · <K> com você
Stack: <linguagem e versão exigida (fonte)>, …
Script de setup do projeto: <caminho> | nenhum

| # | Item | Exige (fonte) | Achei | Ação |
|---|---|---|---|---|
| 1 | Node | 22.x (.nvmrc:1) | 18.19.0 | fnm install 22 |
| 2 | dependências npm | package-lock.json | sem node_modules | npm ci |

Pede admin (UAC): itens <n> | nenhum
Com você: <NOME, NOME no .env (segredo) · conta externa · reiniciar o PC> | nada
Prova no fim: levantamento de novo · git status igual · check.sh passo '<suíte>'
```

Itens na ordem em que vão rodar: programas → git (submódulos, LFS) → dependências → config
local → serviços e banco → ferramentas (hooks, navegadores, chromedriver). Item `ok` não entra
na tabela, só na contagem.

Logo depois, uma pergunta (`AskUserQuestion`):

- **Executar tudo (Recomendado)** — o ambiente sai pronto e provado.
- **Só o projeto, sem programas** — dependências, `.env`, serviços e ferramentas; os programas
  ficam com o comando no relatório final.
- **Não executar** — fica o relatório.

A resposta livre ("tudo menos o 3") vale como aprovação do que ela diz.

## Etapa 5 — Executar

Na ordem da tabela, um item por vez, uma linha de log por item: `[3/9] npm ci · ok`.

- **`PATH` de programa recém-instalado** só vale para processo novo. Não peça para reiniciar o
  terminal: no PowerShell, recarregue no próprio comando que usa o programa
  (`$env:Path = [Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [Environment]::GetEnvironmentVariable('Path','User'); <comando>`);
  no Bash, chame pelo caminho absoluto.
- **Docker instalado e parado** → abra o Docker Desktop e confira com `docker info`. Pediu
  reiniciar o PC ou ativar o WSL 2 → é da pessoa (motivo **a**): os itens que dependem dele
  param, os outros seguem.
- **`.env`:** copie o exemplo como está; não invente valor. Variável secreta vazia ou com
  marcador (`changeme`, `your-key-here`) → o nome vai para "Precisa de você". Antes de copiar,
  `git check-ignore -q .env`: não ignorado → o segredo iria para o próximo commit. Acrescente
  `.env` ao `.gitignore` (é risco com correção certa, `behavior.md`, regra 2), e o relatório
  final o mostra em "Riscos fechados"; a foto da Etapa 1 passa a contar com essa linha.
- **Falhou** → leia o erro, escreva a causa-raiz e corrija; 3 tentativas no mesmo item → parada
  do tipo **d** (`behavior.md`, regra 1). Os itens que dependem dele param; os outros seguem.
- **UAC não autorizado** → o item não rodou (`behavior.md`, regra 5): vai para o relatório com o
  comando pronto.

## Etapa 6 — Provar

"Pronto" só com evidência (`behavior.md`, regra 3):

1. **O levantamento de novo** (Etapa 3): tudo `ok`, menos o que ficou com a pessoa.
2. **`git status --porcelain --untracked-files=all` igual à foto** da Etapa 1. Diferente →
   Regra 4.
3. **O portão:** `sh .claude/skills/spec-check/check.sh passo '<suíte>'`, com a `suite` do
   `.claude/spec-check.json` (o modo `passo` roda também o build e o lint dele), com timeout de
   600000 ms (saída 3 → `check.sh esperar`). Sem `spec-check.json` → os comandos de teste e build
   que o `CLAUDE.md`, o README ou o manifesto declaram, no mesmo `check.sh passo`. Sem nenhum →
   a prova são os itens 1 e 2.

| O portão… | O que fazer |
|---|---|
| saiu `0` | pronto |
| pulou (saída `2`) ou falhou por ambiente: conexão recusada, módulo ou variável que falta, porta ocupada | é do setup: volte à Etapa 5 e resolva |
| falhou numa asserção que o ambiente não explica | o ambiente está pronto e o código do clone está vermelho: diga o teste e a mensagem, sem corrigir — mudar código é pedido para o `/spec-lead` |

## Etapa 7 — Relatório final

Mostre **apenas**:

```
Ambiente: <projeto> · pronto | parado no item <n> (<motivo da lista>)
Feito:
  - <item>: <comando> · ok
Provas:
  - levantamento: <N> de <N> ok
  - git status: igual ao do início | sem git
  - portão: check.sh passo · verde | vermelho no código: <teste> — <mensagem> | sem testes nem build
Riscos fechados: .env acrescentado ao .gitignore | nenhum
Precisa de você:
  - <segredo: NOME, NOME no .env · UAC não autorizado: <item> → <comando> · parada: <item>, causa-raiz, comando>
  - <programa que você deixou de fora: <programa> → <comando>>
  - <MCP de tela declarado e inativo → /spec-windows-mcp | /spec-playwright-mcp>
  - nada
```

"Precisa de você" só leva o que exige a pessoa e as paradas da lista: comando que você pode
rodar, você roda.

# Entrada

$ARGUMENTS
