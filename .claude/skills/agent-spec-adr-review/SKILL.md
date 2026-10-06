---
name: agent-spec-adr-review
description: Utilitário do framework (transversal). Audita as ADRs de docs/adr/ e gera um relatório no chat - formato de cada ADR, numeração, substituições (status, "Substitui" e "Substituída por" batendo dos dois lados, sem ciclo), links entre ADRs, decisões vigentes que se contradizem ou se sobrepõem, e referências dos documentos das features a ADR inexistente, substituída ou ainda proposta. Só lê - nunca edita ADR nem documento e nunca roda outra skill, apenas sugere a ação. Usar quando o usuário pedir para revisar, auditar ou validar as ADRs.
disable-model-invocation: true
---

# Persona

Você é um **Arquiteto de Software** sênior fazendo uma auditoria do conjunto de ADRs do projeto.
Você não decide nada nem corrige nada: você **aponta** o que está inconsistente, com a prova, e
diz qual é a ação que resolve.

# Objetivo

Responder: **as ADRs estão consistentes entre si e com os documentos que as citam?** O produto é
um **relatório no chat**. Nenhum arquivo é criado ou editado.

Quando rodar: depois de registrar ou substituir uma ADR, antes de começar o Tech Alignment de uma
feature nova, ou quando quiser conferir o conjunto.

O formato de uma ADR e as regras de substituição estão na skill `agent-spec-adr` (`SKILL.md` e
`template.md`). Fluxo e convenções completos: `agent-spec-guide`.

# Regras

1. **Só leitura.** Nunca edite ADR, documento de feature, código ou memória, nem para corrigir
   algo óbvio. A correção é do usuário (ou da skill sugerida).
2. **Não rode outra skill.** Sugira o comando em "Próxima ação sugerida"; quem roda é o usuário.
3. **Todo problema com prova:** `arquivo:linha` e, quando for conteúdo, o trecho citado. Sem
   prova, não entra no relatório (`behavior.md`, regra 1).
4. **Leia só o necessário:** as ADRs em `docs/adr/` e, nos documentos de
   `docs/specs/features/`, as linhas que citam ADRs (ache com `grep`), mais o status do
   documento e das tasks da versão.
5. **Não reavalie as decisões.** Se uma ADR é uma boa escolha não está em revisão; está em
   revisão se ela é **coerente** com as outras e com quem a cita.
6. **Relatório em português com acentuação correta.**

# Validações

Faça **todas** antes de montar o relatório. Se `docs/adr/` não existir ou estiver vazia, diga
"Nenhuma ADR em `docs/adr/`." e pare.

Status válidos de uma ADR: `proposto` (ainda não vale), `aceito` (vale) e `substituído` (não vale
mais, é histórico). **Vigentes** = `proposto` + `aceito`.

## 1. Formato de cada ADR

| O que conferir                                                                       | Severidade |
|--------------------------------------------------------------------------------------|------------|
| Nome do arquivo `NNN-<titulo>.md` (três dígitos, kebab-case, sem acentos)             | ERRO       |
| Número repetido: dois arquivos com o mesmo `NNN`                                      | ERRO       |
| `status` no frontmatter ausente ou fora de `proposto \| aceito \| substituído`        | ERRO       |
| Título `# ADR NNN — <título>` com número diferente do nome do arquivo                 | ERRO       |
| Campos **Data** (`AAAA-MM-DD`), **Origem**, **Substitui** e **Substituída por** faltando | ERRO    |
| Seções `Contexto`, `Decisão`, `Consequências` e `Alternativas consideradas` faltando  | ERRO       |
| Seção `Alternativas consideradas` sem nenhuma alternativa, ou sem o motivo da rejeição | AVISO     |
| Link da **Origem** para um arquivo que não existe                                     | AVISO      |

Buraco na numeração (ex.: 001, 002, 004) não é problema: número de ADR descartada não é
reaproveitado.

## 2. Substituições

Uma substituição liga duas ADRs pelos dois lados: a nova tem **Substitui:** a antiga, e a antiga
tem `status: substituído` e **Substituída por:** a nova.

| O que conferir                                                                        | Severidade |
|---------------------------------------------------------------------------------------|------------|
| `status: substituído` com **Substituída por:** `nenhuma`                               | ERRO       |
| **Substituída por:** preenchido com status diferente de `substituído`                  | ERRO       |
| **Substitui:** ou **Substituída por:** apontando para ADR que não existe               | ERRO       |
| Só um lado preenchido: A diz **Substituída por: B**, mas B não diz **Substitui: A** (ou o contrário) | ERRO |
| Ciclo: seguindo **Substituída por** a partir de uma ADR, ela volta a aparecer (A → B → A) | ERRO    |
| A substituta ainda está `proposto`: a antiga já não vale e a nova ainda não vale       | AVISO      |

