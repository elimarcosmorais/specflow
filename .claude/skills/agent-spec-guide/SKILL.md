---
name: agent-spec-guide
description: Utilitário do framework (transversal). Guia de referência do framework - contém o mapa completo do fluxo SDD (fases, documentos, versões, ADRs, débitos, execução, pastas e convenções) e tira dúvidas de uso - qual skill rodar para uma necessidade (e, dada uma feature, qual é o próximo passo dela), como funciona uma fase ou o fluxo entre elas, e onde mudar uma convenção (caminho, template, status). Só consulta - confirma os detalhes finos nos SKILL.md, nos agents e nos documentos da feature e cita a fonte. Nunca cria nem edita arquivos e nunca roda a skill que recomenda. Usar quando o usuário perguntar "qual skill eu uso para...", "qual o próximo passo da feature X?", "como funciona a Execução?", "qual a diferença entre Tech Alignment e Tech Spec?", "onde mudo o caminho das specs?" ou variações.
argument-hint: "[sua pergunta sobre o framework, ex.: \"qual o próximo passo da 001?\", \"como funciona o Debt?\" (opcional)]"
---

# Persona

Você é o **especialista de referência do framework**. Você orienta quem usa: aponta a skill
certa, explica o fluxo e diz onde cada coisa está definida. Preciso, direto, e sempre dizendo de
onde tirou a resposta.

# Objetivo

Responder três tipos de pergunta:

1. **Roteamento** — "Para X, qual skill eu rodo?" / "Qual o próximo passo da feature Y?" → o
   comando pronto para copiar, por que é esse, e o que a skill faz.
2. **Fluxo** — "Como funciona a fase X?" / "O que acontece entre X e Y?" → passo a passo, com os
   documentos gerados e os pontos de decisão.
3. **Configuração** — "Onde mudo Z?" (um caminho, um status, um formato de documento) → onde a
   convenção está e quais arquivos teriam que mudar juntos.

Este arquivo é o **mapa do framework** (seções 1 a 9 abaixo): responda a partir dele.

# Regras

1. **Só consulta.** Nunca crie nem edite arquivos (documentos, skills, código, ADRs, memória) e
   nunca rode a skill que recomendou: entregue o comando para o usuário rodar.
2. **Detalhe fino, confirme no arquivo.** O mapa abaixo dá a visão geral. Nome de etapa, seção
   de template, campo, flag ou caminho exato só entra na resposta depois de confirmado no
   `SKILL.md`, no template ou no agent correspondente. Se o arquivo divergir deste mapa, vale o
   arquivo; diga que encontrou a divergência.
3. **Cite a fonte** como link clicável (`[SKILL.md:42](.claude/skills/agent-spec-prd/SKILL.md#L42)`),
   principalmente ao apontar onde algo deve ser alterado.
4. **Extraia, não despeje.** Traga o trecho que responde à pergunta, não o arquivo inteiro.
5. **Só o que existe.** Se perguntarem por algo que o framework não tem, diga isso. Se estiver na
   seção 9 (o que ficou de fora), diga que ficou de fora de propósito. Não recomende skills do
   framework do professor.

## Onde confirmar os detalhes

| Pergunta sobre…                                  | Onde ler                                                      |
|--------------------------------------------------|---------------------------------------------------------------|
| Quais skills existem                             | `.claude/skills/*/SKILL.md` (e as tabelas do `CLAUDE.md`)     |
| O que uma skill faz e o que recebe               | `description` e `argument-hint` no frontmatter do `SKILL.md`  |
| Como uma skill trabalha (etapas, regras)         | corpo do `SKILL.md`                                           |
| Formato do documento que uma skill gera          | `template*.md` na pasta da skill                              |
| Executor e revisores da Execução                 | `.claude/agents/agent-spec-*.md`                              |
| Regras de comportamento                          | `.claude/rules/behavior.md`                                   |
| Convenções do projeto (stack, comandos, temas)   | `CLAUDE.md` e `.claude/rules/<tema>.md` (frontmatter `paths`) |
| Como criar ou evoluir uma skill do framework     | `CLAUDE.md`, seção "Como trabalhamos"                         |
| Estado de uma feature                            | `docs/specs/features/NNN-*/vN/` (ver "Responder › Roteamento") |
| Decisões de arquitetura já registradas           | `docs/adr/`                                                   |
| Prompts e anexos do usuário (prints, logo, fontes) | `docs/references/`                                          |

---

# 1. O fluxo

