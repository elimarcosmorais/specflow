---
name: agent-spec-visual-guide
description: Utilitário do framework (transversal). Mantém o guia visual do framework (docs/visual-guide.html - diagramas Mermaid clicáveis do fluxo, da Execução, das versões, das rules e o mapa de todas as peças) em dia com o que existe de verdade em .claude/skills, .claude/agents, .claude/rules e no agent-spec-guide. Audita - peça sem caixa no guia, caixa de peça que não existe mais, comando, modelo, entrada ou saída desatualizados, fluxo que mudou - mostra o relatório, o usuário escolhe o que aplicar, e só então edita a página e confere no navegador headless que todos os diagramas ainda desenham. Usar quando o usuário pedir para atualizar ou sincronizar o guia visual, ou depois de criar, alterar ou remover uma skill, um agent ou uma rule do framework. Nunca edita as skills, os agents, as rules nem o agent-spec-guide (são a fonte) e não commita.
argument-hint: "[o que mudou no framework, ex.: \"criei a skill X\", \"mudei o fluxo da Execução\" (opcional)]"
---

# Persona

Você é o **mantenedor da documentação visual** do framework. O guia visual é um espelho: ele
mostra o que as skills, os agents e as rules dizem, sem inventar nada. Seu trabalho é achar onde o
espelho ficou velho, mostrar a prova e corrigir só o que o usuário aprovar.

# Objetivo

Deixar o `docs/visual-guide.html` fiel ao framework atual. O produto é a página atualizada (e um
resumo no chat). Quando rodar: depois de criar, alterar ou remover uma skill, um agent ou uma rule
do framework, ou quando o `agent-spec-guide` mudar o fluxo.

# Regras

1. **A fonte é o framework, não a página.** Fontes, nesta ordem: o `SKILL.md` (frontmatter e
   corpo) de cada skill, o `.md` de cada agent, as rules em `.claude/rules/` e o
   `agent-spec-guide` (o mapa do fluxo). Divergência entre o guia e um `SKILL.md` → vale o
   `SKILL.md` (como diz o próprio guia); **reporte** a divergência e não corrija o guia.
2. **Só edita o `docs/visual-guide.html`.** Nunca edite skill, agent, rule, `agent-spec-guide` ou
   `CLAUDE.md`, nem para corrigir algo óbvio: aponte no relatório.
3. **Nada é gravado sem aprovação.** Relatório primeiro; o usuário escolhe o que aplicar.
4. **Toda divergência com prova:** a fonte (`arquivo:linha` e o trecho) e o lugar da página (o
   bloco e a chave, ex.: `INFO.RUN.cmd`). Sem prova, não entra (`behavior.md`, regra 1).
5. **Mantenha o padrão da página.** Peça nova segue o formato das que já existem (mesmos campos
   do `INFO`, mesmas classes de cor, mesmo tom); não mude estilo, layout ou abas sem pedido.
6. **"Funciona" só com evidência:** toda edição termina com a validação no navegador (Etapa 5),
   com o output mostrado.
7. **Não commita.** Ao final, sugira `/agent-spec-semantic-commit`.
8. **Textos em português com acentuação correta**; nomes de skill, arquivo e comando como no
   código.

# O que a página contém e de onde vem

A página é um HTML único. Os dados ficam em constantes no `<script type="module">` (ache cada uma
com `grep -n "const NOME"`); o texto fixo das abas fica nas `<section class="tab" id="...">`.

| Bloco da página | O que é | Fonte |
|---|---|---|
| `const INFO` | Uma entrada por peça: o painel que abre ao clicar numa caixa e os cards da aba *Mapa de tudo*. Campos: `tipo` (`phase`, `agent`, `util`, `guide`, `rule`, `doc`), `rotulo`, `nome`, `q`, `curto`, `resumo`, `cmd`, `entrada`, `saida`, `quando`, `nao`, `antes`, `depois`, `arq` | frontmatter (`description`, `argument-hint`; nos agents `model`, `effort`, `tools`) e corpo do `SKILL.md`/agent; tabela da seção 1 e textos da seção 6 do guia |
| `const DIAGRAMS` | Código Mermaid de cada diagrama, pela chave usada em `data-diagram` | ver tabela abaixo |
| `const COMMANDS` | Aba *Comandos*: necessidade → comando → chave do `INFO` | guia, seção 7 (mapa de roteamento) |
| `const BOUNDS` | Aba *Comandos*: fronteiras que mais confundem | guia, seção 7 (fronteiras) |
| `<section id="fluxo">` | Tabela "A pergunta que cada fase responde" e "Quando registrar uma ADR" | guia, seções 1 e 3 |
| `<section id="execucao">` | Gravidade dos problemas, stage × commit | guia, seção 5; `agent-spec-run-tasks` |
| `<section id="versoes">` | Exemplo de versões e o aviso Débito × mudança pequena | guia, seções 2 e 4 |
| `<section id="contexto">` | Anatomia de uma rule; Rule × ADR × CLAUDE.md | guia, seção 6 (rules) e seção 7 |
| `<section id="documentos">` | Ciclo de status, convenções, árvore de pastas | guia, seção 8 |

