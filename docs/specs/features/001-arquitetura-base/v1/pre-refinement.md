---
status: aprovado # rascunho | aprovado
---

# Pre-Refinement — Arquitetura base do Dev Tools

> **Legenda:** linha sem marcação = **FATO** (dito pelo usuário) · `[HIPÓTESE]` = inferência
> que precisa ser validada · `[DÚVIDA]` = ponto em aberto (seção 10) ·
> `[fora do escopo do projeto]` = rumo que extrapola o que o projeto se propõe a ser.

- **Feature:** `001-arquitetura-base` · **Versão:** v1 · **Data:** 2026-10-05
- **Fonte da ideia:** `docs/references/prompts/2026-10-05-arquitetura-base.md`

---

## 1. Ideia em uma frase

Criar a fundação do **Dev Tools**, um app desktop para Windows que vai gerenciar ferramentas de
desenvolvimento e servir projetos em localhost (inspirado no FlyEnv, Laragon e Laravel Herd).
A fundação entrega a infraestrutura comum e uma feature real completa (Cadastro de Projetos),
que serve de modelo para todas as próximas.

## 2. Problema

- **Qual é a dor real hoje?** O Dev Tools ainda não existe. Sem uma base definida, cada feature
  futura (PHP, Node.js, Apache, MySQL…) inventaria o seu próprio jeito de tratar logs, erros,
  persistência, idioma e testes, e o sistema ficaria inconsistente e difícil de manter.
  _Exemplo:_ a feature "PHP" grava logs num formato e a feature "Node.js" em outro, e uma
  mensagem de erro aparece traduzida numa tela e só em inglês na outra.
- [HIPÓTESE] A motivação do produto vem da rotina atual do usuário: ele instala cada
  ferramenta à mão em `C:\dev\programs\<programa>\<programa>-vX.X.X` e troca de versão
  recriando um atalho `current`. O Dev Tools automatizaria isso.
- **Quem sente o impacto?** O usuário, como desenvolvedor do Dev Tools (agora) e como usuário
  final dele (depois).
- **Por que resolver agora?** É o primeiro passo do produto: todas as outras features dependem
  dessa base.

## 3. Objetivo

Ao final, existe um Dev Tools **instalável no Windows** que:

- abre como app desktop, sem pedir privilégios de administrador e em instância única;
- mostra o próprio estado de saúde (versão, banco de dados disponível);
- permite trocar tema (Claro, Escuro ou Sistema) e idioma (pt-br ou en), guardando a escolha;
- grava logs úteis tanto para o uso quanto para o desenvolvimento;
- expõe uma API em localhost, documentada e consultável;
- permite cadastrar, listar, editar e remover **Projetos**: a feature-modelo que as próximas
  features copiam (camadas, validação, erros traduzidos, persistência e testes).

## 4. Público

- **Persona primária:** o usuário, desenvolvedor do Dev Tools, que precisa de uma base clara
  para construir as próximas features sempre do mesmo jeito.
- **Persona secundária:** desenvolvedores que usam Windows e querem gerenciar o ambiente local
  de desenvolvimento (público final do Dev Tools; nesta v1, o próprio usuário).
- **Contexto de uso:** computador pessoal com Windows, um único usuário, uso local.
  [HIPÓTESE] Funciona offline.

## 5. Rumos explorados

> Memória do brainstorm: os ramos da Fase 1 e as direções avaliadas na Fase 2, com o que foi
> escolhido e o que foi descartado.
>
> Decisão prévia: este documento fica no **O QUÊ e POR QUÊ**. Stack e arquitetura (linguagem,
> framework HTTP, divisão em camadas) vão para **ADRs** registradas antes do PRD.

### Ramo A — Fatia mínima da fundação

- **A1 — Esqueleto técnico puro:** só infraestrutura, sem funcionalidade de produto.
  _Exemplo:_ tela "Status" com versão, saúde do banco, seletor de tema e idioma.
  _Viabilidade:_ requer algo novo; baixo risco.
- **A2 — Esqueleto + uma feature real fina:** A1 mais uma funcionalidade do Dev Tools que
  atravessa todas as camadas (tela → API → regra → banco).
  _Exemplo:_ A1 + Cadastro de Projetos. _Viabilidade:_ requer algo novo; escopo controlado.
