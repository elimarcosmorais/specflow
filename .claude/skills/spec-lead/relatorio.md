<!-- LLM-ONLY: template do relatório final do /spec-lead. Preencha e mostre no chat, não salve em
arquivo. Seção sem conteúdo some, menos "Precisa de você", que vira "Precisa de você: nada." —
essa linha é a informação mais importante do relatório. Nenhuma linha é livre: toda linha vem de
um campo de um agent, da triagem do débito (Etapa 6) ou da Etapa 7, item 2. -->

## <pedido em uma frase>

**Status:** concluído · <N> passos, <M> planos secundários
<!-- ou: parado no passo <N> de `<plano>` — <qual dos quatro motivos da lista de paradas> -->

### O que foi entregue

- <1 a 4 linhas do que o usuário ganhou, do ponto de vista dele; não a lista de arquivos>
- **Arquivos:** <N> criados, <M> modificados (`git diff --staged --stat`)
- **Specs:** `<arquivo>` (criado | atualizado | subdividido) | nenhuma
- **ADRs:** <NNN — título> (registrada) | nenhuma

### Provas

| O que | Comando | Resultado |
|---|---|---|
| suíte | `<comando>` | <N> passando, 0 falhando, 0 pulados |
| build | `<comando>` | ok |
| lint | `<comando>` | ok |
| tela | <teste de tela ou MCP> | <o que viu> |

**Sem prova:** <o que ficou> — <por quê> | nada

<!-- LLM-ONLY: "0 pulados" é literal (behavior.md, regra 3). Pulo que você não conseguiu
resolver entra em "Sem prova", com o nome do teste e o motivo, e nunca na linha de verde. -->

### Decisões que tomei

| Decisão | Escolhi | Por quê |
|---|---|---|
| <o que precisou ser decidido> | <o que foi escolhido> | <motivo em 1 linha; alternativa descartada> |

<!-- LLM-ONLY: é aqui que o usuário audita a autonomia em vez de ser interrompido por ela. Entra
toda escolha que mudaria o resultado se fosse outra — sua, do executor ou do planner. Não entra
o que o código, uma spec ou uma ADR já respondiam. -->

### Resolvi sozinho, fora do pedido

- <o que estava errado> — `<arquivo>` · <prova> (plano secundário <n> | corrigido no passo <N>)

<!-- LLM-ONLY: o débito que virou trabalho (behavior.md, regra 2, destinos 1 e 2). Sem teto. -->

### Precisa de você

<!-- LLM-ONLY: SÓ o que você não pode fazer — um dos quatro motivos da lista de paradas
(behavior.md, regra 1). Cada item traz o que fazer e o comando pronto. Item que você poderia ter
resolvido e está aqui é falha da triagem (Etapa 8). Nada aqui → escreva literalmente
"Precisa de você: nada." -->

- **Conferir:** <CA-N (manual: …)> → <roteiro em 1 linha>
- **Decidir:** <o que muda no que o produto entrega> · Recomendo: <…> → `/spec-lead "<pedido pronto>"`
- **Fazer:** <exige UAC, 2FA, conta externa ou aparelho> → `<comando pronto>`
- **Risco:** <segredo, credencial, dado pessoal ou de pagamento exposto, `arquivo:linha`> · Recomendo: <…> → `<comando>`

<!-- LLM-ONLY: a linha "Risco" aparece mesmo que já existam outras; risco nunca é descartado
(behavior.md, regra 2). -->

### Estado

Código e specs em **stage, sem commit**. Revise com `git diff --staged` e commite com
`/spec-semantic-commit`.

**Planos:** apagados (remoção em stage) | mantido `<arquivo>` — faltam os passos <N, M>; retome
com `/spec-lead continuar`.
