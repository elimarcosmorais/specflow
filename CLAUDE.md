# SpecFlow

## Sobre o projeto

Um framework enxuto de **skills e agents do Claude Code** para desenvolver software com specs,
em que o usuário conversa com **um** agente — o lead — e ele distribui o trabalho para agents
especializados em sessões isoladas. O usuário é interrompido **uma vez**: para aprovar o plano
principal. Depois disso o lead decide, executa, revisa e resolve o débito que aparece (em planos
secundários, sem aprovação), e só volta pelos quatro motivos da lista fechada de paradas. O que
descreve o sistema fica nas specs vivas em `docs/specs/`, que o agente mantém sozinho.
Visão geral, instalação e uso estão no [README.md](README.md).

Este repositório é o próprio framework: o que se edita aqui são as skills, os agents e as rules.

## Fluxo

Uma skill de conversa, seis agents. O usuário fala **só** com o lead.

| Papel | Quem |
|---|---|
| Ponto único de contato: entende, decide, delega, fecha e relata | skill `spec-lead` |
| Lê o projeto e devolve o mapa (read-only) | `spec-scout` |
| Escreve o plano principal, as partes e os secundários | `spec-planner` |
| Implementa um passo, uma correção ou um item avulso | `spec-executor` |
| Revisa o passo; o de segurança só em `Risco: alto` | `spec-reviewer`, `spec-security-reviewer` |
| Registra a ADR de uma decisão já tomada | `spec-architect` |

Débito achado pelo caminho não vira recado: o executor corrige o que cabe no diff, e o resto
vira **plano secundário**, que o lead escreve e executa sem perguntar. Sem teto de itens.

Skills **utilitárias**:

| Utilitário | Skill |
|---|---|
| Guia do framework | `spec-guide` |
| Decisão de arquitetura (ADR) e o índice delas | `spec-adr` |
| Design system (obrigatório com interface) | `spec-design-system` |
| Commit semântico | `spec-semantic-commit` |
| Versão: SemVer, CHANGELOG, tag e GitHub Release | `spec-release` |
| README do projeto | `spec-readme` |
| CLAUDE.md do projeto | `spec-claude` |
| Rules de um tema | `spec-rule-create` |
| Curadoria das rules | `spec-rule-curate` |
| Windows-MCP: agente vê e opera a tela (só Windows) | `spec-windows-mcp` |
| Playwright MCP: agente vê e opera páginas web | `spec-playwright-mcp` |
| Atualizar o SpecFlow no projeto a partir do GitHub | `spec-upgrade` |

## Ao alterar o framework

- **Autonomia é a feature.** A regra 1 da [`behavior.md`](.claude/rules/behavior.md) — decidir é
  o trabalho, com a lista fechada de quatro paradas — e a regra 2 — nada volta como tarefa —
  são o produto, não detalhe de redação. **Regra, seção ou campo novo que crie mais um motivo
  para o agente parar e perguntar não entra.** Na dúvida, pergunte-se: isto faz o agente
  resolver mais, ou relatar mais?
- **Teto gera fila.** Não reintroduza teto de itens por relatório, por passo ou por agent, nem
  gaveta de "deixei para você": foi o que fez uma demanda simples virar dias de lista. Limite a
  **profundidade** (dois lotes de planos secundários), nunca a **quantidade**.
- **O processo escala com a mudança.** Não reintroduza fase obrigatória, documento por fase,
  versões de documento, arquivo por task nem débito como fase: foi o que tornou o fluxo de 8
  fases inviável. Antes de acrescentar uma skill, seção ou etapa, mostre que ela paga o próprio
  custo. O que já foi tentado e descartado está em `spec-guide`, seção 6.
- **Uma mudança por vez.** Uma skill ou agent por vez, terminada antes de passar à próxima.
- **Explicar as decisões.** Ao criar ou alterar uma skill, dizer o que mudou e por quê.
- **Manter tudo em sincronia** ao criar, alterar ou remover uma skill, agent ou rule: as tabelas
  deste arquivo, as do [README.md](README.md) e o mapa do framework (`spec-guide`).
- **Idioma:** conversar com o usuário e escrever skills e documentos em português.

## Estrutura de arquivos

O fluxo, a estrutura de diretórios e as convenções estão na skill `spec-guide`. Siga-a ao
criar skills e documentos, e atualize-a quando o fluxo mudar.

- Skills: `.claude/skills/<nome-da-skill>/SKILL.md` (templates ao lado: a `spec-lead` tem
  `plano.md` e `relatorio.md`)
- Agents: `.claude/agents/<nome-do-agent>.md`
- Rules: `.claude/rules/<nome>.md` (`behavior.md`, as regras de comportamento; `docs.md`, onde
  cada documento mora e quem mexe nele)
