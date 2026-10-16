# SpecFlow

> Um framework enxuto de skills e agents do Claude Code para desenvolver software com specs,
> em que você conversa com **um** agente, é interrompido **uma** vez, e o resto o time de agents
> resolve sozinho.

## Sobre

Pedir uma feature inteira para um agente de IA de uma vez costuma dar errado: ele preenche as
lacunas com suposições e entrega código que ninguém revisou contra um objetivo claro. Mas o
remédio usual é pior: um fluxo que, a cada dúvida, devolve o problema para você. No fim do dia
você não ganhou um programador — ganhou um estagiário que pergunta tudo e deixa uma lista de
pendências que nunca encolhe.

O SpecFlow inverte isso. **Decidir é o trabalho do agente, não um risco a evitar.**

```
você ──► /spec-lead "<pedido>"
            │
            ▼
      ┌─────────────┐   mapeia o código ─► spec-scout
      │    LEAD     │   escreve o plano ─► spec-planner
      │  (a skill)  │   implementa ──────► spec-executor
      │             │   revisa ──────────► spec-reviewer / spec-security-reviewer
      │ decide,     │   registra ADR ────► spec-architect
      │ nunca       │
      │ escreve     │   débito achado ───► plano secundário, sem te perguntar ──┐
      └─────────────┘                                                           │
            ▲                                                                   │
            └──────────────────── mesmo ciclo ──────────────────────────────────┘
            │
            ▼
      relatório final
```

Você entra **uma** vez: para aprovar o plano. Depois disso o lead decide, delega, revisa,
conserta o que quebra e resolve o débito que aparece — e só te chama de volta por um dos quatro
motivos de uma lista fechada:

| O agente para quando… | Exemplo |
|---|---|
| **exige você fisicamente** | UAC, 2FA, conta externa, testar num aparelho, fornecer imagem ou fonte |
| **é irreversível e fora do plano** | apagar dado de produção, gastar dinheiro, publicar |
| **muda o que o produto entrega** | escolha de produto que o plano aprovado não cobre |
| **falhou 3 vezes no mesmo passo** | já com análise de causa-raiz |

Fora dessa lista, ele decide e segue — e tudo o que decidiu sai no relatório, em "Decisões que
tomei". É ali que você confere de manhã, em vez de ser acordado.

**Débito não vira recado.** Problema achado fora do escopo tem três destinos: o executor corrige
junto quando cabe no mesmo diff; o que não cabe vira **plano secundário**, que o lead escreve e
executa sem pedir nada a você; e só o que bate na lista de paradas chega ao relatório — como
proposta pronta, com recomendação e comando. **Não há teto de itens:** dez problemas, dez
resolvidos.

E o plano só é apagado quando **acaba**: todo checkbox marcado, suíte, build e lint inteiros
verdes, sem um único teste pulado. Faltou prova, o lead monta a precondição e roda de novo;
quebrou, ele conserta. Plano com passo pendente fica, e `/spec-lead continuar` retoma.

O que sobrevive à mudança é a documentação viva em `docs/`, que o agente mantém sozinho, junto
com o código:

| Pasta | Guarda |
|---|---|
| `docs/specs/` | como cada módulo funciona hoje; o agente cria, edita e subdivide sem pedir permissão |
| `docs/architecture.md` | a arquitetura base |
| `docs/adr/` | por que cada decisão transversal foi tomada; índice em `docs/adr/README.md` |
| `docs/design-system.md` + `docs/design/` | a identidade visual, obrigatória em projeto com interface |
| `docs/plans/` | os planos em andamento — temporários, apagados quando cumpridos |

As regras de cada pasta estão na rule [`docs.md`](.claude/rules/docs.md); as de comportamento,
em [`behavior.md`](.claude/rules/behavior.md). Mapa completo:
[spec-guide](.claude/skills/spec-guide/SKILL.md).

## Pré-requisitos

- [Claude Code](https://claude.com/claude-code), com acesso aos modelos Sonnet e Opus (os agents
  rodam neles).
- Git: cada passo aprovado fica em stage, e a revisão compara o `git diff`.

## Instalação

Copie a pasta `.claude/` (skills, agents e rules) para a raiz do seu projeto:

```bash
git clone https://github.com/elimarcosmorais/specflow.git
cp -r specflow/.claude <seu-projeto>/
```

