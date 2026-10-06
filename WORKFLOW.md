# Workflow do framework

Este arquivo descreve **como o nosso framework funciona**: o fluxo das fases do SDD, que
documento cada fase gera e onde cada arquivo fica. É o mapa que as skills seguem.

## Visão geral do fluxo

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
entrada. O Design entra na sequência, entre o PRD e o Tech Alignment, só quando a feature tem
telas. ADR e Debt são **transversais**: não são um passo da fila, são disparadas quando
surge uma decisão de arquitetura ou quando se decide pagar débitos acumulados.

### Regra de ouro: aprovação entre fases

Cada documento nasce como `rascunho`. A próxima fase só começa depois que **o humano revisa e
aprova** o documento anterior (muda para `aprovado`). É isso que impede a IA de construir em
cima de um entendimento errado.

### Versões de uma feature

Cada feature tem versões: `v1`, `v2`, `v3`… A primeira entrega é a `v1`. Uma modificação
posterior na feature **não edita a versão já implementada**: ela abre uma nova versão.

Regras:

- **Uma nova versão só contém os documentos das fases que mudaram.**
- **Cada documento presente é completo**, não um diff: quem lê o `spec.md` da `v2` não precisa
  abrir o da `v1`.
- **O estado atual de cada documento é o da versão mais recente em que ele aparece.**
- **O tamanho da mudança define em que fase ela entra no fluxo:**

| Mudança  | Exemplo                                      | Entra em       | A versão contém                                          |
|----------|----------------------------------------------|----------------|----------------------------------------------------------|
| Grande   | nova capacidade, mudança de regra de negócio | Pre-Refinement | todos os documentos                                      |
| Pequena  | um campo a mais, ajuste de validação         | Tech Spec      | `pre-refinement.md` + `spec.md` + `tasks.md` + `tasks/` (+ `design.md`, se muda telas) |
| Débitos  | remover duplicação, renomear, código morto   | Debt           | `tasks.md` + `tasks/`                                    |
| Trivial  | texto, cor, typo                             | —              | nenhuma versão; implementa direto                        |

Assim, cada versão é um retrato fiel do que foi pedido e construído naquele momento, sem
obrigar a refazer documentos que não mudaram.

### Quando registrar uma ADR

A ADR não tem lugar fixo na fila: ela é registrada **quando a decisão aparece**, e precisa
existir **antes da fase que depende dela**. Na prática:

- **Recomendada pelo Pre-Refinement** (comum em features de fundação: stack, arquitetura) →
  registre **antes do PRD**. O PRD já cita as ADRs em "ADRs relacionadas" e não precisa
  descrever a stack de forma genérica. Prazo limite: antes do Tech Alignment, que trata as ADRs
  `aceito` como restrições e não as reabre.
- **Candidata do Tech Alignment** (ex.: padrão de logging escolhido no `D-02`) → registre
  **antes da Spec**.

Uma ADR nasce `proposto` e passa a valer quando o humano a aprova (`aceito`), como os outros
documentos. Uma decisão nova que contradiz uma ADR aceita não a edita: cria uma ADR nova, e a
antiga fica `substituído`, preservando o histórico.

Para conferir o conjunto, a skill utilitária `agent-spec-adr-review` audita `docs/adr/` e mostra
um relatório no chat: formato de cada ADR, substituições ligadas dos dois lados, links entre
ADRs, decisões vigentes que se contradizem ou se sobrepõem, e documentos de feature que citam ADR
inexistente, substituída ou ainda proposta. Ela só lê: não corrige nada, só sugere a ação. Vale
rodar depois de registrar ou substituir uma ADR e antes do Tech Alignment de uma feature nova.

### Débitos técnicos

Débito técnico é algo que ficou "aceitável, mas não ideal" e deve ser corrigido depois.

1. **Durante a Execução**, cada task passa por revisão. Problemas **críticos ou altos
   bloqueiam** a task (precisam ser corrigidos na hora). Problemas **médios ou baixos não
   bloqueiam**: são anotados no `debts.md` da versão. Também é possível anotar um débito à mão.
2. **Na fase Debt**, a skill lê os débitos `aberto` de todos os `debts.md` da feature, recomenda
   quais pagar agora, o usuário escolhe, e ela abre a próxima versão da feature só com
   `tasks.md` + `tasks/` (uma task por débito, ou por arquivo). Os débitos escolhidos são
   marcados como `pago` no `debts.md` de origem; os outros continuam `aberto` para a próxima
   rodada.

