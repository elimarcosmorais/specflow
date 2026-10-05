---
status: rascunho   # rascunho | aprovado
---

# Tech Spec — <nome da feature>

- **Feature:** `NNN-<nome-da-feature>` · **Versão:** vN · **Data:** AAAA-MM-DD
- **Origem:** [prd.md](prd.md) · [tech-alignment.md](tech-alignment.md)
- **ADRs aplicáveis:** `docs/adr/NNN-<titulo>.md` | nenhuma

<!-- LLM-ONLY: em mudança pequena (versão sem prd.md), a Origem aponta para o pre-refinement.md desta versão e para os documentos mais recentes das versões anteriores (ex.: ../v1/prd.md, ../v1/spec.md). -->

---

## 1. Resumo técnico

<3 a 5 linhas: como a feature será construída, as peças principais e o resultado técnico
esperado.>

**Decisões do Tech Alignment aplicadas:**

| Decisão                 | Onde está nesta spec |
|-------------------------|----------------------|
| D-01 — <ponto>          | §3.1, §8             |

**Mudanças em relação à vN-1:** <o que muda e por quê, em poucas linhas>

<!-- LLM-ONLY: toda D-NN do tech-alignment aparece na tabela. "Mudanças em relação à vN-1" só existe a partir da v2; na v1, remova a linha. -->

## 2. Arquitetura

### 2.1 Componentes

| Componente | Camada | Responsabilidade | Situação |
|------------|--------|------------------|----------|
| <nome>     | <núcleo, infraestrutura, apresentação, front…> | <1 frase> | novo \| modifica \| reusa |

### 2.2 Dependências entre camadas

<Quem chama quem, num diagrama ASCII curto ou em lista. Diga também o que cada camada **não
pode** importar.>

## 3. Contratos

### 3.1 API

| Operação | Método e rota | Entrada | Saída (sucesso) | Erros | Proteção |
|----------|---------------|---------|-----------------|-------|----------|
| <ação>   | `GET /…`      | <campos> | `200` + <campos> | `404 <código>` | <pública \| token…> |

<!-- LLM-ONLY: atualização parcial (PATCH, ou PUT em que campos podem faltar) ganha um exemplo de payload mínimo e a frase "campos ausentes não são alterados e não são obrigatórios", para o executor não copiar a validação do POST. -->

### 3.2 Erros

| Código     | Quando acontece | Status HTTP | Mensagem (pt-br) |
|------------|-----------------|-------------|------------------|
| <código>   | <condição>      | `409`       | <texto>          |

### 3.3 Outras interfaces

<Chamadas internas, eventos, arquivos lidos ou gravados por outros programas, linha de comando…
Formato de cada uma.>

## 4. Dados

### 4.1 Modelo

| Tabela   | Coluna | Tipo | Restrições | Observação |
|----------|--------|------|------------|------------|
| <tabela> | <coluna> | <tipo> | <PK, NOT NULL, UNIQUE…> | <regra que atende> |

### 4.2 Migrações

| Arquivo  | O que faz |
|----------|-----------|
| <arquivo> | <cria a tabela X…> |

### 4.3 Outros dados guardados

<Arquivos de configuração, cache etc.: formato, local e comportamento se ausente ou inválido.>

## 5. Regras e validações

| Regra | Origem | Onde é aplicada | Em caso de falha |
|-------|--------|-----------------|------------------|
| <regra técnica> | RN-01 | <componente> | <código de erro \| comportamento> |

<!-- LLM-ONLY: toda RN do PRD aparece pelo menos uma vez na coluna Origem. Regra técnica sem RN de origem usa "—". -->

## 6. Fluxos

### 6.1 Fluxo principal

1. <camada → camada: o que acontece em cada passo, até a resposta>

### 6.2 Fluxos alternativos e de erro

- **<situação>** → <o que cada camada faz e o que o usuário vê>

## 7. Interface

### 7.1 Telas e navegação

| Tela   | Como se chega | O que mostra | Ações |
|--------|---------------|--------------|-------|
| <tela> | <menu, rota…> | <dados>      | <botões, formulários> |

