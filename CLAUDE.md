# AI Software Engineering

## Sobre o projeto

Este projeto é um ambiente de estudo e trabalho para construir, passo a passo, um
framework/workflow próprio de **Skills e Agents de IA** para desenvolver software de forma
organizada e profissional.

O usuário está fazendo um curso de "AI Software Engineering". No curso, o professor criou um
framework com mais de 35 skills, que o usuário achou complexo de entender. Por isso o objetivo
aqui é montar um framework **nosso**, menor e mais simples, construído aos poucos, para que o
usuário entenda e aprenda cada parte do fluxo.

## Metodologia: SDD (Spec-Driven Development)

O framework segue o SDD. Cada fase do fluxo terá a sua skill:

| # | Fase           | Skill                       | Modelo do professor                  | Status   |
|---|----------------|-----------------------------|--------------------------------------|----------|
| 1 | Pre-Refinement | `agent-spec-pre-refinement` | `agent-spec-pre-refinement`          | criada   |
| 2 | ADR            | `agent-spec-adr`            | `agent-spec-adr-create`              | criada   |
| 3 | PRD            | `agent-spec-prd`            | `agent-spec-sdd-generate-prd`        | criada   |
| — | Design (opcional, só feature com telas) | `agent-spec-design` | `agent-spec-generate-design` | criada |
| 4 | Tech Alignment | `agent-spec-tech-alignment` | `agent-spec-generate-tech-alignment` | criada   |
| 5 | Tech Spec      | `agent-spec-tech-spec`      | `agent-spec-sdd-generate-tech-spec`  | criada   |
| 6 | Task Plan      | `agent-spec-task-plan`      | `agent-spec-sdd-generate-task-plan`  | criada   |
| 7 | Debt           | `agent-spec-debt`           | `agent-spec-debt-resolution`         | criada   |
| 8 | Execução       | `agent-spec-run-tasks`      | `agent-spec-sdd-run-tasks` + agents  | criada   |

Agents da Execução: `agent-spec-executor`, `agent-spec-reviewer` e `agent-spec-security-reviewer`
(só nas tasks de risco alto; modelo: as seções de segurança do `agent-spec-staff-architecture-review`
e do `agent-spec-qa-validator`, mais itens de pagamento nossos).

Além das fases, há skills **utilitárias**, que não geram documento de fase:

| Utilitário       | Skill                        | Modelo do professor          | Status |
|------------------|------------------------------|------------------------------|--------|
| Commit semântico | `agent-spec-semantic-commit` | `agent-spec-semantic-commit` | criada |
| Challenge da spec | `agent-spec-challenge-spec` | `agent-spec-challenge-spec` | criada |
| Guia do framework | `agent-spec-guide` | `agent-spec-guide` | criada |
| Revisão das ADRs | `agent-spec-adr-review` | `agent-spec-adr-review` | criada |
| Design system (opcional) | `agent-spec-design-system` | `agent-spec-design-system-bootstrap` | criada |
| README do projeto | `agent-spec-readme` | `agent-spec-readme-generator` | criada |
| CLAUDE.md do projeto | `agent-spec-claude-md` | `agent-spec-generate-claude-md` | criada |
| Rules de um tema | `agent-spec-rule-create` | `agent-spec-rule-create` | criada |
| Curadoria das rules | `agent-spec-rule-curate` | `agent-spec-curate-project-rules` | criada |

Atualize estas tabelas sempre que uma skill for criada ou alterada (nome da skill e status).

## Como trabalhamos

1. **Claude indica** qual é a próxima skill (ou agent) a criar, e explica para que ela serve e
   onde ela entra no fluxo.
2. **Claude lê o modelo** equivalente no projeto do professor (ver abaixo), ou o usuário indica
   qual usar.
3. **Claude adapta** esse modelo para o nosso framework.

O projeto completo do professor (skills, agents, rules e exemplos de documentos gerados em
`docs/`) fica em `C:\dev\projects\web\professor`. É material de consulta: nunca editar.

Regras:

- **Não criar skills do zero.** Sempre partir de um modelo do professor. Se não houver modelo
  equivalente, avisar o usuário antes de escrever.
- **Um passo de cada vez.** Uma skill por vez, terminada e entendida antes de passar à próxima.
- **Simplicidade acima de completude.** Ao adaptar, remover o que não é necessário para o nosso
  fluxo; o objetivo é um framework fácil de entender, não uma cópia do framework do professor.
- **Explicar as decisões.** O usuário está aprendendo: ao adaptar uma skill, dizer o que foi
  mantido, o que foi removido ou mudado, e por quê.
- **Idioma:** conversar com o usuário em português.

## Estrutura de arquivos

O fluxo completo, a estrutura de diretórios e as convenções estão na skill `agent-spec-guide`
(o mapa do framework). Siga-a ao criar skills e documentos, e atualize-a quando o fluxo mudar.

- Skills: `.claude/skills/<nome-da-skill>/SKILL.md`
- Agents: `.claude/agents/<nome-do-agent>.md`
- Rules: `.claude/rules/<nome>.md` (ex.: `behavior.md`, as regras de comportamento)
