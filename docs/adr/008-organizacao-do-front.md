---
status: aceito   # proposto | aceito | substituído
---

# ADR 008 — Organização do front

- **Data:** 2026-10-05
- **Origem:** [tech-alignment](../specs/features/001-arquitetura-base/v1/tech-alignment.md) — D-10
- **Substitui:** nenhuma
- **Substituída por:** nenhuma

## Contexto

A ADR 002 organiza o motor por feature, mas deixa o front (React + TypeScript + Vite, ADR 001)
de fora. Cada feature nova vai ganhar telas, e copiar uma feature precisa ser mecânico também no
front. As telas da v1 são poucas e simples (lista, formulário, confirmação, status), e o tema
precisa acompanhar o claro/escuro do Windows.

## Decisão

Usar React puro com CSS próprio: pastas por feature espelhando o motor, tema por variáveis CSS
(o tema "Sistema" segue `prefers-color-scheme`), chamadas com `fetch` e tipos gerados do
`openapi.yaml`, sem bibliotecas de estado nem de componentes.

## Consequências

**Ganhos:**

- Criar uma feature no front é copiar a pasta de Projetos, como no motor.
- Zero dependências de UI: nada dita as convenções do front além do próprio React.
- Os tipos das chamadas vêm do contrato da API, e não são escritos à mão.

**Custos:**

- Formulários, modais e confirmações são feitos à mão.
- A geração de tipos a partir do `openapi.yaml` é uma ferramenta nova no build, a validar.
- Se as telas crescerem muito, estado e componentes podem pedir uma biblioteca, e esta ADR
  terá de ser revista.

**Neutras:**

- O tema "Sistema" depende de o WebView2 herdar o claro/escuro do Windows.

## Alternativas consideradas

- **Tailwind + componentes próprios** — utilitários de CSS com modo escuro pronto. Rejeitada
  porque traz mais uma ferramenta e convenção para aprender, sem necessidade nas telas da v1.
- **Biblioteca de componentes (Mantine, MUI…)** — formulários, modais e tema prontos.
  Rejeitada porque é uma dependência grande que passa a ditar as convenções do front.
