# Workflow do framework

Este arquivo descreve **como o nosso framework funciona**: o fluxo das fases do SDD, que
documento cada fase gera e onde cada arquivo fica. É o mapa que as skills seguem.

## Visão geral do fluxo

```
  ideia / pedido
       │
       ▼
  [1] Pre-Refinement ─────► pre-refinement.md   entender o problema
       │
       ▼
  [3] PRD ────────────────► prd.md              O QUÊ construir e POR QUÊ
       │
       ▼
  [4] Tech Alignment ─────► tech-alignment.md   COMO, em alto nível
       │
       ▼
  [5] Tech Spec ──────────► spec.md             COMO, em detalhe
       │
       ▼
  [6] Task Plan ──────────► tasks.md            visão geral das tasks
       │                    tasks/T1.md, T2.md…  uma task por arquivo
       ▼
  [8] Execução ───────────► código              implementa e revisa task a task
                            debts.md            débitos não bloqueantes encontrados

  Transversais (podem acontecer em qualquer momento):
  [2] ADR  ──► docs/adr/     registra uma decisão importante de arquitetura
  [7] Debt ──► nova versão   lê os debts.md e abre a próxima versão (só tasks) para pagá-los
```

As fases 1, 3, 4, 5, 6 e 8 são uma **sequência**: cada uma usa o documento da anterior como
entrada. ADR e Debt são **transversais**: não são um passo da fila, são disparadas quando
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
| Pequena  | um campo a mais, ajuste de validação         | Tech Spec      | `pre-refinement.md` + `spec.md` + `tasks.md` + `tasks/`  |
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
| 1 | Pre-Refinement | Qual é o problema, de verdade?      | ideia ou pedido bruto          | `pre-refinement.md`        |
| 2 | ADR            | Por que decidimos assim?            | uma decisão a registrar        | `docs/adr/NNN-<titulo>.md` |
| 3 | PRD            | O que construir e para quem?        | `pre-refinement.md`            | `prd.md`                   |
| 4 | Tech Alignment | Qual a abordagem técnica?           | `prd.md`                       | `tech-alignment.md`        |
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
  os testes (CTs da spec) que ela implementa, as ADRs que ela segue e como saber que está
  pronta. Cada CT da spec cai em exatamente uma task.

Uma task por arquivo permite entregar ao Claude (ou a um agent) **só a task da vez**, com o
contexto enxuto, em vez do plano inteiro.

### Execução: executor e revisor

A skill da Execução (`agent-spec-run-tasks`) **não escreve código**: ela coordena dois agents,
uma task de cada vez, na ordem das dependências.

```
  task ──► agent-spec-executor ──► agent-spec-reviewer ──┬─ aprovada ───────────► git add, concluída
             (código + testes)      (roda testes,        ├─ só média/baixa ─────► debts.md + git add, concluída
                  ▲                  revisa o diff)       └─ crítica/alta
                  └───────────── correção (até 2) ◄───────┘   3ª reprovação → para e pergunta ao usuário
```

- **Executor** implementa a task seguindo 7 regras (qualidade de sênior sem complexidade
  especulativa, mudanças cirúrgicas, teste primeiro, testes honestos, ADRs como contrato…).
- **Revisor** é a única revisão: roda a suíte, confere cada CT, o escopo, as ADRs, a arquitetura
  e a qualidade, e devolve um JSON com cada problema classificado.
- **Modelo e esforço:** os dois agents rodam em Sonnet com esforço `high`, fixado no frontmatter
  para não depender do esforço da sessão. A última correção do executor (3ª tentativa) roda em
  Opus.
- **Git:** a task aprovada fica em **stage** (`git add`), nunca commitada. Como o trabalho das
  tasks anteriores já está em stage, `git diff` mostra só o que a task da vez mudou: é assim que o
  revisor isola a task.

### Commit

O commit é sempre do usuário. Depois da Execução (ou de qualquer mudança feita fora do fluxo),
a skill utilitária `agent-spec-semantic-commit` lê o que está em stage, propõe a mensagem no
padrão Conventional Commits, em português, e, depois de aprovada, commita e faz push (o usuário
pode escolher só commitar). Quando o commit fecha tasks do framework, o rodapé
`Refs: 001-<feature>/v1 T1, T2` liga o histórico do git aos documentos da feature. Ela nunca usa
`--amend` nem `push --force`, e se o push falhar ela mostra o erro e para.

## Estrutura de diretórios

```
ai-software-engineering/
├── CLAUDE.md                  contexto do projeto (lido pelo Claude em toda sessão)
├── WORKFLOW.md                este arquivo
├── README.md
│
├── .claude/
│   ├── skills/                uma skill por fase
│   │   └── agent-spec-<fase>/
│   │       ├── SKILL.md       instruções da skill
│   │       └── template.md    modelo do documento que a skill gera
│   │                          (o Task Plan tem dois: template-tasks.md e template-task.md;
│   │                          o Debt tem esses dois e mais o template-debts.md)
│   ├── agents/                subagents, chamados pelas skills
│   │   ├── agent-spec-executor.md   implementa uma task (Execução)
│   │   └── agent-spec-reviewer.md   revisa uma task (Execução)
│   └── rules/                 regras que valem sempre, carregadas automaticamente
│       └── behavior.md        como agir: não chutar, escopo, evidência, confirmação
│
└── docs/                      tudo o que as skills produzem
    ├── specs/
    │   └── features/
    │       └── 001-<nome-da-feature>/
    │           ├── v1/        primeira entrega: todos os documentos
    │           │   ├── pre-refinement.md
    │           │   ├── prd.md
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
    └── adr/
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
- **Débito fica dentro da versão.** Ele nasce da execução daquela versão, e é pago abrindo a
  próxima versão da mesma feature, então faz sentido morar junto dela.
- **O template mora dentro da skill**, e não em `docs/templates/`. Assim cada skill é
  autocontida: quem lê a skill acha o modelo do documento ao lado dela.
- **`.claude/rules/`** guarda regras que valem sempre (ex.: convenções de código, padrão de
  commits). O Claude Code carrega esses arquivos automaticamente, como o `CLAUDE.md` (uma rule
  sem `paths` no frontmatter carrega em toda sessão). Os subagents (executor e revisor) também
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

  ADRs usam `proposto | aceito | substituído`, cada task (`T*.md`) usa
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
