---
status: aceito   # proposto | aceito | substituído
---

# ADR 003 — Aplicação portátil

- **Data:** 2026-10-05
- **Origem:** [pre-refinement](../specs/features/001-arquitetura-base/v1/pre-refinement.md) — seção 5 (rumos D3 e E3) e seção 11
- **Substitui:** nenhuma
- **Substituída por:** nenhuma

## Contexto

O Dev Tools precisa rodar sem privilégios de administrador e funcionar igual se a pasta de
instalação for copiada para outro PC. Isso pede uma regra única para onde ficam banco,
configurações e logs, e como o app acha essa pasta, tanto instalado quanto no desenvolvimento
(`wails dev`) e nos testes. Duas instâncias abertas disputariam o mesmo SQLite e a mesma porta
da API.

## Decisão

Guardar tudo (banco, configurações, logs, cache…) dentro da pasta do app, que é a pasta do
executável (`os.Executable`) ou, se definida, a da variável `DEV_TOOLS_HOME`; nunca usar
registro, `%APPDATA%` nem `PATH`, nem pedir admin. O app roda em instância única
(`SingleInstanceLock` do Wails v2: abrir de novo foca a janela aberta), recria ao iniciar as
pastas que faltarem e, se a pasta não permitir escrita, fecha com uma mensagem clara.

## Consequências

**Ganhos:**

- Copiar a pasta de instalação para outro PC leva o app e os dados juntos.
- Nada a limpar no sistema ao desinstalar além da própria pasta.
- `DEV_TOOLS_HOME` isola desenvolvimento e testes da instalação real.
- Sem disputa pelo banco ou pela porta: só existe um processo.

**Custos:**

- Caminhos que apontam para dentro da instalação precisam ser guardados relativos à pasta do
  app, ou a cópia para outro PC quebra.
- Instalar em pasta protegida (ex.: `C:\Program Files`) não funciona; o instalador deve sugerir
  uma pasta do usuário (`C:\dev-tools\`).
- Todo acesso a arquivo passa por um ponto único que resolve a pasta base (em `platform/`,
  conforme a [ADR 002](002-arquitetura-em-camadas-e-padrao-de-feature.md)).

**Neutras:**

- A opção "adicionar ao `PATH` do usuário" (sem admin) fica para uma versão futura, como
  escolha explícita do usuário.

## Alternativas consideradas

- **Integrado ao sistema** — dados em `%APPDATA%`, registro e `PATH` global, estilo Laravel
  Herd. Rejeitada porque costuma exigir admin e quebra a portabilidade.
- **Cair para `%LOCALAPPDATA%` sem permissão de escrita** — o app continua funcionando.
  Rejeitada porque cria um segundo lugar escondido para dados e logs.
- **Diretório de trabalho atual como pasta base** — usa a pasta de onde o app foi chamado.
  Rejeitada porque muda conforme o app é aberto pelo atalho ou pelo terminal.
- **Várias instâncias** — cada janela é um processo. Rejeitada pela disputa do SQLite e da porta.
