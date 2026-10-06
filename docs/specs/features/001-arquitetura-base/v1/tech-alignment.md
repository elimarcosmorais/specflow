---
status: aprovado   # rascunho | aprovado
---

# Tech Alignment — Arquitetura base do Dev Tools

- **Feature:** `001-arquitetura-base` · **Versão:** v1 · **Data:** 2026-10-05
- **Origem:** [prd.md](prd.md) · [pre-refinement.md](pre-refinement.md)
- **ADRs consultadas:**
  - [`docs/adr/001-stack-base-do-dev-tools.md`](../../../../adr/001-stack-base-do-dev-tools.md)
  - [`docs/adr/002-arquitetura-em-camadas-e-padrao-de-feature.md`](../../../../adr/002-arquitetura-em-camadas-e-padrao-de-feature.md)
  - [`docs/adr/003-aplicacao-portatil.md`](../../../../adr/003-aplicacao-portatil.md)

---

## 1. Contexto técnico

Esta é a feature de fundação: ainda não existe código. A stack (Go + Wails v2 + React + SQLite
sem CGO + `net/http` + OpenAPI escrito à mão), a organização do motor em pacotes por feature e a
portabilidade (tudo na pasta do app, instância única, sem admin) já estão decididas nas ADRs
001, 002 e 003, e aqui são tratadas como restrições.

O que ficou em aberto para esta fase são os pontos que o PRD marcou como `[DECISÃO TÉCNICA]` e
os padrões transversais que toda feature futura vai copiar: como o front se conecta ao motor e
prova que é o Dev Tools, como erros e idiomas fluem entre as camadas, onde ficam as
preferências, como os logs e o esquema do banco evoluem, como o app se comporta (e é testado)
com o banco fora do ar, qual ferramenta gera o instalador e como o front é organizado.

Dois casos do PRD moldam várias decisões: com a **porta da API ocupada** (RN-07) o front não
alcança a API, e com o **banco indisponível** (CA-12) a tela de status ainda precisa aparecer no
idioma escolhido. Por isso, inicialização e preferências não podem depender nem da porta TCP nem
do SQLite.

**Invariantes** (o que qualquer implementação precisa preservar):

- Nenhum processo do app, do instalador ou do desinstalador pede elevação de privilégio
  (RN-01, CA-01, CA-04).
- Toda leitura e escrita em disco do app parte da pasta base, inclusive a pasta de dados do
  WebView2, que por padrão iria para `%APPDATA%` (RN-03, CA-10).
- A API escuta só na interface de loopback (RN-08, CA-34).
- Operações que alteram Projetos só são aceitas com prova de que vêm do próprio Dev Tools
  (RN-08, CA-33).
- O app abre e mostra o próprio estado mesmo com a porta da API ocupada ou o banco indisponível;
  só a pasta sem permissão de escrita impede a abertura (RN-04, RN-07, CA-06, CA-07, CA-12).
- Nada em tempo de execução depende da internet: assets, fontes e documentação da API vão
  embutidos no executável (RN-06, CA-08).
- O núcleo de Projetos (entidade e serviço) não importa banco, HTTP nem Wails (CA-38, CA-39).
- Criar, editar ou remover um projeto nunca cria, move, copia ou apaga nada no disco (RN-16,
  CA-19, CA-27, CA-30).
- Toda mensagem mostrada ao usuário, inclusive erros e diálogos nativos, sai no idioma escolhido
  (RN-11, CA-18).

## 2. Decisões

### D-01 — Acesso à API local e conexão do front

