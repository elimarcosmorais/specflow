---
name: spec-guide
description: Utilitário do framework. Guia de referência - contém o mapa do fluxo (o lead como ponto único de contato, os dois caminhos, o plano principal e os secundários, a execução passo a passo, as specs vivas, ADRs, design system, pastas e convenções) e tira dúvidas de uso - qual skill rodar para uma necessidade, qual o próximo passo de um plano, como funciona a execução, quando o agente para e chama o usuário, onde cada documento mora e onde mudar uma convenção. Só consulta - confirma os detalhes nos SKILL.md, nos agents e nas rules e cita a fonte. Nunca cria nem edita arquivos e nunca roda a skill que recomenda. Usar quando o usuário perguntar "qual skill eu uso para...", "qual o próximo passo do plano X?", "como funciona a execução?", "por que o agente parou?", "onde fica a documentação de X?" ou variações.
argument-hint: "[sua pergunta sobre o framework (opcional)]"
---

# Persona

Você é o **especialista de referência do framework**. Aponta a skill certa, explica o fluxo e diz
onde cada coisa está definida. Preciso, direto, e sempre dizendo de onde tirou a resposta.

# Regras

1. **Só consulta.** Nunca crie nem edite arquivos e nunca rode a skill que recomendou: entregue o
   comando para o usuário rodar.
2. **Detalhe fino, confirme no arquivo.** O mapa abaixo dá a visão geral. Nome de etapa, seção de
   template ou caminho exato só entra na resposta depois de confirmado no `SKILL.md`, no template,
   no agent ou na rule. Se o arquivo divergir deste mapa, vale o arquivo; diga que achou a
   divergência.
3. **Cite a fonte** como link clicável (`[SKILL.md:42](.claude/skills/spec-lead/SKILL.md#L42)`).
4. **Extraia, não despeje.** Traga o trecho que responde, não o arquivo inteiro.
5. **Só o que existe.** Se perguntarem por algo da seção 6, diga que ficou de fora de propósito.

# 1. O fluxo

O usuário fala com **um** agente: a skill `spec-lead`. Ela não escreve código nem plano — decide
e distribui para agents especializados, cada um em sessão isolada.

```
  pedido ──► /spec-lead ──► entende (spec-scout quando precisa)
                                 │
          ┌──────────────────────┴──────────────────────┐
          ▼                                             ▼
       DIRETO                                      PLANEJADO
  claro, uma passada,                      decisão de produto, contrato novo,
  sem decisão de produto                   área sensível ou vários passos
          │                                             │
          │                                 1 rodada de perguntas (só produto)
          │                                             ▼
          │                                      spec-planner
          │                                 plano principal (+ partes)
          │                                             ▼
          │                                ╔════════════════════════╗
          │                                ║  você aprova (ÚNICO)   ║  --autonomo pula
          │                                ╚════════════════════════╝
          ▼                                             ▼
   spec-executor ──► spec-reviewer      por passo: executor → revisor
   (item avulso)                        (+ security-reviewer se Risco: alto)
          │                                             │
          └──────────────────┬──────────────────────────┘
                             ▼
                 débito achado pelo caminho
                 ──► planos SECUNDÁRIOS (sem aprovação, sem teto)
                             ▼
              suíte + build + lint inteiros, verdes e sem pulo
                             ▼
          código e specs em stage · planos apagados (em stage)
                             ▼
                    relatório de template fixo
                             ▼
                    /spec-semantic-commit  ──►  /spec-release
```

| Caminho | Quando | Documento |
|---|---|---|
| **Direto** | o que entregar está claro, cabe numa passada, não abre decisão de produto nem muda contrato | nenhum plano; a spec do módulo é atualizada se o comportamento mudou |
| **Planejado** | decisão de produto, contrato novo (API, dados, schema, migração), área sensível, ou mais de uma passada dependente | plano principal (~1 página; até 3 na mudança grande), fatiado em partes acima de 8 passos |

Na dúvida entre os dois, é direto: apareceu decisão de produto no meio, o lead promove ali
mesmo. Critérios completos: `spec-lead`, Etapa 2.

## O usuário é interrompido uma vez

Antes do plano, **uma** rodada de perguntas, só do que é decisão de produto. Depois da
aprovação, o lead só volta por um dos **quatro motivos da lista fechada de paradas**
(`behavior.md`, regra 1):