Pagar débito **não muda comportamento**: o código só passa a seguir melhor o que a spec já diz.
Por isso a versão de débitos não tem `spec.md` (a spec vigente continua a da versão anterior).
Se a correção mudaria API, erros, dados ou regras, não é débito: é uma mudança, e segue o fluxo
normal (Pre-Refinement). O formato do `debts.md` está em
`.claude/skills/agent-spec-debt/template-debts.md`.

## As fases

| # | Fase           | Pergunta que responde               | Entrada                        | Saída (dentro de `vN/`)    |
|---|----------------|-------------------------------------|--------------------------------|----------------------------|
| 1 | Pre-Refinement | Qual é o problema, de verdade?      | ideia ou prompt (+ anexos)     | `pre-refinement.md`        |
| 2 | ADR            | Por que decidimos assim?            | uma decisão a registrar        | `docs/adr/NNN-<titulo>.md` |
| 3 | PRD            | O que construir e para quem?        | `pre-refinement.md`            | `prd.md`                   |
| — | Design (opcional) | O que o usuário vê?              | `prd.md` + `design-system.md`  | `design.md`                |
| 4 | Tech Alignment | Qual a abordagem técnica?           | `prd.md` (+ `design.md`)       | `tech-alignment.md`        |
| 5 | Tech Spec      | Como exatamente vai funcionar?      | `prd.md` + `tech-alignment.md` | `spec.md`                  |
| 6 | Task Plan      | Em que passos implementar?          | `spec.md`                      | `tasks.md` + `tasks/T*.md` |
| 7 | Debt           | Quais débitos pagar agora?          | `debts.md` da feature          | nova versão (`tasks.md` + `tasks/`) |
| 8 | Execução       | Está implementado e revisado?       | `tasks.md` + `tasks/T*.md`     | código + `debts.md`        |

`vN/` é a pasta da versão da feature: `docs/specs/features/NNN-<feature>/vN/`.

### Task Plan: `tasks.md` + `tasks/`

O Task Plan gera dois níveis:

- **`tasks.md`** é o índice: a lista de tasks, a ordem, as dependências entre elas e o status
  de cada uma. É o que se lê para saber "onde estamos".
- **`tasks/T1.md`, `T2.md`…** traz o detalhe de cada task: o que fazer, arquivos envolvidos,
  os testes (CTs da spec) que ela implementa, as ADRs que ela segue, o **risco** e como saber
  que está pronta. Cada CT da spec cai em exatamente uma task.
- **Risco** é `baixo` ou `alto`. É `alto` a task que mexe em autenticação, autorização,
  pagamentos, dados sensíveis, segredos ou entrada externa; na Execução, ela passa também pela
  revisão de segurança. É assim que o fluxo se ajusta a cada projeto sem mudar as skills.

Uma task por arquivo permite entregar ao Claude (ou a um agent) **só a task da vez**, com o
contexto enxuto, em vez do plano inteiro.

### Execução: executor e revisores

A skill da Execução (`agent-spec-run-tasks`) **não escreve código**: ela coordena os agents,
uma task de cada vez, na ordem das dependências.

```
  task ──► agent-spec-executor ──► agent-spec-reviewer ──┬─ aprovada ───────────► git add, concluída
             (código + testes)      (roda testes,        ├─ só média/baixa ─────► debts.md + git add, concluída
                  ▲                  revisa o diff)       └─ crítica/alta
                  └───────────── correção (até 2) ◄───────┘   3ª reprovação → para e pergunta ao usuário
```

- **Executor** implementa a task seguindo 7 regras (qualidade de sênior sem complexidade
  especulativa, mudanças cirúrgicas, teste primeiro, testes honestos, ADRs como contrato…).
- **Revisor** revisa toda task: roda a suíte, confere cada CT, o escopo, as ADRs, a arquitetura,
  a qualidade e as migrações, e devolve um JSON com cada problema classificado.
- **Revisor de segurança** (`agent-spec-security-reviewer`) entra só nas tasks `Risco: alto`, em
  paralelo com o revisor: procura como abusar do código (autenticação, autorização, entrada,
  dados sensíveis, segredos, pagamentos). Não roda testes. Vale o pior dos dois resultados.
