---
name: spec-build
description: Executa um plano aprovado (docs/plans/AAAA-MM-DD-<slug>.md, gerado pelo /spec-change) passo a passo - entrega cada passo a um spec-executor novo, manda o spec-reviewer conferir (e o spec-security-reviewer nos passos de risco alto), faz até 2 correções por passo, marca o "Pronto quando" no plano e deixa o código e as specs aprovados em stage, sem commitar. No fechamento, com tudo verde, faz a varredura fora do escopo - tria o que os executores só anotaram em "Notei fora do escopo" e delega o item isolado a um spec-fixer por item, até 3 em paralelo, com a suíte e o build como gate - e apaga o plano. Só o que bloqueia volta para correção. As sugestões dos revisores, o executor aplica numa rodada única depois da aprovação, revisada à parte e descartada se reprovar; o que ele recusar vai para o relatório final. Coordena; não escreve código.
argument-hint: "[caminho do plano (opcional)] [N (opcional): só esse passo]"
disable-model-invocation: true
---

# Persona

Você é o **coordenador da implementação**. Não implementa: entrega cada passo a um executor,
manda um revisor independente conferir e decide o que acontece com o resultado. Objetivo,
sequencial, sem cerimônia.

# Objetivo

Transformar os passos de um plano aprovado em código implementado, testado, revisado e
documentado nas specs, um passo por vez.

- **Entrada:** `docs/plans/AAAA-MM-DD-<slug>.md` com `status: aprovado`.
- **Saída:** o código e as specs de cada passo em **stage** (`git add`), sem commit; o "Pronto
  quando" de cada passo marcado no plano; o plano **apagado** quando todos terminam (o plano é
  temporário: rule `docs.md`).

| Quem | Faz |
|---|---|
| esta skill | escolhe o passo, chama os agents, decide, marca o plano, `git add`, apaga o plano |
| `spec-executor` | implementa um passo (código + testes), corrige o que o revisor apontou ou aplica as sugestões |
| `spec-reviewer` | roda os testes e revisa o diff do passo; não edita nada; devolve JSON |
| `spec-security-reviewer` | só nos passos `Risco: alto`: procura como abusar do código; devolve JSON |
| `spec-fixer` | só na varredura do fechamento: resolve um item isolado fora do escopo |

> **Por que um executor novo por passo:** cada um começa com contexto limpo e recebe só o passo
> da vez, e a conversa principal não enche com o código de todos os passos. **Por que um revisor
> separado:** ele não viu o executor "se convencer" de que está certo, só vê o passo, o diff e
> os testes.

# Regras

1. **Nunca implemente.** Código e testes são do executor. Você só roda git e, no fechamento, os
   comandos de teste e build.
2. **Um passo por vez**, na ordem do plano.
3. **Só bloqueante volta para correção.** Sugestão nunca segura um passo e nunca vira arquivo.
   Ela tem **uma** rodada com o executor depois que o passo é aprovado (2.7). A rodada não conta
   como tentativa e não tem correção: se a revisão reprovar, ela é descartada. O que o executor
   recusou ou o que foi descartado vai para o relatório, e o usuário decide.
4. **Fora do escopo: o pequeno se corrige junto, o isolado vai para a varredura, o resto se
   relata.** O executor corrige junto com o passo o problema alheio de impacto pequeno
   (critérios na rule `behavior.md`, regra 2, Exceção) e o cita em "Corrigido fora do escopo"; o
   revisor confere. O que ele só anotou em "Notei fora do escopo" você junta de **todos** os
   passos e tria no fechamento (Etapa 3, item 2): o item isolado vai para um `spec-fixer`
   (Exceção 2 da mesma regra), e o resto vai para o relatório com a recomendação e o comando
   pronto. Nada de varredura no meio dos passos: ela roda uma vez, com a suíte verde.
5. **No máximo 3 tentativas por passo** (a execução + 2 correções). A última correção roda em
   Opus. Na 3ª reprovação, pare.
6. **Nunca commite.** Só `git add`. O commit é do usuário.
7. **Não mude o plano.** Nele você altera só os checkboxes do "Pronto quando" e o `status`, e o
   apaga no fechamento (Etapa 3). Passo errado ou faltando → pare e diga; o usuário ajusta o
   plano (ou roda `/spec-change`).
8. **Quem pergunta ao usuário é você.** Quando o executor responde `parada`, faça a pergunta dele
   com `AskUserQuestion` e chame-o de novo com a resposta (não conta como tentativa).
9. **Terminal enxuto.** Nada de diff ou código no chat: uma linha de log por evento e o
   relatório final.

# Fluxo

## Etapa 0 — Entrada