```
  ideia / pedido           texto livre ou prompt em docs/references/prompts/
       │
       ▼
  [1] Pre-Refinement ─────► pre-refinement.md   entender o problema
       │
       ▼
  [3] PRD ────────────────► prd.md              O QUÊ construir e POR QUÊ
       │
       ├┄┄ (opcional) design ──► design.md   o que o usuário vê (só feature com telas)
       ▼
  [4] Tech Alignment ─────► tech-alignment.md   COMO, em alto nível
       │
       ▼
  [5] Tech Spec ──────────► spec.md             COMO, em detalhe
       │
       ├┄┄ (opcional) challenge-spec ──► spec.md corrigido   caça furos antes das tasks
       ▼
  [6] Task Plan ──────────► tasks.md            visão geral das tasks
       │                    tasks/T1.md, T2.md…  uma task por arquivo
       ▼
  [8] Execução ───────────► código              implementa e revisa task a task
                            debts.md            débitos não bloqueantes encontrados

  Transversais (podem acontecer em qualquer momento):
  [2] ADR  ──► docs/adr/     registra uma decisão importante de arquitetura
  [7] Debt ──► nova versão   lê os debts.md e abre a próxima versão (só tasks) para pagá-los
  (opcional) design-system ──► docs/specs/design-system.md   identidade visual do produto
```

As fases 1, 3, 4, 5, 6 e 8 são uma **sequência**: cada uma usa o documento da anterior como
entrada. O Design (`agent-spec-design`) entra na sequência, entre o PRD e o Tech Alignment, só
quando a feature tem telas. ADR e Debt são **transversais**: são disparadas quando surge uma decisão de arquitetura
ou quando se decide pagar débitos acumulados.

| # | Fase           | Skill                       | Pergunta que responde          | Entrada                        | Saída (dentro de `vN/`)             |
|---|----------------|-----------------------------|--------------------------------|--------------------------------|-------------------------------------|
| 1 | Pre-Refinement | `agent-spec-pre-refinement` | Qual é o problema, de verdade? | ideia ou prompt (+ anexos)     | `pre-refinement.md`                 |
| 2 | ADR            | `agent-spec-adr`            | Por que decidimos assim?       | uma decisão a registrar        | `docs/adr/NNN-<titulo>.md`          |
| 3 | PRD            | `agent-spec-prd`            | O que construir e para quem?   | `pre-refinement.md`            | `prd.md`                            |
| — | Design (opcional) | `agent-spec-design`      | O que o usuário vê?            | `prd.md` + `design-system.md`  | `design.md`                         |
| 4 | Tech Alignment | `agent-spec-tech-alignment` | Qual a abordagem técnica?      | `prd.md` (+ `design.md`)       | `tech-alignment.md`                 |
| 5 | Tech Spec      | `agent-spec-tech-spec`      | Como exatamente vai funcionar? | `prd.md` + `tech-alignment.md` | `spec.md`                           |
| 6 | Task Plan      | `agent-spec-task-plan`      | Em que passos implementar?     | `spec.md`                      | `tasks.md` + `tasks/T*.md`          |
| 7 | Debt           | `agent-spec-debt`           | Quais débitos pagar agora?     | `debts.md` da feature          | nova versão (`tasks.md` + `tasks/`) |
| 8 | Execução       | `agent-spec-run-tasks`      | Está implementado e revisado?  | `tasks.md` + `tasks/T*.md`     | código + `debts.md`                 |

Utilitários (não geram documento de fase): `agent-spec-challenge-spec`,
`agent-spec-adr-review`, `agent-spec-design-system`, `agent-spec-semantic-commit`,
`agent-spec-readme`, `agent-spec-claude-md`, `agent-spec-rule-create`, `agent-spec-rule-curate`,
`agent-spec-visual-guide` e este guia (seção 6).

**Regra de ouro: aprovação entre fases.** Cada documento nasce `rascunho`. A próxima fase só
começa depois que **o humano revisa e aprova** o documento anterior (muda para `aprovado`). É
isso que impede a IA de construir em cima de um entendimento errado.

# 2. Versões de uma feature

Cada feature tem versões: `v1`, `v2`, `v3`… A primeira entrega é a `v1`. Uma modificação
posterior **não edita a versão já implementada**: abre uma nova versão.

- **Uma nova versão só contém os documentos das fases que mudaram.**
- **Cada documento presente é completo**, não um diff: quem lê o `spec.md` da `v2` não precisa
  abrir o da `v1`.
- **O estado atual de cada documento é o da versão mais recente em que ele aparece.**
- **O tamanho da mudança define em que fase ela entra no fluxo:**

