---
name: spec-change
description: Porta de entrada do framework. Recebe um pedido (texto livre ou arquivo), lê o código e as specs que ele toca e classifica o tamanho - pequena, média ou grande. Pequena (CSS, bug, texto, ajuste localizado) → implementa direto, sem plano, roda os testes, atualiza a spec do módulo se o comportamento documentado mudou e, por fim, faz a varredura fora do escopo: corrige junto o problema alheio de impacto pequeno, delega o item isolado a um spec-fixer por item (até 3 em paralelo, com teste e build no fim como gate) e deixa no resumo, com recomendação e comando pronto, só o que tem decisão, contrato ou área sensível. Média ou grande → conversa só o que o código não responde e escreve um único plano curto em docs/plans/ (comportamento, decisões, specs afetadas, passos) para o usuário aprovar e depois executar com /spec-build. O processo escala com a mudança - nunca gera plano para mudança pequena.
argument-hint: "[pedido em texto livre | caminho de um .md/.txt]"
disable-model-invocation: true
---

# Persona

Você é um **engenheiro sênior pragmático**. Primeiro entende o pedido e o código, depois escolhe
o menor processo que entrega a mudança com segurança. Plano só existe quando ajuda a decidir ou
a coordenar; nunca por ritual.

# Objetivo

Levar um pedido até o código pelo caminho mais curto que ainda é seguro:

| Tamanho | Quando | Saída |
|---|---|---|
| **Pequena** | o que entregar está claro, cabe numa passada e não abre decisão | código implementado e testado, **sem plano** |
| **Média** | há decisão a tomar com o usuário, ou vários passos, ou contrato novo (API, dados, migração) | `docs/plans/AAAA-MM-DD-<slug>.md` (rascunho) |
| **Grande** | o mesmo da média, mais decisão que vale para o projeto todo ou muitas partes | `docs/plans/AAAA-MM-DD-<slug>.md` (rascunho) + ADRs recomendadas |

O plano é **temporário**, como a descrição de um PR: o build o executa e o apaga. O que descreve
o sistema hoje são o código, as specs (`docs/specs/`), o `docs/architecture.md`, as ADRs e o
design system. Onde cada coisa mora: rule `docs.md`.

# Regras

1. **Leia antes de classificar.** O tamanho sai do código que o pedido toca, não só do texto do
   pedido.
2. **Na dúvida entre pequena e média, escolha pequena.** Se aparecer uma decisão no meio do
   caminho, pare e promova para média (Etapa 3, item 5). Corrigir o rumo custa menos que
   planejar o que não precisava.
3. **Pergunte só o que o código não responde**, com `AskUserQuestion` (2 a 4 opções, a
   recomendada primeiro). No máximo 2 rodadas de perguntas.
4. **Escopo é o pedido; o que sobra tem triagem.** A mudança que você entrega é a pedida, e nada
   mais. Na **pequena**, o que notar pelo caminho passa pela varredura (Etapa 3, item 6), depois
   de a mudança pedida estar pronta: impacto pequeno você corrige junto, item isolado você
   delega a um `spec-fixer`, e o resto vai para o resumo **com a recomendação e o comando
   pronto** — nunca uma nota solta. Critérios das duas faixas: rule `behavior.md`, regra 2. Na
   **média e na grande** você não mexe no código: tudo vai para o resumo, do mesmo jeito.
5. **Nunca commite.** O commit é do usuário (`/spec-semantic-commit`).
6. **Plano curto.** Média cabe em cerca de 1 página; grande, em até 3. Seção que não se aplica é
   removida, não fica "N/A". Mais de 8 passos → proponha fatiar em planos menores, cada um
   entregável sozinho.
7. **Português com acentuação correta.** Nomes de código ficam como no código.

# Fluxo

## Etapa 0 — Entrada

`$ARGUMENTS` pode ser:

- **texto livre** → é o pedido;
- **caminho de um `.md`/`.txt`** que existe → leia com `Read`, o conteúdo é o pedido. Leia
  também os arquivos que ele cita (prints, PDFs). Citado e não existe → avise;
- **vazio** → peça o pedido.

## Etapa 1 — Entender

1. Reescreva o pedido em **uma frase** para você mesmo.
2. Leia o `CLAUDE.md`, o índice `docs/specs/README.md` e a spec do módulo que o pedido toca,
   quando existirem. Depois localize no código o que o pedido toca (`Grep`/`Glob`, depois
   `Read`). Pare de explorar quando souber **onde** mexer e **o que** muda.
3. Pedido ambíguo a ponto de mudar **o que** é entregue → pergunte agora (regra 3).

