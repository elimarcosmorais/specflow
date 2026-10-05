---
status: aprovado # rascunho | aprovado
---

# Task Plan — Pagamento de débitos · Arquitetura base do Dev Tools

- **Feature:** `001-arquitetura-base` · **Versão:** v2 · **Data:** 2026-10-05
- **Tipo:** pagamento de débitos (gerado por `/agent-spec-debt`)
- **Origem:** [../v1/debts.md](../v1/debts.md) · **Spec vigente:** [../v1/spec.md](../v1/spec.md)

---

## 1. Objetivo

Pagar 3 débitos técnicos da v1, sem mudar o comportamento da feature: ao final, os testes
existentes continuam passando.

## 2. Débitos pagos

| Débito   | Severidade | Título                                                      | Task |
| -------- | ---------- | ----------------------------------------------------------- | ---- |
| v1/DT-01 | baixa      | Arquivo gerado pelo Wails fora do .gitignore                | T1   |
| v1/DT-04 | baixa      | Pasta tmp do instalador fora do .gitignore                  | T1   |
| v1/DT-02 | média      | Log config.invalid sem o atributo path exigido pela spec §8 | T2   |

## 3. Tasks

| ID  | Task                                                                     | Fase | Depende de | Status   |
| --- | ------------------------------------------------------------------------ | ---- | ---------- | -------- |
| T1  | [Pagar v1/DT-01 e v1/DT-04: arquivos gerados no .gitignore](tasks/T1.md) | 1    | —          | concluída |
| T2  | [Pagar v1/DT-02: atributo path no log config.invalid](tasks/T2.md)       | 1    | —          | concluída |

### 3.1 Ordem de execução

```
T1 ∥ T2
```

T1 mexe só em `.gitignore` e `README.md`; T2, só no código Go das Preferências e nos seus
chamadores. Nenhum arquivo em comum.

## 4. Critérios de conclusão

- [x] Todas as tasks `concluída`
- [x] Todos os testes automatizados passando: em `dev-tools/`, `go test ./...` e `go test -tags integration ./...`
- [x] Em `dev-tools/`, `go vet -tags integration ./...` sem erros
- [x] Comportamento da feature inalterado (nenhuma API, erro, dado ou mensagem mudou)

## 5. Próximo passo

- **Próxima fase:** Execução → `/agent-spec-run-tasks docs/specs/features/001-arquitetura-base/v2/tasks.md`