- **A3 — Esqueleto + controle de um serviço:** iniciar e parar um serviço real.
  _Exemplo:_ botão "Iniciar servidor web" (Apache) com status verde/vermelho.
  _Viabilidade:_ requer lidar com processos, portas e binários externos.
- **Escolhida:** A2 — o "handler de exemplo" pedido precisa de regra de negócio de verdade para
  servir de modelo.
- **Descartadas / adiadas:** A1 — o exemplo ficaria vazio de regra de negócio. A3 — adiada:
  mistura a fundação com a parte mais difícil do produto.

### Ramo B — Feature-modelo

- **B1 — Health check como modelo:** _Exemplo:_ `GET /health` devolve status do banco e versão.
  _Viabilidade:_ trivial, mas não exercita escrita, validação nem erros de negócio.
- **B2 — Cadastro de Projetos:** criar, listar, editar e remover projetos (nome, pasta, tipo
  web/android/windows), como a lista de projetos do Laragon.
  _Exemplo:_ cadastrar "loja-api" em `projects/web/loja-api`; erro se a pasta não existir ou o
  nome se repetir. _Viabilidade:_ requer algo novo; exercita validação, persistência, erros
  traduzidos e listagem.
- **B3 — Programas instalados:** ler `bin/` e mostrar programas e versões.
  _Exemplo:_ "php: 8.3.12, 8.4.1". _Viabilidade:_ requer algo novo; pouca regra e pouca escrita.
- **Escolhida:** B2 — é o CRUD que o resto do sistema vai copiar. O health check existe à
  parte, porque foi pedido.
- **Descartadas / adiadas:** B1 — fraco como modelo (continua existindo como serviço de saúde).
  B3 — adiada: boa candidata a uma próxima feature.

### Ramo C — Quem consome a API

- **C1 — Só o próprio front, via ponte do app desktop:** sem HTTP exposto.
  _Exemplo:_ a tela chama funções do app diretamente. _Viabilidade:_ o mais simples, mas torna
  sem sentido a documentação de API e o servidor HTTP pedidos.
- **C2 — API HTTP em localhost, usada pelo front e documentada:** escuta só em `127.0.0.1`.
  _Exemplo:_ uma página "/docs" lista os endpoints. _Viabilidade:_ requer algo novo; abre o
  caminho para uma CLI depois.
- **C3 — API + CLI já na fundação:** _Exemplo:_ `dev-tools status` no terminal, como o `herd`
  do Laravel Herd. _Viabilidade:_ mais uma interface para manter desde o início.
- **Escolhida:** C2 — cobre o que foi pedido sem abrir a frente da CLI.
- **Descartadas / adiadas:** C1 — descartada: contradiz o pedido. C3 — CLI adiada.
- **Esclarecimento:** "servir como servidor web em localhost" significa **as duas coisas**: o
  app expõe a própria API (agora) e, no futuro, serve os projetos do usuário (como o Apache do
  Laragon), em outra feature.

### Ramo D — Instalação e layout em disco

- **D1 — Instalador completo:** escolha de pasta, estrutura, atalhos, desinstalador e
  atualização automática. _Exemplo:_ o instalador do Laragon. _Viabilidade:_ requer ferramenta
  de instalador nova; a atualização automática é uma frente grande sozinha.