- **Modelo e esforço:** executor e revisor rodam em Sonnet com esforço `high`, fixado no
  frontmatter para não depender do esforço da sessão. O revisor de segurança roda em Opus. A
  última correção do executor (3ª tentativa) roda em Opus.
- **Git:** a task aprovada fica em **stage** (`git add`), nunca commitada. Como o trabalho das
  tasks anteriores já está em stage, `git diff` mostra só o que a task da vez mudou: é assim que o
  revisão isola a task.

### Challenge da spec (opcional)

Entre a Tech Spec e a aprovação dela, a skill utilitária `agent-spec-challenge-spec` faz um
stress-test do `spec.md`: confronta a spec com o código real, as ADRs aceitas (ao pé da letra),
o PRD e o Tech Alignment, procurando colisões com o que já existe, divergências de ADR,
decisões implícitas e casos de borda descobertos. Ela faz **uma pergunta por vez**, sempre com
recomendação, e corrige o próprio `spec.md` (não gera outro documento).

- **Quando usar:** spec grande, de fundação ou que toca área sensível. Em sessão nova, para não
  ser o mesmo contexto que escreveu a spec.
- **Limites:** só edita o `spec.md`; não cria ADR (aponta candidatas em §14); se o furo exige
  regra de negócio nova ou trocar uma `D-NN`, ela para e avisa que a mudança é na fase anterior.
- Não roda depois que a versão começou a ser implementada (alguma task `em andamento` ou
  `concluída`).

### Design system (opcional)

Projeto com interface tem um `docs/specs/design-system.md`: a identidade visual do produto
inteiro (cores, tipografia, espaçamento, ícones, layout, tema, componentes reutilizáveis, padrões
de feedback e acessibilidade visual). Toda feature com telas parte dele, para que cada uma não
reinvente botões, cores e estados. Projeto sem interface não tem esse arquivo.

A skill utilitária `agent-spec-design-system` cria e atualiza o arquivo. O design de cada
feature fica em outro documento, o `design.md` (ver abaixo).

O design system:

- **Projeto com interface existente:** deriva do código, do Figma e dos documentos do time o que
  já existe e só pergunta o que nenhuma fonte responde, ou onde as fontes divergem.
- **Projeto novo:** propõe o design system inicial, com perguntas de direção (biblioteca de UI,
  estilo, cor da marca, tema), sempre com recomendação. Rode antes da primeira feature com telas.
- **Atualização:** só completa o que está vazio, defasado ou novo, com confirmação.

É global como as ADRs: fica fora das features e não tem versão (o histórico fica no git). Nasce
`rascunho`, é aprovado pelo humano e, a cada atualização, volta para `rascunho`. Decisão
transversal com trade-off real (ex.: a biblioteca de UI) também vira ADR.

### Design da feature (opcional)

Feature com telas ganha, entre o PRD e o Tech Alignment, um `design.md` gerado pela skill
`agent-spec-design`: as telas, o layout, os componentes, os **estados visuais concretos**
(carregando, sucesso, erro, vazio), as mensagens literais, os tamanhos de tela, as interações, a
acessibilidade visual e os assets. Ele parte do `design-system.md` (reuso vence criação) e não
decide regra de negócio (é do PRD) nem mecânica técnica (é da Tech Spec). Feature sem telas pula
esta fase.

Quem consome o `design.md`:

- **Tech Alignment:** as soluções técnicas precisam viabilizar o que ele descreve.
- **Tech Spec:** a §7 (Interface) traz só o técnico de cada tela e aponta para o `design.md`, sem
  redefinir o visual.
- **Task Plan:** as tasks de tela citam na linha **Design** as seções que o executor segue.
- **Execução:** o executor segue essas seções e usa os componentes do design system pela
  referência de import; o revisor confere estados, componentes reusados, mensagens e tokens.

O que nasce reutilizável numa feature pode ser **promovido** ao design system, com confirmação:
token e padrão de feedback na hora; componente só depois de existir no código (a próxima
atualização do design system o encontra).

### Commit

