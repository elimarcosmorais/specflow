# Changelog

Todas as mudanças relevantes deste projeto ficam neste arquivo.
O formato segue o [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/), e o projeto usa
[Versionamento Semântico](https://semver.org/lang/pt-BR/).

## [2.3.0] - 2026-10-27

### Adicionado

- **lead:** leva à spec-lead o pedido de mudança sem o comando (2d934ad)

## [2.2.1] - 2026-10-21

### Alterado

- itera com o teste dirigido e roda a suíte uma vez por passo (14edf8a)

## [2.2.0] - 2026-10-18

### Adicionado

- **behavior:** segura o `RunAs` numa regra `ask` até o usuário aprovar (55ccab3)
- **adr:** passa a criação, atualização e substituição de ADR ao agente (d144f00)

## [2.1.0] - 2026-10-17

### Adicionado

- **spec-lead:** pergunta a fundação e o escopo em projeto novo (26a0b94)
- **behavior:** separa risco fechado de risco aberto no relatório (b522caa)
- **agents:** herda o modelo e o esforço das configurações do usuário (82d7dae)

## [2.0.0] - 2026-10-17

### Incompatível

- **skills:** remove as skills spec-rule-create e spec-rule-curate (2d76484)

  `/spec-rule-create` e `/spec-rule-curate` deixam de existir. Para
  acrescentar uma regra ao projeto, use `/spec-claude`.

### Adicionado

- **behavior:** grava conteúdo com `\` pelo Write/Edit, não pelo Bash (6266a0e)

## [1.0.0] - 2026-10-17

### Incompatível

- reduz o fluxo a dois agents e elimina o débito (8b1bcb6)

  Os agents `spec-scout`, `spec-planner`, `spec-architect` e
  `spec-security-reviewer` foram removidos; o plano perde os campos `tipo` e
  `parte` e não há mais plano secundário nem fatiado. Planos em andamento nesse
  formato devem ser reescritos com `/spec-lead`, e os arquivos dos agents
  removidos apagados do projeto.

## [0.6.0] - 2026-10-16

### Incompatível

- unifica o fluxo no `spec-lead`, com seis agents especializados (3402fca)

  `/spec-change` e `/spec-build` saem. Use `/spec-lead "<pedido>"` para uma mudança
  nova e `/spec-lead continuar` para retomar um plano com passo pendente. O agent
  `spec-fixer` foi removido, e o template do plano passou de
  `.claude/skills/spec-change/template.md` para `.claude/skills/spec-lead/plano.md`.

### Adicionado

- **uac:** trata não autorizado como passo sem prova e proíbe contorno (1e85ac1)
- **varredura:** autoriza a área do plano e isenta risco do teto (1e12243)
- **varredura:** limita o fora do escopo a 2 itens e resolve o resto (4086dc9)
- **adr:** permite a varredura corrigir erro verificável em ADR (ada3b25)

### Corrigido

- **upgrade:** fixa o caminho do clone em `/tmp/specflow-upgrade` (ddf54ba)
- trata teste pulado como falta de prova em todo o fluxo (ea0841d)

## [0.5.2] - 2026-10-14

### Corrigido

- **release:** confere o gh na Etapa 1, antes de publicar a tag (8d8ac66)

## [0.5.1] - 2026-10-14

### Corrigido

- **windows-mcp:** corrige o que o uso.md diz sobre use_ui_tree (2c91af2)

## [0.5.0] - 2026-10-14

### Adicionado

- **mcp:** extrai as dicas de uso para um uso.md ao lado da skill (17bdbd5)

## [0.4.0] - 2026-10-14

### Adicionado

- **windows-mcp:** adota fork pinado por commit como fonte do servidor (273d92f)
- **design-system:** entrega um formulário de especificação (a6d28f3)

## [0.3.0] - 2026-10-13

### Adicionado

- adiciona a varredura fora do escopo com o agent `spec-fixer` (5c398c4)

### Corrigido

- **windows-mcp:** volta para a `0.8.5`, que não congela a tela (30ee7b6)

## [0.2.0] - 2026-10-13

### Incompatível

- **references:** move o material de consulta para `.temp/` (b1104e8)
  projetos que usam o framework devem mover `docs/references/` para `.temp/references/` e
  garantir `.temp/` no `.gitignore` (rodar `/spec-claude` resolve a segunda parte).

### Adicionado

- **build:** tenta as sugestões dos revisores antes do relatório final (fa1a2df)

## [0.1.0] - 2026-10-12

### Incompatível

- **claude:** renomeia a skill spec-claude-md para spec-claude (d6f99f8)
  o comando `/spec-claude-md` passa a ser `/spec-claude`.
- troca o prefixo `agent-` de skills e agents por `spec-` (2fc8c61)
  os comandos mudam de nome (ex.: `/agent-spec-change` → `/spec-change`, `/agent-windows-mcp` →
  `/spec-windows-mcp`); quem copiou a `.claude/` precisa copiá-la de novo e apagar as pastas
  antigas.
- **framework:** troca o spec por plano temporário e specs vivas (b0c7da2)
  docs/specs/NNN-<nome>.md deixa de ser o plano da mudança. Planos novos vão para docs/plans/, o
  status `implementado` deixa de existir, docs/backlog.md não é mais usado e o design system sai
  de docs/specs/design-system.md para docs/design-system.md.
- **framework:** troca as 8 fases pelo fluxo change → spec → build (6e7924e)
  as skills de fase deixam de existir. Mudanças novas começam por /agent-spec-change, e specs
  aprovados rodam com /agent-spec-build.

### Adicionado

- **release:** adiciona a skill spec-release para lançar versões (31c5452)
- **upgrade:** adiciona a skill spec-upgrade para atualizar o SpecFlow (d8a1abf)
- impede pular a conferência de tela quando a ferramenta falha (ad0dae0)
- corrige fora do escopo o problema de impacto pequeno (217742e)
- **semantic-commit:** commita e faz push direto, sem pedir aprovação (7fbab10)
- grava README e CLAUDE.md direto, sem pedir confirmação (7ac8698)
- **adr:** adiciona o índice `docs/adr/README.md` das ADRs em vigor (4e66268)
- **mcp:** registra o servidor de tela na seção Teste visual (6f42c5b)
- **framework:** adiciona a seção Teste visual ao `CLAUDE.md` (4c6fcb7)
- **playwright-mcp:** adiciona a skill `agent-playwright-mcp` (6b3cde6)
- **windows-mcp:** atualiza para a `0.8.8` com Python 3.14 (98351e5)
- **windows-mcp:** adiciona a skill `agent-windows-mcp` (0835bd6)
- **references:** trata `docs/references/` como material temporário (9e02c14)
- **behavior:** adiciona a regra de elevar só o comando que pede admin (42dc9a4)
- **design-system:** adiciona a seção Próximo passo com Design e ADRs (232b942)
- **design:** torna o Design e o design system obrigatórios com telas (b0de082)
- **visual-guide:** adiciona o guia visual e a skill que o mantém (8a903e4)
- **claude-md:** separa convenções de área do CLAUDE.md em rules (ef22672)
- **framework:** adiciona a fase de Design e novas skills utilitárias (bceabf0)
- **framework:** adiciona revisão de segurança e skills utilitárias (09b198d)
- **framework:** adiciona as skills e os agents do fluxo SDD (330e23d)
- **dev-tools:** implementa a arquitetura base do app desktop (1a4ce16)

### Corrigido

- **semantic-commit:** atualiza o rodapé `Refs:` para citar o plano (6dae914)
- **claude-md:** distingue servidor de tela declarado de ativo (e5ac775)
- **mcp:** usa o `.mcp.json` como fonte do MCP de tela no `CLAUDE.md` (36bf744)

[2.3.0]: https://github.com/elimarcosmorais/specflow/compare/v2.2.1...v2.3.0
[2.2.1]: https://github.com/elimarcosmorais/specflow/compare/v2.2.0...v2.2.1
[2.2.0]: https://github.com/elimarcosmorais/specflow/compare/v2.1.0...v2.2.0
[2.1.0]: https://github.com/elimarcosmorais/specflow/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/elimarcosmorais/specflow/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/elimarcosmorais/specflow/compare/v0.6.0...v1.0.0
[0.6.0]: https://github.com/elimarcosmorais/specflow/compare/v0.5.2...v0.6.0
[0.5.2]: https://github.com/elimarcosmorais/specflow/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/elimarcosmorais/specflow/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/elimarcosmorais/specflow/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/elimarcosmorais/specflow/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/elimarcosmorais/specflow/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/elimarcosmorais/specflow/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/elimarcosmorais/specflow/releases/tag/v0.1.0
