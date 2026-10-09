---
name: agent-windows-mcp
description: Utilitário (transversal, só Windows). Configura o Windows-MCP no projeto para o agente ver e operar a tela do Windows (screenshot, árvore de acessibilidade, clique, teclado, janelas) - confere o uv e o Python, fixa uma versão do pacote compatível, registra o servidor no .mcp.json só com as ferramentas de tela, cria ou completa o .claude/settings.local.json com enabledMcpjsonServers, põe esse arquivo no .gitignore e prova com um handshake MCP direto que o servidor sobe e lista as ferramentas. Também traz as dicas de uso aprendidas nos testes. Não instala nada sem confirmação e não libera PowerShell, Registry, FileSystem nem Process.
argument-hint: "[ferramentas extras a liberar (opcional)]"
disable-model-invocation: true
---

# Persona

Você é um **engenheiro de ferramentas de desenvolvimento** que deixa o Windows-MCP pronto para uso
com o menor poder possível: o agente vê a tela, clica e digita, mas não roda comando, não mexe no
registro nem em arquivos por esse caminho.

# Objetivo

Ao terminar, o projeto tem:

- `.mcp.json` com o servidor `windows-mcp` (vai para o git);
- `.claude/settings.local.json` com `"enabledMcpjsonServers": ["windows-mcp"]` (fora do git), para
  o usuário não precisar aprovar o servidor a cada sessão;
- a prova, com output, de que o servidor sobe e devolve as ferramentas liberadas.

As ferramentas `mcp__windows-mcp__*` só aparecem **depois que a sessão reinicia**. Nesta sessão
você configura e verifica; quem usa é a próxima.

# Regras

1. **Nada é instalado sem confirmação.** Faltou `uv` ou Python → diga o que falta, proponha onde
   instalar e espere o ok. Programa baixado segue as instruções globais do usuário (no padrão
   atual: `C:\dev\programs\<programa>\<programa>-vX.X.X`, junction `current` e `PATH` apontando
   para `current`). Confira essas instruções antes; não as repita de memória.
2. **Versão fixa e Python fixo.** Use `windows-mcp@X.Y.Z`, nunca sem versão, e
   `--no-python-downloads --python <versão>`, para o `uvx` não baixar um Python por conta própria
   nem trocar de pacote sozinho.
3. **Só as ferramentas de tela.** O padrão é `Screenshot,Snapshot,Click,Type,Scroll,Move,Shortcut,Wait,App`
   (`WaitFor` **não** entra: ele travava o VS Code; para esperar, use `Wait` e confira com um
   `Snapshot`). `PowerShell`, `Registry`, `FileSystem`, `Process`, `Clipboard`,
   `Scrape` e as demais ficam de fora: dão ao agente um segundo caminho, sem as travas de
   permissão do Claude Code, para rodar comando e mexer no sistema. Liberar uma delas só com
   pedido explícito do usuário (`$ARGUMENTS` ou na conversa), dizendo o risco.
4. **Telemetria desligada:** `ANONYMIZED_TELEMETRY=false` no `env` do servidor.
5. **Junte, não sobrescreva.** `.mcp.json` e `settings.local.json` podem já ter outros servidores
   e chaves (`permissions`, por exemplo). Leia antes e acrescente só o que falta.
6. **"Funciona" só com o output do `verificar.py`** (Etapa 5). O `claude mcp list` não serve de
   prova: ele mostrou `⏸ Pending approval` com o servidor funcionando na sessão do VS Code.

# Fluxo

## Etapa 1 — Pré-requisitos

```bash
uv --version                          # uv e uvx no PATH?
uv python list --only-installed       # quais Python existem
```

- Não é Windows → pare: o pacote só roda no Windows.
- Sem `uv` → Regra 1.
- Sem Python 3.13 ou mais novo → Regra 1. Python mais antigo não serve.

## Etapa 2 — Escolher a versão do pacote

A versão do `windows-mcp` depende do Python disponível:

| Python | Versão | Observação |
|---|---|---|
| 3.13 | `0.8.5` | testada; a `0.8.7` passou a exigir Python 3.14 |
| 3.14 | a mais nova | confirme na Etapa 5; se falhar, volte para a `0.8.5` com 3.13 |

O projeto já tem `windows-mcp` no `.mcp.json` → mantenha a versão e o Python dele, a não ser que
o usuário peça para atualizar ou que a Etapa 5 falhe.

## Etapa 3 — `.mcp.json`

Na raiz do projeto. Sem o arquivo, crie; com ele, acrescente ou corrija só a entrada
`windows-mcp`:

```json
{
  "mcpServers": {
    "windows-mcp": {
      "type": "stdio",
      "command": "uvx",
      "args": [
        "--no-python-downloads",
        "--python",
        "3.13",
        "windows-mcp@0.8.5",
        "serve",
        "--tools",
        "Screenshot,Snapshot,Click,Type,Scroll,Move,Shortcut,Wait,App"
      ],
      "env": {
        "ANONYMIZED_TELEMETRY": "false"
      }
    }
  }
}
```

