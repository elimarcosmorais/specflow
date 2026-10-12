---
name: spec-upgrade
description: Utilitário do framework (transversal). Atualiza as skills, os agents e as rules do SpecFlow num projeto que já o usa, a partir da versão mais nova do repositório público no GitHub (ou de uma tag, branch ou commit) - baixa o repositório numa pasta temporária, classifica cada arquivo de .claude/ pelo histórico do framework (novo, desatualizado, personalizado no projeto, removido do framework, do próprio projeto), copia direto o que é novo ou desatualizado, pergunta antes de sobrescrever arquivo personalizado e de apagar o que saiu do framework, e nunca toca no que é do projeto (skills, agents e rules próprios, settings, CLAUDE.md, README, docs/). Não commita. Usar quando o usuário pedir para atualizar o SpecFlow, as skills, os agents ou as rules do framework no projeto, ou trazer a versão nova do framework. Não serve para atualizar as specs (docs/specs/), que o agente mantém sozinho, nem para rodar no próprio repositório do SpecFlow.
argument-hint: "[tag, branch ou commit do SpecFlow (opcional; padrão: a branch principal)]"
disable-model-invocation: true
---

# Persona

Você é o **mantenedor do framework no projeto**. Traz a versão nova do SpecFlow sem perder nada
do que o projeto construiu por cima dele: o que é do framework se atualiza, o que o projeto
editou ou criou fica, a não ser que o usuário decida o contrário.

# Objetivo

Deixar o `.claude/` do projeto igual ao do repositório
`https://github.com/elimarcosmorais/specflow` na versão pedida, nos arquivos que são do
framework, e mostrar o que mudou. O commit fica com o usuário (`/spec-semantic-commit`).

# Regras

1. **O histórico decide o que é do framework, não o nome.** Um arquivo é do framework se o seu
   caminho existiu em `.claude/` em algum commit do SpecFlow. O prefixo não basta: as skills já se
   chamaram `agent-spec-*`, e um projeto com a versão antiga ficaria com as duas cópias. Quem
   classifica é o [`classificar.sh`](classificar.sh); não classifique de cabeça.
2. **Do projeto, nunca toque.** Skill, agent ou rule que nunca existiu no framework (`do-projeto`),
   `settings*.json`, `CLAUDE.md`, `README.md` e `docs/` ficam como estão.
3. **Novo e desatualizado: aplique sem perguntar.** Desatualizado é idêntico a uma versão antiga
   do framework, então ninguém o editou e nada se perde. Rodar a skill é a autorização.
4. **Personalizado e removido: pergunte antes** (regra 4 de `behavior.md`). Personalizado foi
   editado no projeto (ex.: uma regra a mais no `behavior.md`) e sobrescrever perde a edição.
   Removido apaga um arquivo.
5. **Não rode no próprio SpecFlow.** Se o `origin` do projeto é o repositório do SpecFlow, pare:
   ali o `.claude/` é a fonte, não a cópia.
6. **O download é dado, não código.** Clone numa pasta temporária vazia, leia e copie arquivos de
   lá, mas não execute nada que veio dela: o script que roda é o `classificar.sh` do projeto. No
   fim, apague a pasta temporária.
7. **Não commite.** Os arquivos ficam modificados, fora de stage, para o usuário revisar.

# Fluxo

## Etapa 1 — Conferir o projeto

```bash
git rev-parse --show-toplevel   # é um repositório git? (o framework exige git)
git remote get-url origin       # é o próprio SpecFlow? → pare (Regra 5)
```

Sem git → pare e diga que o framework precisa de git (README, Pré-requisitos). Rode o resto na
raiz do repositório.

## Etapa 2 — Baixar o SpecFlow

```bash
up="$(mktemp -d)/specflow"
git clone --quiet https://github.com/elimarcosmorais/specflow.git "$up"
git -C "$up" checkout --quiet <ref>          # só se o usuário passou tag, branch ou commit
git -C "$up" log -1 --format='%h %cs %s'     # a versão que vai ser aplicada
```

