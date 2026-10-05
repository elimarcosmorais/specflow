---
status: aprovado # rascunho | aprovado
---

# Task Plan — Arquitetura base do Dev Tools

- **Feature:** `001-arquitetura-base` · **Versão:** v1 · **Data:** 2026-10-05
- **Origem:** [spec.md](spec.md) · [prd.md](prd.md)

---

## 1. Objetivo

O Dev Tools existe em `dev-tools/`:

- o motor Go/Wails com a infraestrutura comum;
- as features Projetos (feature-modelo), Preferências e Sistema, com a API local protegida por
  token e documentada em `/docs`;
- o front React com as telas de Projetos, Preferências e Status;
- o instalador NSIS, que roda sem admin.

Tudo coberto pelos testes unitários e de integração da spec e pelo roteiro manual.

## 2. Fases

| Fase                           | Objetivo                                                                                                       | Tasks      |
| ------------------------------ | -------------------------------------------------------------------------------------------------------------- | ---------- |
| 1 — Esqueleto                  | Projeto Wails + React que compila e abre a janela; valida a hipótese do `wails.json` em `cmd/dev-tools/` (§13) | T1         |
| 2 — Infraestrutura comum       | Caminhos, configuração, logs, idiomas do motor e banco com migrações                                           | T2, T3, T4 |
| 3 — Erros e servidor HTTP      | Erro de domínio, servidor HTTP com token, CORS e erros traduzidos, e `/docs` com o `openapi.yaml`              | T5         |
| 4 — Features no motor          | Preferências completa, núcleo de Projetos e Sistema                                                            | T6, T7, T8 |
| 5 — Projetos no banco e na API | Repositório SQLite e rotas de Projetos, com a integração completa                                              | T9         |
| 6 — Montagem do app            | Camada desktop e `main.go` com a sequência de inicialização; rotas × contrato                                  | T10        |
| 7 — Front base                 | Shell, cliente da API, idioma, tema, telas de Status e Preferências                                            | T11        |
| 8 — Front de Projetos          | Lista, formulário e confirmação de remoção                                                                     | T12        |
| 9 — Instalador e entrega       | Instalador NSIS ajustado, `README.md` e roteiro manual do build de release                                     | T13        |

## 3. Tasks

| ID  | Task                                                                 | Fase | Depende de                         | Status   |
| --- | -------------------------------------------------------------------- | ---- | ---------------------------------- | -------- |
| T1  | [Esqueleto Wails (motor + front)](tasks/T1.md)                       | 1    | —                                  | concluída |
| T2  | [Configuração, caminhos e logs](tasks/T2.md)                         | 2    | T1                                 | concluída |
| T3  | [Idiomas do motor](tasks/T3.md)                                      | 2    | T1                                 | concluída |
| T4  | [Banco de dados e migrações](tasks/T4.md)                            | 2    | T1                                 | concluída |
| T5  | [Erros de domínio, servidor HTTP e documentação da API](tasks/T5.md) | 3    | T3                                 | concluída |
| T6  | [Preferências](tasks/T6.md)                                          | 4    | T2, T5                             | concluída |
| T7  | [Projetos: núcleo](tasks/T7.md)                                      | 4    | T5                                 | concluída |
| T8  | [Sistema: saúde e pasta de logs](tasks/T8.md)                        | 4    | T2, T4, T5                         | concluída |
| T9  | [Projetos: repositório SQLite e API](tasks/T9.md)                    | 5    | T2, T4, T5, T6, T7                 | concluída |
| T10 | [Camada desktop e montagem do motor](tasks/T10.md)                   | 6    | T1, T2, T3, T4, T5, T6, T7, T8, T9 | concluída |
| T11 | [Front: shell, Status e Preferências](tasks/T11.md)                  | 7    | T1, T5, T10                        | concluída |
| T12 | [Front: telas de Projetos](tasks/T12.md)                             | 8    | T9, T11                            | concluída |
| T13 | [Instalador, README e entrega](tasks/T13.md)                         | 9    | T1, T12                            | concluída |

### 3.1 Ordem de execução

```
T1 -> T2 -> T3 -> T4 -> T5 -> (T6 ∥ T7 ∥ T8) -> T9 -> T10 -> T11 -> T12 -> T13
```

T2, T3 e T4 não dependem umas das outras, mas não estão marcadas com `∥`: todas alteram
`dev-tools/go.mod` e `go.sum`.

## 4. Rastreabilidade

### 4.1 Histórias → tasks

| História | Tasks                          |
| -------- | ------------------------------ |
| US-01    | T13                            |
| US-02    | T13                            |
| US-03    | T1, T2, T5, T10, T11, T12, T13 |
| US-04    | T2                             |
| US-05    | T2, T7, T9, T13                |
| US-06    | T4, T5, T8, T10, T11, T12      |
| US-07    | T6, T10, T11                   |
| US-08    | T3, T5, T6, T10, T11, T12      |
| US-09    | T5, T7, T9, T12                |
| US-10    | T7, T9, T12                    |
| US-11    | T7, T9, T12                    |
| US-12    | T7, T9, T12                    |
| US-13    | T5, T6, T9, T10                |
| US-14    | T2, T5, T7, T8, T9, T11        |
| US-15    | T7                             |
| US-16    | T4, T7, T8, T9                 |

US-03, US-06, US-08 e US-14 passam por muitas tasks porque são transversais: abrir o app,
saúde, idioma e logs atravessam todas as camadas, do motor ao front.