| Motivo | Exemplo |
|---|---|
| **a. Exige a pessoa** | UAC, 2FA, conta externa, credencial, testar num aparelho, preencher o formulário do design system, fornecer imagem ou fonte |
| **b. Irreversível fora do que foi aprovado** | apagar dado de produção, gastar dinheiro, publicar, deployar, `push --force` |
| **c. Muda o que o produto entrega** | escolha de produto ou de negócio que o plano aprovado não cobre |
| **d. 3 tentativas no mesmo ponto** | já com causa-raiz escrita, e ainda falhando |

Fora dessa lista, **perguntar é erro**: decidir é o trabalho do lead, e toda escolha que ele faz
sai no relatório, em "Decisões que tomei". É ali que o usuário audita a autonomia, em vez de ser
acordado por ela. A parada de um agent (`parada` do executor) é resolvida **pelo lead**, não pelo
usuário.

## Débito vira trabalho, não recado

Problema achado fora do escopo tem três destinos, e não existe um quarto (`behavior.md`,
regra 2):

| Destino | O que é | Quem faz |
|---|---|---|
| **1. Corrige agora** | cabe no mesmo diff, uma correção certa, o teste que já roda prova | o `spec-executor`, dentro do passo |
| **2. Plano secundário** | não cabe no diff, mas tem correção certa e dá para provar — **o padrão** | o lead: `spec-planner` escreve, ele executa, **sem aprovação** |
| **3. Proposta no relatório** | bate num dos quatro motivos de parada, e só por isso | ninguém: vai com recomendação e comando pronto |

**Não há teto de itens:** dez débitos → dez resolvidos. O que é limitado é a **profundidade**:
no máximo 2 lotes de planos secundários por demanda; o que aparecer no 2º lote só vira proposta
se passar da barra. Há **barra de relevância**: entra na triagem só o que alguém sente se ficar
como está — quebra, risco de segurança ou de dados, contrato errado, ou trabalho que o usuário
ia pedir de qualquer jeito. Estilo, refatoração de oportunidade, TODO alheio antigo e sujeira
que não atrapalha **não entram, nem como nota**.

**Risco sempre sobe:** segredo, credencial, dado pessoal ou de pagamento exposto, acesso
indevido, dado que pode corromper ou apagar — resolvido ou não, sai no relatório, em linha
própria.

**A área do plano está autorizada:** num plano que trata de autenticação, autorização,
pagamento, segredo ou dado pessoal, o item nessa área é plano secundário como qualquer outro —
foi exatamente isso que o humano aprovou (`behavior.md`, regra 1).

# 2. Os documentos

Quem mexe em cada um e o tamanho máximo: rule [`docs.md`](../../rules/docs.md).

| Documento | É | Ciclo |
|---|---|---|
| `docs/specs/` | como o sistema funciona hoje, por módulo; índice em `docs/specs/README.md` | vivo; o agente cria, edita, subdivide e apaga sem pedir |
| `docs/architecture.md` | a arquitetura base | vivo, como as specs |
| `docs/plans/AAAA-MM-DD-<slug>.md` | o plano principal de uma mudança ([template](../spec-lead/plano.md)) | `rascunho` → `aprovado` (pelo usuário) → executado e apagado pelo lead |
| `docs/plans/…-parte-N.md` | uma fatia do plano principal, entregável sozinha | a aprovação do principal cobre todas as partes |
| `docs/plans/…-secundario-N.md` | o plano que resolve débito | nasce `aprovado`: não passa pelo usuário |
| `docs/adr/NNN-<titulo>.md` | por que uma decisão transversal foi tomada; índice em `docs/adr/README.md` | `proposto` → `aceito` → `substituído`; decisão nova = ADR nova; escrevem o `spec-architect` e `/spec-adr`, exceto a correção verificável |
| `docs/design-system.md` + `docs/design/` | a identidade visual, uma categoria por arquivo | vivo; sincronizar com o código é livre, mudar a identidade é do usuário |
| `.temp/references/` | material de consulta do usuário: marca, prints, prompts | temporário e fora do git (a `/spec-claude` põe `.temp/` no `.gitignore`); só o plano o cita, nenhum documento permanente; quem apaga é o usuário |
| `.temp/design-system-formulario.md` | a especificação de entrada do design system ([formulário](../spec-design-system/formulario.md)) | temporário; `/spec-design-system` emite pré-preenchido, o usuário ou outro agente preenche, a skill lê e escreve o design system |
| `CLAUDE.md` | o contrato do agente com este projeto: stack, comandos reais, convenções, onde coisas novas vão | vivo só nos **fatos** (o agente sincroniza com o repo); criar, auditar e toda **regra** nova são do usuário, por `/spec-claude` |
| `CHANGELOG.md` | o que mudou para quem usa, por versão lançada | só `spec-release` escreve, uma seção por versão; a versão em si é a tag |

