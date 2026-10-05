---
status: rascunho   # rascunho | aprovado
---

# Tech Alignment — <nome da feature>

- **Feature:** `NNN-<nome-da-feature>` · **Versão:** vN · **Data:** AAAA-MM-DD
- **Origem:** [prd.md](prd.md) · [pre-refinement.md](pre-refinement.md)
- **ADRs consultadas:** `docs/adr/NNN-<titulo>.md` | nenhuma

---

## 1. Contexto técnico

<1 a 3 parágrafos curtos: o problema reescrito do ponto de vista técnico, o que o projeto já
decide (stack, ADRs) e o que ficou em aberto para esta fase.>

**Invariantes** (o que qualquer implementação precisa preservar):

- <invariante> (RN-01, CA-03)

<!-- LLM-ONLY: declarações técnicas diretas, sem narrativa ("o usuário deve poder..."). Cada invariante cita os IDs do PRD de onde vem. -->

## 2. Decisões

### D-01 — <ponto de decisão>

- **Decisão:** <solução escolhida> — <por quê>.
- **Alternativas avaliadas:**
  - **<alternativa>** — <exemplo concreto>. _Prós:_ … _Contras:_ … Não escolhida porque <motivo>.
- **Viabilidade:** <reusa X | requer Y novo `[HIPÓTESE]` | conflita com ADR-NNN>
- **Trade-off aceito:** <o que se ganha e o que se abre mão>
- **Atende:** US-01, RN-02, CA-05
- **ADR:** candidata — "<título da decisão>" | registrada em `docs/adr/NNN-<titulo>.md` | não

<!-- LLM-ONLY: decisão direta (cravada pelo usuário ou óbvia) não tem leque: troque "Alternativas avaliadas" por "Decisão direta — <motivo>". Nada de endpoints, tabelas, campos ou arquivos: isso é da Spec. Sem pontos de decisão, escreva "N/A — a feature segue os padrões <X/Y>; ver Restrições herdadas". -->

## 3. Restrições herdadas

<O que a implementação precisa respeitar e que não foi decidido aqui: stack, ADRs, padrões do
projeto e imposições do PRD.>

- <restrição> — <origem: ADR-NNN | CLAUDE.md | PRD seção 10 | padrão da feature NNN>

## 4. Pontos em aberto

**Para a Spec** (detalhe de implementação, fora da altitude desta fase):

- <ponto> (D-01)

**Dependências de produto** (o PRD não decidiu; não foram decididas aqui):

- <dependência>

**Fora do escopo** (observado, mas não vira proposta nesta feature):

- <observação>

## 5. Próximo passo

- **Próxima fase:** Tech Spec → `/agent-spec-tech-spec docs/specs/features/NNN-<feature>/vN/tech-alignment.md`
- **ADR a registrar antes:**
  - `/agent-spec-adr "D-NN de docs/specs/features/NNN-<feature>/vN/tech-alignment.md"` — <título da decisão>

<!-- LLM-ONLY: uma linha por decisão marcada como "ADR: candidata", com o caminho real deste arquivo. Sem candidatas → "**ADR a registrar antes:** nenhuma". -->

