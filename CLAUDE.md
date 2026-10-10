# AI Software Engineering

## Sobre o projeto

Um framework enxuto de **skills e agents do Claude Code** para desenvolver software com specs,
em que o processo escala com o tamanho da mudança: pequena → implementa direto, sem plano;
média ou grande → um único plano em `docs/plans/`, aprovado pelo humano, executado passo a passo
com revisão e apagado no fim. O que descreve o sistema fica nas specs vivas em `docs/specs/`, que
o agente mantém sozinho.
Visão geral, instalação e uso estão no [README.md](README.md).

Este repositório é o próprio framework: o que se edita aqui são as skills, os agents e as rules.

## Fluxo

| Etapa | Skill |
|---|---|
| Entrada: classifica, implementa o pequeno, escreve o plano | `spec-change` |
| Executa um plano aprovado | `spec-build` |

Agents do build: `spec-executor`, `spec-reviewer` e `spec-security-reviewer`
(só nos passos de risco alto).

Skills **utilitárias**:

| Utilitário | Skill |
|---|---|
| Guia do framework | `spec-guide` |
| Decisão de arquitetura (ADR) | `spec-adr` |
| Design system (obrigatório com interface) | `spec-design-system` |
| Commit semântico | `spec-semantic-commit` |
| README do projeto | `spec-readme` |
| CLAUDE.md do projeto | `spec-claude-md` |
| Rules de um tema | `spec-rule-create` |
| Curadoria das rules | `spec-rule-curate` |
| Windows-MCP: agente vê e opera a tela (só Windows) | `spec-windows-mcp` |
| Playwright MCP: agente vê e opera páginas web | `spec-playwright-mcp` |

## Ao alterar o framework

- **O processo escala com a mudança.** Não reintroduza fase obrigatória, documento por fase,
  versões de documento, arquivo por task nem débito como fase: foi o que tornou o fluxo antigo
  inviável. Antes de acrescentar uma skill, seção ou etapa, mostre que ela paga o próprio custo.
- **Uma mudança por vez.** Uma skill ou agent por vez, terminada antes de passar à próxima.
- **Explicar as decisões.** Ao criar ou alterar uma skill, dizer o que mudou e por quê.
- **Manter tudo em sincronia** ao criar, alterar ou remover uma skill, agent ou rule: as tabelas
  deste arquivo, as do [README.md](README.md) e o mapa do framework (`spec-guide`).
- **Idioma:** conversar com o usuário e escrever skills e documentos em português.

## Estrutura de arquivos

O fluxo, a estrutura de diretórios e as convenções estão na skill `spec-guide`. Siga-a ao
criar skills e documentos, e atualize-a quando o fluxo mudar.

- Skills: `.claude/skills/<nome-da-skill>/SKILL.md`
- Agents: `.claude/agents/<nome-do-agent>.md`
- Rules: `.claude/rules/<nome>.md` (`behavior.md`, as regras de comportamento; `docs.md`, onde
  cada documento mora e quem mexe nele)