- **D2 — Sem instalador na v1:** descompactar e rodar; o app cria as pastas na primeira
  execução. _Exemplo:_ extrair o zip em `C:\dev-tools\programs\` e abrir.
  _Viabilidade:_ o mais simples, mas deixa de fora algo pedido.
- **D3 — Instalador mínimo + estrutura garantida pelo app:** o instalador escolhe a pasta
  (sugerindo `C:\dev-tools\`), copia o executável, cria um atalho e o desinstalador; ao iniciar,
  o app verifica a estrutura (`data/`, `logs/`, `bin/`…) e recria o que faltar.
  _Exemplo:_ apagar `logs/` sem querer e o app recriar no próximo início.
  _Viabilidade:_ requer um instalador novo, porém pequeno.
- **Escolhida:** D3 — atende ao pedido e permite rodar o app sem instalador durante o
  desenvolvimento.
- **Descartadas / adiadas:** D1 — atualização automática e atalhos extras adiados.
  D2 — descartada: o instalador foi pedido explicitamente.

### Ramo E — Relação com o Windows

- **E1 — Portátil:** banco, configurações e logs ficam dentro da pasta do app; nada em
  registro, `%APPDATA%` ou `PATH`. _Exemplo:_ copiar `C:\dev-tools\` para outro PC e funcionar
  igual. _Viabilidade:_ coerente com o esquema de pastas sugerido.
- **E2 — Integrado ao sistema:** dados em `%APPDATA%`, `PATH` global e registro.
  _Exemplo:_ estilo do Laravel Herd. _Viabilidade:_ costuma exigir admin; conflita com o
  esquema de pastas.
- **E3 — Portátil, com integração opcional depois:** E1 por padrão; numa versão futura, opção
  "adicionar ao PATH do usuário" (sem admin). _Viabilidade:_ igual à E1 agora.
- **Escolhida:** E3 — fundação simples, sem fechar a porta para a integração. Junto: o app
  **nunca pede admin** e roda em **instância única** (abrir de novo foca a janela já aberta).
- **Descartadas / adiadas:** E2 — descartada: conflita com a portabilidade e exige admin.
  E1 "pura" — substituída pela E3, que é igual agora e mais flexível depois.

## 6. Escopo

**Dentro** (direções escolhidas):

- [ ] App desktop para Windows que abre sem pedir admin e em instância única.
- [ ] Serviço de saúde da aplicação (versão, banco disponível), visível numa tela e pela API.
- [ ] Persistência local em banco de dados, dentro da pasta do app.
- [ ] Logs do sistema gravados em `logs/`, úteis para o uso e para o desenvolvimento.
- [ ] Idiomas pt-br e en, selecionáveis pelo usuário e guardados entre execuções.
- [ ] Tema Claro, Escuro ou Sistema, selecionável e guardado entre execuções.
- [ ] API HTTP em localhost (só `127.0.0.1`), consumida pelo front e documentada.
- [ ] Cadastro de Projetos (criar, listar, editar, remover) como feature-modelo, com validação
      e mensagens de erro traduzidas.
- [ ] Testes unitários e de integração, com a feature-modelo coberta como exemplo.
- [ ] Separação de responsabilidades entre regra de negócio, infraestrutura e testes (a forma
      exata fica na ADR de arquitetura).
- [ ] Instalador mínimo: escolhe a pasta (sugerindo `C:\dev-tools\`), copia o app, cria atalho
      e desinstalador.
- [ ] O app verifica a estrutura de pastas ao iniciar e recria o que faltar.

**Fora** (com motivo):

- Instalar e gerenciar programas (PHP, Node.js, Apache, MySQL…) — adiado: próximas features.
- Servir os projetos do usuário em localhost (Apache/Nginx, hosts virtuais) — adiado: feature
  futura.
- Iniciar e parar serviços (A3) — adiado: puxa processos, portas e binários externos.
- Listagem de programas instalados em `bin/` (B3) — adiado: candidata a próxima feature.
- CLI no terminal (C3) — adiado.
- Atualização automática do app e atalhos extras (D1) — adiado.
- Opção "adicionar ao PATH do usuário" (E3) — adiado para uma versão futura.
- Dados em `%APPDATA%`, registro, `PATH` global e execução como admin (E2) — descartado:
  conflita com a portabilidade.
- Idiomas além de pt-br e en — não pedido.
- Escolha de stack e de arquitetura — não é descartado: vai para as ADRs (seção 11).

## 7. Ancoramento no projeto

- **O que o projeto é** (CLAUDE.md / README): ambiente de estudo para construir um framework
  próprio de skills e agents no fluxo SDD. O Dev Tools é o produto especificado com esse
  framework.
- **Features relacionadas** (`docs/specs/features/`): nenhuma (esta é a primeira).
- **Capacidades reutilizáveis:** nenhuma ainda; o código do Dev Tools começa do zero.
- **Conflitos ou duplicações:** nenhum.

## 8. Premissas e restrições

**Premissas** (suposições para a ideia fazer sentido):

- [HIPÓTESE] Uso local por um único usuário, em computador pessoal; a API não precisa ser
  acessível pela rede.
- [HIPÓTESE] Um projeto cadastrado apenas **referencia** uma pasta; o cadastro não move nem
  copia arquivos.
- [HIPÓTESE] Instalar em `C:\dev-tools\` não exige admin, porque o Windows permite que um
  usuário comum crie pastas na raiz do `C:\` por padrão. Precisa ser confirmado na máquina-alvo.
- [HIPÓTESE] Tema padrão "Sistema" e idioma padrão igual ao do Windows (en se não for pt-br).

**Restrições já decididas** (prazo, tecnologia imposta, integração obrigatória):

- Plataforma: Windows, como programa executável.
- Stack **sugerida** pelo usuário, a confirmar nas ADRs: Go 1.25+ (motor e CLI), Wails v2 (app
  desktop, WebView2), React 19 + TypeScript + Vite, SQLite via `modernc.org/sqlite` (sem CGO),
  testes com `go test`.
- Esquema de diretórios **sugerido** após a instalação: `C:\dev-tools\projects\{web,android,windows}`
  e `C:\dev-tools\programs\` (com `dev-tools.exe`, `backup/`, `cache/`, `data/`, `lang/`,
  `logs/`, `temp/`, `sessions/` e `bin/<programa>/<programa>-vX.X.X`).
- Idiomas: pt-br e en. Temas: Claro, Escuro e Sistema.
- Instalador com escolha de pasta, sugerindo `C:\dev-tools\`, como o do Laragon.
- Decidido neste refinamento: app portátil, sem admin, instância única.

## 9. Riscos

- **Produto / adoção:** construir infraestrutura demais antes de existir uma feature de
  produto → a base só inclui o que a feature-modelo (Projetos) realmente usa.
- **Escopo** (pode crescer demais?): são 12 itens de escopo; a fundação pode demorar a ficar
  pronta → o Task Plan deve fatiar em entregas que já rodam; se apertar, o instalador pode
  virar a `v2`.
- **Técnico / segurança:** uma API em `127.0.0.1` pode ser chamada por outros programas e por
  páginas abertas no navegador → escutar só em localhost e definir no Tech Alignment como
  restringir quem chama (ex.: token ou checagem de origem).
- **Técnico:** a porta da API pode já estar em uso por outro programa → ver dúvida 2.
- **Técnico:** um instalador sem assinatura digital dispara o aviso do SmartScreen do Windows
  → aceitável na fase de uso pessoal; reavaliar antes de distribuir.
- **Técnico:** desinstalar pode apagar dados do usuário → o desinstalador nunca apaga
  `projects/` sem confirmação (ver dúvida 4).

## 10. Dúvidas em aberto

1. [DÚVIDA] Wails v2 (sugerido) ou v3? Avaliar a maturidade de cada uma na ADR de stack.
2. [DÚVIDA] A porta da API é fixa ou configurável? O que o app faz se ela estiver ocupada?
3. [DÚVIDA] Cadastro de Projetos: o app só registra pastas que já existem ou também cria a
   pasta? Aceita projetos fora de `C:\dev-tools\projects\`?
4. [DÚVIDA] O desinstalador remove `data/`, `logs/` e `projects/`, ou pergunta ao usuário?
5. [DÚVIDA] Logs: o usuário final vê os logs dentro do app ou só nos arquivos? Por quanto
   tempo (ou até que tamanho) eles são guardados?
6. [DÚVIDA] "Testes contra um servidor local que injeta falhas": quais falhas se quer simular
   (banco indisponível, disco cheio, porta ocupada)? Isso entra já na fundação?
7. [DÚVIDA] Onde fica o código do Dev Tools: num repositório próprio, separado deste projeto
   (que guarda só a documentação), ou aqui dentro?

## 11. Próximo passo

- **Tamanho da mudança:** grande — é a primeira capacidade do produto e define o padrão de
  todas as features seguintes.
- **Próxima fase:** PRD → `/agent-spec-prd docs/specs/features/001-arquitetura-base/v1/pre-refinement.md`
- **ADR a registrar antes:**
  - `/agent-spec-adr "Stack base do Dev Tools"` — linguagem, app desktop, front, banco,
    framework HTTP e documentação da API;
  - `/agent-spec-adr "Arquitetura em camadas e padrão de feature"` — negócio, infraestrutura e
    testes; como uma feature é organizada;
  - `/agent-spec-adr "Aplicação portátil"` — dados e caminhos relativos à pasta do app, sem
    admin, instância única.
