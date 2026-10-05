---
status: aprovado   # rascunho | aprovado
---

# PRD — Arquitetura base do Dev Tools

> **Legenda:** `[DECISÃO TÉCNICA]` = ponto que envolve o COMO; fica registrado aqui de forma
> abstrata e é decidido no Tech Alignment.

- **Feature:** `001-arquitetura-base` · **Versão:** v1 · **Data:** 2026-10-05
- **Origem:** [pre-refinement.md](pre-refinement.md)
- **ADRs relacionadas:**
  - [`docs/adr/001-stack-base-do-dev-tools.md`](../../../../adr/001-stack-base-do-dev-tools.md)
  - [`docs/adr/002-arquitetura-em-camadas-e-padrao-de-feature.md`](../../../../adr/002-arquitetura-em-camadas-e-padrao-de-feature.md)
  - [`docs/adr/003-aplicacao-portatil.md`](../../../../adr/003-aplicacao-portatil.md)

---

## 1. Contexto e motivação

- **Qual é o problema hoje?** O Dev Tools ainda não existe. Sem uma base definida, cada feature
  futura (PHP, Node.js, Apache, MySQL…) inventaria o seu próprio jeito de tratar logs, erros,
  dados guardados, idioma e testes, e o sistema ficaria inconsistente e difícil de manter.
  Exemplo: a feature "PHP" grava logs num formato e a feature "Node.js" em outro, e uma mensagem
  de erro aparece traduzida numa tela e só em inglês na outra.
- **Quem sofre o impacto?** O desenvolvedor do Dev Tools, que precisa de um padrão único para
  construir as próximas features, e o usuário final, que recebe um app inconsistente se esse
  padrão não existir.
- **Por que resolver agora?** É o primeiro passo do produto: todas as outras features dependem
  dessa base. A motivação vem da rotina atual de instalar cada ferramenta à mão e trocar de
  versão manualmente, que o Dev Tools vai automatizar nas próximas features.

## 2. Objetivo

- **Resultado esperado, do ponto de vista do usuário:** existe um Dev Tools instalável no
  Windows que abre sem pedir administrador, mostra o próprio estado de saúde, permite escolher
  tema e idioma, e permite cadastrar, listar, editar e remover Projetos. Para o desenvolvedor,
  existe uma feature completa (Projetos) que serve de modelo para todas as próximas, com API
  local documentada, logs e testes.
- **O que muda no comportamento de quem usa:** o usuário final passa a ter um lugar único para
  registrar os seus projetos locais. O desenvolvedor passa a criar uma feature nova copiando o
  modelo de Projetos, em vez de decidir de novo como tratar erros, idioma, dados, logs e testes.

## 3. Personas

- **Primária:** desenvolvedor do Dev Tools — quer uma base clara para construir as próximas
  features sempre do mesmo jeito, com testes e documentação desde o início.
- **Secundária:** usuário final, desenvolvedor que usa Windows — quer instalar o app sem
  complicação e organizar o seu ambiente local de desenvolvimento (nesta v1, o próprio
  desenvolvedor do Dev Tools).

## 4. Escopo

**Dentro:**

- App desktop para Windows que abre sem pedir administrador e em instância única.
- Funcionamento portátil: todos os dados ficam dentro da pasta de instalação.
- Estrutura de pastas verificada ao iniciar, recriando o que faltar.
- Saúde da aplicação (versão e disponibilidade do banco de dados), visível numa tela e
  consultável pela API local.
- Dados guardados localmente em banco de dados, dentro da pasta do app.
- Logs gravados em arquivo, úteis para o uso e para o desenvolvimento, com descarte automático
  dos antigos e atalho para abrir a pasta de logs.
- Idiomas pt-br e en, selecionáveis e guardados entre execuções.
- Tema Claro, Escuro ou Sistema, selecionável e guardado entre execuções.
- API local, acessível só a partir do próprio computador, com documentação consultável.
- Cadastro de Projetos (criar, listar, editar, remover) como feature-modelo, com validação e
  mensagens de erro traduzidas.
- Separação entre regra de negócio, infraestrutura e testes, seguida pela feature-modelo.
- Testes unitários e de integração cobrindo a feature-modelo, incluindo a simulação de banco
  indisponível.