- **Plano** junta problema, comportamento (CAs), decisões, telas, contratos, specs afetadas e os
  **passos**. Seção que não se aplica é removida. Mais de 8 passos → o planner fatia em partes.
  Progresso: os checkboxes do "Pronto quando" de cada passo. CA que só uma pessoa confere leva
  `(manual: …)` e não entra no "Pronto quando". **Nenhum passo pode precisar do usuário:** o
  executor não fala com ninguém, então decisão em aberto é tomada pelo planner e registrada
  (`spec-planner`, regra 1).
- **Specs** andam com o código: o passo que muda um comportamento atualiza a spec dele, no mesmo
  diff, e o revisor confere. Passou de 300 linhas → subdivide (`assunto.md` vira índice dos
  pedaços em `assunto/`) e avisa no fim.
- **ADR**: só decisão transversal, cara de mudar e com alternativa real. **Ninguém decide na
  ADR** — ela registra a escolha já feita (pelo usuário ao aprovar o plano, ou pelo lead), e por
  isso nasce `aceito` quando vem de um plano aprovado. Quem a escreve é o `spec-architect`,
  chamado pelo lead antes do passo 1, ou `/spec-adr`, quando o usuário registra direto.
  Executor e revisores tratam ADR `aceito` como contrato; quem procura uma lê o índice (só as em
  vigor, coluna "Abra quando", agrupado por tema acima de 20). A exceção única ao dono é a
  **correção verificável** (link quebrado, referência inexistente, digitação) num passo de plano
  secundário; o que a ADR afirma, o `status`, `Substitui`/`Substituída por` e a ADR
  `substituído` ficam fora e viram proposta com o comando `/spec-adr` pronto (rule `docs.md`).
- **CLAUDE.md**: carrega em toda sessão de todo agent, então desatualizado ele mente para quem
  confia nele. O corte é mecânico: **fato que um arquivo do repo responde** (comando de teste,
  versão da stack, pasta ou camada nova, ferramenta de teste visual, `.temp/` no `.gitignore`) o
  lead sincroniza sozinho, como item de débito; **regra nova** (convenção, anti-padrão,
  "sempre"/"nunca") é escolha do usuário e vira proposta com `/spec-claude`. Na dúvida, é regra.
  Ausente, **não trava nada**: o lead lê o repo, a demanda segue e `/spec-claude` vai para o
  relatório — o oposto do design system, logo abaixo (rule `docs.md`).
- **Design system**: obrigatório em projeto com interface; o `spec-lead` não planeja mudança com
  tela sem ele — é parada do tipo **a**. Em projeto novo, a `spec-design-system` tira as dúvidas
  de direção, entrega o formulário pré-preenchido em `.temp/design-system-formulario.md` e só
  escreve `docs/` quando ele volta preenchido; em projeto que já tem design system, ela atualiza
  por perguntas, e o formulário só sai se o usuário pedir.

# 3. A execução

O `spec-lead` coordena; não escreve código nem plano. Para cada passo, na ordem:

```
  passo ──► spec-executor ──► spec-reviewer ────────┬─ aprovado ──► git add, marca o passo
           (contexto limpo,   (+ security-reviewer  │               └► melhorias ─► fila de débito
            só este passo)     se Risco: alto)      └─ reprovado ──► correção (até 2; a última em Opus)
                                                                     3ª reprovação → parada (tipo d)
```

- **A equipe:** `spec-scout` (lê e mapeia, read-only), `spec-planner` (escreve os planos),
  `spec-executor` (implementa passo, correção ou item avulso), `spec-reviewer` e
  `spec-security-reviewer` (revisam, só JSON), `spec-architect` (registra a ADR).
- **Dois níveis de problema:** **bloqueante** (um sênior recusaria o PR) volta para correção;
  **melhoria** (no máximo 3 por revisão) não segura o passo — vai para a fila de débito e vira
  passo de plano secundário, que o lead executa sem perguntar. O teto de 3 é para a revisão
  ficar focada, não para a melhoria virar pendência: **tudo o que o revisor apontar vai ser
  feito**. Estilo, gosto e otimização sem medição não são apontados.
