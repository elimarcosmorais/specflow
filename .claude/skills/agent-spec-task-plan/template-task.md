---
status: pendente   # pendente | em andamento | concluída
---

# T1 — <nome da task>

- **Fase:** 1 — <nome da fase>
- **Depende de:** — | T1 (usa `<símbolo>`), T2 (usa `<símbolo>`)
- **Entrega para outras tasks:** `<tipo, função ou tabela que outras tasks usam>` | nada
- **Histórias:** US-01 · **Critérios de aceite:** CA-01, CA-02
- **Spec:** [spec.md](../spec.md) §3.1, §5, §6.1
- **Design:** [design.md](../design.md) §3.1, §4 · [design-system.md](../../../../design-system.md) §7, §8 | —
- **Risco:** baixo | alto — <áreas: autenticação, autorização, pagamentos…>

<!-- LLM-ONLY: "Depende de" e "Entrega para outras tasks" são o que prova a ordem (regra 5): todo símbolo que esta task usa de outra task aparece em "Depende de", com a task que o cria. "Spec" lista só as seções que o executor precisa consultar para esta task; "Design" idem, só em task que cria ou altera tela ou componente visual (as outras: "—"). "Risco" segue as "Regras de risco" do SKILL.md: `baixo` fica sozinho; `alto` traz as áreas que o justificam. -->

---

## 1. Objetivo

<1 a 3 linhas: o resultado técnico ao final desta task.>

## 2. O que fazer

- <passo concreto: criar `X` com `Y`, seguindo §3.1 da spec>
- <passo concreto: alterar `Z` para…>
- <regra ou decisão que vale aqui, citando o ID: RN-02, D-01, §5>

<!-- LLM-ONLY: nível de spec, não de código: nomes, assinaturas, códigos de erro, ordem dos passos. Não escreva o corpo das funções. -->

## 3. Arquivos

| Arquivo | Ação | O que muda |
|---------|------|------------|
| `<caminho>` | criar | <o que contém> |
| `<caminho>` | modificar | <o que muda> |
| `<caminho>` | consultar | <por que ler: padrão a seguir, schema…> |

<!-- LLM-ONLY: se esta task muda a assinatura de algo público, liste aqui os chamadores e testes que precisam acompanhar: tocar neles faz parte da task. Inclua os arquivos de teste. -->

## 4. Testes

**Arquivo de teste:** `<caminho>` · **Padrão:** <framework e convenção do projeto, da spec §10>

| CT    | Camada   | O que verifica | Entrada / cenário | Resultado esperado | CA    |
|-------|----------|----------------|-------------------|--------------------|-------|
| CT-01 | unitário | <regra>        | <dados>           | <saída exata>      | CA-01 |

<!-- LLM-ONLY: copie as linhas da spec §10.2 dos CTs desta task (a task é autocontida). Resultado esperado sempre concreto: valor, código de erro ou status, nunca "funciona" ou "correto". CT criado nesta fase leva "(novo)" depois do ID. -->

## 5. ADRs que se aplicam

- [ADR NNN](../../../../../adr/NNN-<titulo>.md) — <a regra concreta que esta task deve seguir>

<!-- LLM-ONLY: só as ADRs que governam os arquivos da seção 3, com a regra concreta tirada da seção Decisão da ADR. Nenhuma → "Nenhuma". -->

## 6. Pronta quando

- [ ] Arquivos da seção 3 criados e alterados conforme a seção 2
- [ ] Testes da seção 4 implementados e passando: `<comando de testes do projeto>`
- [ ] `<comando de build/lint do projeto>` sem erros
- [ ] Nenhum teste existente quebrado