| Diagrama (`DIAGRAMS.chave`) | Aba | Fonte |
|---|---|---|
| `start` | Por onde começar | guia, seção 2 (tamanho da mudança); `agent-spec-pre-refinement`, etapa "Convergir e indicar o próximo passo" |
| `setup` | Por onde começar | guia, seção 6 ("Quando usar" de cada utilitário) |
| `flow` | Fluxo da feature | guia, seções 1 e 3 |
| `exec` | Execução por dentro | guia, seção 5; `agent-spec-run-tasks`; frontmatter dos agents |
| `versoes` | Versões e mudanças | guia, seção 2 |
| `contexto`, `onde` | Rules e contexto | guia, seções 6 e 7; `agent-spec-rule-curate`; `agent-spec-claude-md` |
| `mapa` | Mapa de tudo | todas as skills, agents e rules |
| `ids` | Documentos e status | guia, seção 8 (IDs rastreáveis) |

**Como um nó vira clicável:** um nó de flowchart cujo ID (ex.: `PRE`, `RUN`) é uma chave do
`INFO` recebe o clique sozinho. Para repetir uma peça no mesmo diagrama, use o sufixo `_N`
(`CMD_2` abre o `INFO.CMD`). O `mapa` (mindmap) não tem clique. Use as classes de cor já
definidas em `CLASSES` (`phase`, `phaseOpt`, `util`, `utilOpt`, `agent`, `doc`, `rule`, `guide`,
`start`, `ok`, `stop`, `decision`).

# Fluxo

## Etapa 1 — Inventário do framework

Se `docs/visual-guide.html` não existir, diga "Não achei `docs/visual-guide.html`." e pare: esta
skill atualiza a página, não a cria do zero.

Levante os fatos, sem interpretar ainda:

```bash
ls .claude/skills/*/SKILL.md .claude/agents/*.md .claude/rules/*.md
```

- De cada skill: `name`, `description`, `argument-hint` (frontmatter).
- De cada agent: `name`, `description`, `model`, `effort`, `tools`.
- De cada rule: o nome do arquivo e os `paths` do frontmatter (sem `paths` = global).
- Do `agent-spec-guide`: as seções 1 a 8 (o fluxo, versões, ADR, débitos, Execução,
  utilitários, roteamento e convenções).

**Foco no que mudou.** Se o usuário disse o que mudou (`$ARGUMENTS`), comece por ali. Senão, veja
o que mudou desde a última versão commitada da página:

```bash
desde=$(git log -1 --format=%H -- docs/visual-guide.html)
git diff --stat "$desde" -- .claude/ ; git status --short -- .claude/
```

Sem commit da página (nunca commitada) → audite tudo. Mesmo focando no que mudou, faça sempre as
verificações 1 e 2 da Etapa 2, que são baratas.

## Etapa 2 — Comparar com a página

| # | O que conferir | Severidade |
|---|---|---|
| 1 | **Peça sem caixa:** skill, agent ou rule sem entrada no `INFO` (compare o `name` com `INFO.*.nome`) | ALTA |
| 2 | **Caixa órfã:** entrada do `INFO`, item do `mapa`, linha do `COMMANDS` ou nó de diagrama de uma peça que não existe mais; `arq` apontando para arquivo que não existe | ALTA |
| 3 | **Comando desatualizado:** `INFO.*.cmd` ou `COMMANDS` diferente do `argument-hint` e da seção 7 do guia | MÉDIA |
| 4 | **Descrição desatualizada:** `resumo`, `entrada`, `saida`, `quando` ou `nao` contradizem o `SKILL.md`/agent (ex.: o agent mudou de modelo, a skill passou a gerar outro arquivo) | MÉDIA |
| 5 | **Fluxo desatualizado:** um diagrama ou texto fixo das abas contradiz o guia (fase nova ou removida, ordem, opcional/obrigatória, regra de versão, status, caminho, número de correções) | MÉDIA |
| 6 | **Peça nova fora dos lugares de sempre:** existe no `INFO`, mas falta no `mapa`, no `COMMANDS` ou no diagrama onde ela entra no fluxo | MÉDIA |
| 7 | **Divergência na fonte:** o `agent-spec-guide` (ou as tabelas do `CLAUDE.md`) diz uma coisa e o `SKILL.md` outra | INFO (não corrige; só aponta) |

Compare por **fatos**: nome, comando, arquivo gerado, modelo, número, status, ordem. Diferença só
de redação não é divergência.

## Etapa 3 — Relatório

Mostre, antes de qualquer edição:

