---
status: aceito   # proposto | aceito | substituído
---

# ADR 006 — Padrão de logging

- **Data:** 2026-10-05
- **Origem:** [tech-alignment](../specs/features/001-arquitetura-base/v1/tech-alignment.md) — D-04
- **Substitui:** nenhuma
- **Substituída por:** nenhuma

## Contexto

O motor precisa registrar o que acontece (inicialização, porta ocupada, falhas do banco) em
arquivos dentro da pasta do app (ADR 003), e toda feature nova vai logar do mesmo jeito. Hoje
quem lê os logs é uma pessoa, abrindo o arquivo à mão. Os arquivos não podem crescer sem limite
num app que fica instalado por meses.

## Decisão

Usar `log/slog` da biblioteca padrão em formato texto (`chave=valor`), com rotação e descarte
por `lumberjack`: arquivos de até 10 MB, no máximo 5 antigos e descarte após 14 dias (no máximo
~60 MB em disco).

## Consequências

**Ganhos:**

- Logs legíveis no Bloco de Notas e ainda fáceis de buscar por chave.
- Rotação e descarte por tamanho e idade prontos e testados, sem código nosso.
- `slog` é da biblioteca padrão: as features logam pela mesma API, sem depender de `lumberjack`.

**Custos:**

- Uma dependência nova (`lumberjack`).
- A rotação no Windows com o arquivo aberto e o descarte na abertura são hipóteses a validar
  na implementação.
- Processar os logs por programa é mais difícil que com JSON.

**Neutras:**

- A configuração do logger fica na infraestrutura comum (`platform/log`, ADR 002).

## Alternativas consideradas

- **`slog` em JSON + `lumberjack`** — uma linha JSON por evento. Rejeitada porque é ruim de ler
  à mão, e a tela de logs, que se beneficiaria do JSON, está adiada.
- **`slog` em texto + rotação própria** — um arquivo por dia, com limpeza ao iniciar e a cada
  24 h. Rejeitada porque seriam ~60 linhas nossas para manter e testar, reimplementando algo
  maduro.