O commit é sempre do usuário. Depois da Execução (ou de qualquer mudança feita fora do fluxo),
a skill utilitária `agent-spec-semantic-commit` coloca em stage todas as mudanças do projeto
(`git add -A`, respeitando o `.gitignore`), propõe a mensagem no padrão Conventional Commits, em português, e, depois de aprovada, commita e faz push (o usuário
pode escolher só commitar). Quando o commit fecha tasks do framework, o rodapé
`Refs: 001-<feature>/v1 T1, T2` liga o histórico do git aos documentos da feature. Ela nunca usa
`--amend` nem `push --force`, e se o push falhar ela mostra o erro e para.

### README

A skill utilitária `agent-spec-readme` gera ou atualiza o `README.md` de um projeto, a
documentação para pessoas (o `CLAUDE.md` e as rules são para o agente). Parte de um template
único, com seções por tipo de projeto (lib, cli, app, service), extrai do repositório os comandos
reais e pergunta só o que o repo não responde (sempre o idioma). Se o README já existe, ela
complementa com edições pontuais e mostra o diff antes de gravar. Não cria `LICENSE` nem outros
arquivos de suporte: só avisa que faltam.

### CLAUDE.md

A skill utilitária `agent-spec-claude-md` gera ou audita o `CLAUDE.md` de um projeto, o contrato
do agente: visão, stack, comandos reais, convenções e onde coisas novas vão, em até ~80 linhas.
Todo fato sai de um arquivo do repo ou do usuário; o que não tem fonte é perguntado ou fica de
fora. Ela não repete as regras de comportamento (ficam em `.claude/rules/behavior.md`) e aponta
para as ADRs e o design system em vez de copiá-los. Se o `CLAUDE.md` já existe, ela enxuga e
complementa com edições pontuais e mostra o diff, com a fonte de cada fato, antes de gravar.

### Rules do projeto

Além do `CLAUDE.md`, um projeto pode ter **rules por área** em `.claude/rules/<tema>.md`: as
convenções de um tema (banco, erros HTTP, testes, um módulo). Elas têm `paths` no frontmatter e
só carregam quando o agente lê ou edita um arquivo que casa com eles, então não ocupam o contexto
fora da sua área. Uma rule sem `paths` é global, como o `behavior.md`, que é do framework e não
recebe convenção de projeto.

Uma rule que guia **como o código é escrito** (naming, camadas, erros, testes, contratos) inclui
`docs/specs/**` nos `paths`: a Tech Spec e o Task Plan já decidem nomes e assinaturas, e todas as
fases leem os documentos da feature, então assim a rule carrega na geração da spec e na execução.
Cada regra tem uma linha, o porquê e um micro-exemplo ✅/❌ escrito nela.

Duas skills utilitárias cuidam delas:

- **`agent-spec-rule-create`** estabelece as convenções de **um tema inteiro**: divide o tema em
  3 a 5 eixos, propõe 3 alternativas por eixo com recomendação, deriva do código o que já existe
  ou, em projeto sem código, propõe e marca a rule como `provisória`. Completa uma rule que já
  existe e aponta as decisões que valem ADR.
- **`agent-spec-rule-curate`** decide se **um item** (convenção, decisão, feedback) vira regra,
  com o teste de fricção (alguém faria errado sem ela? dá para descobrir em 1 minuto? fica velha
  em 3 meses? tem o porquê?), e onde ela fica: `CLAUDE.md`, rule global ou rule por área. Também
  audita as rules (desatualizadas, duplicadas, óbvias, sem exemplo, escopo errado) e promove a
  `estável` a rule provisória cujo tema já tem código.

As duas mostram o diff antes de gravar e não tocam código.

### Dúvidas sobre o framework

A skill utilitária `agent-spec-guide` responde dúvidas de uso: qual skill rodar (e, dada uma
feature, qual é o próximo passo dela, pelo status dos documentos), como funciona uma fase e onde
mudar uma convenção. Ela só consulta: não edita nada e não roda a skill que recomenda. Ela
carrega o próprio mapa do framework (fluxo, versões, ADRs, débitos, execução, pastas e
convenções), para continuar funcionando quando este arquivo for removido. Por isso, enquanto ele
existir, uma mudança no fluxo precisa ser feita aqui **e** no guia.

## Estrutura de diretórios