O `CLAUDE.md` deste repositório descreve este repositório e não vai junto: gere o do seu projeto
com `/spec-claude`.

Para atualizar depois, rode `/spec-upgrade` no seu projeto: ele baixa a versão nova deste
repositório, atualiza o que é do framework e pergunta antes de sobrescrever um arquivo que o
projeto editou. Skills, agents e rules próprios do projeto ficam intocados.

## Uso

Dentro do Claude Code, na pasta do seu projeto, descreva a mudança:

```
/spec-lead "o botão de salvar está azul, deveria ser verde"
/spec-lead "quero que o usuário faça login com o Google"
```

O primeiro é direto: o lead manda um executor implementar, um revisor conferir, e entrega. O
segundo é planejado: ele pergunta de uma vez só o que é decisão de produto, manda escrever o
plano e te mostra:

```
Plano: docs/plans/2026-10-16-login-com-google.md · 6 passos · 1 parte · risco alto no passo 3
Aprovado? Respondo a ajuste e executo em seguida; não preciso de você até o relatório.
```

Aprovou, ele executa tudo — passo a passo, com revisão independente, ADRs registradas, specs
atualizadas e os débitos resolvidos em planos secundários — e volta com o relatório.

Vai dormir e quer o sistema pronto de manhã? Pule até a aprovação:

```
/spec-lead "<pedido>" --autonomo
```

Nada é commitado sozinho: no fim, `/spec-semantic-commit`. Para lançar uma versão com o que já
foi commitado, `/spec-release`.

## Skills e agents

**O fluxo é uma skill só** — você não precisa conhecer as outras para trabalhar:

| Skill | Faz |
|---|---|
| `/spec-lead` | Ponto único de contato: entende, decide, delega, executa, resolve o débito e entrega o relatório |

Ela coordena seis agents, cada um em sessão isolada com contexto limpo:

| Agent | Faz |
|---|---|
| `spec-scout` | Lê o projeto e devolve o mapa do que o pedido toca; não escreve nada |
| `spec-planner` | Escreve o plano principal, as partes e os planos secundários |
| `spec-executor` | Implementa um passo, uma correção ou um item avulso, com os testes que provam |
| `spec-reviewer` | Roda os testes e revisa o diff do passo; não edita nada |
| `spec-security-reviewer` | Nos passos de risco alto, procura como alguém abusaria do código |
| `spec-architect` | Registra a ADR de uma decisão já tomada e mantém o índice |

Na revisão, só o que um sênior recusaria num PR volta para correção. O resto vira **melhoria**,
e melhoria não é pendência: entra na fila de débito e é feita num plano secundário. O que é
decisão de produto, e só isso, chega a você.

**Utilitários:**

| Skill | Para quê |
|---|---|
| `/spec-guide` | Tirar dúvidas: qual skill rodar, próximo passo de um plano, onde fica cada documento |
| `/spec-adr` | Registrar direto uma decisão de arquitetura; `/spec-adr índice` organiza o índice |
| `/spec-design-system` | Criar ou atualizar o design system do produto (obrigatório com interface); em projeto novo, entrega um formulário técnico já pré-preenchido e escreve o design system quando ele volta |
| `/spec-semantic-commit` | Commitar tudo com mensagem em Conventional Commits |
| `/spec-release` | Lançar uma versão: calcula o número pelos commits, atualiza o `CHANGELOG.md`, cria a tag e o GitHub Release |
| `/spec-readme` | Gerar ou completar o `README.md` |
| `/spec-claude` | Gerar, auditar ou enxugar o `CLAUDE.md`, e acrescentar uma regra a ele (e pôr `.temp/` no `.gitignore`). Manter os **fatos** dele em dia — comando de teste, stack, pastas — é do próprio lead, sem você |
| `/spec-rule-create` | Definir as convenções de um tema (banco, erros HTTP…) em uma rule |
| `/spec-rule-curate` | Decidir se um item vira regra, e auditar as rules |
| `/spec-windows-mcp` | Configurar o Windows-MCP para o agente ver e operar a tela (só Windows) |
| `/spec-playwright-mcp` | Configurar o Playwright MCP para o agente ver e operar páginas web |
| `/spec-upgrade` | Trazer as skills, os agents e as rules da versão nova do SpecFlow, preservando o que o projeto editou ou criou |

Argumentos de cada comando: [mapa de roteamento](.claude/skills/spec-guide/SKILL.md#5-mapa-de-roteamento).
