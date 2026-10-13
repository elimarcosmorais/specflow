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
2. **Versão fixa do pacote, Python pinado mas provisionável.** Use `windows-mcp@X.Y.Z`, nunca sem
   versão, e `--python <versão>`. O pin do pacote é o que impede o `uvx` de trocar de versão
   sozinho; o `--python` diz qual interpretador usar. **Não use `--no-python-downloads`:** ele não
   acrescenta determinismo (quem o dá é o pin do pacote) e amarra o `.mcp.json` a um Python
   pré-instalado — numa máquina sem ele o servidor não sobe
   (`error: No interpreter found for Python <versão>` … `Python downloads are set to 'never'`).
   Sem o flag, o `uv` baixa o interpretador na área dele, sem admin, e o mesmo `.mcp.json`
   funciona em qualquer máquina.
3. **Só as ferramentas de tela.** O padrão é `Screenshot,Snapshot,Click,Type,Scroll,Move,Shortcut,Wait,App`
   (`WaitFor` **não** entra: ele travava o VS Code; para esperar, use `Wait` e confira com um
   `Snapshot`). `PowerShell`, `Registry`, `FileSystem`, `Process`, `Clipboard`,
   `Scrape` e as demais ficam de fora: dão ao agente um segundo caminho, sem as travas de
   permissão do Claude Code, para rodar comando e mexer no sistema. Liberar uma delas só com
   pedido explícito do usuário (`$ARGUMENTS` ou na conversa), dizendo o risco. Numa `0.8.8` ou
   mais nova o pacote acrescenta o `ControlStatus` mesmo com `--tools`: é só leitura e não fere
   esta regra. Na `0.8.5` ele não existe, porque o subsistema de controle que ele relata também
   não existe (ver Etapa 2).
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
- **Python não bloqueia.** A versão que a Etapa 2 escolher pede Python `>=3.12`, e sem o
  `--no-python-downloads` (Regra 2) o `uv` provisiona o interpretador sozinho. Isso não cai na
  Regra 1: não é programa instalado na máquina, é um interpretador descartável no cache do `uv`
  (`%LOCALAPPDATA%\uv`), sem admin e sem entrar no `PATH`. Só avise que a primeira execução baixa
  uns 20 MB. Prefira um Python que `uv python list --only-installed` já mostre, para nem baixar;
  entre os instalados, prefira os gerenciados pelo `uv` aos que vêm embutidos em outro programa
  (ex.: o Python do laragon sai se o laragon sair).

## Etapa 2 — Escolher a versão do pacote

**Use a `0.8.5`. Não use a `0.8.8`**, e não pegue "a mais nova" sem fazer a checagem abaixo.

| Versão | Python exigido | Usar? |
|---|---|---|
| `0.8.5` | `>=3.12` | **sim**, é o padrão desta skill |
| `0.8.8` | `>=3.14` | **não**: congela a máquina e não tem como desligar |
| qualquer outra (`0.8.6`, `0.8.7`, `0.8.9`+) | ver no PyPI | só se passar na checagem do `--no-control` |

A `0.8.8` instala um hook global de mouse e teclado
(`SetWindowsHookExW(14/13, …, 0)`, em `desktop/control_hooks.py`) que **engole o input físico do
Windows inteiro** enquanto o agente opera a tela, renovado de 1 em 1 segundo durante toda a
operação. O sintoma é "o VS Code travou", mas ele só parou de receber mouse e teclado. Não é
problema de um cliente: o hook é do processo do servidor e vale para a área de trabalho toda, então
o VS Code, o app do Claude Code e o Claude Desktop travam igual. A `0.8.5` não tem nada desse
código — nenhum arquivo `control*.py`.

