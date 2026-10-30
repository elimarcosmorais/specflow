# SpecFlow

> Um framework enxuto de skills e agents do Claude Code para desenvolver software com specs,
> em que você conversa com **um** agente, só entra **antes** da execução, e recebe a tarefa
> **pronta**: implementada, testada, revisada e documentada — sem lista de pendências.

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
      ┌─────────────┐   lê o código e decide
      │    LEAD     │
      │  (a skill)  │   pequeno e claro ──► faz ele mesmo e prova
      │             │
      │             │   com plano: cada passo ──► spec-executor (sessão limpa, até o fim)
      │             │              no fim ──────► spec-reviewer: aprova, ou devolve
      │             │                             os bloqueantes ao executor
      └─────────────┘
            │
            ▼
      relatório final: pronto, ou parado por um motivo da lista
```

Você só entra **antes** da execução: numa rodada de perguntas, cada uma com a opção que o lead
recomenda (sempre em projeto novo, para a plataforma, a stack e o escopo da primeira entrega), e
para aprovar o plano. Depois disso o lead decide, delega, confere e
conserta o que quebra — e só te chama de volta por um dos quatro motivos de uma lista fechada:

| O agente para quando… | Exemplo |
|---|---|
| **exige você fisicamente** | UAC, 2FA, conta externa, testar num aparelho, fornecer imagem ou fonte |
| **é irreversível e fora do plano** | apagar dado de produção, gastar dinheiro, publicar |
| **muda o que o produto entrega** | escolha de produto que o plano aprovado não cobre |
| **falhou 3 vezes no mesmo passo** | já com análise de causa-raiz |

Fora dessa lista, ele decide e segue — e tudo o que decidiu sai no relatório, em "Decisões que
tomei". É ali que você confere de manhã, em vez de ser acordado.

**Tarefa dada é tarefa cumprida.** Não existe débito. O problema que aparece pelo caminho ou
atrapalha a tarefa — e aí é parte dela, e o executor resolve — ou não atrapalha, e não é assunto
de ninguém: não vira nota, fila nem plano para depois. A revisão é binária: aprova, ou devolve
o que um sênior recusaria num PR; "melhoria" e "sugestão" não existem. A única coisa que sobe
ao relatório sem ser pedida é **risco** — segredo, dado pessoal ou de pagamento exposto.

E o plano só é apagado quando **acaba**: todo checkbox marcado, suíte, build e lint inteiros
verdes, sem um único teste pulado, e a revisão final aprovada. Faltou prova, o lead monta a
precondição e roda de novo; quebrou, ele conserta. Plano com passo pendente fica, e `/spec-lead continuar` retoma.

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

- [Claude Code](https://claude.com/claude-code). Os agents usam o modelo e o esforço escolhidos
  nas configurações dele.
- Git: cada passo aprovado fica em stage, e a revisão compara o `git diff`.

## Instalação

Copie a pasta `.claude/` (skills, agents, rules e o `settings.json`) para a raiz do seu projeto,
sem sobrescrever o que ele já tem:

```bash
git clone https://github.com/elimarcosmorais/specflow.git
cp -rn specflow/.claude <seu-projeto>/
```

O projeto já tinha um `.claude/settings.json` → rode `/spec-upgrade` em seguida: ele junta no
seu as regras de permissão do framework, sem mexer no resto.

Uma delas segura o comando que abre o UAC do Windows (`RunAs`) até você aprovar no Claude Code,
mesmo em bypass e mesmo vindo de um subagent: o diálogo do UAC expira em ~2 minutos, e assim ele
só abre com você na frente do PC. Para receber esse aviso no celular, ligue o
[Remote Control](https://code.claude.com/docs/en/remote-control) e ative **Push when actions
required** no `/config`. Aprove só quando estiver no PC: aprovar de longe abre um UAC sem
ninguém para responder.

O `CLAUDE.md` deste repositório descreve este repositório e não vai junto: gere o do seu projeto
com `/spec-claude`.

Para atualizar depois, rode `/spec-upgrade` no seu projeto: ele baixa a versão nova deste
repositório, atualiza o que é do framework e pergunta antes de sobrescrever um arquivo que o
projeto editou. Skills, agents e rules próprios do projeto ficam intocados, e o
`.claude/settings.json` só recebe as regras de permissão do framework que faltam.

## Uso

Dentro do Claude Code, na pasta do seu projeto, descreva a mudança:

```
/spec-lead "o botão de salvar está azul, deveria ser verde"
/spec-lead "quero que o usuário faça login com o Google"
```

Esqueceu o `/spec-lead` ou digitou o nome errado? O pedido vai para ele do mesmo jeito: o
agente carrega a skill sozinho, e a mudança nunca sai sem passar pelo lead.

O primeiro é direto: o lead faz ele mesmo — troca a cor, roda os testes que cobrem aquele ponto
e entrega, sem plano e sem agents. O segundo é planejado: ele pergunta de uma vez só o que é decisão de produto, com a recomendação dele em cada pergunta, escreve o plano e
te mostra:

```
Plano: docs/plans/2026-10-16-login-com-google.md · 6 passos · risco alto no passo 3
Aprovado? Respondo a ajuste e executo em seguida; não preciso de você até o relatório.
```

Aprovou, ele executa tudo — um executor por passo, a prova de cada um conferida, ADRs
registradas e mantidas em dia sem te consultar, specs atualizadas e uma revisão independente da mudança inteira no fim — e volta
com o relatório.

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
| `/spec-lead` | Ponto único de contato: lê o código e decide; a mudança pequena faz sozinha, a maior planeja, delega, confere e entrega o relatório |

Na mudança planejada, ela coordena dois agents, cada um em sessão isolada com contexto limpo:

| Agent | Faz |
|---|---|
| `spec-executor` | Implementa um passo do plano ou a correção dos bloqueantes, até o fim: código, testes que provam, specs e as ADRs que o passo cria, atualiza ou substitui |
| `spec-reviewer` | Revisão final da mudança inteira: roda a suíte, o build, o lint e a integração do que a mudança alcança, confere o plano e, com risco alto, a segurança; aprova ou reprova, sem "melhoria" |

Na revisão, só o que um sênior recusaria num PR volta para correção; o resto ele aceita calado.
O que é decisão de produto, e só isso, chega a você.

**Utilitários:**

| Skill | Para quê |
|---|---|
| `/spec-guide` | Tirar dúvidas: qual skill rodar, próximo passo de um plano, onde fica cada documento |
| `/spec-adr` | Registrar direto uma decisão de arquitetura; `/spec-adr índice` organiza o índice |
| `/spec-design-system` | Criar ou atualizar o design system do produto (obrigatório com interface); em projeto novo, entrega um formulário técnico já pré-preenchido e escreve o design system quando ele volta; também guarda modelos de design systems que já existem (ex.: o do GitHub) e aplica um deles com `modelo <slug>`, sem formulário |
| `/spec-semantic-commit` | Commitar tudo com mensagem em Conventional Commits |
| `/spec-release` | Lançar uma versão: roda a integração inteira, calcula o número pelos commits, atualiza o `CHANGELOG.md`, cria a tag e o GitHub Release, com o instalável (`.exe`, `.apk`) anexado quando o projeto tem um |
| `/spec-readme` | Gerar ou completar o `README.md` |
| `/spec-claude` | Gerar, auditar ou enxugar o `CLAUDE.md`, e acrescentar uma regra a ele (e pôr `.temp/` no `.gitignore`). Manter os **fatos** dele em dia — comando de teste, stack, pastas — é do próprio lead, sem você |
| `/spec-windows-mcp` | Configurar o Windows-MCP para o agente ver e operar a tela (só Windows) |
| `/spec-playwright-mcp` | Configurar o Playwright MCP para o agente ver e operar páginas web |
| `/spec-upgrade` | Trazer as skills, os agents e as rules da versão nova do SpecFlow, preservando o que o projeto editou ou criou |

Argumentos de cada comando: [mapa de roteamento](.claude/skills/spec-guide/SKILL.md#5-mapa-de-roteamento).
