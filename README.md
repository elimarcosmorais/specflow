# AI Software Engineering

> Um framework enxuto de skills e agents do Claude Code para desenvolver software com
> Spec-Driven Development: da ideia ao código implementado e revisado, uma fase por vez.

## Sobre

Pedir uma feature inteira para um agente de IA de uma vez costuma dar errado: ele preenche as
lacunas com suposições, decide a arquitetura sozinho e entrega código que ninguém revisou contra
um objetivo claro. Este framework divide o trabalho em fases curtas (Spec-Driven Development):
cada fase gera um documento, **um humano revisa e aprova**, e só então a próxima fase começa,
usando esse documento como entrada. A IA nunca constrói em cima de um entendimento que ninguém
conferiu.

São skills e agents do [Claude Code](https://claude.com/claude-code), de propósito poucos e
simples: cada um faz uma coisa, com um template próprio, e o fluxo todo cabe num mapa só.

```
ideia ─► Pre-Refinement ─► PRD ─► (Design) ─► Tech Alignment ─► Tech Spec ─► Task Plan ─► Execução
          o problema       o quê   as telas    o como, alto nível  o como, detalhe  tasks      código revisado

Transversais: ADR (decisões de arquitetura) · Debt (paga os débitos que a Execução anotou)
```

- **[Guia visual](docs/visual-guide.html)**: diagramas clicáveis do fluxo, da Execução, das
  versões e de todas as peças (abra o arquivo no navegador).
- **[Mapa completo do framework](.claude/skills/agent-spec-guide/SKILL.md)**: fases, versões,
  ADRs, débitos, pastas e convenções.

## Pré-requisitos

- [Claude Code](https://claude.com/claude-code), com acesso aos modelos Sonnet e Opus (os agents
  da Execução rodam neles).
- Git: a Execução deixa cada task aprovada em stage, e a revisão compara o `git diff`.
- Chrome ou Edge: só para a `agent-spec-visual-guide`, que confere os diagramas no navegador
  headless.

## Instalação

Copie a pasta `.claude/` (skills, agents e rules) para a raiz do seu projeto:

```bash
git clone https://github.com/elimarcosmorais/ai-software-engineering.git
cp -r ai-software-engineering/.claude <seu-projeto>/
```

O `CLAUDE.md` deste repositório descreve este repositório e não vai junto: gere o do seu projeto
com `/agent-spec-claude-md`.

## Uso

Dentro do Claude Code, na pasta do seu projeto, comece pela ideia:

```
/agent-spec-pre-refinement "quero que o usuário faça login com o Google"
```

A skill gera `docs/specs/features/001-<nome>/v1/pre-refinement.md`. Revise, mude o status para
`aprovado`, e rode o comando que o próprio documento traz na seção "Próximo passo". Perdeu-se no
caminho? Pergunte ao guia:

```
/agent-spec-guide "qual o próximo passo da 001?"
```

Exemplo completo do que o fluxo gera (de um app de teste): a feature
[`001-arquitetura-base`](docs/specs/features/001-arquitetura-base/) (v1 com todas as fases, v2 de
débitos) e as [ADRs](docs/adr/).

## Skills e agents

**Fases** (cada uma usa o documento aprovado da anterior):

| # | Skill | Gera |
|---|-------|------|
| 1 | `/agent-spec-pre-refinement` | `pre-refinement.md`: o problema, de verdade |
| 2 | `/agent-spec-adr` | `docs/adr/NNN-<titulo>.md`: uma decisão de arquitetura (transversal) |
| 3 | `/agent-spec-prd` | `prd.md`: o que construir e para quem |
| — | `/agent-spec-design` | `design.md`: as telas (obrigatório em feature com interface) |
| 4 | `/agent-spec-tech-alignment` | `tech-alignment.md`: a abordagem técnica |
| 5 | `/agent-spec-tech-spec` | `spec.md`: como vai funcionar, em detalhe |
| 6 | `/agent-spec-task-plan` | `tasks.md` + `tasks/T*.md`: os passos de implementação |
| 7 | `/agent-spec-debt` | nova versão com as tasks que pagam os débitos (transversal) |
| 8 | `/agent-spec-run-tasks` | código + `debts.md`: implementa e revisa task a task |

A Execução coordena três agents: `agent-spec-executor` (implementa a task),
`agent-spec-reviewer` (roda os testes e revisa o diff) e `agent-spec-security-reviewer` (só nas
tasks de risco alto).

**Utilitários:**

| Skill | Para quê |
|-------|----------|
| `/agent-spec-guide` | Tirar dúvidas: qual skill rodar, próximo passo de uma feature, como funciona uma fase |
| `/agent-spec-challenge-spec` | Caçar furos no `spec.md` antes do Task Plan |
| `/agent-spec-adr-review` | Conferir se as ADRs são consistentes entre si e com as features |
| `/agent-spec-design-system` | Criar ou atualizar o design system do produto (obrigatório em projeto com interface, antes do primeiro Design) |
| `/agent-spec-semantic-commit` | Commitar tudo com mensagem em Conventional Commits |
| `/agent-spec-readme` | Gerar ou completar o `README.md` |
| `/agent-spec-claude-md` | Gerar ou enxugar o `CLAUDE.md` |
| `/agent-spec-rule-create` | Definir as convenções de um tema (banco, erros HTTP…) em uma rule |
| `/agent-spec-rule-curate` | Decidir se um item vira regra, e auditar as rules |
| `/agent-spec-visual-guide` | Manter o guia visual em dia com as skills, agents e rules |

Argumentos de cada comando: [mapa de roteamento](.claude/skills/agent-spec-guide/SKILL.md#7-mapa-de-roteamento--necessidade--skill).