O projeto `Windows-MCP` resolveu isso, mas **depois** da `0.8.8`: o `main` documenta
`--no-control`, `WINDOWS_MCP_CONTROL=false` e `[control] enabled = false` ("to run background tasks
while using the computer yourself, disable it"), e nada disso existe na `0.8.8`. Ou seja, a `0.8.8`
é a única versão com o bloqueio e sem a chave.

**Checagem antes de subir para uma `0.8.9` ou mais nova.** Pergunte ao próprio CLI, que é a fonte
autoritativa: o README documenta o `main`, não a versão lançada, e por isso descreve o
`--no-control` como se já existisse.

```bash
uvx --python 3.14 windows-mcp@<versão> serve --help | grep -- --no-control
```

Imprimiu a opção → pode usar essa versão, acrescentando `--no-control` depois do `serve`
(mantenha o `--tools`). Não imprimiu nada → a chave não está comprovada nessa versão: **fique na
`0.8.5`**. Na `0.8.8` esse comando não imprime nada.

Versão mais nova e Python exigido: `https://pypi.org/pypi/windows-mcp/json` (`info.version` e
`info.requires_python`).

O projeto já tem `windows-mcp` no `.mcp.json` → mantenha a versão dele, **exceto se for a `0.8.8`
ou mais nova sem `--no-control`**: nesse caso avise o usuário do congelamento e volte para a
`0.8.5`.

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
        "--python",
        "3.14",
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

Saída esperada (exit 0, com a `0.8.5`):

```
server: windows-mcp 4.1.0
tools (9): App, Click, Move, Screenshot, Scroll, Shortcut, Snapshot, Type, Wait
```

Nove ferramentas, sem `ControlStatus`: essa ausência é a prova de que o subsistema de controle —
o do hook que congela a tela (Etapa 2) — não está no pacote que subiu.

- Apareceu `ControlStatus` (10 ferramentas) → você está numa `0.8.8` ou mais nova. Confira a
  Etapa 2: ou a versão passa na checagem do `--no-control` e você o acrescenta depois do `serve`,
  ou volte para a `0.8.5`.
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
    conseguiu operar (ferramentas ausentes) → pergunte no chat que precisa conferir com o
    Windows-MCP.
  ```

  Só se o `.mcp.json` estiver numa `0.8.8` ou mais nova, acrescente "mouse em uso depois de
  esperar com `Wait`, " à lista de motivos: a disputa pelo mouse e teclado só existe a partir
  dela. Na `0.8.5`, que é o padrão, não entra.

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
- **Tela congelada: `Ctrl+Alt+Shift+Backspace`.** Só acontece numa `0.8.8` ou mais nova (na
  `0.8.5`, o padrão, não existe). Se o mouse e o teclado pararem de responder enquanto o agente
  opera, a máquina não travou: o hook do pacote está engolindo o input (Etapa 2). Esse atalho
  devolve o controle na hora. Mover bastante o mouse também devolve, e o controle volta ao
  agente depois de ~10 s sem uso. **Diga isso ao usuário** ao configurar uma versão dessas.
- **Disputa pelo mouse e teclado (`0.8.8`+).** Enquanto o agente opera, o pacote assume o mouse e
  o teclado físicos e mostra um aviso na tela. Se o usuário mexer, o agente perde o controle e só
  o retoma depois de uma espera sem uso. Ação recusada nesse meio-tempo → consulte o
  `ControlStatus` (dono atual e segundos que faltam), espere com `Wait` e tente de novo.
  **Mouse em uso não é motivo para pular a conferência**, nem para passá-la ao usuário: depois de
  umas 3 esperas sem conseguir, pergunte no chat ("Preciso conferir a tela com o Windows-MCP.
  Solte o mouse e o teclado e me avise.") e, com a resposta, confira. O próprio projeto avisa que
  essa retomada é "best-effort, not a security boundary", e que o Windows pode remover um hook de
  baixo nível sem avisar o processo.
- **`Snapshot` no VS Code é a travessia mais pesada.** Ele é Electron e expõe a árvore de
  acessibilidade inteira. Numa `0.8.8`+ isso é o pior caso: o congelamento dura toda a travessia.
  Numa `0.8.8`+ dá para excluí-lo com `WINDOWS_MCP_EXCLUDE_PROCESSES=Code.exe` no `env` (a
  variável não existe na `0.8.5`, e lá não faz falta).
- **Coordenadas.** Se a imagem vier reduzida, multiplique as coordenadas da imagem pela razão
  entre o tamanho real da tela e o da imagem antes de clicar.
- **Mostrar o print ao usuário.** O resultado do `Snapshot`/`Screenshot` diz onde a imagem foi
  salva (`tool-results\...png`). Mande esse arquivo para o usuário em vez de só descrevê-lo.

# Entrada

$ARGUMENTS
