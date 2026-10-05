---
status: aceito   # proposto | aceito | substituído
---

# ADR 007 — Migrações do banco de dados

- **Data:** 2026-10-05
- **Origem:** [tech-alignment](../specs/features/001-arquitetura-base/v1/tech-alignment.md) — D-05
- **Substitui:** nenhuma
- **Substituída por:** nenhuma

## Contexto

O esquema do SQLite (ADR 001) vai mudar a cada feature e a cada versão de feature. Como a
aplicação é portátil (ADR 003), uma instalação antiga pode ser copiada ou atualizada e abrir um
banco criado por uma versão anterior: o esquema precisa se ajustar sozinho ao iniciar, sem
ferramenta externa nem intervenção do usuário.

## Decisão

Usar migrações SQL numeradas, embutidas no executável (`embed`) e aplicadas ao iniciar, com a
versão aplicada guardada no próprio arquivo do banco (`PRAGMA user_version`); cada feature nova
acrescenta uma migração.

## Consequências

**Ganhos:**

- Uma instalação antiga copiada ou atualizada ajusta o banco sozinha ao abrir.
- O executável continua autossuficiente: as migrações vão dentro dele.
- Mecanismo de ~30 linhas, sem dependência e fácil de entender.

**Custos:**

- Não há rollback de migração; corrigir um erro exige uma migração nova.
- O código de aplicar migrações, a ordem e as falhas no meio do caminho ficam por nossa conta.

**Neutras:**

- Toda mudança de esquema passa a ser uma migração nova; nenhuma tabela é criada fora delas.

## Alternativas consideradas

- **`CREATE TABLE IF NOT EXISTS` ao iniciar** — criar as tabelas que faltam. Rejeitada porque
  não altera tabela existente e quebra na v2 de qualquer feature.
- **Biblioteca de migração (`goose` ou `golang-migrate`)** — completa, com rollback. Rejeitada
  porque traz dependência e conceitos a mais para um app local de um usuário, onde o rollback
  não é necessário.