### 7.2 Estados

| Tela   | Estado | O que o usuário vê |
|--------|--------|--------------------|
| <tela> | carregando \| vazio \| erro \| <outro> | <descrição> |

<!-- LLM-ONLY: feature sem interface → "N/A — <motivo>" nas duas subseções. -->

## 8. Aspectos transversais

- **Segurança:** <autenticação, autorização, validação de entrada, exposição de rede>
- **Logs:** <eventos registrados, nível, dados que não podem aparecer>
- **Idiomas:** <o que é traduzido e onde fica cada texto>
- **Configuração:** <valores configuráveis e padrões>
- **Inicialização e encerramento:** <o que acontece ao abrir e fechar>

<!-- LLM-ONLY: só os aspectos que a feature toca; remova os outros. Acrescente um aspecto se a feature tocar algo não listado (ex.: desempenho). -->

## 9. Cobertura do PRD

| História | Atendida por (componentes, operações, telas) | Seções |
|----------|----------------------------------------------|--------|
| US-01    | <…>                                          | §3.1, §7 |

<!-- LLM-ONLY: toda US do PRD tem uma linha. -->

## 10. Estratégia de testes

> **Resumo:** X casos · unitários Y · integração Z · e2e W · manuais M
> **Padrão:** <framework, onde ficam os arquivos, como se nomeiam, como se separam as camadas>

### 10.1 Critérios de aceite → testes

| CA    | Resumo   | Testes       |
|-------|----------|--------------|
| CA-01 | <1 linha> | CT-01, CT-02 |

### 10.2 Casos de teste

#### <Componente> (`<arquivo de teste>`)

| CT    | Camada   | O que verifica | Entrada / cenário | Resultado esperado | CA |
|-------|----------|----------------|-------------------|--------------------|----|
| CT-01 | unitário | <regra>        | <dados>           | <saída>            | CA-01 |

<!-- LLM-ONLY: preencha 10.1 primeiro: ela guia os casos. Todo CA tem pelo menos um CT. Camadas: unitário | integração | e2e | manual. Manual só quando automatizar não é viável (ex.: instalador sem admin), com o roteiro em "Entrada / cenário". Agrupe por componente; teste manual fica num grupo "Manuais". -->

## 11. Arquivos

```
<raiz>/
├── <pasta>/
│   ├── <arquivo>        [N] <o que contém>
│   └── <arquivo>        [M] <o que muda>
└── <arquivo>            [R] <por que consultar>
```

Legenda: `[N]` novo · `[M]` modificado · `[R]` referência (só leitura)

<!-- LLM-ONLY: todos os arquivos da feature, inclusive testes, migrações e documentação (ex.: openapi.yaml). Nada que já existe no código aparece como [N]. -->

## 12. Conformidade com ADRs

| ADR | Aplica? | Como esta spec cumpre |
|-----|---------|-----------------------|
| [ADR 001](../../../../adr/001-<titulo>.md) | sim \| parcial \| não | "<trecho da Decisão da ADR>" → <onde e como a spec cumpre, §> |

<!-- LLM-ONLY: só ADRs com status aceito. Em "sim" e "parcial", cite o texto real da seção Decisão da ADR. Sem ADRs aceitas no projeto → "N/A — nenhuma ADR aceita". -->

## 13. Riscos e hipóteses

| Risco ou `[HIPÓTESE]` | Impacto | Como validar ou mitigar |
|-----------------------|---------|-------------------------|
| <…>                   | <…>     | <teste, protótipo, fallback> |

## 14. Próximo passo

- **Próxima fase:** Task Plan → `/agent-spec-task-plan docs/specs/features/NNN-<feature>/vN/spec.md`
- **ADR a registrar antes:**
  - `/agent-spec-adr "<título da decisão>"` — <o que ela decide, em poucas palavras>

<!-- LLM-ONLY: um comando por decisão nova desta spec que vale para o projeto todo. Sem candidatas → "**ADR a registrar antes:** nenhuma". -->
