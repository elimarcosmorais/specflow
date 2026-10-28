# SpecFlow

## Sobre o projeto

Um framework enxuto de **skills e agents do Claude Code** para desenvolver software com specs,
em que o usuário conversa com **um** agente — o lead — e ele distribui o trabalho para agents
especializados em sessões isoladas. O usuário só entra **antes** da execução: numa rodada de
perguntas, cada uma com a opção recomendada (obrigatória em projeto novo, onde cobre a stack e o
escopo), e para aprovar o plano. Depois disso o lead decide, executa e revisa, e entrega a tarefa **pronta** — não
existe débito: o que atrapalha a tarefa é parte dela, o que não atrapalha não existe. Ele só
volta pelos quatro motivos da lista fechada de paradas. O que descreve o sistema fica nas specs
vivas em `docs/specs/`, que o agente mantém sozinho.
Visão geral, instalação e uso estão no [README.md](README.md).

Este repositório é o próprio framework: o que se edita aqui são as skills, os agents e as rules.

## Fluxo

Uma skill de conversa, dois agents. O usuário fala **só** com o lead, e o pedido de mudança que
chega sem o `/spec-lead` vai para ele do mesmo jeito (regra 7 da `behavior.md`). Por isso a
`spec-lead` não leva `disable-model-invocation`: com ele, o agente não conseguiria carregá-la.

| Papel | Quem |
|---|---|
| Ponto único de contato: lê o código e decide; faz sozinho a mudança pequena; na maior, escreve o plano, confere cada passo, fecha e relata | skill `spec-lead` |
| Implementa um passo do plano ou a correção dos bloqueantes; cria, atualiza e substitui as ADRs que o passo pede | `spec-executor` |
| Revisão final, binária, da mudança inteira (segurança quando há `Risco: alto`) | `spec-reviewer` |

Tarefa dada é tarefa cumprida: não há fila de débito, plano secundário nem "melhoria".

Skills **utilitárias**:

| Utilitário | Skill |
|---|---|
| Guia do framework | `spec-guide` |
| Decisão de arquitetura (ADR) e o índice delas | `spec-adr` |
| Design system (obrigatório com interface) | `spec-design-system` |
| Commit semântico | `spec-semantic-commit` |
| Versão: integração inteira, SemVer, CHANGELOG, tag e GitHub Release com os instaláveis | `spec-release` |
| README do projeto | `spec-readme` |
| CLAUDE.md do projeto | `spec-claude` |
| Windows-MCP: agente vê e opera a tela (só Windows) | `spec-windows-mcp` |
| Playwright MCP: agente vê e opera páginas web | `spec-playwright-mcp` |
| Atualizar o SpecFlow no projeto a partir do GitHub | `spec-upgrade` |

## Ao alterar o framework

- **Autonomia é a feature.** A regra 1 da [`behavior.md`](.claude/rules/behavior.md) — decidir é
  o trabalho, com a lista fechada de quatro paradas — e a regra 2 — a tarefa é a entrega, não
  existe débito — são o produto, não detalhe de redação. **Regra, seção ou campo novo que crie mais um motivo
  para o agente parar e perguntar não entra.** Na dúvida, pergunte-se: isto faz o agente
  resolver mais, ou relatar mais?
- **Não fabrique débito.** Não reintroduza campo de "notei fora do escopo", categoria
  "melhoria", fila, plano secundário nem gaveta de "deixei para você": revisor pedido para achar
  sugestões sempre acha, e foi isso que fez uma demanda simples virar dias de lista. Pergunte-se:
  isto dá ao agente um lugar para deixar trabalho para depois? Então não entra.
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
- Settings: `.claude/settings.json`, só com as regras `permissions.ask` de que as rules dependem
  (ex.: segurar o `RunAs` da regra 5). No projeto o arquivo é dele: a `spec-upgrade` só
  acrescenta as regras que faltam, e o `classificar.sh` o deixa fora da cópia.
