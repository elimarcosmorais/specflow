---
status: aceito   # proposto | aceito | substituído
---

# ADR 005 — Padrão de erros e idiomas

- **Data:** 2026-10-05
- **Origem:** [tech-alignment](../specs/features/001-arquitetura-base/v1/tech-alignment.md) — D-02
- **Substitui:** nenhuma
- **Substituída por:** nenhuma

## Contexto

Toda mensagem mostrada ao usuário precisa sair no idioma escolhido, inclusive erros da API e
diálogos nativos que o motor mostra antes de existir janela. O front precisa reagir a cada erro
(ex.: destacar o campo inválido) sem depender do texto, e os testes Go precisam conferir a
mensagem traduzida. Com a porta da API ocupada, o front continua precisando dos próprios textos.

## Decisão

Cada camada traduz o que ela mesma produz: o motor devolve todo erro com um código estável e a
mensagem já traduzida no idioma salvo, usando o próprio catálogo (erros e diálogos nativos); o
front tem o catálogo das telas (rótulos, títulos, textos de status).

## Consequências

**Ganhos:**

- O front decide o que fazer pelo código, e mostra a mensagem sem traduzir nada.
- Os testes Go verificam o texto traduzido de verdade.
- O front tem todos os textos das telas mesmo sem API.
- Toda feature nova segue a mesma fronteira: o motor traduz o que diz, o front o que desenha.

**Custos:**

- Dois catálogos para manter; um idioma novo exige traduzir os dois.
- Todo erro novo exige um código estável e a tradução em todos os idiomas do catálogo do motor.
- O código do erro vira contrato da API: renomear quebra o front.

**Neutras:**

- O catálogo do motor fica na infraestrutura comum (`platform/i18n`, ADR 002), e o idioma vem
  das preferências entregues no bootstrap (ADR 004).

## Alternativas consideradas

- **Só o front traduz** — o motor devolve apenas o código. Rejeitada porque os testes Go não
  verificariam texto traduzido, e os diálogos nativos exigiriam catálogo no motor do mesmo jeito.
- **Tudo no motor** — o front baixa o dicionário das telas pela API. Rejeitada porque, com a
  porta da API ocupada, o front ficaria sem nenhum texto.
