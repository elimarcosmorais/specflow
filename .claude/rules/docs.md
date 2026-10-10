# Documentação do projeto

Onde cada documento mora, quem pode mexer nele e quando. Vale em qualquer tarefa.

## Onde cada coisa mora

| Lugar | Guarda | Quem mexe |
|---|---|---|
| `docs/specs/` | **como o sistema funciona hoje**, por módulo ou assunto | o agente, livremente |
| `docs/architecture.md` | a arquitetura base: camadas, módulos, fluxos e integrações principais | o agente, livremente |
| `docs/adr/NNN-<titulo>.md` | **por que** uma decisão transversal foi tomada | só `/spec-adr`, com o usuário |
| `docs/design-system.md` + `docs/design/` | a identidade visual: tokens, componentes e padrões | `/spec-design-system`; sincronizar com o código é livre |
| `docs/plans/AAAA-MM-DD-<slug>.md` | o plano de uma mudança média ou grande, até ser executado | `/spec-change` cria, o usuário aprova, `/spec-build` executa e apaga |
| `docs/references/` | material de consulta temporário do usuário: marca, prints, prompts | o usuário |

**Não duplique.** A spec diz **como funciona** e aponta para a ADR que diz **por quê**. O plano
diz o que uma mudança vai fazer e some depois dela: o que sobrevive vai para a spec (como
funciona), para uma ADR (por quê) ou para o `CLAUDE.md`/rules (regra para todo código novo).
Pastas e arquivos só são criados quando há conteúdo para eles.

## Specs e `architecture.md`: livres

Crie, edite, mova, subdivida e apague **sem pedir permissão**. Ao fim da resposta, uma linha:
`Specs: <arquivo> (criado | atualizado | subdividido | apagado), …`.

Atualize quando:

- uma mudança altera o comportamento que uma spec descreve: na **mesma** mudança, nunca depois;
- um assunto novo passa de um parágrafo;
- você descobre algo que custaria tempo para redescobrir: armadilha, regra de negócio implícita,
  detalhe de integração.

Não entra: o que o código diz sozinho, o porquê de uma decisão (ADR), plano ou pendência.

**Índice.** `docs/specs/README.md` lista cada arquivo com a coluna "Abra quando" (a situação em
que ele deve ser lido). Arquivo criado entra no índice; apagado sai. Antes de mexer num módulo,
leia o índice e abra só a spec dele.

## Tamanho máximo: 300 linhas

Vale para spec, `architecture.md`, `design-system.md` e os arquivos de `docs/design/`. Passou de
300 linhas (`wc -l`) → **subdivida na mesma tarefa, sem pedir permissão**, e diga no fim:

1. `<assunto>.md` vira índice: um resumo de poucas linhas e a tabela "Abra quando" dos pedaços.
2. O conteúdo vai para `<assunto>/<sub-assunto>.md`, um sub-assunto por arquivo.
3. Os links que apontavam para o trecho movido passam a apontar para o arquivo novo.

A profundidade é livre: um pedaço que passar de 300 linhas se subdivide do mesmo jeito.

## Design system: livre para sincronizar

Registrar o que o código **já implementa** é livre, como nas specs. Mudar a identidade (cor,
token, componente ou padrão que o código ainda não tem) é decisão do usuário: pergunte antes, ou
rode `/spec-design-system`.

## Planos: temporários

Um plano por mudança, nome `AAAA-MM-DD-<slug>.md`. `rascunho` → `aprovado` (pelo usuário) →
executado e **apagado pelo build** no fechamento: a remoção fica em stage junto com o código e
entra no mesmo commit. Plano não vira documentação.

## Referências: temporárias

`docs/references/` é material de consulta do usuário e pode sumir a qualquer momento. Leia à
vontade, mas **só o plano** o cita, porque também é temporário. Nenhum documento permanente
(specs, `architecture.md`, ADR, design system, `CLAUDE.md`, README, rules) cita, linka ou
documenta esses arquivos. O que ele precisar de uma referência (uma cor, uma regra, um fluxo) é
escrito nele como fato do projeto, sem o caminho. Quem apaga é o usuário; o agente não apaga
nem move nada ali.

## Nomes

kebab-case, minúsculas, sem acentos. Conteúdo em português com acentuação correta; nomes de
código ficam como no código.