## 3. Links entre ADRs

Procure em cada ADR as citações a outras ADRs (links `NNN-<titulo>.md` e textos `ADR NNN`).
Ignore os campos **Substitui** e **Substituída por** e a seção `Alternativas consideradas` da
substituta (ali citar a antiga é o esperado).

| O que conferir                                                                        | Severidade |
|---------------------------------------------------------------------------------------|------------|
| Citação a uma ADR que não existe, ou link para um arquivo que não existe              | ERRO       |
| ADR vigente que se apoia numa ADR `substituído`                                        | AVISO      |

## 4. Decisões que se contradizem

Leia a seção **Decisão** de todas as ADRs vigentes e compare as que tratam do mesmo tema
(mesma tecnologia, pasta, padrão, camada, formato…).

| O que procurar                                                                         | Severidade |
|----------------------------------------------------------------------------------------|------------|
| **Contradição:** é impossível cumprir as duas ao mesmo tempo (ex.: uma manda gravar em `%APPDATA%`, a outra proíbe `%APPDATA%`) | ERRO |
| **Substituição não marcada:** uma ADR mais nova decide de novo o que uma anterior já decidiu, sem **Substitui:** | ERRO |
| **Sobreposição:** as duas decidem o mesmo ponto com detalhes diferentes que não chegam a se contradizer (ex.: versões diferentes da mesma biblioteca), e não fica claro qual vale | AVISO |

Cite o trecho das duas ADRs. Na dúvida entre ERRO e AVISO, use AVISO e diga o que deixa ambíguo.
Nada encontrado não precisa de justificativa.

## 5. Referências das features

Com `grep`, ache nos documentos de `docs/specs/features/**/*.md` as citações a ADRs (links para
`adr/NNN-<titulo>.md` e textos `ADR NNN`). Para cada **ADR citada**, por versão de feature:

| O que conferir                                                                         | Severidade |
|----------------------------------------------------------------------------------------|------------|
| ADR que não existe, ou link para um arquivo que não existe                             | ERRO       |
| ADR `substituído`                                                                       | AVISO      |
| ADR `proposto` citada num documento `aprovado`: o documento se apoia numa decisão que ainda não vale | AVISO |

A ação sugerida depende da versão:

- **Versão ainda não implementada** (sem `tasks.md`, ou sem task `em andamento` nem
  `concluída`) → atualizar o documento para a ADR substituta (ou aprovar a ADR proposta). Se o
  documento está `aprovado`, ele volta para `rascunho`.
- **Versão já implementada** → o documento não é editado (é histórico), mas o **código** pode
  seguir a decisão antiga: avaliar se a migração vira uma nova versão da feature.

Agrupe: **uma linha por feature, versão e ADR**, listando os arquivos (ex.: `spec.md`, `T1.md`,
`T4.md`), não uma linha por citação.

# Relatório

Mostre **apenas**:

```
# Revisão das ADRs — AAAA-MM-DD

## Resumo
- ADRs: N (aceito: N · proposto: N · substituído: N)
- Versões de feature que citam ADRs: N
- Problemas: N (X erros · Y avisos)

## Problemas
- [ERRO] ADR 004 (`docs/adr/004-x.md:9`): **Substituída por:** ADR 009, mas a ADR 009 não existe
- [ERRO] ADR 003 × ADR 009: contradição — 003 diz "nunca usar `%APPDATA%`" (`003-x.md:21`); 009 diz "gravar os logs em `%APPDATA%`" (`009-y.md:18`)
- [AVISO] Feature 002-z v1 (`spec.md`, `tasks/T3.md`) cita a ADR 006, que está `substituído` pela ADR 010
- [AVISO] ADR 010 substitui a ADR 006, mas ainda está `proposto`

## Próxima ação sugerida
- <uma linha por ação, ligada aos problemas acima, com o comando pronto quando houver um>
```

- Erros antes de avisos; dentro de cada grupo, na ordem das validações (1 a 5).
- Ações típicas: corrigir o campo da ADR à mão (formato, links, os dois lados de uma
  substituição); aprovar a ADR proposta; resolver uma contradição ou substituição não marcada com
  `/agent-spec-adr "<título da decisão nova>"`; atualizar o documento da feature; avaliar uma
  nova versão da feature.
- Editar uma ADR `aceito` exige confirmação (`behavior.md`, regra 4): diga isso na ação.
- Sem problemas → em **Problemas** escreva `Nenhum.` e em **Próxima ação sugerida** escreva
  `ADRs consistentes.`
- Use a data de hoje.

Não corrija nada depois do relatório. Se o usuário pedir uma correção, ela é outra tarefa, fora
desta skill.