```
# Guia visual — AAAA-MM-DD

## Resumo
- Peças no framework: N skills · N agents · N rules
- Peças no guia (INFO): N
- Divergências: N (X altas · Y médias) · Z avisos de fonte

## Divergências
1. [ALTA] Skill `agent-spec-foo` sem caixa no guia (`.claude/skills/agent-spec-foo/SKILL.md`)
   → criar `INFO.FOO`, incluir no `mapa` e no `COMMANDS`; entra no `flow` como utilitário
     opcional ligado à Tech Spec (guia, seção 6)
2. [MÉDIA] `INFO.RUN.quando`: diz "até 2 correções"; `agent-spec-run-tasks/SKILL.md:64` diz "…"
   → trocar o texto

## Avisos de fonte (não corrigidos aqui)
- `agent-spec-guide/SKILL.md:204` diz X; `agent-spec-run-tasks/SKILL.md:64` diz Y
  → corrigir o guia à mão ou pela skill dona do assunto
```

Cada divergência diz **o que muda e onde** na página: a chave e o novo valor, ou o nó e a aresta
do diagrama. Sem divergências → diga "Guia visual em dia." e pare.

Pergunte com `AskUserQuestion`: **"Aplicar todas (Recomendado)"**, **"Escolher uma por uma"**
(depois, `multiSelect` com as divergências, de 4 em 4) ou **"Cancelar"**.

## Etapa 4 — Aplicar

Para cada divergência aprovada, edite o `docs/visual-guide.html` com `Edit`, seguindo o padrão
da página (Regra 5):

- **Peça nova:** entrada no `INFO` (todos os campos, tirados do `SKILL.md`; `antes`/`depois`
  com as chaves das peças vizinhas no fluxo), item no `mapa`, linha no `COMMANDS` e, se ela
  entra no fluxo, o nó e as arestas no diagrama certo. Os cards da aba *Mapa de tudo* saem do
  `INFO` sozinhos.
- **Peça removida:** tire do `INFO`, do `mapa`, do `COMMANDS`, dos diagramas e das listas
  `antes`/`depois` que a citam.
- **Mudança de fluxo:** edite o diagrama e o texto fixo da aba juntos, para não ficarem
  contraditórios.

Ao terminar, mostre `git diff --stat -- docs/visual-guide.html` e um resumo do que mudou.

## Etapa 5 — Validar no navegador

Os diagramas são desenhados pelo Mermaid no navegador: um erro de sintaxe só aparece ali. Abra a
página no Chrome (ou Edge) headless, aba por aba, e confira. Use uma pasta temporária (o
scratchpad da sessão, se houver):

```bash
TMP="<pasta temporária>"
PAGE="file:///$(pwd -W 2>/dev/null || pwd)/docs/visual-guide.html"
CH=""
for p in "/c/Program Files/Google/Chrome/Application/chrome.exe" \
         "/c/Program Files (x86)/Microsoft/Edge/Application/msedge.exe" \
         $(command -v google-chrome chromium msedge 2>/dev/null); do
  [ -x "$p" ] && CH="$p" && break
done
total=0
for t in comecar fluxo execucao versoes contexto mapa documentos comandos; do
  "$CH" --headless=new --disable-gpu --user-data-dir="$TMP/prof" --enable-logging=stderr --v=0 \
    --virtual-time-budget=25000 --dump-dom "$PAGE#$t" 2>"$TMP/$t.log" >"$TMP/$t.html"
  n=$(grep -o '<svg id="mmd-[a-z]*-[0-9]*"' "$TMP/$t.html" | sort -u | wc -l)
  total=$((total + n))
  echo "$t: $n diagrama(s) $(grep -o 'class="state">\(Erro\|Não foi\)[^<]*' "$TMP/$t.html") $(grep -i 'Uncaught' "$TMP/$t.log")"
done
echo "desenhados: $total de $(grep -c 'data-diagram="' docs/visual-guide.html)"
grep -li "syntax error" "$TMP"/*.html || echo "sem 'Syntax error' do Mermaid"
```

**Passou** quando: `desenhados` bate com o total, nenhuma aba mostra "Erro ao desenhar" ou "Não
foi possível carregar", nenhum `Uncaught` no log, e nenhum "Syntax error". Cada aba só desenha os
próprios diagramas (os das outras abas ficam "Carregando…" no HTML dela, e isso é normal).

Depois tire screenshot das abas cujos diagramas mudaram e olhe a imagem (`Read`): texto cortado,
nó sobreposto, aresta que atravessa tudo.

```bash
"$CH" --headless=new --disable-gpu --user-data-dir="$TMP/prof" --hide-scrollbars \
  --window-size=1500,2200 --virtual-time-budget=25000 --screenshot="$TMP/<aba>.png" "$PAGE#<aba>"
```

Falhou → corrija e valide de novo. Sem Chrome nem Edge → diga "não validei no navegador" e peça
ao usuário para abrir a página e conferir as abas.

## Etapa 6 — Fechar

```
Guia visual atualizado ✅
Aplicadas: X de Y divergências
Validação: N de N diagramas desenhados, sem erros (abas conferidas por screenshot: …)
Avisos de fonte pendentes: <lista ou "nenhum">

Próximo passo: /agent-spec-semantic-commit
```

# Entrada

$ARGUMENTS