### 4.2 Testes → tasks

| CT           | CA                                | Task |
| ------------ | --------------------------------- | ---- |
| CT-01        | CA-10                             | T2   |
| CT-02        | CA-10                             | T2   |
| CT-03        | CA-09                             | T2   |
| CT-04        | CA-06                             | T2   |
| CT-05        | — (D-08)                          | T2   |
| CT-06        | CA-16, CA-17                      | T3   |
| CT-07        | CA-18                             | T3   |
| CT-08        | CA-40                             | T4   |
| CT-09        | CA-12                             | T4   |
| CT-10        | CA-41                             | T4   |
| CT-11        | CA-35                             | T2   |
| CT-12        | CA-37                             | T2   |
| CT-13        | CA-37                             | T2   |
| CT-14        | CA-33                             | T5   |
| CT-15        | CA-13                             | T5   |
| CT-16        | — (D-01)                          | T5   |
| CT-17        | CA-18, CA-23                      | T5   |
| CT-18        | CA-35                             | T5   |
| CT-19        | —                                 | T5   |
| CT-20        | CA-08, CA-32                      | T5   |
| CT-21        | CA-07                             | T5   |
| CT-22        | CA-34                             | T5   |
| CT-23        | CA-32                             | T10  |
| CT-24        | CA-33                             | T10  |
| CT-25        | CA-07                             | T10  |
| CT-26        | CA-07, CA-12                      | T10  |
| CT-27        | CA-05                             | T10  |
| CT-28        | CA-14, CA-18                      | T10  |
| CT-29        | CA-14, CA-16, CA-17               | T6   |
| CT-30        | CA-14                             | T6   |
| CT-31        | CA-15, CA-18                      | T6   |
| CT-32        | — (RN-09, RN-10)                  | T6   |
| CT-33        | CA-15, CA-18                      | T6   |
| CT-34        | CA-33                             | T6   |
| CT-35        | CA-10, CA-19                      | T7   |
| CT-36        | CA-10                             | T7   |
| CT-37        | CA-23                             | T7   |
| CT-38        | CA-39                             | T7   |
| CT-39        | CA-20                             | T7   |
| CT-40        | CA-20                             | T7   |
| CT-41        | CA-21                             | T7   |
| CT-42        | CA-22                             | T7   |
| CT-43        | CA-25, CA-26                      | T7   |
| CT-44        | CA-27                             | T7   |
| CT-45        | CA-28                             | T7   |
| CT-46        | CA-29                             | T7   |
| CT-47        | CA-39                             | T7   |
| CT-48        | CA-30                             | T7   |
| CT-49        | CA-19, CA-27, CA-30               | T7   |
| CT-50        | CA-35                             | T7   |
| CT-51        | CA-38                             | T7   |
| CT-52        | CA-19, CA-25, CA-27, CA-30, CA-40 | T9   |
| CT-53        | CA-20, CA-40                      | T9   |
| CT-54        | CA-21, CA-22, CA-40               | T9   |
| CT-55        | CA-23, CA-40                      | T9   |
| CT-56        | CA-28, CA-29, CA-40               | T9   |
| CT-57        | CA-26, CA-40                      | T9   |
| CT-58        | CA-33, CA-40                      | T9   |
| CT-59        | CA-24, CA-41                      | T9   |
| CT-60        | CA-10                             | T9   |
| CT-61        | CA-35                             | T9   |
| CT-62        | — (RN-13, RN-14)                  | T9   |
| CT-63        | CA-11, CA-12                      | T8   |
| CT-64        | CA-36                             | T8   |
| CT-65        | CA-11, CA-13, CA-40, CA-41        | T8   |
| CT-66        | CA-01                             | T13  |
| CT-67        | CA-02                             | T13  |
| CT-68        | CA-03                             | T13  |
| CT-69        | CA-04                             | T13  |
| CT-70        | CA-05                             | T10  |
| CT-71        | CA-06                             | T10  |
| CT-72        | CA-07                             | T11  |
| CT-73        | CA-08                             | T12  |
| CT-74        | CA-10                             | T13  |
| CT-75        | CA-11, CA-12                      | T12  |
| CT-76        | CA-14, CA-15                      | T11  |
| CT-77        | CA-16, CA-17, CA-18               | T12  |
| CT-78        | CA-23                             | T12  |
| CT-79        | CA-30, CA-31                      | T12  |
| CT-80        | CA-36                             | T11  |
| CT-81        | CA-32                             | T10  |
| CT-82        | CA-34                             | T10  |
| CT-83 (novo) | —                                 | T1   |

## 5. Critérios de conclusão da feature

- [x] Todas as tasks `concluída`
- [x] Todos os testes automatizados passando: em `dev-tools/`, `go test ./...` e `go test -tags integration ./...`
- [ ] Testes manuais executados conforme o roteiro (CT-66 a CT-82 e CT-83), os de CT-66 a CT-82 no build de release (`wails build -nsis`)
- [x] `go vet -tags integration ./...` sem erros; em `dev-tools/frontend`, `npm run typecheck` sem erros; em `dev-tools/cmd/dev-tools`, `wails build -nsis` sem erros

## 6. Próximo passo

- **Próxima fase:** Execução → `/agent-spec-run-tasks docs/specs/features/001-arquitetura-base/v1/tasks.md`