- **Decisão:** token aleatório gerado a cada execução, mantido só em memória e exigido nas
  operações que alteram dados; o front o recebe por uma **única chamada interna do Wails de
  bootstrap**, que entrega também o endereço da API, o estado da inicialização (ex.: "porta
  41790 ocupada") e as preferências. Depois do bootstrap, toda comunicação segue pela API HTTP.
  Por quê: o token barra tanto outro programa quanto páginas do navegador, e o canal interno é
  o único jeito de o front saber da porta ocupada (RN-07) sem depender da própria porta.
- **Alternativas avaliadas:**
  - **Checar o cabeçalho `Origin`** — recusar requisições que não venham da origem do Wails.
    _Prós:_ dispensa token. _Contras:_ qualquer programa forja o `Origin` (um `curl` passa).
    Não escolhida porque não atende o CA-33.
  - **Rotas de Projetos só pelo AssetServer do Wails** (em processo), deixando na porta TCP só
    saúde e documentação. _Prós:_ sem token, nada externo alcança as rotas. _Contras:_ conflita
    com a ADR 001 (API na porta, consumida pelo front) e depende de `[HIPÓTESE]` sobre o
    funcionamento no `wails dev` e com POST no WebView2. Não escolhida pelo conflito e pelo
    risco.
- **Viabilidade:** reusa Wails e `net/http`. Conflita com a consequência neutra da ADR 001
  ("o front fala com o motor pela API HTTP, não pelos bindings do Wails"): o bootstrap é uma
  exceção pontual e explícita, limitada a essa única chamada.
- **Trade-off aceito:** ganha-se proteção contra programas e páginas locais e um canal que
  funciona com a porta ocupada; abre-se mão da pureza "só HTTP" da ADR 001 em um ponto.
- **Atende:** US-03, US-13, RN-07, RN-08, CA-07, CA-33, CA-34
- **ADR:** registrada em `docs/adr/004-acesso-a-api-local-token-por-execucao-e-bootstrap-do-front.md`

### D-02 — Erros e idiomas

- **Decisão:** cada camada traduz o que ela mesma produz. O motor devolve todo erro com um
  **código estável** e a **mensagem já traduzida** no idioma salvo, e tem o próprio catálogo
  (erros da API e diálogos nativos mostrados antes de existir janela). O front tem o catálogo
  das telas (rótulos, títulos, textos de status). Por quê: o teste Go do CA-41 enxerga o texto
  traduzido de verdade, a mensagem nativa da RN-04 já exige catálogo no motor, e o código
  estável permite ao front reagir (ex.: destacar o campo que falta, CA-23).
- **Alternativas avaliadas:**
  - **Só o front traduz** — motor devolve apenas o código (ex.: `project.name_taken`).
    _Prós:_ um único catálogo. _Contras:_ `go test` não verifica texto traduzido (CA-41) e o
    diálogo nativo da RN-04 precisaria de catálogo no motor do mesmo jeito. Não escolhida.
  - **Tudo no motor** — o front baixa o dicionário das telas pela API. _Prós:_ fonte única.
    _Contras:_ com a porta ocupada (RN-07), o front fica sem nenhum texto. Não escolhida.
- **Viabilidade:** reusa a biblioteca padrão; nenhuma dependência nova.
- **Trade-off aceito:** dois catálogos, com fronteira clara (o motor traduz o que diz, o front
  traduz o que desenha), em troca de testabilidade e de funcionar com a API fora do ar.
- **Atende:** US-08, US-16, RN-04, RN-11, CA-12, CA-18, CA-20, CA-21, CA-22, CA-23, CA-24, CA-41
- **ADR:** registrada em `docs/adr/005-padrao-de-erros-e-idiomas.md`

### D-03 — Preferências (tema e idioma)

- **Decisão:** arquivo de configuração na pasta do app, **fora do SQLite**, lido pelo motor ao
  iniciar e entregue ao front no bootstrap (D-01); alterações passam pela API. O motor detecta o
  idioma do Windows na primeira execução. Por quê: as preferências precisam sobreviver ao banco
  indisponível (CA-12) e estar disponíveis antes da API (RN-07).
- **Alternativas avaliadas:**
  - **Tabela no SQLite** — junto com os demais dados. _Prós:_ um só lugar de persistência.
    _Contras:_ com o banco fora do ar o idioma se perde. Não escolhida porque quebra o CA-12.
  - **`localStorage` do WebView2** — preferência só no front. _Prós:_ zero código no motor.
    _Contras:_ o motor não enxerga o idioma (não traduz erros nem o diálogo nativo) e os dados
    dependem de detalhes internos do WebView2. Não escolhida.
- **Viabilidade:** reusa a biblioteca padrão.
- **Trade-off aceito:** dois lugares de persistência (arquivo e banco), em troca de as
  preferências independerem do banco.
- **Atende:** US-07, US-08, RN-03, RN-09, RN-10, CA-10, CA-12, CA-14, CA-15, CA-16, CA-17, CA-18
- **ADR:** não

### D-04 — Logs

- **Decisão:** `log/slog` da biblioteca padrão em formato **texto** (`chave=valor`), com
  rotação e descarte por `lumberjack`: arquivos de até **10 MB**, no máximo **5** antigos e
  descarte após **14 dias** (no máximo ~60 MB em disco). Por quê: legível no Bloco de Notas e
  ainda fácil de buscar, com descarte por idade e tamanho já pronto e testado.
- **Alternativas avaliadas:**
  - **`slog` em JSON + `lumberjack`** — uma linha JSON por evento. _Prós:_ fácil de processar,
    útil para uma futura tela de logs. _Contras:_ ruim de ler à mão. Não escolhida porque a
    tela de logs está adiada e hoje quem lê é uma pessoa.
  - **`slog` em texto + rotação própria** — um arquivo por dia e limpeza ao iniciar e a cada
    24 h. _Prós:_ zero dependências. _Contras:_ ~60 linhas nossas para manter e testar. Não
    escolhida por ser reimplementar algo maduro.
- **Viabilidade:** `slog` é biblioteca padrão; `lumberjack` é dependência nova
  `[HIPÓTESE]` — validar a rotação no Windows com o arquivo aberto e o descarte na abertura.
- **Trade-off aceito:** uma dependência pequena em troca de não manter código de rotação.
- **Atende:** US-14, RN-19, CA-07, CA-35, CA-37
- **ADR:** registrada em `docs/adr/006-padrao-de-logging.md`

### D-05 — Esquema do banco

- **Decisão:** migrações SQL numeradas, **embutidas no executável** e aplicadas ao iniciar; a
  versão aplicada fica no próprio arquivo do banco (`PRAGMA user_version`). Cada feature nova
  acrescenta uma migração. Por quê: uma instalação antiga copiada ou atualizada se ajusta sozinha,
  com ~30 linhas e sem dependência.
- **Alternativas avaliadas:**
  - **`CREATE TABLE IF NOT EXISTS` ao iniciar** — _Prós:_ o mais simples. _Contras:_ não
    altera tabela existente; quebra na v2 de qualquer feature. Não escolhida.
  - **Biblioteca (`goose` ou `golang-migrate`)** — _Prós:_ completa, com rollback. _Contras:_
    dependência e conceitos a mais para um app local de um usuário. Não escolhida.
- **Viabilidade:** reusa `modernc.org/sqlite` (ADR 001) e `embed` da biblioteca padrão.
- **Trade-off aceito:** sem rollback de migração (desnecessário num app local), em troca de
  simplicidade.
- **Atende:** US-05, US-15, CA-10, CA-38
- **ADR:** registrada em `docs/adr/007-migracoes-do-banco-de-dados.md`

### D-06 — Banco indisponível (em execução e nos testes)

- **Decisão:** o app sobe em **modo degradado** quando o banco não abre ou cai: a janela abre,
  a saúde reporta o banco como indisponível e as operações de Projetos respondem com erro
  traduzido. No teste de integração do CA-41, a indisponibilidade é simulada **fechando a
  conexão real do SQLite** antes de chamar a API. Por quê: reproduz a falha pelo caminho real
  (SQLite → serviço → API) sem ferramenta nova.
- **Alternativas avaliadas:**
  - **Fake do repositório que devolve erro** — _Prós:_ rápido. _Contras:_ não passa pela
    saúde nem pela API. Não escolhida para o CA-41 (continua sendo o padrão dos testes
    unitários, pela ADR 002).
  - **Servidor que injeta falhas** (ideia do pre-refinement) — _Contras:_ infraestrutura
    inteira para simular uma única falha. Não escolhida.
- **Viabilidade:** reusa `httptest` e SQLite real em pasta temporária (ADR 002).
- **Trade-off aceito:** simula "conexão perdida", não todos os modos de falha do arquivo;
  outras falhas ficaram fora do escopo no PRD.
- **Atende:** US-06, US-16, CA-11, CA-12, CA-24, CA-41
- **ADR:** não

### D-07 — Caminhos guardados

- **Decisão:** decisão direta — a ADR 003 já obriga caminho relativo para o que fica dentro da
  instalação. Fica fixado aqui o resto: banco, logs e configuração **nunca são guardados**
  (derivam da pasta base); a pasta de um projeto **dentro** da instalação é guardada relativa à
  raiz da instalação, e **fora** dela, absoluta (caminho relativo não atravessa unidades, ex.:
  `C:` → `D:`). A unicidade (RN-14) compara o caminho resolvido e normalizado, sem diferenciar
  maiúsculas. Copiada a instalação, um projeto absoluto que não exista no outro computador
  aparece como "pasta não encontrada" (RN-18).
- **Viabilidade:** segue a ADR 003; biblioteca padrão.
- **Trade-off aceito:** projetos fora da instalação não "viajam" com a cópia, o que é coerente
  com o cadastro apenas referenciar pastas.
- **Atende:** US-05, RN-03, RN-14, RN-15, RN-18, CA-10, CA-22, CA-26
- **ADR:** registrada em `docs/adr/003-aplicacao-portatil.md`

### D-08 — Porta fixa da API

- **Decisão:** decisão direta — porta **41790**. Fica fora da faixa dinâmica do Windows
  (49152–65535), da faixa reservada na máquina de desenvolvimento (50000–50059), das portas
  comuns de ferramentas de desenvolvimento (3000, 5173, 8000, 8080, 3306…) e da 47001 (WinRM),
  e estava livre na verificação. Porta ocupada segue a RN-07, e o estado chega ao front pelo
  bootstrap (D-01).
- **Viabilidade:** `net/http`; nada novo.
- **Trade-off aceito:** uma porta fixa pode colidir com algum programa raro; nesse caso o app
  avisa e continua aberto, e a escolha pelo usuário fica para uma versão futura (PRD).
- **Atende:** RN-07, CA-07, CA-13, CA-32
- **ADR:** não

### D-09 — Instalador

- **Decisão:** **NSIS pela integração do Wails v2** (`wails build -nsis`), com o script
  ajustado para rodar como usuário comum, sugerir `C:\dev-tools\`, criar atalho e desinstalador,
  e perguntar no desinstalador sobre dados e logs sem nunca tocar na pasta de projetos. Por quê:
  reusa o que a stack já oferece, inclusive a verificação do WebView2.
- **Alternativas avaliadas:**
  - **Inno Setup** — o do Laragon; sem admin numa linha. _Prós:_ script mais legível.
    _Contras:_ fora do fluxo de build do Wails; verificação do WebView2 feita à mão. Não
    escolhida por não reusar a stack.
  - **MSI com WiX** — _Contras:_ complexo demais para um instalador mínimo. Não escolhida.
- **Viabilidade:** reusa Wails v2; requer o NSIS instalado na máquina de build. O modelo padrão
  do Wails instala em `Program Files` como administrador: ajustá-lo para usuário comum é
  `[HIPÓTESE]` a validar.
- **Trade-off aceito:** a linguagem de script do NSIS é árida, em troca de ficar no fluxo do
  Wails.
- **Atende:** US-01, US-02, RN-01, RN-20, CA-01, CA-02, CA-03
- **ADR:** não

### D-10 — Organização do front

- **Decisão:** **React puro com CSS próprio**: pastas por feature espelhando o motor (ADR 002),
  tema por variáveis CSS (o tema "Sistema" acompanha `prefers-color-scheme`, que o WebView2
  herda do Windows), chamadas com `fetch` e tipos gerados do `openapi.yaml` (já previsto na
  ADR 001), sem bibliotecas de estado nem de componentes. Por quê: é o mínimo que atende às
  telas da v1 e mantém mecânica a cópia de uma feature.
- **Alternativas avaliadas:**
  - **Tailwind + componentes próprios** — _Prós:_ estilização rápida, modo escuro pronto.
    _Contras:_ mais uma ferramenta e convenção. Não escolhida.
  - **Biblioteca de componentes (Mantine, MUI…)** — _Prós:_ formulários, modais e tema prontos.
    _Contras:_ dependência grande que passa a ditar as convenções do front. Não escolhida.
- **Viabilidade:** reusa React + TypeScript + Vite (ADR 001); a geração de tipos a partir do
  OpenAPI é ferramenta nova de build `[HIPÓTESE]`.
- **Trade-off aceito:** formulário e confirmação de remoção feitos à mão, em troca de zero
  dependências de UI.
- **Atende:** US-07, US-09, US-10, US-11, US-12, US-15, RN-10, RN-17, CA-14, CA-15, CA-31, CA-38
- **ADR:** registrada em `docs/adr/008-organizacao-do-front.md`

### D-11 — Onde fica o código do Dev Tools

- **Decisão:** decisão direta (escolha do usuário) — numa **subpasta deste repositório**, ao lado
  de `docs/`. Por quê: specs, skills e código no mesmo workspace e histórico, o que simplifica a
  fase de Execução; se o produto crescer, dá para extrair depois preservando o histórico.
- **Viabilidade:** nada novo.
- **Trade-off aceito:** o histórico do repositório mistura framework e produto.
- **Atende:** US-15, US-16
- **ADR:** não

## 3. Restrições herdadas

- Go 1.25+, Wails v2 (WebView2), React 19 + TypeScript + Vite, SQLite via `modernc.org/sqlite`
  sem CGO, `net/http` com o `ServeMux` padrão, contrato `openapi.yaml` escrito à mão servido numa
  página de documentação, `go test` — ADR 001.
- Um teste confere que as rotas registradas batem com o `openapi.yaml` — ADR 001.
- Motor organizado por feature, núcleo sem importar infraestrutura, infraestrutura comum
  separada, montagem num único ponto de entrada — ADR 002.
- Testes unitários com fakes ao lado do código; testes de integração na própria feature,
  separados por build tag, com SQLite real em pasta temporária e `httptest` — ADR 002.
- Pasta base = pasta do executável ou `DEV_TOOLS_HOME`; nada em registro, `%APPDATA%` ou `PATH`;
  instância única pelo `SingleInstanceLock`; recriar pastas que faltam; fechar com mensagem se a
  pasta não permitir escrita — ADR 003.
- Plataforma Windows, como executável; o app roda sem o instalador durante o desenvolvimento —
  PRD seção 10.
- Instalador sem assinatura digital na v1; se o escopo apertar, o instalador pode ir para a v2
  — PRD seção 10.
- Ferramentas de desenvolvimento (Go, Node.js, NSIS…) instaladas em
  `C:\dev\programs\<programa>\<programa>-vX.X.X`, com o atalho `current` no `PATH` — CLAUDE.md
  global do usuário.

## 4. Pontos em aberto

**Para a Spec** (detalhe de implementação, fora da altitude desta fase):

- Como o token viaja na requisição, quais rotas o exigem e a política de CORS para a origem do
  Wails (D-01).
- Campos exatos entregues pelo bootstrap e como o front reage a cada estado de inicialização
  (D-01).
- Formato do erro na API (código, mensagem, campo, parâmetros) e a lista de códigos da v1 (D-02).
- Como o motor detecta o idioma do Windows na primeira execução (D-03, RN-09).
- Formato do arquivo de configuração e o comportamento se ele estiver ausente ou corrompido
  (D-03).
- Nível de log padrão, eventos registrados e saída no console durante o desenvolvimento (D-04).
- Conteúdo da primeira migração (D-05).
- Se o app tenta reabrir o banco depois de uma falha, ou só ao reiniciar (D-06).
- Layout de pastas da instalação, e qual pasta é a "raiz da instalação" em relação à pasta do
  executável (D-07).
- Como rodar o app em desenvolvimento ao mesmo tempo que o instalado, já que a porta é fixa
  (ex.: sobrescrita só para desenvolvimento e testes) (D-08).
- Ajustes do script NSIS e o tratamento do WebView2 ausente (D-09).
- Estrutura de pastas do front, navegação entre telas e o mecanismo de tradução das telas
  (D-10).
- Nome da subpasta do código (D-11).

**Dependências de produto** (o PRD não decidiu; não foram decididas aqui):

- Com a porta da API ocupada (RN-07), o app continua aberto, mas sem API: o PRD não diz o que o
  usuário pode fazer nesse estado (só ver a mensagem? tentar de novo? trocar tema e idioma?).

**Fora do escopo** (observado, mas não vira proposta nesta feature):

- A futura CLI pode reusar o token da D-01 e o catálogo do motor da D-02.
- Uma tela de logs (adiada) pode pedir que o formato de log mude para JSON.
- A migração para Wails v3 afetaria o bootstrap (D-01) e o instalador (D-09).

## 5. Próximo passo

- **Próxima fase:** Tech Spec → `/agent-spec-tech-spec docs/specs/features/001-arquitetura-base/v1/tech-alignment.md`
- **ADR a registrar antes:**
  - `/agent-spec-adr "D-01 de docs/specs/features/001-arquitetura-base/v1/tech-alignment.md"` — Acesso à API local: token por execução e bootstrap do front
  - `/agent-spec-adr "D-02 de docs/specs/features/001-arquitetura-base/v1/tech-alignment.md"` — Padrão de erros e idiomas
  - `/agent-spec-adr "D-04 de docs/specs/features/001-arquitetura-base/v1/tech-alignment.md"` — Padrão de logging
  - `/agent-spec-adr "D-05 de docs/specs/features/001-arquitetura-base/v1/tech-alignment.md"` — Migrações do banco de dados
  - `/agent-spec-adr "D-10 de docs/specs/features/001-arquitetura-base/v1/tech-alignment.md"` — Organização do front
