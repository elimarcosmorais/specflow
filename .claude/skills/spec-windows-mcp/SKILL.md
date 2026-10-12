---
name: spec-windows-mcp
description: Utilitário (transversal, só Windows). Configura o Windows-MCP no projeto para o agente ver e operar a tela do Windows (screenshot, árvore de acessibilidade, clique, teclado, janelas) - confere o uv e o Python, fixa uma versão do pacote compatível, registra o servidor no .mcp.json só com as ferramentas de tela, cria ou completa o .claude/settings.local.json com enabledMcpjsonServers, põe esse arquivo no .gitignore, prova com um handshake MCP direto que o servidor sobe e lista as ferramentas, e registra o servidor na seção Teste visual do CLAUDE.md do projeto (sem precisar voltar ao spec-claude). Também traz as dicas de uso aprendidas nos testes. Não instala nada sem confirmação e não libera PowerShell, Registry, FileSystem nem Process.
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
   pedido explícito do usuário (`$ARGUMENTS` ou na conversa), dizendo o risco. A partir da
   `0.8.8` o pacote acrescenta o `ControlStatus` mesmo com `--tools`: é só leitura e não fere
   esta regra.
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
| 3.14 | a mais nova | `0.8.8` testada; confirme na Etapa 5; se falhar, volte para a `0.8.5` com 3.13 |

Versão mais nova e Python exigido: `https://pypi.org/pypi/windows-mcp/json` (`info.version` e
`info.requires_python`).

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
        "3.14",
        "windows-mcp@0.8.8",
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

O `.mcp.json` vai para o git. Se o `.gitignore` o ignora (o starter kit do Laravel traz
`/.mcp.json`), avise o usuário e tire a linha só com o ok dele. Confira com
`git check-ignore -v .mcp.json`.

Se o projeto exige uma quebra de linha (ex.: LF), grave os arquivos novos nela.

## Etapa 5 — Verificar

O `verificar.py` (na pasta desta skill) lê o `.mcp.json`, sobe o servidor como o Claude Code
sobe, faz o handshake MCP (`initialize` e `tools/list`) e imprime as ferramentas. A primeira
execução pode demorar: o `uvx` baixa o pacote para o cache (`%LOCALAPPDATA%\uv`).

```bash
python -I .claude/skills/spec-windows-mcp/verificar.py .mcp.json
```

Saída esperada (exit 0, com a `0.8.8`):

```
server: windows-mcp 4.1.0
tools (10): App, Click, ControlStatus, Move, Screenshot, Scroll, Shortcut, Snapshot, Type, Wait
```

A `0.8.5` lista as mesmas, sem o `ControlStatus` (9 ferramentas).

- `ControlStatus` a mais é esperado a partir da `0.8.8`: o pacote o mantém ligado mesmo com
  `--tools` e não há opção para tirá-lo. Ele só lê quem controla a área de trabalho (ver
  "Disputa pelo mouse e teclado" nas Dicas de uso). Avise o usuário.
- Fora o `ControlStatus`, a lista não bate com o `--tools` → corrija o `.mcp.json` e rode de novo.
- O projeto já tinha `WaitFor` no `--tools` → tire (Regra 3).
- `No solution found when resolving tool dependencies` → a versão do pacote não aceita esse
  Python: volte à Etapa 2.
- Timeout sem erro → rode o comando do `.mcp.json` à mão num terminal e leia o erro.

## Etapa 6 — Registrar no `CLAUDE.md`

Esta skill fecha o próprio trabalho: não mande o usuário rodar `/spec-claude` por causa do
MCP.

- Sem `CLAUDE.md` no projeto → não crie. O `/spec-claude`, quando rodar, acha o servidor no
  `.mcp.json`.
- O `CLAUDE.md` já cita o Windows-MCP, em qualquer seção → não mova a regra para a seção Teste
  visual: reorganizar o arquivo é trabalho do `/spec-claude`. Só complete o bullet existente
  com a frase "Não pule a conferência…" do bloco abaixo, se ela faltar.
- Senão → acrescente este bullet à seção Teste visual, no idioma do arquivo. Sem a seção, crie-a
  na posição do [template](../spec-claude/template.md) (depois de UI, antes de "Pronto
  significa"):

  ```
  - Ver a tela (conferir layout, reproduzir bug) → Windows-MCP (`windows-mcp` no `.mcp.json`):
    abra o app e confira com `Snapshot`. Não pule a conferência nem a passe ao usuário: não
    conseguiu operar (mouse em uso depois de esperar com `Wait`, ferramentas ausentes) →
    pergunte no chat que precisa conferir com o Windows-MCP.
  ```

  O trecho "mouse em uso depois de esperar com `Wait`, " só entra com a `0.8.8` ou mais nova
  no `.mcp.json` (a disputa pelo mouse e teclado começou nela); numa versão anterior, tire-o.

  Mostre o diff e grave só depois do OK. Mexa só nessa linha, mesmo que o arquivo tenha outras
  alterações do usuário.

## Etapa 7 — Entregar

Mostre o output da Etapa 5 e o que mudou no `CLAUDE.md` (ou que o projeto não tem um). Diga o
que o usuário precisa fazer, porque só ele pode:

1. recarregar o VS Code (`Ctrl+Shift+P` → "Developer: Reload Window") ou abrir uma conversa nova;
2. aceitar o pedido de confiança no projeto ou de aprovação do `windows-mcp`, se aparecer;
3. se nada aparecer e as ferramentas não vierem, rodar `claude` num terminal na pasta do projeto
   e aceitar por lá.

Diga também como testar na sessão nova, em termos do usuário (ele não chama a ferramenta; pede
ao agente): basta escrever algo como "tire um print da minha tela e liste as janelas abertas".
O agente chama o `Snapshot` com `use_vision: true` (inclui a imagem) e `use_annotation: false`
(sem os retângulos coloridos que o padrão desenha em volta dos elementos). Se o modo de
permissão da sessão pedir aprovação, o Claude Code pergunta antes de cada ferramenta nova.

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
- **Disputa pelo mouse e teclado (a partir da `0.8.8`).** Enquanto o agente opera, o pacote
  assume o mouse e o teclado físicos e mostra um aviso na tela. Se o usuário mexer, o agente
  perde o controle e só o retoma depois de uma espera sem uso. Ação recusada nesse meio-tempo →
  consulte o `ControlStatus` (dono atual e segundos que faltam), espere com `Wait` e tente de
  novo. **Mouse em uso não é motivo para pular a conferência**, nem para passá-la ao usuário:
  depois de umas 3 esperas sem conseguir, pergunte no chat ("Preciso conferir a tela com o
  Windows-MCP. Solte o mouse e o teclado e me avise.") e, com a resposta, confira.
- **Coordenadas.** Se a imagem vier reduzida, multiplique as coordenadas da imagem pela razão
  entre o tamanho real da tela e o da imagem antes de clicar.
- **Mostrar o print ao usuário.** O resultado do `Snapshot`/`Screenshot` diz onde a imagem foi
  salva (`tool-results\...png`). Mande esse arquivo para o usuário em vez de só descrevê-lo.

# Entrada

$ARGUMENTS