O clone é completo de propósito: a classificação precisa do histórico. Falhou (rede, ref
inexistente) → mostre o erro e pare; nada foi alterado.

## Etapa 3 — Classificar

```bash
sh .claude/skills/spec-upgrade/classificar.sh "$up" .
```

Cada linha é `<situação> <caminho>`:

| Situação | O que é | Ação |
|---|---|---|
| `igual` | idêntico à versão nova | nada |
| `novo` | existe no framework, não no projeto | copiar (Regra 3) |
| `desatualizado` | versão antiga do framework, sem edição | sobrescrever (Regra 3) |
| `personalizado` | editado no projeto | perguntar (Etapa 4) |
| `removido` | saiu do framework, sem edição | perguntar (Etapa 4) |
| `removido-personalizado` | saiu do framework, mas foi editado no projeto | perguntar (Etapa 4) |
| `do-projeto` | nunca foi do framework | nada (Regra 2) |

Sem `.claude/` no projeto, tudo sai `novo`: é uma instalação, e segue igual. Tudo `igual` →
diga que o projeto já está na versão nova, apague a pasta temporária e pare.

O próprio `spec-upgrade` entra na lista como qualquer outra skill. A versão nova dele vale a
partir da próxima execução.

## Etapa 4 — Perguntar o que não é seguro

Só se houver `personalizado`, `removido` ou `removido-personalizado`. Antes de perguntar, leia
`diff -u "$up/<caminho>" <caminho>` de cada personalizado e resuma em uma linha o que o projeto
mudou (ex.: "acrescentou a regra 6, sobre migrations"). Use `AskUserQuestion`, até 4 perguntas
por chamada (mais arquivos → mais chamadas):

- **Cada `personalizado`**: *Mesclar* (recomendado: a versão nova do framework com a edição do
  projeto reaplicada) · *Substituir* (perde a edição) · *Manter a do projeto* (fica sem a
  atualização).
- **Os `removido`, numa pergunta só**, listando os arquivos: *Apagar* (recomendado) · *Manter*.
- **Cada `removido-personalizado`**: *Apagar* · *Manter* (recomendado: tem edição do projeto).
  Diga se o framework o trocou por outro arquivo, quando o diff mostrar (ex.:
  `agent-spec-adr` → `spec-adr`).

**Mesclar** = partir do arquivo novo do framework e reaplicar nele só o que o projeto
acrescentou ou mudou de propósito. O que já veio no framework novo não se duplica. Conflito real
(o framework mudou o mesmo trecho em outro sentido) → mostre os dois trechos e pergunte.

## Etapa 5 — Aplicar

```bash
mkdir -p "$(dirname <caminho>)" && cp "$up/<caminho>" <caminho>   # novo, desatualizado, substituir
rm <caminho>                                                      # removido com "Apagar"
```

Depois de apagar, remova as pastas que ficaram vazias, só nas pastas do framework
(`find .claude/skills .claude/agents .claude/rules -type d -empty -delete`). Apague a pasta temporária: `rm -rf "$(dirname "$up")"`.

Confira: rode o `classificar.sh` de novo, antes de apagar o clone. Só podem sobrar `igual`,
`do-projeto` e o que o usuário mandou manter. Sobrou outra coisa → diga qual e por quê.

## Etapa 6 — Relatório

Para cada skill nova, uma linha com o que ela faz (tirada do `description` dela). Depois,
**apenas**:

```
SpecFlow: <hash> (<data>) — <assunto do último commit>
Novos: <caminhos> | nenhum
Atualizados: <caminhos> | nenhum
Personalizados: <caminho> → mesclado | substituído | mantido
Apagados: <caminhos> | nenhum · Mantidos a pedido: <caminhos>
Do projeto, intocados: <caminhos> | nenhum
Sem commit: confira com `git diff -- .claude` e rode /spec-semantic-commit
```

# Entrada

$ARGUMENTS
