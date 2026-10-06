# AI Software Engineering

## Sobre o projeto

Um framework enxuto de **skills e agents do Claude Code** para desenvolver software com
Spec-Driven Development (SDD): cada fase gera um documento, um humano revisa e aprova, e só então
a próxima fase começa, usando esse documento como entrada. Visão geral, instalação e uso estão no
[README.md](README.md).

Este repositório é o próprio framework: o que se edita aqui são as skills, os agents e as rules.
Os documentos em `docs/specs/` e `docs/adr/` são o exemplo do que o fluxo gera.

## Fases

| # | Fase           | Skill                       |
|---|----------------|-----------------------------|
| 1 | Pre-Refinement | `agent-spec-pre-refinement` |
| 2 | ADR            | `agent-spec-adr`            |
| 3 | PRD            | `agent-spec-prd`            |
| — | Design (obrigatório em feature com telas) | `agent-spec-design` |
| 4 | Tech Alignment | `agent-spec-tech-alignment` |
| 5 | Tech Spec      | `agent-spec-tech-spec`      |
| 6 | Task Plan      | `agent-spec-task-plan`      |
| 7 | Debt           | `agent-spec-debt`           |
| 8 | Execução       | `agent-spec-run-tasks`      |

Agents da Execução: `agent-spec-executor`, `agent-spec-reviewer` e `agent-spec-security-reviewer`
(só nas tasks de risco alto).

Além das fases, há skills **utilitárias**, que não geram documento de fase:

| Utilitário                             | Skill                        |
|----------------------------------------|------------------------------|
| Commit semântico                       | `agent-spec-semantic-commit` |
| Challenge da spec                      | `agent-spec-challenge-spec`  |
| Guia do framework                      | `agent-spec-guide`           |
| Revisão das ADRs                       | `agent-spec-adr-review`      |
| Design system (projeto com interface)  | `agent-spec-design-system`   |
| README do projeto                      | `agent-spec-readme`          |
| CLAUDE.md do projeto                   | `agent-spec-claude-md`       |
| Rules de um tema                       | `agent-spec-rule-create`     |
| Curadoria das rules                    | `agent-spec-rule-curate`     |
| Guia visual (`docs/visual-guide.html`) | `agent-spec-visual-guide`    |

## Ao alterar o framework

- **Simplicidade acima de completude.** Cada skill faz uma coisa, com um template próprio. Antes
  de acrescentar uma fase, skill ou seção, confira se ela é necessária para o fluxo.
- **Uma mudança por vez.** Uma skill ou agent por vez, terminada antes de passar à próxima.
- **Explicar as decisões.** Ao criar ou alterar uma skill, dizer o que mudou e por quê.
- **Manter tudo em sincronia** ao criar, alterar ou remover uma skill, agent ou rule:
  - as tabelas deste arquivo e as do [README.md](README.md);
  - o mapa do framework (`agent-spec-guide`);
  - o guia visual (rodar `/agent-spec-visual-guide`).
- **Idioma:** conversar com o usuário e escrever skills e documentos em português.

## Estrutura de arquivos

O fluxo completo, a estrutura de diretórios e as convenções estão na skill `agent-spec-guide`
(o mapa do framework). Siga-a ao criar skills e documentos, e atualize-a quando o fluxo mudar.

- Skills: `.claude/skills/<nome-da-skill>/SKILL.md`
- Agents: `.claude/agents/<nome-do-agent>.md`
- Rules: `.claude/rules/<nome>.md` (ex.: `behavior.md`, as regras de comportamento)
