---
status: rascunho   # rascunho | aprovado
---

# Task Plan — Pagamento de débitos · <nome da feature>

- **Feature:** `NNN-<nome-da-feature>` · **Versão:** vN · **Data:** AAAA-MM-DD
- **Tipo:** pagamento de débitos (gerado por `/agent-spec-debt`)
- **Origem:** [../v1/debts.md](../v1/debts.md) · **Spec vigente:** [../v1/spec.md](../v1/spec.md)

<!-- LLM-ONLY: a Origem lista todos os debts.md de onde vieram os débitos escolhidos. A spec vigente é a da versão mais recente que tem spec.md. -->

---

## 1. Objetivo

Pagar K débitos técnicos de <versões>, sem mudar o comportamento da feature: ao final, os
testes existentes continuam passando.

## 2. Débitos pagos

| Débito   | Severidade | Título    | Task |
|----------|------------|-----------|------|
| v1/DT-01 | média      | <título>  | T1   |

<!-- LLM-ONLY: só os débitos escolhidos, cada um em exatamente uma task. Os que ficaram abertos continuam no debts.md de origem e não aparecem aqui. -->

## 3. Tasks

| ID | Task | Fase | Depende de | Status |
|----|------|------|------------|--------|
| T1 | [Pagar v1/DT-01: <título>](tasks/T1.md) | 1 | — | pendente |

<!-- LLM-ONLY: a fonte do status é o frontmatter de cada tasks/TN.md; esta coluna é um espelho, atualizado junto pela Execução. Status: pendente | em andamento | concluída. -->

### 3.1 Ordem de execução

```
T1 ∥ T2 -> T3
```

<!-- LLM-ONLY: "∥" só entre tasks que não tocam nenhum arquivo em comum. Na dúvida, sequencial. -->

## 4. Critérios de conclusão

- [ ] Todas as tasks `concluída`
- [ ] Todos os testes automatizados passando: `<comando de testes do projeto>`
- [ ] `<comando de build/lint do projeto>` sem erros
- [ ] Comportamento da feature inalterado (nenhuma API, erro, dado ou mensagem mudou)

<!-- LLM-ONLY: use os comandos reais do projeto. -->

## 5. Próximo passo

- **Próxima fase:** Execução → `/agent-spec-run-tasks docs/specs/features/NNN-<feature>/vN/tasks.md`
