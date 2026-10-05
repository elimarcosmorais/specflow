---
status: rascunho   # rascunho | aprovado
---

# PRD — <nome da feature>

> **Legenda:** `[DECISÃO TÉCNICA]` = ponto que envolve o COMO; fica registrado aqui de forma
> abstrata e é decidido no Tech Alignment.

- **Feature:** `NNN-<nome-da-feature>` · **Versão:** vN · **Data:** AAAA-MM-DD
- **Origem:** [pre-refinement.md](pre-refinement.md)
- **ADRs relacionadas:** `docs/adr/NNN-<titulo>.md` | nenhuma

---

## 1. Contexto e motivação

- **Qual é o problema hoje?**
- **Quem sofre o impacto?**
- **Por que resolver agora?**

## 2. Objetivo

- **Resultado esperado, do ponto de vista do usuário:**
- **O que muda no comportamento de quem usa:**

## 3. Personas

- **Primária:** <quem> — <o que quer ao usar a feature>
- **Secundária:** <quem> — <o que quer> | nenhuma

## 4. Escopo

**Dentro:**

- <comportamento esperado>

**Fora** (com motivo):

- <item> — <motivo: descartado, adiado para outra versão ou fora do escopo do projeto>

## 5. Histórias de usuário

- **US-01:** Como <persona>, quero <ação> para <resultado>.

<!-- LLM-ONLY: uma história por comportamento que entrega valor. Todo item de "Dentro" (seção 4) deve estar coberto por pelo menos uma história. -->

## 6. Regras de negócio

- **RN-01:** <condição do domínio que o sistema deve respeitar>.

<!-- LLM-ONLY: só regras do domínio ("um pedido cancelado não pode ser reaberto"), nunca lógica técnica. -->

## 7. Fluxo de uso

**Principal:**

1. O usuário <ação>.
2. O sistema <resposta>.

**Alternativos:**

- **<situação: erro, dado faltando, caso-limite>** → <o que o sistema faz>.

## 8. Critérios de aceite

- [ ] **CA-01:** DADO <situação> QUANDO <ação> ENTÃO <resultado observável>.

<!-- LLM-ONLY: todo CA em DADO/QUANDO/ENTÃO, com resultado que dá para verificar. Cubra o fluxo principal e cada fluxo alternativo. -->

## 9. Rastreabilidade

| História | Resumo | Critérios de aceite |
|----------|--------|---------------------|
| US-01    | <resumo em poucas palavras> | CA-01, CA-02 |

<!-- LLM-ONLY: toda US com pelo menos um CA; todo CA presente na tabela. -->

## 10. Premissas e restrições

**Premissas** (confirmadas com o usuário):

- <premissa>

**Restrições** (prazo, dependência, imposição externa):

- <restrição> | `[DECISÃO TÉCNICA]` <imposição técnica descrita sem nome de tecnologia>

## 11. Métricas de sucesso

- <como saber, depois de entregue, que a feature resolveu o problema>

## 12. Próximo passo

- **Próxima fase:** Tech Alignment → `/agent-spec-tech-alignment docs/specs/features/NNN-<feature>/vN/prd.md`
- **ADR a registrar antes:**
  - `/agent-spec-adr "<título da decisão>"` — <o que ela decide, em poucas palavras>

<!-- LLM-ONLY: os comandos ficam prontos para copiar, com o caminho real deste arquivo. ADR aqui = recomendada pelo pre-refinement e ainda ausente em docs/adr/; sem nenhuma → "**ADR a registrar antes:** nenhuma". -->