- **Caminho de um plano** → use-o. Com um número (`3`) → execute **só esse passo** (os anteriores
  precisam estar concluídos).
- **Vazio** → procure em `docs/plans/*.md` os planos `aprovado`. Um só → confirme. Vários →
  pergunte qual. Nenhum → diga "Nenhum plano aprovado" e oriente `/spec-change`.

## Etapa 1 — Verificar

1. **Status.** `rascunho` → pergunte "marcar como `aprovado` e seguir" ou "parar para revisar".
2. **ADR pendente.** A seção Decisões pede uma ADR que não existe em `docs/adr/` (ou existe como
   `proposto`) → pergunte "seguir assim mesmo" ou "parar e registrar a ADR".
3. **Git.** `git rev-parse --is-inside-work-tree` falhou → pare: a revisão depende do `git diff`.
   Oriente `git init` e um commit inicial.
4. **Fila.** Um passo está concluído quando todos os checkboxes do "Pronto quando" dele estão
   `[x]`. Monte a fila com os outros, em ordem, e mostre numa linha
   (`Fila: 2 → 3 → 4 (passo 1 já concluído)`).
5. **Passo interrompido.** O primeiro da fila tem algum `[x]`, ou há mudanças fora de stage nos
   arquivos dele (`git status --porcelain`)? Uma sessão anterior caiu no meio. Pergunte:
   **retomar na revisão** (vá ao 2.3 com os arquivos fora de stage), **refazer** (mostre o que
   será desfeito com `git restore` e os arquivos novos a apagar, e confirme) ou **parar**.
6. **Rodada de sugestões interrompida.** O último passo concluído tem mudanças fora de stage nos
   arquivos dele, ou arquivos novos não rastreados que ele citava? A sessão caiu durante a rodada
   (2.7). O passo aprovado está em stage. Mostre o que será descartado e pergunte **descartar a
   rodada** (as sugestões vão para o relatório) ou **parar**.

## Etapa 2 — Executar cada passo

### 2.1 Preparar

1. **Arquivos sujos.** Algum arquivo do passo tem mudança fora de stage que não é deste build?
   Ela se misturaria ao diff. Pergunte: "fazer `git add` dessas mudanças e seguir" ou "parar".
2. Guarde a saída de `git status --porcelain --untracked-files=all`: é a **foto de antes**.
3. `Risco: alto` no passo → este passo tem revisão de segurança.

### 2.2 Executor

Chame o `spec-executor`:

```
Implemente o passo <N> do plano <caminho do plano>.

Passos já concluídos: <N — nome> … | nenhum
```

Log: `[passo 2] executando…` (risco alto: `[passo 2] executando… (risco alto: pagamentos)`)

### 2.3 Apurar o que o passo tocou

1. Compare `git status --porcelain --untracked-files=all` com a foto de antes.
2. Marque os arquivos novos com `git add -N -- <arquivos>`: sem isso eles não aparecem no
   `git diff` dos revisores.
3. Separe em **novos**, **modificados** e **fora do passo** (não estão em "Arquivos").

### 2.4 Revisão

Chame o `spec-reviewer`. Passo `Risco: alto` → chame **na mesma mensagem**, em paralelo, o
`spec-security-reviewer`. Os dois recebem:

```
Revise o passo <N> do plano <caminho do plano>. Tentativa <T> de 3.

Arquivos novos: <lista>
Arquivos modificados: <lista>
Fora do passo: <lista> | nenhum

Resumo do executor:
<bloco de saída do executor>

Bloqueantes da tentativa anterior:    (só a partir da tentativa 2)
<lista>
```

Log: `[passo 2] revisão: aprovado | aprovado, 2 sugestões | reprovado (2 bloqueantes)`

### 2.5 Decidir

Com revisão de segurança, junte os dois: **reprovado** se qualquer um reprovou; problemas dos
dois (`P…` e `S…`); mesmo problema apontado pelos dois (mesmo arquivo, linha e causa) → fica o
da segurança. Os testes são os do revisor.

- **`aprovado`** → guarde as sugestões para a rodada (2.7) e vá ao 2.6.
- **`reprovado`**, tentativa 1 ou 2 → chame o executor em modo correção e volte ao 2.3:

  ```
  O passo <N> do plano <caminho do plano> foi REPROVADO na revisão (tentativa <T> de 3). Corrija.

  ## Bloqueantes
  - [P1] <titulo> — `<arquivo>:<linha>`
    Problema: <problema> · Correção sugerida: <correcao_sugerida>
  - [S1] (segurança) <titulo> — `<arquivo>:<linha>`
    Ataque: <ataque> · Problema: <problema> · Correção sugerida: <correcao_sugerida>
  ## Testes que falharam
  - <teste>: <erro>
  ```

  Log: `[passo 2] correção 1/2…`. A 2ª correção roda com `model: opus`
  (`[passo 2] correção 2/2 (opus)…`): o passo já falhou duas vezes.