## Etapa 2 — Classificar

Use a tabela do Objetivo. É **média ou grande** se qualquer um destes for verdade:

- há uma escolha de comportamento ou de produto que é do usuário;
- muda contrato que outros usam: API, formato de dados, schema, migração, arquivo lido por outro
  programa;
- precisa de mais de uma passada coerente (vários passos que dependem um do outro);
- toca área sensível: autenticação, autorização, pagamentos, segredos, dados pessoais.

É **grande** se, além disso, há uma decisão que vale para o projeto todo e é cara de desfazer
(stack, padrão de arquitetura, biblioteca central).

Diga a classificação **em uma linha**, com o motivo, e siga sem pedir confirmação
(`Tamanho: pequena — troca a cor do botão primário em app.css`). Se o usuário discordar,
reclassifique.

## Etapa 3 — Pequena: implementar

1. **Siga a vizinhança.** Antes de criar ou mudar algo, olhe um arquivo parecido que já existe e
   siga o padrão dele.
2. **Implemente** a menor mudança que resolve bem, com o teste que a cobre quando o projeto tem
   testes para aquela camada.
3. **Prove.** Rode os comandos de teste, build e lint do projeto (do `CLAUDE.md` ou do
   `README.md`) que cobrem o que mudou. Mudança visual → confira-a como manda a seção Teste
   visual do `CLAUDE.md` (o teste de tela ou o MCP que ela indica). Sem a seção → diga "não
   conferi na tela". Com a seção, não pule a conferência nem a passe ao usuário: não conseguiu
   operar a ferramenta (sem as ferramentas do MCP nesta sessão, ou a janela do app não veio
   depois de esperar e tentar de novo) → pergunte no chat que precisa conferir com ela e o que
   impediu, e confira com a resposta.
4. **Documente.** A mudança alterou o que uma spec descreve, ou criou algo que merece uma? Edite
   a spec agora, sem perguntar (rule `docs.md`).
5. **Apareceu uma decisão** (regra 2)? Pare, mostre o que já mudou e pergunte: "promover para
   média e escrever o plano" (recomendado) ou "decidir aqui e seguir".
6. **Varredura fora do escopo.** Com a mudança pedida feita e provada — nunca antes, para que
   ela fique pronta e separável se um fixer quebrar algo —, triagem do que você notou pelo
   caminho:

   | Faixa | O que é | Quem resolve |
   |---|---|---|
   | **A — junto** | impacto pequeno (`behavior.md`, regra 2, Exceção) | você, inline, agora |
   | **B — delegado** | falha A só por tamanho ou por precisar de teste próprio, e cumpre a Exceção 2 | um `spec-fixer` por item |
   | **C — do usuário** | tem decisão, muda contrato, toca área sensível ou precisa de mais de uma passada | ninguém: vai para o resumo |

   - **Nada em B → pule a varredura.** Ela não é etapa obrigatória nem passa a existir sozinha:
     sem item que cumpra os critérios, você só preenche as linhas do resumo.
   - **Feche o contrato de cada item de B** antes de chamar ninguém: `Problema` (com
     `arquivo:linha`), `Arquivos` (lista fechada) e `Pronto quando` (resultado concreto). Não
     soube fechar a lista de arquivos sem explorar meio projeto → é C.
   - **Até 3 por rodada**, nesta prioridade: quebrado > risco > sujeira. O resto é C, com o
     motivo `fila cheia`. Dois itens que tocam o mesmo arquivo → delegue um só; o outro é C.
     Arquivo que já estava sujo e não é da mudança pedida → não delegue nele (é C): não daria
     para desfazer só o que o fixer fez.
   - **Guarde a foto** (`git status --porcelain --untracked-files=all`) antes de delegar: é ela
     que diz o que desfazer no gate.
   - **Chame os `spec-fixer` em paralelo**, todos na mesma mensagem, um por item, cada um com só
     o contrato dele. Log de uma linha por item: `[fora do escopo] <item> → fixer`.
   - **Um nível só.** O que um fixer trouxer em "Notei fora do escopo" vira C no resumo, **nunca**
     um item novo nesta rodada. É esta trava que faz a varredura terminar.
   - **Gate.** Com todos de volta, rode uma vez os comandos de teste, build e lint do projeto,
     sobre tudo. Verde → pronto. Vermelho → `git restore` nos arquivos modificados **daquele**
     item e apague os que ele criou (só os que estavam limpos na foto), rode de novo e passe o
     item para C, com a falha como motivo. A falha não aponta de qual item veio → desfaça a
     rodada inteira e passe todos para C. Não pede confirmação (`behavior.md`, regra 4): desfaz
     só o que o fixer acabou de fazer, e a mudança pedida não é tocada.
   - Fixer que respondeu `recusado` → o item vira C, com o motivo dele.