```
ai-software-engineering/
├── CLAUDE.md                  contexto do projeto (lido pelo Claude em toda sessão)
├── WORKFLOW.md                este arquivo
├── README.md
│
├── .claude/
│   ├── skills/                uma skill por fase, mais as utilitárias (ex.: semantic-commit)
│   │   └── agent-spec-<fase>/
│   │       ├── SKILL.md       instruções da skill
│   │       └── template.md    modelo do documento que a skill gera
│   │                          (o Task Plan tem dois: template-tasks.md e template-task.md;
│   │                          o Debt tem esses dois e mais o template-debts.md; a Execução
│   │                          e as utilitárias não geram documento, então não têm template,
│   │                          exceto o design-system, que gera o design-system.md)
│   ├── agents/                subagents, chamados pelas skills
│   │   ├── agent-spec-executor.md   implementa uma task (Execução)
│   │   ├── agent-spec-reviewer.md   revisa uma task (Execução)
│   │   └── agent-spec-security-reviewer.md   revisa a segurança de uma task de risco alto (Execução)
│   └── rules/                 regras carregadas automaticamente
│       ├── behavior.md        como agir: não chutar, escopo, evidência, confirmação (sempre)
│       └── <tema>.md          convenções de uma área do projeto (só ao tocar os paths dela)
│
└── docs/
    ├── references/            o que o usuário escreve ou anexa (entrada do fluxo)
    │   ├── prompts/
    │   │   └── AAAA-MM-DD-<assunto>.md   pedido para o Pre-Refinement
    │   ├── screens/           prints de telas de outros sistemas
    │   └── brand/             logo, fontes, paleta
    │       ├── logo/
    │       └── fonts/
    ├── specs/                 o que as skills produzem (saída do fluxo)
    │   ├── design-system.md   só em projeto com interface: identidade visual do produto
    │   └── features/
    │       └── 001-<nome-da-feature>/
    │           ├── v1/        primeira entrega: todos os documentos
    │           │   ├── pre-refinement.md
    │           │   ├── prd.md
    │           │   ├── design.md      só em feature com telas
    │           │   ├── tech-alignment.md
    │           │   ├── spec.md
    │           │   ├── tasks.md
    │           │   ├── tasks/
    │           │   │   ├── T1.md
    │           │   │   ├── T2.md
    │           │   │   └── T3.md
    │           │   └── debts.md   só existe se a execução deixou débitos
    │           └── v2/        mudança pequena (ou pagamento de débitos): só o que mudou
    │               ├── pre-refinement.md   só na mudança pequena: descreve a mudança
    │               ├── spec.md             só na mudança pequena
    │               ├── tasks.md
    │               └── tasks/
    │                   └── T1.md
    └── adr/                   também saída do fluxo
        └── 001-<titulo-da-decisao>.md
```

### Por que assim

- **Docs agrupados por feature, não por tipo.** Com `docs/prd/`, `docs/specs/` etc., uma
  mesma feature fica espalhada em cinco pastas. Agrupando por feature, toda a história dela
  fica junta, e basta dizer ao Claude "leia `.../001-<feature>/v1/`" para ele ter o contexto
  inteiro daquela versão.
- **`features/` dentro de `specs/`.** É a convenção usada no curso e em outros materiais de
  SDD: `specs/` reúne tudo o que especifica o produto, e `features/` é uma das partes dele.
