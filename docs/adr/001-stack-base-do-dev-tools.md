---
status: aceito   # proposto | aceito | substituído
---

# ADR 001 — Stack base do Dev Tools

- **Data:** 2026-10-05
- **Origem:** [pre-refinement](../specs/features/001-arquitetura-base/v1/pre-refinement.md) — seções 8, 10 (dúvida 1) e 11
- **Substitui:** nenhuma
- **Substituída por:** nenhuma

## Contexto

O Dev Tools é um executável para Windows que abre como app desktop e expõe uma API HTTP só em
`127.0.0.1`, consumida pelo próprio front e documentada. Precisa de persistência local num
arquivo dentro da pasta do app, sem pedir admin e sem depender de compilador C na máquina de
build. A stack escolhida aqui vale para todas as features seguintes.

## Decisão

Usar Go 1.25+ com Wails v2 (WebView2) no app desktop, React 19 + TypeScript + Vite no front,
SQLite via `modernc.org/sqlite` (sem CGO), `net/http` da biblioteca padrão na API, contrato
OpenAPI escrito à mão (`openapi.yaml`) servido numa página `/docs`, e `go test` nos testes do
motor.

## Consequências

**Ganhos:**

- Um único executável, sem CGO: build simples no Windows e nada para instalar além do WebView2.
- Wails v2 é a versão estável, com documentação e exemplos de sobra.
- A API não depende de framework: rotas com método e parâmetro (`GET /projects/{id}`) vêm do
  `ServeMux` do Go 1.22+, e middlewares são funções comuns.
- O `openapi.yaml` é o contrato da API: dá para ler, revisar e gerar os tipos TS do front a
  partir dele.

**Custos:**

- O `openapi.yaml` precisa ser atualizado junto com os handlers; um teste deve conferir que as
  rotas registradas batem com o contrato.
- Migrar para Wails v3 quando ela ficar estável exigirá ajustes na camada desktop.
- `modernc.org/sqlite` é mais lento que o driver com CGO; irrelevante no volume de um app local.
- Agrupamento de rotas e middlewares prontos (que um chi ou Gin dariam) ficam por nossa conta.

**Neutras:**

- O front fala com o motor pela API HTTP, não pelos bindings do Wails (rumo C2 do
  pre-refinement).

## Alternativas consideradas

- **Wails v3** — nova versão, com múltiplas janelas e tray nativo. Rejeitada porque ainda está
  em beta (API pode mudar) e a fundação não precisa desses recursos.
- **chi** — roteador leve compatível com `net/http`. Rejeitada porque o `ServeMux` atual já
  cobre o que a API local precisa; seria uma dependência por pouco ganho.
- **Gin ou Echo** — frameworks HTTP completos. Rejeitados pelo acoplamento a um contexto
  próprio e pela complexidade extra.
- **swaggo (anotações)** — gera Swagger a partir de comentários. Rejeitado porque gera
  OpenAPI 2.0 e os comentários também saem de sincronia.
- **Huma** — gera OpenAPI 3.1 dos structs Go. Rejeitado porque passa a ditar a forma dos
  handlers (mais um framework para aprender).
- **SQLite com CGO (`mattn/go-sqlite3`)** — driver mais rápido. Rejeitado porque exige
  compilador C no build.