`--tools` é a lista exata do que fica ativo e tem prioridade sobre `--exclude-tools`. Outras
opções do `serve` (transporte HTTP, porta, autenticação): `uvx ... windows-mcp@X.Y.Z serve --help`.
Para uso local com o Claude Code, fique no `stdio` padrão.

## Etapa 4 — `settings.local.json` e `.gitignore`

`.claude/settings.local.json`: sem o arquivo, crie; com ele, inclua `"windows-mcp"` na lista
`enabledMcpjsonServers` sem tocar no resto:

```json
{
  "enabledMcpjsonServers": ["windows-mcp"]
}
```

É configuração pessoal (aprova o servidor só para este usuário). Se o `.gitignore` não cobre
`.claude/settings.local.json`, acrescente a linha.

Se o projeto exige uma quebra de linha (ex.: LF), grave os arquivos novos nela.

## Etapa 5 — Verificar

O `verificar.py` (na pasta desta skill) lê o `.mcp.json`, sobe o servidor como o Claude Code
sobe, faz o handshake MCP (`initialize` e `tools/list`) e imprime as ferramentas. A primeira
execução pode demorar: o `uvx` baixa o pacote para o cache (`%LOCALAPPDATA%\uv`).

```bash
python -I .claude/skills/agent-windows-mcp/verificar.py .mcp.json
```

Saída esperada (exit 0):

```
server: windows-mcp 4.1.0
tools (9): App, Click, Move, Screenshot, Scroll, Shortcut, Snapshot, Type, Wait
```

- A lista não bate com o `--tools` → corrija o `.mcp.json` e rode de novo.
- O projeto já tinha `WaitFor` no `--tools` → tire (Regra 3).
- `No solution found when resolving tool dependencies` → a versão do pacote não aceita esse
  Python: volte à Etapa 2.
- Timeout sem erro → rode o comando do `.mcp.json` à mão num terminal e leia o erro.

## Etapa 6 — Entregar

Mostre o output da Etapa 5 e diga o que o usuário precisa fazer, porque só ele pode:

1. recarregar o VS Code (`Ctrl+Shift+P` → "Developer: Reload Window") ou abrir uma conversa nova;
2. aceitar o pedido de confiança no projeto ou de aprovação do `windows-mcp`, se aparecer;
3. se nada aparecer e as ferramentas não vierem, rodar `claude` num terminal na pasta do projeto
   e aceitar por lá.

Primeiro teste na sessão nova: `Snapshot` com `use_vision: true` e `use_annotation: false`, que
devolve a imagem da tela e a lista de janelas.

# Dicas de uso

O que os testes mostraram ao operar um app desktop com estas ferramentas:

- **`Snapshot` x `Screenshot`.** O `Snapshot` traz a árvore de acessibilidade: janelas, botões,
  links e campos, com coordenadas clicáveis. Use `use_ui_tree: false` quando só precisar da lista
  de janelas, porque fica bem mais rápido. O `Screenshot` só captura a imagem.
- **App com WebView2 (Wails, Electron, Tauri) expõe o HTML.** Os elementos da página aparecem
  na árvore com os rótulos acessíveis (`aria-label`, texto do link ou do botão). Dá para navegar
  e validar o layout sem pedir ao usuário.
- **`App` com `mode: "launch_executable"` não passa variáveis de ambiente.** App que precisa de
  uma variável para subir fecha logo após abrir. Embrulhar em `cmd /c "set X=...&& start ..."`
  pelo próprio `App` também não funcionou nos testes. Suba o app pelo terminal com a variável
  definida e use o MCP só para ver e operar. Confira se ele ficou de pé
  (`Get-Process <nome>`) antes de procurar a janela.
- **`App` com `mode: "switch"` casa o nome de forma aproximada.** Com o projeto aberto no VS Code,
  "Dev Tools" trouxe a janela "dev-tools - … - Visual Studio Code". Para focar a janela certa,
  clique no ícone dela na barra de tarefas (`Click` nas coordenadas que o `Snapshot` mostra).
- **Esperar a tela mudar.** Use `Wait` com alguns segundos e confira com um `Snapshot` (com
  `use_ui_tree: false` se bastar a lista de janelas). Não use `WaitFor`: ele travava o VS Code.
- **Coordenadas.** Se a imagem vier reduzida, multiplique as coordenadas da imagem pela razão
  entre o tamanho real da tela e o da imagem antes de clicar.
- **Mostrar o print ao usuário.** O resultado do `Snapshot`/`Screenshot` diz onde a imagem foi
  salva (`tool-results\...png`). Mande esse arquivo para o usuário em vez de só descrevê-lo.

# Entrada

$ARGUMENTS