- **Versões em vez de edição.** Ver [Versões de uma feature](#versões-de-uma-feature).
- **ADR fica fora das features.** Uma decisão de arquitetura costuma valer para o projeto todo
  e sobreviver à feature que a criou. Quando nasce dentro de uma feature, o documento da
  feature apenas aponta para ela.
- **`references/` separa a entrada da saída.** O prompt e os anexos são escritos antes do
  Pre-Refinement, quando a pasta da feature ainda não existe (é a skill que define o `NNN`).
  Por isso ficam fora das features, num lugar só: um print usado por duas features é guardado
  uma vez e citado pelos dois prompts. O prompt leva data no nome, não número de feature, e a
  ponte para a feature é o campo "Fonte da ideia" do `pre-refinement.md`. Logo e fontes que o
  app carrega de fato vão para os assets do código; `brand/` guarda só o material de
  referência.
- **Débito fica dentro da versão.** Ele nasce da execução daquela versão, e é pago abrindo a
  próxima versão da mesma feature, então faz sentido morar junto dela.
- **O template mora dentro da skill**, e não em `docs/templates/`. Assim cada skill é
  autocontida: quem lê a skill acha o modelo do documento ao lado dela.
- **`.claude/rules/`** guarda regras (ex.: convenções de código, padrão de commits). O Claude
  Code carrega esses arquivos automaticamente, como o `CLAUDE.md`: uma rule sem `paths` no
  frontmatter carrega em toda sessão; uma com `paths`, só ao tocar os arquivos dela (ver
  [Rules do projeto](#rules-do-projeto)). Os subagents (executor e revisor) também
  recebem o `CLAUDE.md` e as rules, a menos que o frontmatter do agent tenha
  `omitClaudeMd: true`. As regras de comportamento ficam aqui, e
  não no `CLAUDE.md`, porque fazem parte do framework: o `CLAUDE.md` descreve este projeto, e
  quem copiar a `.claude/` para outro projeto leva as regras junto.
- **Pastas só são criadas quando usadas.** Nada de pastas vazias "para o futuro".

## Convenções

- **Nomes** em kebab-case, sem acentos: `001-login-com-google`.
- **Numeração** sequencial de três dígitos, separada por pasta (`features/`, `adr/`).
  O número nunca é reaproveitado, mesmo se o item for descartado.
- **Versões** da feature: `v1`, `v2`, `v3`… Uma versão implementada não é mais editada.
- **Tasks**: `T1.md`, `T2.md`…, numeradas na ordem do `tasks.md`, recomeçando em `T1` a cada
  versão.
- **IDs rastreáveis**: o PRD numera histórias de usuário (`US-01`), regras de negócio
  (`RN-01`) e critérios de aceite (`CA-01`), o Tech Alignment numera decisões (`D-01`) e a Tech
  Spec numera casos de teste (`CT-01`) e o `debts.md` numera débitos (`DT-01`; fora do próprio
  arquivo, citado com a versão: `v1/DT-01`). As fases seguintes citam esses IDs para mostrar que
  cada requisito foi atendido e de qual decisão uma solução veio. Numeração sequencial, sem
  pular números.
- **Nome da skill**: prefixo `agent-spec-` + nome da fase em kebab-case (ex.:
  `agent-spec-pre-refinement`). O prefixo agrupa as skills do framework na listagem. Exceção: a
  Execução é `agent-spec-run-tasks`, que diz o que ela faz (o nome do professor). Skills
  utilitárias, que não são fase, usam o prefixo + o que fazem (ex.: `agent-spec-semantic-commit`).
- **Nome do agent**: prefixo `agent-spec-` + papel (ex.: `agent-spec-reviewer`).
- **Skills, agents, rules e commands têm nome em inglês** (ex.: `agent-spec-run-tasks`,
  `agent-spec-reviewer`, `behavior.md`). O conteúdo deles continua em português.
- **Status** no topo de cada documento, em frontmatter:

  ```markdown
  ---
  status: rascunho   # rascunho | aprovado
  ---
  ```

  O `design-system.md` usa os mesmos dois. ADRs usam `proposto | aceito | substituído`, cada task (`T*.md`) usa
  `pendente | em andamento | concluída`, e cada débito do `debts.md` usa `aberto | pago`.
- **Próximo passo:** todo documento de fase (pre-refinement, PRD, Tech Alignment…) termina com
  a seção "Próximo passo", com o **comando pronto para copiar** da próxima fase e de cada ADR a
  registrar antes dela. O caminho do fluxo fica no próprio documento, não só no chat. A ADR é a
  exceção: ela vale por anos e não aponta para fase nenhuma.

  ```markdown
  - **Próxima fase:** Tech Alignment → `/agent-spec-tech-alignment docs/specs/features/001-x/v1/prd.md`
  - **ADR a registrar antes:**
    - `/agent-spec-adr "Padrão de logging"` — formato e destino dos logs
  ```
- **Idioma** dos documentos: português, com acentuação correta. Nomes de código ficam como
  estão no código.

## Pontos em aberto

- **Execução — o que ficou de fora do modelo do professor** (trazer só se fizer falta):
  execução de tasks em paralelo; escolha de modelo por task (`model`, `risk`, `gates` no Task
  Plan) e escalonamento por área sensível (do escalonamento ficou só "a última correção roda em
  Opus"); uma segunda revisão separada (QA + Tech Review); doutrina de testes
  completa (uma skill própria, com 29 antipadrões); mineração de regras a partir das revisões.
