# SpecFlow

> Um framework enxuto de skills e agents do Claude Code para desenvolver software com specs,
> onde o processo cresce com o tamanho da mudança: nada de plano para mudar uma linha de CSS.

## Sobre

Pedir uma feature inteira para um agente de IA de uma vez costuma dar errado: ele preenche as
lacunas com suposições e entrega código que ninguém revisou contra um objetivo claro. O outro
extremo também falha: um fluxo com uma fase e um documento para cada etapa faz você escrever mais
documentação do que código.

Este framework escolhe o processo pelo tamanho da mudança:

| Tamanho | Exemplo | O que acontece |
|---|---|---|
| **Pequena** | CSS, bug, texto, ajuste de validação | implementa direto, roda os testes, mostra o resultado. Nenhum plano |
| **Média** | feature contida, endpoint + tela | um plano curto → você aprova → implementação passo a passo, com revisão |
| **Grande** | fundação, decisão de arquitetura | o mesmo, com ADR para as decisões que valem para o projeto todo |

```
pedido ─► /spec-change ─┬─ pequena ─► código testado
                        └─ média/grande ─► plano ─► você aprova ─► /spec-build ─► código revisado
```

Na revisão, só o que um sênior recusaria num PR volta para correção. O resto vira sugestão: o
executor a aplica numa rodada revisada à parte, e o que ele recusar fica no relatório. Mapa completo: [spec-guide](.claude/skills/spec-guide/SKILL.md).

O plano é temporário: o build o apaga depois de executar. O que fica é a documentação viva em
`docs/`, que o agente mantém sozinho, junto com o código:

| Pasta | Guarda |
|---|---|
| `docs/specs/` | como cada módulo funciona hoje; o agente cria, edita e subdivide sem pedir permissão |
| `docs/architecture.md` | a arquitetura base |
| `docs/adr/` | por que cada decisão transversal foi tomada (com a sua aprovação); índice em `docs/adr/README.md` |
| `docs/design-system.md` + `docs/design/` | a identidade visual, obrigatória em projeto com interface |
| `docs/plans/` | o plano da mudança em andamento |

As regras de cada pasta estão na rule [`docs.md`](.claude/rules/docs.md).

## Pré-requisitos

- [Claude Code](https://claude.com/claude-code), com acesso aos modelos Sonnet e Opus (os agents
  do build rodam neles).
- Git: o build deixa cada passo aprovado em stage, e a revisão compara o `git diff`.

## Instalação

Copie a pasta `.claude/` (skills, agents e rules) para a raiz do seu projeto:

```bash
git clone https://github.com/elimarcosmorais/specflow.git
cp -r specflow/.claude <seu-projeto>/
```

O `CLAUDE.md` deste repositório descreve este repositório e não vai junto: gere o do seu projeto
com `/spec-claude`.

Para atualizar depois, rode `/spec-upgrade` no seu projeto: ele baixa a versão nova deste
repositório, atualiza o que é do framework e pergunta antes de sobrescrever um arquivo que o
projeto editou. Skills, agents e rules próprios do projeto ficam intocados.

## Uso

Dentro do Claude Code, na pasta do seu projeto, descreva a mudança:

```
/spec-change "o botão de salvar está azul, deveria ser verde"
/spec-change "quero que o usuário faça login com o Google"
```

O primeiro é pequeno: a skill muda o CSS e mostra o resultado. O segundo é médio ou grande: ela
pergunta o que o código não responde e salva `docs/plans/2026-10-08-login-com-google.md`.
Revise, mude o status para `aprovado` e rode:

```
/spec-build docs/plans/2026-10-08-login-com-google.md
```

O build implementa, atualiza as specs dos módulos tocados e apaga o plano. Nada é commitado
sozinho: no fim, `/spec-semantic-commit`. Para lançar uma versão com o que já foi commitado,
`/spec-release`.

## Skills e agents

**Fluxo:**

| Skill | Faz |
|---|---|
| `/spec-change` | Porta de entrada: classifica o pedido; implementa o pequeno, escreve o plano do médio e do grande |
| `/spec-build` | Implementa um plano aprovado, passo a passo, e o apaga no fim |

O build coordena três agents: `spec-executor` (implementa um passo),
`spec-reviewer` (roda os testes e revisa o diff) e `spec-security-reviewer` (só nos
passos de risco alto).

**Utilitários:**

| Skill | Para quê |
|---|---|
| `/spec-guide` | Tirar dúvidas: qual skill rodar, próximo passo de um plano, onde fica cada documento |
| `/spec-adr` | Registrar uma decisão de arquitetura que vale para o projeto todo; `/spec-adr índice` organiza o índice das ADRs |
| `/spec-design-system` | Criar ou atualizar o design system do produto (obrigatório com interface) |
| `/spec-semantic-commit` | Commitar tudo com mensagem em Conventional Commits |
| `/spec-release` | Lançar uma versão: calcula o número pelos commits, atualiza o `CHANGELOG.md`, cria a tag e o GitHub Release |
| `/spec-readme` | Gerar ou completar o `README.md` |
| `/spec-claude` | Gerar ou enxugar o `CLAUDE.md` |
| `/spec-rule-create` | Definir as convenções de um tema (banco, erros HTTP…) em uma rule |
| `/spec-rule-curate` | Decidir se um item vira regra, e auditar as rules |
| `/spec-windows-mcp` | Configurar o Windows-MCP para o agente ver e operar a tela (só Windows) |
| `/spec-playwright-mcp` | Configurar o Playwright MCP para o agente ver e operar páginas web |
| `/spec-upgrade` | Trazer as skills, os agents e as rules da versão nova do SpecFlow, preservando o que o projeto editou ou criou |

Argumentos de cada comando: [mapa de roteamento](.claude/skills/spec-guide/SKILL.md#5-mapa-de-roteamento).