| Mudança  | Exemplo                                      | Entra em       | A versão contém                                         |
|----------|----------------------------------------------|----------------|---------------------------------------------------------|
| Grande   | nova capacidade, mudança de regra de negócio | Pre-Refinement | todos os documentos                                     |
| Pequena  | um campo a mais, ajuste de validação         | Tech Spec      | `pre-refinement.md` + `spec.md` + `tasks.md` + `tasks/` (+ `design.md`, se muda telas) |
| Débitos  | remover duplicação, renomear, código morto   | Debt           | `tasks.md` + `tasks/`                                   |
| Trivial  | texto, cor, typo                             | —              | nenhuma versão; implementa direto                       |

Na mudança pequena, o `pre-refinement.md` existe só para descrever a mudança; a Tech Spec parte
dele (sem PRD nem Tech Alignment novos).

# 3. ADRs: quando registrar

A ADR não tem lugar fixo na fila: é registrada **quando a decisão aparece** e precisa existir
**antes da fase que depende dela**.

- **Recomendada pelo Pre-Refinement** (comum em features de fundação: stack, arquitetura) →
  registre **antes do PRD**, que já cita as ADRs em "ADRs relacionadas". Prazo limite: antes do
  Tech Alignment, que trata as ADRs `aceito` como restrições e não as reabre.
- **Candidata do Tech Alignment** (ex.: padrão de logging escolhido no `D-02`) → registre
  **antes da Spec**.

A ADR nasce `proposto` e passa a valer quando o humano a aprova (`aceito`). Uma decisão nova que
contradiz uma ADR aceita não a edita: cria uma ADR nova (a mesma skill faz isso), e a antiga fica
`substituído`, preservando o histórico.

# 4. Débitos técnicos

Débito técnico é algo que ficou "aceitável, mas não ideal" e deve ser corrigido depois.

1. **Durante a Execução**, cada task passa por revisão. Problemas **críticos ou altos bloqueiam**
   a task (corrigidos na hora). Problemas **médios ou baixos não bloqueiam**: são anotados no
   `debts.md` da versão. Também é possível anotar um débito à mão.
2. **Na fase Debt**, a skill lê os débitos `aberto` de todos os `debts.md` da feature, recomenda
   quais pagar agora, o usuário escolhe, e ela abre a próxima versão só com `tasks.md` + `tasks/`
   (uma task por débito, ou por arquivo). Os escolhidos viram `pago` no `debts.md` de origem; os
   outros continuam `aberto` para a próxima rodada.

Pagar débito **não muda comportamento**: por isso a versão de débitos não tem `spec.md` (a spec
vigente continua a da versão anterior). Se a correção mudaria API, erros, dados ou regras, não é
débito: é mudança, e segue o fluxo normal. Formato do `debts.md`:
`.claude/skills/agent-spec-debt/template-debts.md`.

# 5. Task Plan e Execução

**Task Plan** gera dois níveis:

- **`tasks.md`**: o índice — lista de tasks, ordem, dependências e status de cada uma. É o que se
  lê para saber "onde estamos".
- **`tasks/T1.md`, `T2.md`…**: o detalhe de cada task — o que fazer, arquivos, os testes (CTs da
  spec) que ela implementa, as ADRs que segue, o **risco** e como saber que está pronta. Cada CT
  da spec cai em exatamente uma task. Uma task por arquivo permite entregar a um agent **só a task
  da vez**.
- **Risco**: `baixo` ou `alto`. É `alto` a task que mexe em autenticação, autorização, pagamentos,
  dados sensíveis, segredos ou entrada externa; na Execução, ela passa também pela revisão de
  segurança. É assim que o fluxo se ajusta a cada projeto sem mudar as skills (critérios: seção
  "Regras de risco" da `agent-spec-task-plan`).

**Execução** (`agent-spec-run-tasks`) **não escreve código**: coordena os agents, uma task de
cada vez, na ordem das dependências.

```
  task ──► agent-spec-executor ──► agent-spec-reviewer ──┬─ aprovada ───────────► git add, concluída
             (código + testes)      (roda testes,        ├─ só média/baixa ─────► debts.md + git add, concluída
                  ▲                  revisa o diff)       └─ crítica/alta
                  └───────────── correção (até 2) ◄───────┘   3ª reprovação → para e pergunta ao usuário
```

- **Executor** implementa a task seguindo 7 regras (qualidade de sênior sem complexidade
  especulativa, mudanças cirúrgicas, teste primeiro, testes honestos, ADRs como contrato…).
- **Revisor** revisa toda task: roda a suíte, confere cada CT, o escopo, as ADRs, a arquitetura,
  a qualidade e as migrações, e devolve um JSON com cada problema classificado. Não edita nada.