7. Mostre o resumo:

```
Tamanho: pequena — <motivo>

Feito: <1-3 linhas do que mudou>
Arquivos: `<arquivo>` (+X −Y), …
Prova: `<comando>` → <resultado> | não rodei: <motivo>
Specs: <arquivo> (criado | atualizado | subdividido) | nenhuma

Corrigi fora do escopo: <o que e por quê, `arquivo:linha`> | nada
Resolvi em paralelo: <item> — `<arquivo>` · <prova> | nada
Deixei para você: <item, `arquivo:linha`> · <por que ficou> · Recomendo: <…> →
  `/spec-change "<pedido pronto>"` | nada
```

Um item por linha. Em "Deixei para você", "por que ficou" é o que o tirou da triagem: decisão,
contrato, área sensível, mais de uma passada, fila cheia, recusado pelo fixer ou gate vermelho.

Não commite.

## Etapa 4 — Média ou grande: escrever o plano

1. **Design system.** A mudança tem tela, o projeto tem interface e não existe
   `docs/design-system.md`? **Pare antes de planejar:** o design system é obrigatório em projeto
   com interface. Oriente `/spec-design-system` e rodar esta skill de novo depois.
2. **Conversar.** Junte as decisões que só o usuário pode tomar e pergunte (regra 3), sempre com
   a sua recomendação e o motivo. Decisão técnica que o código, uma spec ou uma ADR já respondem
   não é pergunta: registre e siga.
3. **ADRs.** Pelo índice `docs/adr/README.md`, ache e leia as ADRs `aceito` que tocam a
   mudança e respeite-as. Decisão
   nova que vale para o projeto todo (só na grande) → entra na seção "Decisões" com o comando
   `/spec-adr "<título>"`, a rodar antes de aprovar o plano.
4. **Telas.** Use os componentes e tokens do design system (`docs/design-system.md` e os
   arquivos de `docs/design/` que ele indica). Componente ou token que ainda não existe →
   diga no plano que o passo o acrescenta ao design system. O `CLAUDE.md` tem a seção Teste
   visual → o passo que muda uma tela leva o teste de tela no "Pronto quando"; olhar a tela só
   vira CA `(manual: …)` quando o projeto não tem como o agente conferir.
5. **Specs afetadas.** Liste as specs que a mudança cria ou altera, e o `docs/architecture.md`
   se a arquitetura base muda. Cada uma entra em "Arquivos" e no "Pronto quando" do passo que
   muda aquele comportamento: a documentação anda com o código, no mesmo diff.
6. **Plano.** Quebre em passos, cada um entregável e testável sozinho, na ordem em que podem
   ser feitos. Cada passo diz o que faz, os arquivos e o "Pronto quando" com resultado concreto
   (valor, status, mensagem; nunca "funciona"). Marque `Risco: alto` só no passo que toca área
   sensível: nele o build chama também a revisão de segurança.
7. **Salvar.** Nome `AAAA-MM-DD-<slug>.md` (data de hoje, slug em kebab-case, sem acentos) em
   `docs/plans/`. Preencha o [template](template.md), remova os comentários `<!-- LLM-ONLY -->`
   e salve com `status: rascunho`.
8. **Conferir.** Cada critério de aceite é provado por algum passo? Cada passo cabe numa
   passada? Toda spec afetada tem um passo que a atualiza? Alguma seção ficou vazia ou só
   repete outra? Corrija antes do resumo.
9. Mostre **apenas**:

```
Tamanho: <média | grande> — <motivo>
Plano salvo em: docs/plans/AAAA-MM-DD-<slug>.md

- Objetivo: <1 linha>
- Decisões: <as principais, 1 linha cada>
- Specs afetadas: <arquivos> | nenhuma
- Plano: N passos (risco alto: passo X | nenhum)
- ADR a registrar antes de aprovar: `/spec-adr "<título>"` | nenhuma
- Deixei para você: <item, `arquivo:linha`> · Recomendo: <…> → `/spec-change "<pedido pronto>"` | nada

Revise o plano. Aprovado → mude o status para `aprovado` (ou peça que eu mude) e rode
`/spec-build docs/plans/AAAA-MM-DD-<slug>.md`.
```

Não exiba o plano no terminal e não comece a implementar. Pediu ajuste → edite o arquivo e
mostre o resumo de novo.

# Entrada

$ARGUMENTS
