---
status: aceito   # proposto | aceito | substituído
---

# ADR 004 — Acesso à API local: token por execução e bootstrap do front

- **Data:** 2026-10-05
- **Origem:** [tech-alignment](../specs/features/001-arquitetura-base/v1/tech-alignment.md) — D-01
- **Substitui:** nenhuma
- **Substituída por:** nenhuma

## Contexto

A API do motor escuta em `127.0.0.1` numa porta fixa (ADR 001), e qualquer programa local ou
página aberta no navegador consegue chegar a ela. As operações que alteram dados só podem ser
aceitas se vierem do próprio Dev Tools. Além disso, o front precisa saber o estado da
inicialização (ex.: porta ocupada) e as preferências mesmo quando a porta não está disponível,
ou seja, sem depender da própria API.

## Decisão

Gerar a cada execução um token aleatório, mantido só em memória e exigido nas operações que
alteram dados. O front recebe o token por uma única chamada interna do Wails de bootstrap, que
entrega também o endereço da API, o estado da inicialização e as preferências; depois dela,
toda comunicação segue pela API HTTP.

## Consequências

**Ganhos:**

- Outros programas e páginas do navegador não conseguem alterar dados: não conhecem o token.
- O front mostra o estado da inicialização (inclusive porta ocupada) sem precisar da API.
- O token não vai para o disco: nada a proteger, rotacionar ou limpar entre execuções.

**Custos:**

- Abre uma exceção à regra "o front fala com o motor só pela API HTTP" (ADR 001): o bootstrap
  é a única chamada pelos bindings do Wails, e nenhuma outra deve ser criada por esse caminho.
- Toda rota nova que altera dados precisa exigir o token; esquecer abre uma brecha.
- Testes de integração da API precisam obter e enviar o token.

**Neutras:**

- A futura CLI terá de obter um token por outro meio, a definir quando ela existir.

## Alternativas consideradas

- **Checar o cabeçalho `Origin`** — recusar requisições que não venham da origem do Wails.
  Rejeitada porque qualquer programa forja o `Origin` (um `curl` passa).
- **Rotas de Projetos só pelo AssetServer do Wails** — atender as rotas em processo e deixar na
  porta TCP só saúde e documentação. Rejeitada porque contraria a ADR 001 (API na porta,
  consumida pelo front) e depende de hipóteses não validadas sobre `wails dev` e POST no
  WebView2.