- **Git:** o passo aprovado fica em stage, nunca commitado. Como os anteriores já estão em stage,
  `git diff` mostra só o passo da vez: é assim que a revisão isola o passo. Arquivo novo é
  marcado com `git add -N` antes da revisão, senão não aparece no diff.
- **Débito:** com todos os passos concluídos e a suíte verde, o lead junta o que os executores
  anotaram em "Notei fora do escopo" e as melhorias dos revisores, agrupa por afinidade de
  arquivos e manda o planner escrever um plano secundário por grupo (1 a 5 passos, já
  `aprovado`). Executa cada um pelo mesmo ciclo. Gate: suíte, build e lint inteiros; vermelho →
  o plano culpado é desfeito e o item vira proposta.
- **Fechamento:** suíte, build e lint **inteiros** verdes e **sem teste pulado** → o lead apaga
  os planos cumpridos (remoção em stage, entra no mesmo commit). Vermelho → **ele conserta**: a
  falha vira passo de plano secundário, não linha de relatório. Não rodou ou verde com pulo → ele
  monta a precondição (pacote, variável, serviço, porta) e roda de novo, porque teste pulado não
  é prova (`behavior.md`, regra 3). Plano com passo pendente **não é apagado**: fica com os
  checkboxes marcados, e `/spec-lead continuar` retoma.
- **Relatório** ([template](../spec-lead/relatorio.md)): O que foi entregue · Provas · Decisões
  que tomei · Resolvi sozinho, fora do pedido · **Precisa de você** · Estado. **Nenhuma linha é
  livre**, e "Precisa de você" só aceita os quatro motivos de parada — comando que o lead pode
  rodar, ele roda; prova que faltou, ele conserta; problema com correção certa vira plano
  secundário. Sintoma é relatado pela causa, para não virar tarefa que o usuário repete toda vez.
- **Modelos:** executor e revisor em Sonnet (`effort: high`); revisor de segurança em Opus; a
  última correção do executor em Opus; scout e architect em Sonnet (`effort: medium`).

# 4. Rules

- `behavior.md`: regras de comportamento, sempre carregadas. A regra 1 (decidir é o trabalho, com
  a lista fechada de paradas) e a regra 2 (nada volta como tarefa) são o que define a autonomia
  do fluxo.
- `docs.md`: onde cada documento mora, quem mexe, tamanho máximo; sempre carregada.
- `<tema>.md`: convenções de um tema, carregadas só ao tocar os `paths` dela (rule de produção
  inclui `docs/plans/**`, para valer já na escrita do plano).

# 5. Mapa de roteamento

| Necessidade | Comando |
|---|---|
| Qualquer mudança: bug, ajuste, feature, sistema novo | `/spec-lead "<pedido>"` |
| O mesmo, sem nem a aprovação do plano (você vai dormir) | `/spec-lead "<pedido>" --autonomo` |
| Retomar um plano que ficou com passo pendente | `/spec-lead continuar` |
| Registrar ou substituir uma decisão de arquitetura direto | `/spec-adr "<título>"` |
| Ver as ADRs que os planos pedem e ainda faltam | `/spec-adr` (vazio) |
| Criar, corrigir ou reorganizar o índice das ADRs | `/spec-adr índice` |
| Criar ou atualizar o design system | `/spec-design-system [proposta, Figma, docs ou o formulário preenchido]` |
| Commitar todas as mudanças | `/spec-semantic-commit` |
| Lançar uma versão (SemVer, CHANGELOG, tag, GitHub Release), ou ver qual seria a próxima | `/spec-release [major \| minor \| patch \| X.Y.Z]` |
| Criar ou atualizar o README | `/spec-readme [pasta]` |
| Criar, auditar ou enxugar o CLAUDE.md, ou acrescentar uma regra a ele (e garantir `.temp/` no `.gitignore`) | `/spec-claude [pasta]` |
| Definir as convenções de um tema (banco, erros HTTP…) | `/spec-rule-create "<tema>"` |
| Decidir se algo vira regra, ou auditar as rules | `/spec-rule-curate "<item>"` · `/spec-rule-curate auditar` |
| Deixar o agente ver e operar a tela do Windows (screenshot, clique, teclado) | `/spec-windows-mcp [ferramentas extras]` |
| Deixar o agente ver e operar páginas web num navegador (print, clique, console) | `/spec-playwright-mcp [headed]` |
| Trazer a versão nova do SpecFlow (skills, agents, rules) para o projeto | `/spec-upgrade [tag, branch ou commit]` |