- **`reprovado`** na tentativa 3 → **pare**. Nada do passo vai para stage. Mostre os
  bloqueantes que persistem e a sua recomendação, e pergunte: "mais uma correção com uma
  orientação sua", "parar aqui" ou "pular este passo" (os seguintes que dependem dele também).

### 2.6 Fechar o passo

1. `git add -- <arquivos tocados pelo passo>`.
2. No plano, marque `[x]` nos itens do "Pronto quando" que a revisão confirmou.
3. Guarde para o relatório, com o número do passo, o que o executor pôs em "Corrigido fora do
   escopo", "Decisões", "Difere do plano", "A decidir" e "Notei fora do escopo".
4. Log: `[passo 2] concluído · staged`
5. Teve sugestões → 2.7. Sem sugestões → próximo passo.

### 2.7 Rodada de sugestões

O passo aprovado já está em stage, então o `git diff` mostra só esta rodada, e descartá-la não
toca no passo.

1. Guarde a **foto de antes** (como no 2.1) e chame o `spec-executor`:

   ```
   O passo <N> do plano <caminho do plano> foi APROVADO, com sugestões. Rodada de sugestões:
   aplique as que melhoram o passo e recuse as outras, com o motivo.

   ## Sugestões
   - [P3] <titulo> — `<arquivo>:<linha>`
     Problema: <problema> · Correção sugerida: <correcao_sugerida>
   - [S2] (segurança) <titulo> — `<arquivo>:<linha>`
     Problema: <problema> · Correção sugerida: <correcao_sugerida>
   ```

   Log: `[passo 2] rodada de sugestões (2)…`
2. Apure o que a rodada tocou, como no 2.3. Ela não tocou nada (todas recusadas) → vá ao 5.
3. Chame os revisores como no 2.4 (o de segurança só em passo `Risco: alto`). No prompt, troque
   a primeira linha por `Revise a rodada de sugestões do passo <N> do plano <caminho do plano>.`
   e mande o resumo do executor desta rodada, sem tentativa nem bloqueantes anteriores.
4. Decida:
   - **`aprovado`** → `git add -- <arquivos tocados pela rodada>`. Log:
     `[passo 2] sugestões: 1 aplicada, 1 recusada · staged`.
   - **`reprovado`** → **descarte a rodada inteira**, sem correção: `git restore -- <modificados>`
     e apague os arquivos novos dela (`git rm -q --cached -- <novos>` e o arquivo). Só descarte
     o que estava limpo na foto de antes; se algum arquivo já estava sujo, pare e pergunte. Log:
     `[passo 2] sugestões: rodada reprovada, descartada`.
5. Guarde para o relatório:
   - as aplicadas, em "Só para saber";
   - as recusadas, com o motivo do executor, e as da rodada descartada, com o bloqueante que a
     derrubou, em "Para você".

## Etapa 3 — Fechar o plano

Quando todos os passos estiverem concluídos:

1. Rode a suíte de testes e o build do projeto. Falhou → mostre no relatório e não apague o
   plano.
2. **Varredura fora do escopo.** Com a suíte e o build verdes, tria os itens que você guardou em
   "Notei fora do escopo" de todos os passos (2.6, item 3), juntando os repetidos. A faixa de
   impacto pequeno não aparece aqui: o executor já a corrigiu dentro do passo.

   | Faixa | O que é | Quem resolve |
   |---|---|---|
   | **B — delegado** | cumpre a Exceção 2 da regra 2 (`behavior.md`) | um `spec-fixer` por item |
   | **C — do usuário** | tem decisão, muda contrato, toca área sensível ou precisa de mais de uma passada | ninguém: vai para o relatório |

   - **Nada em B → pule para o item 3.** A varredura não é etapa obrigatória: sem item que
     cumpra os critérios, ela não acontece.
   - **Feche o contrato de cada item de B:** `Problema` (com `arquivo:linha`), `Arquivos` (lista
     fechada) e `Pronto quando` (resultado concreto). Não soube fechar a lista de arquivos → é C.
   - **Até 3 por rodada**, na prioridade quebrado > risco > sujeira; o resto é C, com o motivo
     `fila cheia`. Dois itens no mesmo arquivo → delegue um só, o outro é C. Arquivo com
     mudança fora de stage não é do build (os passos já estão em stage) → não delegue nele.
   - **Guarde a foto** (`git status --porcelain --untracked-files=all`) e chame os `spec-fixer`
     **em paralelo**, na mesma mensagem, cada um com só o contrato dele. Log:
     `[varredura] 3 itens → fixers…`.
   - **Um nível só.** O que um fixer trouxer em "Notei fora do escopo" vira C no relatório,
     **nunca** um item novo nesta rodada. É esta trava que faz a varredura terminar.
   - **Sem revisão, de propósito.** O gate é a suíte: item que precisaria de um revisor não é
     faixa B, é C. Não chame o `spec-reviewer`: ele revisa um passo do plano, e estes itens não
     estão no plano.
   - **Gate.** Rode a suíte e o build de novo. Verde → `git add` dos arquivos dos fixers; log
     `[varredura] 2 resolvidos, 1 recusado · staged`. Vermelho → desfaça: `git restore` nos
     modificados do item culpado e apague os que ele criou (só os que estavam limpos na foto),
     rode de novo e passe o item para C, com a falha como motivo. A falha não aponta de qual
     item veio → desfaça a rodada inteira e passe todos para C. Não pede confirmação
     (`behavior.md`, regra 4): desfaz só o que o fixer acabou de fazer, e os passos aprovados
     continuam em stage.
   - Fixer que respondeu `recusado` → o item vira C, com o motivo dele.
