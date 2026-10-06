---
status: rascunho   # rascunho | aprovado
---

# Task Plan — <nome da feature>

- **Feature:** `NNN-<nome-da-feature>` · **Versão:** vN · **Data:** AAAA-MM-DD
- **Origem:** [spec.md](spec.md) · [prd.md](prd.md)

<!-- LLM-ONLY: em mudança pequena (versão sem prd.md), a Origem aponta para o spec.md e o pre-refinement.md desta versão e para o prd.md mais recente das versões anteriores (ex.: ../v1/prd.md). -->

---

## 1. Objetivo

<2 a 3 linhas: o que estará entregue e funcionando quando todas as tasks estiverem concluídas.>

## 2. Fases

| Fase | Objetivo | Tasks |
|------|----------|-------|
| 1 — <nome> | <o que esta fase deixa pronto> | T1 |
| 2 — <nome> | <…> | T2, T3 |

## 3. Tasks

| ID | Task | Fase | Depende de | Risco | Status |
|----|------|------|------------|-------|--------|
| T1 | [<nome>](tasks/T1.md) | 1 | — | baixo | pendente |
| T2 | [<nome>](tasks/T2.md) | 2 | T1 | alto | pendente |

<!-- LLM-ONLY: a fonte é cada tasks/TN.md (status no frontmatter, risco no cabeçalho); estas colunas são espelho. O Status é atualizado junto pela Execução. Status: pendente | em andamento | concluída. Risco: baixo | alto. -->

### 3.1 Ordem de execução

```
T1 -> (T2 ∥ T3) -> T4
```

<!-- LLM-ONLY: "∥" marca tasks que podem ser feitas em qualquer ordem: nenhuma depende da outra (nem indiretamente) e elas não tocam nenhum arquivo em comum. Na dúvida, sequencial. -->

## 4. Rastreabilidade

### 4.1 Histórias → tasks

| História | Tasks |
|----------|-------|
| US-01    | T1, T2 |

### 4.2 Testes → tasks

| CT    | CA    | Task |
|-------|-------|------|
| CT-01 | CA-01 | T1   |

<!-- LLM-ONLY: toda US do PRD tem uma linha em 4.1. Todo CT da spec aparece exatamente uma vez em 4.2. CT criado nesta fase (fora da spec) leva "(novo)" depois do ID. -->

## 5. Critérios de conclusão da feature

- [ ] Todas as tasks `concluída`
- [ ] Todos os testes automatizados passando: `<comando de testes do projeto>`
- [ ] Testes manuais executados conforme o roteiro (CT-XX, CT-YY)
- [ ] `<comando de build/lint do projeto>` sem erros

<!-- LLM-ONLY: use os comandos reais do projeto. Sem testes manuais, remova a linha. -->

## 6. Próximo passo

- **Próxima fase:** Execução → `/agent-spec-run-tasks docs/specs/features/NNN-<feature>/vN/tasks.md`
