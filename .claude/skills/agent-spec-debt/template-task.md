---
status: pendente   # pendente | em andamento | concluída
---

# T1 — Pagar v1/DT-01: <título>

- **Fase:** 1 — Pagamento de débitos
- **Depende de:** — | T1 (toca `<arquivo>` também)
- **Débitos:** v1/DT-01 (média, origem: revisão da T3 da v1) — [debts.md](../../v1/debts.md)
- **Spec vigente:** [spec.md](../../v1/spec.md) §<seções que a correção precisa respeitar>

<!-- LLM-ONLY: "Débitos" lista todos os débitos desta task (vários quando estão no mesmo arquivo). Os links apontam para a versão de origem do débito e para a versão da spec vigente. -->

---

## 1. Objetivo

<1 a 2 linhas: o que fica melhor no código ao final desta task, sem mudar comportamento.>

## 2. O que fazer

- <a correção, concreta: remover o teste `X` de `Y`, que duplica o CT-14>

**Não fazer:**

- Corrigir outros problemas do mesmo arquivo que não estão nesta task. Achou outro? Cite-o nas
  "Pendências" da sua resposta: a revisão o registra como novo débito no `debts.md` desta versão.
- Alterar comportamento: resposta de API, erro, mensagem, dado gravado. Se a correção exigir
  isso, pare e avise: não é débito, é mudança.

<!-- LLM-ONLY: a correção vem do "Correção sugerida" do débito, refinada pela leitura do código. Nível de spec, não de código: o que remover, renomear ou mover, e para onde. -->

## 3. Arquivos

| Arquivo | Ação | O que muda |
|---------|------|------------|
| `<caminho>` | modificar | <a correção neste arquivo> |

<!-- LLM-ONLY: se a correção muda algo usado em outros arquivos (ex.: renomear uma função pública), os chamadores e testes afetados entram aqui. -->

## 4. Testes

Nenhum teste novo: a suíte existente prova que o comportamento não mudou.

<!-- LLM-ONLY: quando o débito é um teste faltando, troque a frase acima por "**Arquivo de teste:** `<caminho>` · **Padrão:** <da spec §10>" e uma tabela com o CT novo, no mesmo formato da spec §10.2 (CT | Camada | O que verifica | Entrada / cenário | Resultado esperado | CA), com "(novo)" depois do ID. -->

## 5. ADRs que se aplicam

- [ADR NNN](../../../../../adr/NNN-<titulo>.md) — <a regra concreta que esta task deve seguir>

<!-- LLM-ONLY: só as ADRs que governam os arquivos da seção 3, com a regra concreta tirada da seção Decisão da ADR. Nenhuma → "Nenhuma". -->

## 6. Pronta quando

- [ ] Correção da seção 2 aplicada, só nos arquivos da seção 3
- [ ] Testes passando: `<comando de testes do projeto>` (nenhum teste alterado além dos citados
      na seção 2)
- [ ] `<comando de build/lint do projeto>` sem erros