- Instalador mínimo (escolha de pasta sugerindo `C:\dev-tools\`, atalho e desinstalador).

**Fora** (com motivo):

- Instalar e gerenciar programas (PHP, Node.js, Apache, MySQL…) — adiado: próximas features.
- Servir os projetos do usuário em localhost (servidor web, hosts virtuais) — adiado: feature
  futura.
- Iniciar e parar serviços — adiado: envolve processos, portas e programas externos.
- Listagem de programas instalados — adiado: candidata a próxima feature.
- CLI no terminal — adiado.
- Atualização automática do app e atalhos extras no instalador — adiado.
- Opção "adicionar ao PATH do usuário" — adiado para uma versão futura.
- Dados fora da pasta do app (pastas do sistema, registro, PATH global) e execução como
  administrador — descartado: conflita com a portabilidade.
- Escolha da porta da API pelo usuário — adiado: na v1 a porta é fixa.
- Tela de logs dentro do app — adiado: na v1 os logs ficam só nos arquivos.
- Criar a pasta do projeto pelo app, ou apagá-la ao remover o projeto — descartado: o cadastro
  só referencia pastas existentes.
- Simulação de falhas além do banco indisponível (disco cheio, porta ocupada) nos testes —
  adiado.
- Idiomas além de pt-br e en — não pedido.

## 5. Histórias de usuário

- **US-01:** Como usuário final, quero instalar o Dev Tools escolhendo a pasta de instalação
  para ter o app no meu computador sem precisar de administrador.
- **US-02:** Como usuário final, quero desinstalar o Dev Tools decidindo se meus dados e logs
  vão junto para não perder nada sem querer.
- **US-03:** Como usuário final, quero abrir o Dev Tools sem pedido de administrador e sempre
  na mesma janela para usá-lo sem fricção e sem instâncias duplicadas.
- **US-04:** Como usuário final, quero que o app recrie sozinho as pastas que faltarem na
  instalação para que ele continue funcionando se eu apagar alguma sem querer.
- **US-05:** Como usuário final, quero levar o Dev Tools para outro computador copiando a pasta
  de instalação para manter o app e os meus dados juntos.
- **US-06:** Como usuário final, quero ver o estado de saúde do app (versão e banco de dados)
  para saber se ele está funcionando corretamente.
- **US-07:** Como usuário final, quero escolher o tema (Claro, Escuro ou Sistema) para usar o
  app do jeito mais confortável, sem precisar escolher de novo a cada execução.
- **US-08:** Como usuário final, quero escolher o idioma (pt-br ou en) para ler todas as telas e
  mensagens no idioma que prefiro, sem precisar escolher de novo a cada execução.
- **US-09:** Como usuário final, quero cadastrar um projeto (nome, pasta e tipo) para registrar
  no Dev Tools os projetos que tenho no computador.
- **US-10:** Como usuário final, quero ver a lista dos projetos cadastrados para saber o que
  tenho e onde cada um está.
- **US-11:** Como usuário final, quero editar um projeto cadastrado para corrigir ou atualizar o
  seu nome, pasta ou tipo.
- **US-12:** Como usuário final, quero remover um projeto do cadastro para tirar da lista o que
  não uso mais, sem perder os arquivos dele.
- **US-13:** Como desenvolvedor do Dev Tools, quero consultar a documentação da API local para
  saber quais operações existem e como usá-las.
- **US-14:** Como desenvolvedor do Dev Tools, quero que o app grave logs em arquivo e ofereça um
  atalho para a pasta de logs para entender o que aconteceu quando algo der errado.
- **US-15:** Como desenvolvedor do Dev Tools, quero uma feature-modelo com regra de negócio,
  infraestrutura e testes separados para criar as próximas features seguindo um padrão único.
- **US-16:** Como desenvolvedor do Dev Tools, quero rodar testes unitários e de integração,
  incluindo a simulação de banco indisponível, para saber que a base e a feature-modelo
  funcionam antes de seguir.

## 6. Regras de negócio

- **RN-01:** O app nunca pede privilégios de administrador, nem para instalar nem para rodar.
- **RN-02:** Só uma instância do app roda por vez; abrir o app de novo traz para a frente a
  janela já aberta.
- **RN-03:** Todos os dados do app (banco, configurações, logs, cache…) ficam dentro da pasta de
  instalação; nada é gravado em pastas do sistema, no registro ou no PATH.
- **RN-04:** Se a pasta de instalação não permitir escrita, o app não abre e mostra uma mensagem
  clara explicando o problema.
- **RN-05:** Ao iniciar, o app recria as pastas da sua estrutura que estiverem faltando, sem
  alterar o conteúdo das que já existem.
- **RN-06:** O app funciona sem conexão com a internet.
- **RN-07:** A API local usa uma porta fixa. Se a porta estiver ocupada por outro programa, o app
  mostra uma mensagem clara informando qual porta está em uso, registra o ocorrido no log e
  continua aberto.
- **RN-08:** A API local só é acessível a partir do próprio computador. As operações de Projetos
  só são aceitas quando vêm do próprio Dev Tools; a saúde e a documentação podem ser
  consultadas pelo navegador.
- **RN-09:** Os idiomas disponíveis são pt-br e en. Na primeira execução, o idioma é pt-br se o
  Windows estiver em português do Brasil, e en em qualquer outro caso. A escolha do usuário é
  guardada entre execuções.
- **RN-10:** Os temas disponíveis são Claro, Escuro e Sistema. Na primeira execução, o tema é
  Sistema, que acompanha o tema do Windows. A escolha do usuário é guardada entre execuções.
- **RN-11:** Toda mensagem exibida ao usuário, inclusive as de erro, aparece no idioma escolhido.
- **RN-12:** Um projeto tem nome, pasta e tipo, todos obrigatórios. O tipo é um entre web,
  android e windows.
- **RN-13:** Não pode haver dois projetos com o mesmo nome, ignorando maiúsculas e minúsculas.
- **RN-14:** Não pode haver dois projetos apontando para a mesma pasta, ignorando maiúsculas e
  minúsculas.
- **RN-15:** A pasta de um projeto precisa existir no momento de cadastrar ou editar, e pode
  estar em qualquer lugar do disco.
- **RN-16:** O cadastro só referencia a pasta: criar, editar ou remover um projeto nunca cria,
  move, copia ou apaga arquivos ou pastas no disco.
- **RN-17:** Remover um projeto exige confirmação do usuário.
- **RN-18:** Um projeto cuja pasta deixou de existir continua cadastrado e aparece na lista
  marcado como "pasta não encontrada".
- **RN-19:** Os logs antigos são descartados automaticamente ao passar de um limite de idade ou
  de tamanho.
- **RN-20:** O desinstalador sempre remove o app, pergunta se deve remover também os dados e os
  logs, e nunca apaga a pasta de projetos.

## 7. Fluxo de uso

**Principal:**

1. O usuário executa o instalador, aceita a pasta sugerida (`C:\dev-tools\`) ou escolhe outra.
2. O instalador copia o app e cria o atalho e o desinstalador, sem pedir administrador.
3. O usuário abre o app pelo atalho.
4. O app verifica a estrutura de pastas, recria o que falta e abre a janela no tema e no idioma
   padrão (ou nos já escolhidos antes).
5. O usuário consulta a tela de status: versão do app e banco de dados disponível.
6. O usuário escolhe o tema e o idioma; a tela muda na hora e a escolha fica guardada.
7. O usuário cadastra um projeto informando nome, pasta existente e tipo.
8. O app valida os dados, salva o projeto e o mostra na lista.
9. O usuário edita um projeto da lista; o app valida e mostra os dados atualizados.
10. O usuário remove um projeto, confirma, e o projeto sai da lista; a pasta continua no disco.
11. O desenvolvedor abre a documentação da API no navegador e consulta as operações
    disponíveis.

**Alternativos:**

- **App já aberto** → abrir de novo não cria uma segunda janela; a janela existente vem para a
  frente.
- **Pasta de instalação sem permissão de escrita** → o app não abre e mostra uma mensagem clara
  explicando o problema.
- **Porta da API ocupada** → o app mostra uma mensagem informando qual porta está em uso,
  registra o ocorrido no log e continua aberto.
- **Pasta da estrutura apagada** → o app a recria ao iniciar.
- **Banco de dados indisponível** → a tela de status mostra o banco como indisponível e as
  operações de Projetos respondem com um erro no idioma escolhido.
- **Cadastro ou edição com pasta inexistente, nome repetido, pasta já usada ou campo
  obrigatório vazio** → o app recusa, mostra o motivo no idioma escolhido e não salva nada.
- **Pasta de um projeto apagada ou movida depois do cadastro** → o projeto aparece na lista
  marcado como "pasta não encontrada" e pode ser editado ou removido.
- **Usuário cancela a remoção** → o projeto continua na lista.
- **Operação de Projetos chamada por outro programa ou página do navegador** → a operação é
  recusada e nada é alterado.
- **Desinstalação** → o app é removido; o usuário decide se os dados e logs vão junto; a pasta
  de projetos fica intacta.

## 8. Critérios de aceite

- [ ] **CA-01:** DADO um Windows com uma conta sem privilégios de administrador QUANDO o usuário
      executa o instalador ENTÃO o instalador sugere `C:\dev-tools\`, permite escolher outra
      pasta, copia o app e cria o atalho e o desinstalador, sem pedir administrador.
- [ ] **CA-02:** DADO o app instalado, com projetos cadastrados QUANDO o usuário desinstala e
      responde "sim" à pergunta sobre dados e logs ENTÃO o app, os dados e os logs são
      removidos, e a pasta de projetos e o seu conteúdo continuam no disco.
- [ ] **CA-03:** DADO o app instalado QUANDO o usuário desinstala e responde "não" à pergunta
      sobre dados e logs ENTÃO o app é removido, e os dados, os logs e a pasta de projetos
      continuam no disco.
- [ ] **CA-04:** DADO o app instalado e uma conta sem privilégios de administrador QUANDO o
      usuário abre o app ENTÃO a janela abre sem nenhum pedido de administrador.
- [ ] **CA-05:** DADO o app já aberto QUANDO o usuário abre o app de novo ENTÃO nenhuma segunda
      janela é criada e a janela existente vem para a frente.
- [ ] **CA-06:** DADO uma pasta de instalação sem permissão de escrita QUANDO o usuário abre o
      app ENTÃO o app não abre e mostra uma mensagem explicando que a pasta não permite escrita.
- [ ] **CA-07:** DADO outro programa usando a porta da API QUANDO o usuário abre o app ENTÃO o
      app mostra uma mensagem informando qual porta está em uso, o ocorrido fica registrado no
      log e a janela do app continua aberta.
- [ ] **CA-08:** DADO um computador sem conexão com a internet QUANDO o usuário abre o app e
      cadastra, lista, edita e remove projetos ENTÃO todas as operações funcionam normalmente.
- [ ] **CA-09:** DADO que a pasta de logs foi apagada QUANDO o usuário abre o app ENTÃO a pasta
      de logs é recriada, e as outras pastas da estrutura continuam com o mesmo conteúdo.
- [ ] **CA-10:** DADO uma instalação com projetos cadastrados, tema e idioma escolhidos QUANDO a
      pasta de instalação é copiada para outro computador Windows e o app é aberto lá ENTÃO os
      mesmos projetos, tema e idioma aparecem, e nada é gravado fora da pasta copiada.
- [ ] **CA-11:** DADO o app aberto e o banco de dados disponível QUANDO o usuário abre a tela de
      status ENTÃO ela mostra a versão do app e o banco como disponível.
- [ ] **CA-12:** DADO o banco de dados indisponível QUANDO o usuário abre a tela de status ENTÃO
      ela mostra o banco como indisponível, no idioma escolhido.
- [ ] **CA-13:** DADO o app aberto QUANDO o desenvolvedor consulta a saúde pela API no navegador
      ENTÃO recebe a versão do app e o estado do banco de dados.
- [ ] **CA-14:** DADO a primeira execução do app QUANDO a janela abre ENTÃO o tema é Sistema e
      acompanha o tema claro ou escuro do Windows.
- [ ] **CA-15:** DADO o usuário escolheu o tema Escuro QUANDO fecha e abre o app de novo ENTÃO o
      app abre no tema Escuro.
- [ ] **CA-16:** DADO a primeira execução em um Windows em português do Brasil QUANDO a janela
      abre ENTÃO o app está em pt-br.
- [ ] **CA-17:** DADO a primeira execução em um Windows em outro idioma (ex.: espanhol) QUANDO a
      janela abre ENTÃO o app está em en.
- [ ] **CA-18:** DADO o usuário escolheu o idioma en QUANDO fecha e abre o app de novo ENTÃO
      todas as telas e mensagens, inclusive as de erro, aparecem em en.
- [ ] **CA-19:** DADO uma pasta existente que nenhum projeto usa QUANDO o usuário cadastra um
      projeto com nome inédito, essa pasta e o tipo web ENTÃO o projeto aparece na lista, e o
      conteúdo da pasta no disco não muda.
- [ ] **CA-20:** DADO uma pasta que não existe no disco QUANDO o usuário tenta cadastrar um
      projeto com ela ENTÃO o app recusa, informa no idioma escolhido que a pasta não existe e
      nenhum projeto é salvo.
- [ ] **CA-21:** DADO um projeto cadastrado com o nome "loja-api" QUANDO o usuário tenta
      cadastrar outro com o nome "Loja-API" ENTÃO o app recusa, informa que o nome já está em
      uso e nenhum projeto é salvo.
- [ ] **CA-22:** DADO um projeto cadastrado com uma pasta QUANDO o usuário tenta cadastrar outro
      projeto com a mesma pasta ENTÃO o app recusa, informa que a pasta já está cadastrada e
      nenhum projeto é salvo.
- [ ] **CA-23:** DADO o formulário de cadastro QUANDO o usuário tenta salvar sem nome, sem pasta
      ou sem tipo ENTÃO o app recusa, indica qual campo falta no idioma escolhido e nenhum
      projeto é salvo.
- [ ] **CA-24:** DADO o banco de dados indisponível QUANDO o usuário tenta cadastrar um projeto
      ENTÃO o app mostra um erro no idioma escolhido e nenhum projeto é salvo.
- [ ] **CA-25:** DADO projetos cadastrados QUANDO o usuário abre a lista de projetos ENTÃO cada
      projeto aparece com o seu nome, tipo e pasta.
- [ ] **CA-26:** DADO um projeto cuja pasta foi apagada depois do cadastro QUANDO o usuário abre
      a lista ENTÃO o projeto aparece marcado como "pasta não encontrada" e continua disponível
      para editar e remover.
- [ ] **CA-27:** DADO um projeto cadastrado QUANDO o usuário altera o seu nome, pasta e tipo para
      valores válidos ENTÃO a lista mostra os dados novos, e nenhuma pasta é criada, movida ou
      apagada no disco.
- [ ] **CA-28:** DADO dois projetos cadastrados QUANDO o usuário edita um deles para usar o nome
      ou a pasta do outro ENTÃO o app recusa, informa o motivo e o projeto mantém os dados
      anteriores.
- [ ] **CA-29:** DADO um projeto cadastrado QUANDO o usuário edita a pasta para uma que não
      existe ENTÃO o app recusa, informa que a pasta não existe e o projeto mantém os dados
      anteriores.
- [ ] **CA-30:** DADO um projeto cadastrado QUANDO o usuário pede para removê-lo e confirma
      ENTÃO o projeto sai da lista, e a pasta e os arquivos dele continuam no disco.
- [ ] **CA-31:** DADO um projeto cadastrado QUANDO o usuário pede para removê-lo e cancela a
      confirmação ENTÃO o projeto continua na lista.
- [ ] **CA-32:** DADO o app aberto QUANDO o desenvolvedor abre a documentação da API no navegador
      ENTÃO vê todas as operações de saúde e de Projetos, com o que cada uma recebe e responde.
- [ ] **CA-33:** DADO o app aberto QUANDO outro programa ou uma página no navegador tenta criar,
      editar ou remover um projeto pela API ENTÃO a operação é recusada e os projetos
      cadastrados não mudam.
- [ ] **CA-34:** DADO o app aberto QUANDO outro computador da rede tenta acessar a API ENTÃO não
      recebe resposta.
- [ ] **CA-35:** DADO o app em uso QUANDO um projeto é cadastrado ou ocorre um erro ENTÃO um
      registro com data, hora e descrição do evento é gravado no arquivo de log.
- [ ] **CA-36:** DADO a tela de status aberta QUANDO o usuário clica em "abrir pasta de logs"
      ENTÃO o Explorer do Windows abre na pasta de logs.
- [ ] **CA-37:** DADO arquivos de log que passaram do limite de idade ou de tamanho QUANDO o app
      está em execução ENTÃO esses logs são descartados automaticamente e os recentes são
      mantidos.
- [ ] **CA-38:** DADO a feature-modelo de Projetos QUANDO o desenvolvedor a consulta ENTÃO
      encontra a regra de negócio, a infraestrutura e os testes separados, e a regra de negócio
      pode ser verificada sem banco de dados e sem a API.
- [ ] **CA-39:** DADO o código da base QUANDO o desenvolvedor roda os testes unitários ENTÃO as
      regras de Projetos (RN-12 a RN-16) são verificadas sem banco de dados e sem a API, e todos
      os testes passam.
- [ ] **CA-40:** DADO o código da base QUANDO o desenvolvedor roda os testes de integração ENTÃO
      cadastro, listagem, edição, remoção e saúde são verificados com banco de dados real e pela
      API, e todos os testes passam.
- [ ] **CA-41:** DADO um teste que simula o banco de dados indisponível QUANDO o desenvolvedor o
      roda ENTÃO o teste confirma que a saúde reporta o banco como indisponível e que as
      operações de Projetos respondem com erro traduzido.

## 9. Rastreabilidade

| História | Resumo                                   | Critérios de aceite                       |
|----------|------------------------------------------|-------------------------------------------|
| US-01    | Instalar escolhendo a pasta, sem admin   | CA-01                                     |
| US-02    | Desinstalar decidindo sobre dados e logs | CA-02, CA-03                              |
| US-03    | Abrir sem admin, instância única         | CA-04, CA-05, CA-06, CA-07, CA-08         |
| US-04    | Recriar pastas que faltam                | CA-09                                     |
| US-05    | Levar a instalação para outro PC         | CA-10                                     |
| US-06    | Ver a saúde do app                       | CA-11, CA-12, CA-13                       |
| US-07    | Escolher o tema                          | CA-14, CA-15                              |
| US-08    | Escolher o idioma                        | CA-16, CA-17, CA-18                       |
| US-09    | Cadastrar projeto                        | CA-19, CA-20, CA-21, CA-22, CA-23, CA-24  |
| US-10    | Listar projetos                          | CA-25, CA-26                              |
| US-11    | Editar projeto                           | CA-27, CA-28, CA-29                       |
| US-12    | Remover projeto                          | CA-30, CA-31                              |
| US-13    | Consultar a documentação da API          | CA-32, CA-33, CA-34                       |
| US-14    | Logs em arquivo e atalho para a pasta    | CA-35, CA-36, CA-37                       |
| US-15    | Feature-modelo com camadas separadas     | CA-38                                     |
| US-16    | Testes unitários e de integração         | CA-39, CA-40, CA-41                       |

## 10. Premissas e restrições

**Premissas** (confirmadas com o usuário):

- Uso local por um único usuário, em computador pessoal; a API não precisa ser acessível pela
  rede.
- O app funciona offline.
- Um projeto cadastrado apenas referencia uma pasta; o cadastro não move, copia nem apaga
  arquivos.
- Instalar em `C:\dev-tools\` não exige administrador, porque o Windows permite que um usuário
  comum crie pastas na raiz do `C:\` por padrão (a confirmar na máquina-alvo durante a
  implementação).
- Tema padrão Sistema; idioma padrão igual ao do Windows (en se não for pt-br).

**Restrições** (prazo, dependência, imposição externa):

- Plataforma: Windows, como programa executável.
- Stack, organização em camadas e portabilidade já decididas nas ADRs 001, 002 e 003.
- O app também roda sem o instalador durante o desenvolvimento.
- O instalador da v1 não tem assinatura digital; o aviso de segurança do Windows é aceitável no
  uso pessoal e deve ser reavaliado antes de distribuir o app.
- Se o escopo apertar, o instalador pode virar a v2.
- `[DECISÃO TÉCNICA]` Número da porta fixa da API.
- `[DECISÃO TÉCNICA]` Mecanismo que garante que só o Dev Tools executa operações de Projetos
  pela API (RN-08).
- `[DECISÃO TÉCNICA]` Limite de idade e de tamanho para o descarte de logs (RN-19).
- `[DECISÃO TÉCNICA]` Como os caminhos guardados se comportam quando a instalação é copiada
  para outro computador (RN-03, CA-10).
- `[DECISÃO TÉCNICA]` Como simular o banco de dados indisponível nos testes (CA-41).
- `[DECISÃO TÉCNICA]` Onde fica o código do Dev Tools: repositório próprio ou dentro deste
  projeto.

## 11. Métricas de sucesso

- A próxima feature do Dev Tools é criada seguindo o modelo de Projetos, sem inventar um jeito
  novo de tratar logs, erros, idioma, dados ou testes.
- Instalado em `C:\dev-tools\` numa conta comum, o app abre sem pedir administrador e cria a
  estrutura de pastas.
- Todos os testes unitários e de integração passam, inclusive o de banco indisponível.
- Trocando entre pt-br e en, nenhuma tela ou mensagem de erro aparece no outro idioma.

## 12. Próximo passo

- **Próxima fase:** Tech Alignment → `/agent-spec-tech-alignment docs/specs/features/001-arquitetura-base/v1/prd.md`
- **ADR a registrar antes:** nenhuma (as ADRs 001, 002 e 003 recomendadas pelo pre-refinement
  já estão registradas).