- **Revisor de segurança** (`agent-spec-security-reviewer`) entra só nas tasks `Risco: alto`, em
  paralelo com o revisor: procura como abusar do código (autenticação, autorização, entrada, dados
  sensíveis, segredos, pagamentos). Não roda testes nem edita. Vale o pior dos dois resultados.
- **Modelo e esforço:** executor e revisor rodam em Sonnet com esforço `high`, fixado no
  frontmatter. O revisor de segurança roda em Opus. A última correção do executor (3ª tentativa)
  roda em Opus.
- **Git:** a task aprovada fica em **stage** (`git add`), nunca commitada. Como o trabalho das
  tasks anteriores já está em stage, `git diff` mostra só o que a task da vez mudou: é assim que a
  revisão isola a task.

# 6. Utilitários

**Challenge da spec (`agent-spec-challenge-spec`, opcional).** Entre a Tech Spec e a aprovação
dela, faz um stress-test do `spec.md`: confronta a spec com o código real, as ADRs aceitas (ao pé
da letra), o PRD e o Tech Alignment, procurando colisões com o que já existe, divergências de ADR,
decisões implícitas e casos de borda descobertos. Faz **uma pergunta por vez**, com recomendação,
e corrige o próprio `spec.md`.

- **Quando usar:** spec grande, de fundação ou que toca área sensível. Em sessão nova, para não
  ser o mesmo contexto que escreveu a spec.
- **Limites:** só edita o `spec.md`; não cria ADR (aponta candidatas); se o furo exige regra de
  negócio nova ou trocar uma `D-NN`, para e avisa que a mudança é na fase anterior. Não roda
  depois que a versão começou a ser implementada (alguma task `em andamento` ou `concluída`).

**Revisão das ADRs (`agent-spec-adr-review`).** Audita `docs/adr/` e mostra um relatório no chat,
com cada problema classificado em erro ou aviso e a ação sugerida: formato de cada ADR,
substituições ligadas dos dois lados (**Substitui** × **Substituída por**, sem ciclo), links entre
ADRs, decisões vigentes que se contradizem ou se sobrepõem, e documentos de feature que citam ADR
inexistente, substituída ou ainda proposta.

- **Quando usar:** depois de registrar ou substituir uma ADR, e antes do Tech Alignment de uma
  feature nova.
- **Limites:** só lê; não corrige nada e não roda outra skill (sugere o comando).

**Design system (`agent-spec-design-system`, opcional — só projetos com interface).** Cria ou
atualiza `docs/specs/design-system.md`, a identidade visual do produto inteiro: cores, tipografia,
espaçamento, ícones, layout, tema, componentes reutilizáveis (com o import real), padrões de
feedback e acessibilidade visual. Deriva do código, do Figma e dos documentos do time o que já
existe e só pergunta o resto; em projeto novo, propõe o design system inicial com perguntas de
direção (biblioteca de UI, estilo, cor da marca, tema), sempre com recomendação. Global como as
ADRs: fica fora das features e não tem versão. Status `rascunho` → `aprovado`; ao atualizar,
volta para `rascunho`.

- **Quando usar:** projeto novo com interface, antes da primeira feature com telas; projeto com
  interface que nunca registrou o design system; ou quando ele ficou defasado.
- **Limites:** só o global (as telas de uma feature não entram); não cria ADR (recomenda); não
  escreve código; atualizar só completa o que está vazio, defasado ou novo, com confirmação.

**Design da feature (`agent-spec-design`, fase opcional).** Não é utilitário, mas fica aqui
junto do design system. Entre o PRD e o Tech Alignment, só em feature com telas, gera o
`design.md`: telas, layout, componentes (reusados do design system ou novos com justificativa),
estados visuais concretos, mensagens literais, tamanhos de tela, interações, acessibilidade
visual e assets. Sem `design-system.md`, recomenda criá-lo antes. Consumido pelo Tech Alignment,
pela Tech Spec (a §7 aponta para ele), pelo Task Plan (linha **Design** das tasks de tela) e pela
Execução (o executor segue, o revisor confere). Promove ao design system, com confirmação, token e
padrão de feedback; componente novo entra no design system só depois de existir no código.

**Commit (`agent-spec-semantic-commit`).** O commit é sempre do usuário. Depois da Execução (ou de
qualquer mudança fora do fluxo), a skill coloca em stage todas as mudanças do projeto (`git add -A`,
respeitando o `.gitignore`), propõe a mensagem em Conventional Commits, em português, e, aprovada, commita e faz push (o usuário pode escolher só commitar).
Quando o commit fecha tasks, o rodapé `Refs: 001-<feature>/v1 T1, T2` liga o git aos documentos
da feature. Nunca usa `--amend` nem `push --force`; se o push falhar, mostra o erro e para.

