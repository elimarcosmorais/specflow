---
name: spec-check
description: Utilitário do framework. O portão de testes - um script (check.sh) que decide o que rodar a partir do git diff e do mapa em .claude/spec-check.json, roda fora da sessão do modelo, guarda pelo hash do código o que já passou e grava o resultado. Tem três modos - passo (testes do passo, build e lint), mudanca (suíte padrão, build e lint inteiros, integração e tela só do que mudou) e versao (tudo inteiro). Chamado pelo spec-executor, pelo spec-lead e pela spec-release, e usado para gerar o spec-check.json de um projeto. Usar para rodar o portão na mão ou para criar ou corrigir o mapa de impacto.
argument-hint: "[passo <comandos> | mudanca | versao | esperar]"
---

# Objetivo

Tirar dos agents a execução dos testes. Quem escreve código não escolhe o que rodar, não vigia
log e não repete o que já passou: chama o script e lê o resultado.

# Como chamar

Sempre da raiz do projeto:

```
sh .claude/skills/spec-check/check.sh passo '<comando de teste>' ['<comando>' …]
sh .claude/skills/spec-check/check.sh mudanca [--arquivos <arquivo> …]
sh .claude/skills/spec-check/check.sh versao
sh .claude/skills/spec-check/check.sh esperar
```

| Modo | Quem chama | Roda |
|---|---|---|
| `passo` | executor, lead | os comandos dados, build e lint |
| `mudanca` | lead, antes da revisão final | lint, build e suíte padrão inteiros; verdes, integração e tela **por impacto** |
| `versao` | `/spec-release` | lint, build, suíte, integração e tela **inteiros**, sem cache |

`--arquivos` é da tentativa 2 da revisão: integração e tela rodam só o que esses arquivos
alcançam, mais o que falhou na tentativa anterior.

| Saída | Significa | O que fazer |
|---|---|---|
| `0` | verde | siga |
| `1` | vermelho | leia o fim da saída, impresso pelo script, e corrija |
| `2` | passou, mas pulou teste | monte a precondição (`behavior.md`, regra 3) e rode de novo |
| `3` | ainda rodando | `check.sh esperar`, quantas vezes precisar |
| `4` | uso errado ou config inválida | corrija o comando ou o `spec-check.json` |

O script espera até 9 min por chamada, dentro do limite da ferramenta. Chame-o com timeout de
600000 ms, sem embrulhar em `timeout`, e sem laço de `grep` ou `sleep` (`behavior.md`, regra 8).

**Cache.** Comando que passou sem pulo não roda de novo sobre o mesmo código. O hash vem da
árvore de trabalho, com o que não está em stage, e ignora `docs/` e `*.md`: marcar um checkbox
não invalida a prova. Falha e pulo nunca entram no cache. O modo `versao` não o usa.

**Resultado.** `.git/spec-check/ultimo.json` traz `status` (`verde`, `vermelho` ou `pulou`) e, por
etapa, `resultado`, `duracao_s`, `alcance`, `comando`, `log` e `pulos`. Resultados de etapa:
`ok`, `cache`, `falhou`, `pulou`, `nenhum_alvo`, `nao_rodou` e `n/a`.

# O `.claude/spec-check.json`

É do projeto: a `/spec-upgrade` não o toca. Tudo nele é fato tirado do repo (`docs.md`), então o
agente o gera e corrige sem perguntar. Só texto e lista de textos; chave desconhecida é erro.

```json
{
  "suite": "go tool task test",
  "build": "go tool task build",
  "lint": "go tool task lint",
  "tela": "go tool task test:e2e",
  "tela_impacto": "go tool task test:e2e -- {alvos}",
  "tela_mapa": [
    "frontend/src/features/{x}/ => e2e/{x}*.spec.ts",
    "e2e/*.spec.ts => e2e/*.spec.ts"
  ],
  "tela_inteira": ["frontend/src/components/ui/", "frontend/src/styles/", "playwright.config.ts", "e2e/support/", "package.json", "package-lock.json"]
}
```

| Chave | O que é |
|---|---|
| `suite`, `build`, `lint` | o comando de cada um; o que o projeto não tem fica de fora |
| `integracao`, `tela` | o comando **inteiro** |
| `integracao_impacto`, `tela_impacto` | o comando com `{alvos}`. Precisa do mesmo ambiente do inteiro: se a task exporta variáveis ou sobe serviços, passe por ela |
| `*_mapa` | `padrão => alvo`. No padrão: `{x}` captura um segmento, `*` não cruza `/`, `**` cruza, e `/` no fim é a pasta inteira. No alvo: `{x}` e `*`; alvo que não existe é ignorado |
| `*_alvos` | comando que lê os arquivos alterados na entrada e escreve um alvo por linha, para o que o mapa não expressa |
| `*_inteira` | padrões que pedem o comando inteiro: componente compartilhado, estilos, config do runner, apoio dos testes, manifesto, lockfile |

Arquivo que não casa com nada não tem alvo: backend sem spec de tela não dispara o teste de tela.
Sem `*_impacto`, o comando inteiro roda na `mudanca`, uma vez.

**Como gerar.** Leia o `CLAUDE.md` (a linha `Integração:` de projeto antigo vai para cá e sai de
lá), os scripts, o `Taskfile` ou `Makefile`, a config do runner e a pasta dos testes.

- O alvo do mapa são os specs de cada tela, e os próprios specs de tela também entram no mapa.
- Toda imagem de referência fotografa a moldura (sidebar, cabeçalho)? → o arquivo da moldura vai em
  `tela_inteira` até as imagens serem recortadas. Sem isso, a imagem velha só apareceria na
  release e travaria a versão.
- **Jest:** `integracao_impacto` com `--findRelatedTests {alvos}` e `integracao_alvos` igual a `cat`.
- **Go:** `integracao_alvos` escreve os pacotes cujo `go list -tags=integration -test -deps` inclui
  o pacote de um arquivo alterado. O `go test` só mostra pulado com `-v` ou `-json`: use um deles
  nos comandos, senão o script não vê o pulo.
- **Dusk:** `php artisan dusk {alvos}`, com o mapa apontando os arquivos de `tests/Browser/`.

Não há como achar os alvos → deixe só o comando inteiro.

# Entrada

$ARGUMENTS
