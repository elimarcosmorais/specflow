---
status: aceito   # proposto | aceito | substituído
---

# ADR 002 — Arquitetura em camadas e padrão de feature

- **Data:** 2026-10-05
- **Origem:** [pre-refinement](../specs/features/001-arquitetura-base/v1/pre-refinement.md) — seções 6 e 11
- **Substitui:** nenhuma
- **Substituída por:** nenhuma

## Contexto

O motor em Go ([ADR 001](001-stack-base-do-dev-tools.md)) vai crescer feature a feature (PHP,
Node.js, Apache…), e cada uma precisa separar regra de negócio, infraestrutura e testes do
mesmo jeito. A feature-modelo (Projetos) é o molde que as próximas copiam, então criar uma
feature nova tem que ser mecânico. A regra de negócio precisa ser testável sem banco nem HTTP.

## Decisão

Organizar o motor em pacotes por feature (`internal/features/<feature>/`), cada um com entidade,
serviço, interface de repositório, implementação SQLite e handler HTTP, e a infraestrutura comum
em `internal/platform/` (db, log, i18n, httpserver, config). O núcleo (entidade e serviço) nunca
importa infraestrutura: depende de interfaces que a infra implementa, e `cmd/dev-tools/main.go`
monta tudo.

Testes unitários ficam no `_test.go` ao lado do código, com fakes das interfaces; testes de
integração ficam em `*_integration_test.go` na própria feature, com `//go:build integration`,
SQLite real em pasta temporária e `httptest`.

## Consequências

**Ganhos:**

- Criar uma feature é copiar a pasta de `projects/` e trocar os nomes.
- Regra de negócio testável em milissegundos, sem banco e sem servidor.
- Trocar uma peça de infra (ex.: o banco) mexe só no adaptador, não no serviço.
- `go test ./...` roda só os unitários; `go test -tags integration ./...` roda tudo.

**Custos:**

- Uma interface e um fake a mais por repositório, mesmo quando só existe a implementação SQLite.
- É preciso vigiar a regra de dependência nas revisões: nada impede, no compilador, que o
  serviço importe `platform/db`.
- O que for comum a várias features precisa ir para `platform/`, e não ser copiado entre elas.

**Neutras:**

- A organização do front (React) não é coberta por esta ADR.

## Alternativas consideradas

- **Camadas horizontais (Clean clássica)** — pastas globais `domain/`, `application/`,
  `infrastructure/` e `http/`. Rejeitada porque espalha cada feature em quatro pastas e torna
  copiar o modelo mais trabalhoso, com a mesma regra de dependência.
- **Handler → service → repository sem interfaces** — três pastas acopladas diretamente.
  Rejeitada porque prende o serviço ao SQLite e obriga os testes unitários a usar banco.
- **Testes de integração numa pasta `tests/` separada** — testam a API só por HTTP. Rejeitada
  porque ficam longe da feature que testam e não alcançam o que não é exportado.