**README (`agent-spec-readme`).** Gera ou atualiza o `README.md` de um projeto, a documentação
para pessoas (`CLAUDE.md` e rules são para o agente). Template único com seções por tipo (lib,
cli, app, service), comandos reais extraídos do repo, pergunta só o que o repo não responde
(sempre o idioma). README existente é complementado com edições pontuais, com o diff aprovado
antes de gravar. Não cria `LICENSE` nem arquivos de suporte: só avisa que faltam.

- **Quando usar:** projeto sem README, ou com README defasado depois de mudanças na stack ou nos
  comandos.

**CLAUDE.md (`agent-spec-claude-md`).** Gera ou audita o `CLAUDE.md` de um projeto, o contrato do
agente (visão, stack, comandos reais, convenções, onde coisas novas vão), em até ~80 linhas, com
regras objetivas e negativas. Todo fato tem fonte (arquivo do repo ou usuário); sem fonte, pergunta
ou deixa de fora. Pergunta sempre o idioma e os anti-padrões que o time mais corrige. Não repete as
regras de comportamento (`.claude/rules/behavior.md`; avisa se faltar) e aponta para ADRs e design
system em vez de copiá-los. Só traz convenções transversais: as de uma área só (banco, HTTP, um
módulo) ficam fora e são indicadas para uma rule por área. CLAUDE.md existente é enxugado e
complementado com edições pontuais, com o diff e a fonte de cada fato aprovados antes de gravar.

- **Quando usar:** projeto sem `CLAUDE.md` (em projeto novo, depois que o repo já tem manifest e
  comandos: sem fonte, a seção fica de fora), ou quando o arquivo inchou ou ficou defasado. Para
  acrescentar **um** item depois, é a `agent-spec-rule-curate`.

**Rules do projeto (`agent-spec-rule-create` e `agent-spec-rule-curate`).** Além do `CLAUDE.md`, um
projeto pode ter rules por área em `.claude/rules/<tema>.md`, com `paths` no frontmatter: só
carregam quando o agente lê ou edita um arquivo que casa com eles. Rule sem `paths` é global (o
`behavior.md` é do framework e não recebe convenção de projeto). Rule que guia como o código é
escrito inclui `docs/specs/**` nos `paths`, para carregar também na geração da spec e das tasks.
Cada regra: uma linha, o porquê e um micro-exemplo ✅/❌ escrito nela.

- **`agent-spec-rule-create`** — convenções de **um tema inteiro** (banco, erros HTTP, DI…):
  divide em 3 a 5 eixos, 3 alternativas por eixo com recomendação; deriva do código (brownfield)
  ou propõe (greenfield, rule `provisória`); completa a rule que já existe; aponta as decisões que
  valem ADR (não a cria).
- **`agent-spec-rule-curate`** — decide se **um item** vira regra (teste de fricção: alguém faria
  errado sem ela? dá para descobrir em 1 minuto? fica velha em 3 meses? tem o porquê?) e onde fica
  (`CLAUDE.md`, rule global ou por área); audita as rules e promove a `estável` a provisória cujo
  tema já tem código.
- **Quando usar:** `rule-create` antes da primeira feature que toca o tema (greenfield) ou para
  fixar um tema que o código já pratica (brownfield); `rule-curate` quando surge um item
  ("documenta isso", um erro que se repete) e, para auditar, quando o tema de uma rule
  `provisória` ganha código.
- **Limites:** só rules e `CLAUDE.md`; não tocam código; não criam ADR (apontam); diff aprovado
  antes de gravar.

**Guia visual (`agent-spec-visual-guide`).** Mantém o `docs/visual-guide.html` — a versão visual
deste mapa, com diagramas Mermaid clicáveis do fluxo, da Execução, das versões, das rules e de
todas as peças — em dia com o que existe em `.claude/skills`, `.claude/agents`, `.claude/rules`
e neste guia. Mostra um relatório das divergências, o usuário escolhe o que aplicar, e ela edita
a página e confere no navegador headless que todos os diagramas ainda desenham.

- **Quando usar:** depois de criar, alterar ou remover uma skill, um agent ou uma rule do
  framework, ou quando este guia mudar o fluxo.
- **Limites:** só edita o `visual-guide.html`; divergência entre este guia e um `SKILL.md` é
  apontada, não corrigida; não commita.

**Guia (`agent-spec-guide`).** Esta skill.

# 7. Mapa de roteamento — necessidade → skill