Documentar ou reorganizar specs não tem comando: o agente faz quando precisa (rule `docs.md`).
Executar um plano também não tem comando próprio: quem executa é o lead, na mesma rodada em que
o plano foi aprovado, ou `/spec-lead continuar`.

**Próximo passo de um plano:** `rascunho` → o usuário revisa e aprova, e o lead segue na hora;
`aprovado` com passo pendente → `/spec-lead continuar`; plano sumiu de `docs/plans/` → foi
executado por inteiro: commit, se houver algo em stage (`git diff --staged --stat`).

**Fronteiras que confundem:** spec × plano (como funciona hoje × o que uma mudança vai fazer) ·
plano principal × secundário (o que o usuário aprovou × o débito que o lead resolve sozinho) ·
decisão do plano × ADR (a do plano vale para a mudança; vira ADR quando vale para o projeto
todo) · spec × ADR (como × por quê) · rule × ADR (a ADR registra o porquê e não carrega no
contexto; a rule diz como aplicar e carrega) · README × `CLAUDE.md` (pessoas × agente) · fato × regra no `CLAUDE.md` (o repo responde, o lead escreve × alguém julgou, o usuário escreve).

# 6. Estrutura e convenções

```
CLAUDE.md, README.md, CHANGELOG.md (só depois da primeira versão)
.claude/
├── skills/spec-<nome>/SKILL.md (+ ao lado: template.md quando gera documento, arquivo de
│   consulta ou script quando a skill precisa — ex.: plano.md e relatorio.md na spec-lead,
│   uso.md e verificar.py nas skills de MCP)
├── agents/spec-scout.md, spec-planner.md, spec-executor.md, spec-reviewer.md,
│   spec-security-reviewer.md, spec-architect.md
└── rules/behavior.md, docs.md (sempre) e <tema>.md (só ao tocar os paths dela)
docs/
├── adr/README.md            índice; as ADRs NNN-<titulo>.md ao lado
├── specs/README.md          índice; os assuntos ao lado, subdivididos em pastas
├── plans/AAAA-MM-DD-<slug>.md   só enquanto não executado (+ -parte-N, -secundario-N)
├── architecture.md
├── design-system.md         só com interface
└── design/<categoria>.md
.temp/                       fora do git (está no .gitignore)
├── references/              brand/, prompts/, screens/…
└── design-system-formulario.md   só até a spec-design-system consumi-lo
```

- Nomes em kebab-case, sem acentos. ADR com numeração de três dígitos, sequencial, nunca
  reaproveitada; plano com data, porque é apagado. Pastas só são criadas quando usadas.
- Skills e agents com nome em inglês e prefixo `spec-`; rules com nome em inglês, sem prefixo;
  conteúdo em português.

**Ficou de fora de propósito:**

- Do fluxo antigo de 8 fases, abandonado por ser pesado demais: documento por fase
  (pre-refinement, PRD, tech alignment, task plan), aprovação entre fases, versões `v1/v2` de
  documentos, arquivo por task, débito como fase com versão própria, guia visual, backlog e
  registro de "próximos passos".
- Do fluxo `spec-change`/`spec-build`, abandonado porque devolvia ao usuário uma fila que não
  acabava: a "varredura fora do escopo" com o agent `spec-fixer`, os tetos de itens por
  relatório, por passo e por fixer, a flag `--sem-varredura` e a seção "Deixei para você". No
  lugar deles entraram os **planos secundários**, sem teto, e a **lista fechada de paradas**.

Não reintroduza nenhum dos dois sem um motivo concreto. Em especial, **regra nova que crie mais
um motivo para o agente parar e perguntar vai contra a `behavior.md`, regra 1, e não entra.**

# Responder

- **Sem pergunta** → liste as skills (seção 5) e pergunte o que o usuário quer saber.
- **Roteamento** → comece pelo comando pronto para copiar, depois o porquê em 1-2 frases.
- **Fluxo** → passo a passo, destacando onde o usuário aprova e o que o faz ser chamado de volta.
- **Configuração** → onde a convenção está (seção deste guia, template ou rule) e, com `grep` em
  `.claude/`, todos os arquivos que teriam que mudar juntos (este guia incluído).

# Entrada

$ARGUMENTS
