---
name: spec-check
description: Utilitário do framework. O portão de testes - um script (check.sh) que decide o que rodar a partir do git diff e do mapa em .claude/spec-check.json, roda fora da sessão do modelo, guarda pelo hash do código o que já passou e grava o resultado. Tem três modos - passo (testes do passo, build e lint), mudanca (suíte padrão, build e lint inteiros, integração e tela só do que mudou) e versao (tudo inteiro) - e o auditar, que confere o spec-check.json com o repo (comando inexistente, linha morta do mapa, teste que nenhuma linha alcança). Chamado pelo spec-executor, pelo spec-lead e pela spec-release, e usado para gerar o spec-check.json de um projeto. Usar para rodar o portão na mão ou para criar ou corrigir o mapa de impacto.
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
sh .claude/skills/spec-check/check.sh auditar
```

| Modo | Quem chama | Roda |
|---|---|---|
| `passo` | executor, lead | os comandos dados, build e lint |
| `mudanca` | lead, antes da revisão final | o mapa dos testes do diff; lint, build e suíte padrão inteiros; verdes, integração e tela **por impacto** |
| `versao` | `/spec-release` | lint, build, suíte, integração e tela **inteiros**, sem cache |
| `auditar` | quem gera ou edita o `spec-check.json` | confere o arquivo com o repo (abaixo); não roda teste |

`--arquivos` é da tentativa 2 da revisão: integração e tela rodam só o que esses arquivos
alcançam, mais o que falhou na tentativa anterior.

| Saída | Significa | O que fazer |
|---|---|---|
| `0` | verde | siga |
| `1` | vermelho (no `auditar`: o arquivo tem problema) | leia o fim da saída, impresso pelo script, e corrija |
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
`ok`, `cache`, `falhou`, `pulou`, `nenhum_alvo`, `nao_rodou` e `n/a`. No `mudanca`, a primeira etapa
é `mapa` (abaixo).

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
    "e2e/{x}.spec.ts => e2e/{x}.spec.ts"
  ],
  "tela_testes": ["e2e/*.spec.ts"],
  "tela_inteira": ["frontend/src/components/ui/", "frontend/src/styles/", "playwright.config.ts", "e2e/support/", "package.json", "package-lock.json"]
}
```

| Chave | O que é |
|---|---|
| `suite`, `build`, `lint` | o comando de cada um; o que o projeto não tem fica de fora |
| `integracao`, `tela` | o comando **inteiro** |
| `integracao_impacto`, `tela_impacto` | o comando com `{alvos}`. Precisa do mesmo ambiente do inteiro: se a task exporta variáveis ou sobe serviços, passe por ela |
| `*_mapa` | `padrão => alvo`. No padrão: `{x}` captura um segmento (um por linha), `*` não cruza `/`, `**` cruza, e `/` no fim é a pasta inteira. No alvo: `{x}` é o capturado e `*` é glob; alvo que não existe é ignorado |
| `*_testes` | onde ficam os arquivos de teste daquele tipo. É o que permite achar o teste que nenhuma linha do mapa alcança |
| `*_alvos` | comando que lê os arquivos alterados na entrada e escreve um alvo por linha, para o que o mapa não expressa |
| `*_inteira` | padrões que pedem o comando inteiro: componente compartilhado, estilos, config do runner, apoio dos testes, manifesto, lockfile |

Arquivo que não casa com nada não tem alvo: backend sem spec de tela não dispara o teste de tela.
Sem `*_impacto`, o comando inteiro roda na `mudanca`, uma vez. O próprio `spec-check.json` no
diff também pede integração e tela inteiras: mapa mexido não prova a si mesmo.

# Quem confere o arquivo

O script, nos dois pontos em que um erro do mapa custaria caro:

- **`auditar`**, rodado por quem gera ou edita o arquivo, sai `1` com a lista de problemas:
  comando cujo programa não existe na máquina; linha de mapa que não casa com nenhum arquivo do
  repo, ou que não leva a nenhum teste que existe; padrão de `*_testes` vazio; e **teste órfão**,
  que nenhuma linha do mapa alcança. A linha que liga um teste só a ele mesmo
  (`e2e/{x}.spec.ts => e2e/{x}.spec.ts`) não conta: com ela, o teste só roda quando é editado.
  Com `*_alvos`, quem acha os testes é a ferramenta do ecossistema, e não há órfão a apontar.
- **`mudanca`** confere os testes **do diff**: teste novo ou alterado que o mapa não alcança deixa
  o portão vermelho (etapa `mapa`), antes da suíte. A correção é a linha que falta no mapa. Órfão
  antigo, fora do diff, não trava a mudança: aparece no `auditar`.

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
- **.NET** (WPF, WinForms, WinUI, MAUI, ASP.NET): o alvo é a **pasta** do projeto de teste, e o
  `dotnet test` aceita um projeto por vez, então o comando por impacto é um laço:

  ```json
  {
    "suite": "DOTNET_CLI_UI_LANGUAGE=en dotnet test",
    "build": "dotnet build",
    "lint": "dotnet format --verify-no-changes",
    "integracao": "DOTNET_CLI_UI_LANGUAGE=en dotnet test --filter Category=Integration",
    "integracao_impacto": "for p in {alvos}; do DOTNET_CLI_UI_LANGUAGE=en dotnet test \"$p\" --filter Category=Integration || exit 1; done",
    "integracao_mapa": ["src/{x}/ => tests/{x}.Tests/", "tests/{x}.Tests/ => tests/{x}.Tests/"],
    "integracao_testes": ["tests/*.Tests/**"],
    "integracao_inteira": ["Directory.Build.props", "Directory.Packages.props", "*.sln", "global.json"]
  }
  ```

  Os nomes de pasta e a categoria vêm do repo, não deste exemplo. O `DOTNET_CLI_UI_LANGUAGE=en`
  deixa a saída em inglês (`Skipped:`), que é o que o script reconhece como pulado; em
  português ela sai traduzida e o pulo passaria despercebido. Teste de tela com FlaUI é um
  projeto de teste como outro qualquer: entra em `tela`, `tela_impacto` e `tela_mapa` do mesmo
  jeito.
- **Delphi** (DUnitX): não há como achar os testes de uma unit, então fica só o comando inteiro,
  o que compila e roda o projeto de teste, tirado do script de build do repo. Falha é saída
  diferente de 0; confira se a saída do runner diz os ignorados num formato que o script
  reconhece (`N skipped`, `Skipped: N`) e, se não, configure o logger para isso.

Não há como achar os alvos → deixe só o comando inteiro.

Gravou ou editou o arquivo → `check.sh auditar` até sair `0`.

# Entrada

$ARGUMENTS