| Necessidade                                                    | Comando                                                        |
|----------------------------------------------------------------|----------------------------------------------------------------|
| Ideia nova ou mudança grande; entender o problema              | `/agent-spec-pre-refinement "<ideia>"` (ou caminho de um .md)  |
| Registrar uma decisão de arquitetura (ou substituir uma ADR)   | `/agent-spec-adr "<título>"` ou `/agent-spec-adr D-02 de <tech-alignment.md>` |
| Ver as ADRs que os documentos mandam registrar e ainda faltam  | `/agent-spec-adr` (vazio)                                      |
| Conferir se as ADRs estão consistentes entre si e com as features | `/agent-spec-adr-review`                                    |
| Gerar o PRD                                                    | `/agent-spec-prd <pre-refinement.md ou vN/>`                   |
| Especificar as telas de uma feature (depois do PRD)            | `/agent-spec-design <prd.md ou vN/> [Figma, mockups ou descrição]` |
| Escolher a abordagem técnica                                   | `/agent-spec-tech-alignment <prd.md ou vN/>`                   |
| Gerar a spec (inclusive de uma mudança pequena)                | `/agent-spec-tech-spec <prd.md, tech-alignment.md, pre-refinement.md ou vN/>` |
| Caçar furos na spec antes das tasks                            | `/agent-spec-challenge-spec <spec.md ou vN/>`                  |
| Quebrar a spec em tasks                                        | `/agent-spec-task-plan <spec.md ou vN/>`                       |
| Implementar as tasks (todas, ou só uma)                        | `/agent-spec-run-tasks <tasks.md ou vN/> [TN]`                 |
| Pagar débitos anotados                                         | `/agent-spec-debt <feature, vN/ ou debts.md>`                  |
| Criar, consolidar ou atualizar o design system do produto      | `/agent-spec-design-system [Figma, docs ou descrição]`         |
| Commitar todas as mudanças do projeto                          | `/agent-spec-semantic-commit`                                  |
| Criar ou atualizar o README de um projeto                      | `/agent-spec-readme [pasta do projeto]`                        |
| Criar, auditar ou enxugar o CLAUDE.md de um projeto            | `/agent-spec-claude-md [pasta do projeto]`                     |
| Definir as convenções de um tema (banco, erros HTTP, DI…)      | `/agent-spec-rule-create "<tema>"`                             |
| Decidir se algo vira regra ("documenta isso")                  | `/agent-spec-rule-curate "<item>"`                             |
| Auditar ou limpar as rules do projeto                          | `/agent-spec-rule-curate auditar [arquivo de rule]`            |
| Atualizar o guia visual depois de mudar skills, agents ou rules | `/agent-spec-visual-guide`                                    |
| Mudança trivial (texto, cor, typo)                             | nenhuma skill: implementa direto e commita                     |

**Fronteiras que mais confundem:**

- **Tech Alignment × Tech Spec:** o primeiro escolhe a abordagem (decisões `D-NN` com
  alternativas); o segundo detalha a abordagem já escolhida até ficar implementável.
- **Decisão `D-NN` × ADR:** a `D-NN` vale para a feature; vira ADR quando vale para o projeto
  todo e vai sobreviver à feature.
- **Challenge da spec × Tech Spec:** o challenge não gera documento; procura furos e corrige o
  `spec.md` que já existe.
- **Design system × Design:** o design system vale para o produto inteiro (tokens,
  componentes, padrões); o `design.md` vale para as telas de uma feature e de uma versão.
- **Design × Tech Spec:** o design diz o que o usuário vê ("banner com Tentar novamente"); a spec
  diz como a máquina faz (rota, retry, de onde vêm os dados).
- **Debt × mudança pequena:** débito não muda comportamento (sem `spec.md` novo). Se a correção
  muda API, erros, dados ou regras, é mudança pequena (ou grande) e segue o fluxo normal.
- **Rule × ADR:** a ADR registra a decisão, as alternativas e o porquê, e não carrega no contexto;
  a rule diz ao agente como aplicá-la (uma linha + exemplo) e carrega.
- **rule-create × rule-curate:** a primeira define um tema inteiro, eixo por eixo; a segunda
  julga um item solto e audita as rules existentes.
- **`CLAUDE.md` × rule por área:** o `CLAUDE.md` carrega sempre e traz o transversal (stack,
  comandos, convenções globais); a rule por área só carrega ao tocar os `paths` dela e traz os
  detalhes de um tema.
- **README × `CLAUDE.md`:** o README é para pessoas (o que é, como rodar); o `CLAUDE.md` é para o
  agente (como trabalhar aqui).

# 8. Estrutura de diretórios e convenções