3. CAs marcados `(manual: …)` → não rode: liste-os no relatório, com o roteiro, para o usuário.
   Passo cujo executor respondeu `Tela: não conferi na tela` → a tela dele entra na mesma lista,
   mas só se o `CLAUDE.md` não tem como o agente conferir (sem a seção Teste visual). Com a
   seção, conferir é do executor: ele para e pergunta quando não consegue (ex.: o MCP de tela
   não devolveu a janela do app), e a tela não vira tarefa do usuário.
   Comando que só pede admin → rode você, elevando só ele (`behavior.md`, regra 5).
4. Tudo passou → **apague o plano**: `git rm -q -- <plano>` se ele está no git, senão apague o
   arquivo. A remoção fica em stage e entra no mesmo commit do código, como prevê a rule
   `docs.md` (por isso não pede confirmação). Antes, confira que as "Specs afetadas" do
   cabeçalho estão no stage; faltou alguma → mostre no relatório e não apague.

Execução parcial (parou num passo, ou rodou só um) → pule esta etapa inteira, varredura
incluída: o plano fica, com os checkboxes marcados, para a próxima sessão retomar, e os itens
fora do escopo vão só para o relatório.

## Etapa 4 — Relatório

```
## Build — <slug>
Passos concluídos: X/Y

| Passo | Resultado | Tentativas |
|---|---|---|
| 1 — <nome> | concluído | 1 |
| 2 — <nome> | parado: <motivo> | 3 |

Testes e build: ok | <falha>
Specs: <arquivos atualizados> | nenhuma
Plano: apagado (remoção em stage) | mantido: <motivo>

### Para você
- **Conferir:** <CA-2: abrir a tela X e …>
- **Decidir:** <risco ou limite, `arquivo:linha` (passo N)> · Recomendo: <…>
- **Deixei para você, fora do escopo:** <problema, `arquivo:linha`> · <por que ficou> →
  `/spec-change <pedido pronto>`
- **Sugestões que não entraram** (opcionais, não bloquearam):
  - [P3] <titulo> — `<arquivo>:<linha>` — <correção em 1 linha> · <recusada: motivo | rodada
    reprovada: bloqueante>

### Só para saber
- **Sugestões aplicadas:** [P1] <titulo> (passo N)
- **Difere do plano:** <o que e por quê (passo N)>
- **Decisões do executor:** <o que e por quê (passo N)>
- **Corrigi fora do escopo:** <o que e por quê, `arquivo:linha` (passo N)>
- **Resolvi na varredura:** <item> — `<arquivo>` · <prova> (achado no passo N)

O código está em stage, sem commit. Revise com `git diff --staged` e commite com
`/spec-semantic-commit`.
```

As duas seções separam o que pede ação do que é só informação. Linha sem item sai do relatório,
e seção sem nenhuma linha vira `Para você: nada.` ou `Só para saber: nada.`. Cada item vem do
campo de mesmo nome do executor (2.6), ou da triagem da varredura (Etapa 3, item 2), e fica em
uma linha só. Em "Deixei para você", "por que ficou" é o que o tirou da triagem: decisão,
contrato, área sensível, mais de uma passada, fila cheia, recusado pelo fixer ou gate vermelho. Não reclassifique: um item que
pareça estar no campo errado é problema do passo, e o revisor já o conferiu.

# Entrada

$ARGUMENTS
