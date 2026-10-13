---
status: rascunho   # rascunho | aprovado
tamanho: média     # média | grande
---

# <Título da mudança>

- **Data:** AAAA-MM-DD · **Pedido:** <texto curto do pedido | `.temp/references/...`>
- **ADRs que se aplicam:** [ADR NNN](../adr/NNN-<titulo>.md) | nenhuma
- **Specs afetadas:** `docs/specs/<assunto>.md` (cria | altera) · `docs/architecture.md` (altera) | nenhuma

<!-- LLM-ONLY: este é o plano temporário de uma mudança: o build o executa e o apaga. O que vale depois dela (como funciona) vai para as specs afetadas, nos passos. Seção que não se aplica é removida (não fica "N/A"). Média cabe em cerca de 1 página; grande, em até 3. -->

## 1. Problema e objetivo

<2 a 5 linhas: o que está errado ou faltando hoje, para quem, e o que muda quando estiver pronto.>

## 2. Comportamento esperado

- **CA-1** — <comportamento verificável, do ponto de vista de quem usa: dado X, quando Y, então Z>
- **CA-2** — <caso de erro ou de borda: o que o usuário vê, com a mensagem literal quando importa>

**Fora do escopo:** <o que parece fazer parte, mas fica de fora, e por quê> | nada

<!-- LLM-ONLY: cada CA é verificável por um teste ou por um roteiro manual curto. CA que só uma pessoa confere (olhar a tela num projeto sem a seção Teste visual no CLAUDE.md, usar outra máquina) termina com "(manual: <roteiro em 1 linha>)": ele não entra no "Pronto quando" dos passos, e o build o lista no relatório para o usuário. Só os comportamentos que mudam; o que continua igual não entra. -->

## 3. Decisões

| Decisão | Escolha | Por quê |
|---|---|---|
| <o que precisou ser decidido> | <o que foi escolhido> | <motivo em 1 linha; alternativa descartada> |

**ADR a registrar antes de aprovar:** `/spec-adr "<título>"` | nenhuma

<!-- LLM-ONLY: só decisões não óbvias, tomadas com o usuário ou que fogem do padrão do código. O que o código, uma spec ou uma ADR já respondem não entra. ADR só na grande, para decisão que vale para o projeto todo. -->

## 4. Telas

<Por tela: como se chega, o que mostra, ações, estados (carregando, vazio, erro) com as mensagens
literais. Componentes e tokens do design system; o que for novo nele, dito aqui.>

## 5. Contratos e dados

<API (método, rota, entrada, saída, erros), formato de arquivo, tabelas e migrações (com o que
desfaz cada uma, ou "irreversível: <plano>").>

## 6. Passos

### Passo 1 — <nome>

- **Faz:** <o que muda, em nível de plano: nomes, assinaturas, códigos de erro; não o corpo das
  funções>
- **Arquivos:** `<caminho>` (novo) · `<caminho>` (modifica) · `docs/specs/<assunto>.md` (modifica)
- **Pronto quando:**
  - [ ] <teste: entrada → resultado concreto> (CA-1)
  - [ ] `docs/specs/<assunto>.md` descreve <o comportamento novo>
  - [ ] `<comando de testes e build do projeto>` sem erros

### Passo 2 — <nome>

- **Faz:** …
- **Arquivos:** …
- **Risco:** alto — <área: autenticação, pagamentos…>
- **Pronto quando:**
  - [ ] …

<!-- LLM-ONLY: passos na ordem de execução; cada um entregável e testável sozinho. "Pronto quando" só tem o que um agente prova (teste com resultado concreto, comando, trecho da spec); todo CA que não é manual aparece em algum "Pronto quando", e toda spec afetada é atualizada pelo passo que muda aquele comportamento. "Risco" só aparece no passo de risco alto. Mais de 8 passos → fatiar em planos menores. -->

## 7. Riscos e dúvidas

- <risco ou `[HIPÓTESE]`> → <como validar ou mitigar>