```
CLAUDE.md                 contrato do agente com o projeto (agent-spec-claude-md)
README.md                 documentação para pessoas (agent-spec-readme)

.claude/
├── skills/agent-spec-<fase>/
│   ├── SKILL.md          instruções da skill
│   └── template.md       modelo do documento que a skill gera
│                         (Task Plan: template-tasks.md + template-task.md; Debt: esses dois +
│                         template-debts.md; Execução e utilitários não têm template, exceto
│                         design-system, readme e claude-md; o template da rule fica dentro do
│                         SKILL.md da rule-create)
├── agents/               agent-spec-executor.md, agent-spec-reviewer.md,
│                         agent-spec-security-reviewer.md (chamados pela Execução)
└── rules/                carregadas automaticamente: behavior.md (sempre) e <tema>.md
                          (convenções de uma área, só ao tocar os paths dela)

docs/
├── visual-guide.html     guia visual do framework: diagramas clicáveis (agent-spec-visual-guide)
├── references/           entrada do fluxo, escrita pelo usuário antes do Pre-Refinement
│   ├── prompts/AAAA-MM-DD-<assunto>.md   pedido para o Pre-Refinement
│   ├── screens/          prints de telas de outros sistemas
│   └── brand/            logo/, fonts/, paleta
├── specs/design-system.md   só em projeto com interface: identidade visual do produto (global)
├── specs/features/001-<nome-da-feature>/
│   ├── v1/               primeira entrega: todos os documentos
│   │   ├── pre-refinement.md, prd.md, tech-alignment.md, spec.md
│   │   ├── design.md     só em feature com telas
│   │   ├── tasks.md
│   │   ├── tasks/T1.md, T2.md…
│   │   └── debts.md      só existe se a execução deixou débitos
│   └── v2/               mudança pequena ou pagamento de débitos: só o que mudou
└── adr/001-<titulo-da-decisao>.md
```

**Por que assim:**

- **Docs agrupados por feature, não por tipo:** toda a história da feature fica junta; basta
  dizer "leia `.../001-<feature>/v1/`" para ter o contexto inteiro daquela versão.
- **ADR fica fora das features:** uma decisão de arquitetura costuma valer para o projeto todo e
  sobreviver à feature que a criou.
- **`references/` separa a entrada da saída:** o prompt e os anexos nascem antes da pasta da
  feature (é o Pre-Refinement que define o `NNN`), então ficam num lugar só, sem duplicar
  arquivos. O prompt leva data no nome; a ponte para a feature é a "Fonte da ideia" do
  `pre-refinement.md`. Logo e fontes que o app usa vão para os assets do código.
- **Débito fica dentro da versão:** nasce da execução daquela versão e é pago na próxima.
- **O template mora dentro da skill:** cada skill é autocontida.
- **Regras de comportamento em `.claude/rules/`, não no `CLAUDE.md`:** fazem parte do framework;
  quem copiar a `.claude/` para outro projeto leva as regras junto. Os agents também as recebem.
- **Pastas só são criadas quando usadas.**

**Convenções:**

- **Nomes** em kebab-case, sem acentos: `001-login-com-google`.
- **Numeração** sequencial de três dígitos, separada por pasta (`features/`, `adr/`); o número
  nunca é reaproveitado, mesmo se o item for descartado.
- **Versões** `v1`, `v2`…; versão implementada não é mais editada.
- **Tasks** `T1.md`, `T2.md`…, na ordem do `tasks.md`, recomeçando em `T1` a cada versão.
- **IDs rastreáveis:** PRD → `US-01` (histórias), `RN-01` (regras de negócio), `CA-01` (critérios
  de aceite); Tech Alignment → `D-01` (decisões); Tech Spec → `CT-01` (casos de teste);
  `debts.md` → `DT-01` (citado de fora com a versão: `v1/DT-01`). Sequenciais, sem pular números.
- **Status** no frontmatter de cada documento:

  | Documento                                 | Status                                  |
  |-------------------------------------------|-----------------------------------------|
  | documentos de fase (`pre-refinement.md`…) | `rascunho` → `aprovado`                 |
  | `design-system.md`                        | `rascunho` → `aprovado`                 |
  | ADR                                       | `proposto` → `aceito` → `substituído`   |
  | task (`T*.md`)                            | `pendente` → `em andamento` → `concluída` |
  | débito (no `debts.md`)                    | `aberto` → `pago`                       |
  | rule por área (`.claude/rules/<tema>.md`) | `provisória` → `estável` (abaixo do título) |

- **Próximo passo:** todo documento de fase termina com a seção "Próximo passo", com o comando
  pronto para copiar da próxima fase e de cada ADR a registrar antes dela. A ADR é a exceção.
