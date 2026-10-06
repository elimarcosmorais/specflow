---
status: rascunho   # rascunho | aprovado
---

# Pre-Refinement — <nome da feature>

> **Legenda:** linha sem marcação = **FATO** (dito pelo usuário) · `[HIPÓTESE]` = inferência
> que precisa ser validada · `[DÚVIDA]` = ponto em aberto (seção 10) ·
> `[fora do escopo do projeto]` = rumo que extrapola o que o projeto se propõe a ser.

- **Feature:** `NNN-<nome-da-feature>` · **Versão:** vN · **Data:** AAAA-MM-DD
- **Fonte da ideia:** texto livre | `<path do arquivo lido>`
- **Anexos:** nenhum | `<path de cada arquivo citado na ideia>`

---

## 1. Ideia em uma frase

<!-- LLM-ONLY: a ideia do usuário reescrita em UMA frase clara, sem ambiguidade. -->

## 2. Problema

- **Qual é a dor real hoje?** (com um exemplo concreto do dia a dia)
- **Quem sente o impacto?**
- **Por que resolver agora?**

## 3. Objetivo

<!-- LLM-ONLY: o resultado esperado ao final (O QUÊ), nunca como implementar. -->

## 4. Público

- **Persona primária:** quem usa ou recebe o valor
- **Persona secundária:** (se houver)
- **Contexto de uso:** onde, quando, em qual dispositivo

## 5. Rumos explorados

> Memória do brainstorm: os ramos da Fase 1 e as direções avaliadas na Fase 2, com o que foi
> escolhido e o que foi descartado.

### Ramo A — <título>

- **A1 — <direção>:** <descrição em 1 frase>.
  _Exemplo:_ <exemplo concreto>. _Viabilidade:_ <reusa X / requer algo novo / `[fora do escopo do projeto]`>
- **A2 — <direção>:** ...
- **Escolhida:** A2 — <motivo, decidido com o usuário>
- **Descartadas / adiadas:** A1 — <motivo>

<!-- LLM-ONLY: repita o bloco para cada ramo aprovado na Fase 1. Toda direção tem exemplo concreto. -->

## 6. Escopo

**Dentro** (direções escolhidas):

- [ ] <item>

**Fora** (com motivo):

- <item> — <motivo: descartado, adiado para outra versão ou fora do escopo do projeto>

## 7. Ancoramento no projeto

- **O que o projeto é** (CLAUDE.md / README): <propósito e domínio relevantes>
- **Features relacionadas** (`docs/specs/features/`): <`NNN-feature/vN` — cobre parte / adjacente / conflita> ou "nenhuma"
- **Capacidades reutilizáveis:** <módulos ou dependências que tornam algum rumo viável>
- **Conflitos ou duplicações:** <descrição> ou "nenhum"

## 8. Premissas e restrições

**Premissas** (suposições para a ideia fazer sentido):

- [HIPÓTESE] <premissa> — <por que foi assumida>

**Restrições já decididas** (prazo, tecnologia imposta, integração obrigatória):

- <restrição, no texto do usuário> ou "nenhuma declarada"

## 9. Riscos

- **Produto / adoção:** <risco> → <mitigação>
- **Escopo** (pode crescer demais?): <risco> → <mitigação>
- **Técnico / segurança:** <risco> → <mitigação>

## 10. Dúvidas em aberto

1. [DÚVIDA] <pergunta objetiva a responder antes da próxima fase>

## 11. Próximo passo

- **Tamanho da mudança:** grande | pequena | trivial — <justificativa em 1-2 frases>
- **Próxima fase:** PRD → `/agent-spec-prd docs/specs/features/NNN-<feature>/vN/pre-refinement.md`
- **ADR a registrar antes:**
  - `/agent-spec-adr "<título da decisão>"` — <o que ela decide, em poucas palavras>

<!-- LLM-ONLY: os comandos ficam prontos para copiar, com o caminho real deste arquivo. Pequena → "Tech Spec → `/agent-spec-tech-spec <caminho deste pre-refinement.md>`"; trivial → "implementar direto", sem comando. Sem ADR → "**ADR a registrar antes:** nenhuma". -->