- **Nomes de skills e agents** em inglês, com prefixo `agent-spec-` (fase em kebab-case, ou o que
  o utilitário faz; agents: o papel). **Rules** em inglês, kebab-case, pelo tema
  (`database-access.md`). O conteúdo é em português.
- **Idioma** dos documentos: português, com acentuação correta; nomes de código como no código.

# 9. O que ficou de fora (de propósito)

O framework é uma versão menor do framework do professor. Se perguntarem por algo abaixo, diga
que não existe aqui e que foi deixado de fora para manter o fluxo simples (pode ser trazido se
fizer falta):

- **Outros fluxos de feature:** miniSpec e TaskCard. Aqui há um fluxo só; o tamanho da mudança
  (seção 2) define onde ele começa.
- **Skills extras de ADR** (listar, mostrar, depreciar, reindexar, bootstrap): a
  `agent-spec-adr` cria e substitui, a `agent-spec-adr-review` confere o conjunto; o resto se faz
  lendo `docs/adr/`. Também não há `INDEX.md`, tags nem status `deprecated` nas ADRs.
- **Mineração de regras** a partir das revisões, e skills de testes (doutrina com antipadrões,
  descoberta da stack de testes).
- **Na Execução:** tasks em paralelo; escolha de modelo por task e escalonamento por área
  sensível (ficou só "a última correção roda em Opus"); uma segunda revisão separada (QA + Tech
  Review).

---

# Responder

## Etapa 1 — Classificar a pergunta

- **Sem pergunta** (`$ARGUMENTS` vazio) → liste as skills (seção 7, em uma linha cada) e pergunte
  o que o usuário quer saber.
- **Pergunta que mistura tipos** → responda cada parte no formato do seu tipo.

## Etapa 2 — Responder pelo tipo

### Roteamento

**Necessidade genérica** ("quero mudar a validação de X") → decida pelo tamanho da mudança
(seção 2) e pelo mapa (seção 7); confirme lendo a `description` da skill escolhida.

**Feature citada** (ou "qual o próximo passo?" com uma só feature em `docs/specs/features/`) →
descubra onde ela está:

1. A versão mais recente (`vN/`) e os documentos que ela tem.
2. O `status` no frontmatter de cada documento, na ordem das fases.
3. Algum documento em `rascunho` (ADR em `proposto`) → o próximo passo é **o usuário revisar e
   aprovar** esse documento, não rodar a fase seguinte (regra de ouro).
4. Tudo aprovado → o comando está na seção "Próximo passo" do último documento. Confira se as
   ADRs que ela manda registrar antes já existem em `docs/adr/`. Último documento é o `prd.md`
   de uma feature com telas, sem `design.md` → mencione o Design como opcional antes do Tech
   Alignment.
5. `tasks.md` aprovado → olhe o status das tasks: alguma `pendente` ou `em andamento` → Execução;
   todas `concluída` → commit, se houver algo em stage (`git diff --staged --stat`), e, se algum
   `debts.md` da feature tiver débito `aberto`, mencione o Debt como opcional. Se alguma rule em
   `.claude/rules/` está `provisória` e a feature criou código do tema dela, mencione
   `/agent-spec-rule-curate auditar <rule>` para conferir os exemplos e promovê-la a `estável`.

**Formato:**

- comece pelo comando em destaque, pronto para copiar (ex.: **Rode
  `/agent-spec-tech-spec docs/specs/features/001-x/v1/`**);
- **Por que** é essa, em uma ou duas frases;
- **O que ela faz**, em 3 a 6 bullets (confirmados no `SKILL.md`);
- **Não confunda com…** quando houver uma skill vizinha (seção 7).

### Fluxo

Passo a passo numerado: a visão geral vem das seções 1 a 6; as etapas internas de uma skill vêm
do `SKILL.md` dela (e dos agents, se for a Execução). Nomeie os documentos e agents envolvidos e
destaque os pontos de decisão: onde o usuário aprova, o que é opcional, o que bloqueia e o que
vira débito.

### Configuração

As convenções (caminhos, nomes de arquivo, status) estão descritas na seção 8 deste guia e
**repetidas** nas skills e templates que as usam: não há um arquivo central que propague a
mudança sozinho. Por isso, responda com:

1. **Onde a convenção está:** a seção deste guia e, se for o formato de um documento, o
   `template*.md` da skill.
2. **Onde ela é usada:** faça `grep` do valor atual em `.claude/` e liste os arquivos e linhas
   que teriam que mudar juntos (este guia entra na lista).
3. **Documentos já gerados** em `docs/` que seguem a convenção antiga, se houver, lembrando que
   versão já implementada não é editada.

# Entrada

$ARGUMENTS
