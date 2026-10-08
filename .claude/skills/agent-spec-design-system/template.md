---
status: rascunho   # rascunho | aprovado
---

# Design System — <nome do produto>

> Identidade visual do produto: tokens, componentes e padrões que **toda feature com interface
> segue**. Gerado e atualizado por `/agent-spec-design-system`.
>
> **Origem de cada valor:** `[código]` já implementado (com o caminho) · `[referência]` Figma ou
> documento do time (com o link) · `[decisão]` decidido com o usuário (proposto e aceito).

- **Data:** AAAA-MM-DD
- **Plataformas:** web | mobile | desktop
- **Stack de UI:** <!-- LLM-ONLY: detectada no código ou na ADR. Ex.: React + Tailwind, Flutter (Material 3), WPF. Não invente. -->
- **Biblioteca de componentes base:** <!-- LLM-ONLY: ex.: shadcn/ui, MUI, Material 3, nenhuma (componentes próprios) -->
- **Fontes consultadas:** <!-- LLM-ONLY: caminhos do código, links de Figma, documentos do time, ADRs -->
- **ADRs relacionadas:** `docs/adr/NNN-<titulo>.md` | nenhuma

---

## 1. Cores

<!-- LLM-ONLY: se o código já declara o token (variável CSS, chave do tema), use o nome dele.
  Mínimo num projeto novo: primária, neutros (fundo, superfície, borda, texto, texto secundário)
  e os semânticos (sucesso, aviso, erro, informação). -->

| Token | Valor | Uso | Origem |
|-------|-------|-----|--------|
|       |       |     |        |

## 2. Tipografia

- **Fonte(s):** <!-- LLM-ONLY: família e de onde vem (Google Fonts, do sistema, pacote) -->

| Estilo | Tamanho / Peso / Altura de linha | Uso | Origem |
|--------|----------------------------------|-----|--------|
|        |                                  |     |        |

## 3. Espaçamento, raios e sombras

- **Escala de espaçamento:** <!-- LLM-ONLY: ex.: base 4 (4, 8, 12, 16, 24, 32, 48) -->

| Token | Valor | Uso | Origem |
|-------|-------|-----|--------|
|       |       |     |        |

## 4. Ícones

- **Biblioteca:** <!-- LLM-ONLY: ex.: lucide, Material Symbols, SF Symbols -->
- **Tamanhos:**

## 5. Layout

<!-- LLM-ONLY: as faixas de tamanho de tela da plataforma. Web: breakpoints. Mobile: celular,
  tablet e orientação. Desktop: tamanho mínimo da janela e o que muda ao redimensionar. -->

| Faixa | Valor | O que muda |
|-------|-------|------------|
|       |       |            |

- **Densidade:** <!-- LLM-ONLY: ex.: confortável; compacta só em tabelas -->

## 6. Tema

- **Modos:** claro | escuro | claro e escuro | segue o sistema operacional
- **Como o tema é aplicado:** <!-- LLM-ONLY: ex.: tokens semânticos (fundo, texto) sobre a paleta, trocados por classe/tema -->

## 7. Componentes

<!-- LLM-ONLY: só componentes REUTILIZÁVEIS entre features. A Referência é o import real que o
  código usa (`@/components/ui/button`, `import { Button } from '@acme/ui'`): é por ela que o
  executor reusa o componente em vez de criar outro. Componente próprio sem referência que
  exista no repo não entra. Projeto novo, sem componentes ainda: liste os da biblioteca base que
  serão o padrão, ou escreva a frase abaixo da tabela. -->

| Componente | Referência (import) | Variantes | Quando usar |
|------------|---------------------|-----------|-------------|
|            |                     |           |             |

Componente reutilizável novo entra aqui antes de ser usado por uma segunda feature.

## 8. Padrões de feedback

<!-- LLM-ONLY: comportamento CONCRETO (o quê, onde, com qual ação), nunca "mostrar loading".
  É o padrão que toda tela herda; a feature que fugir dele diz por quê. -->

| Situação | Padrão | Componente | Origem |
|----------|--------|------------|--------|
| Carregando |      |            |        |
| Erro     |        |            |        |
| Vazio    |        |            |        |
| Sucesso  |        |            |        |
| Ação destrutiva |  |           |        |

## 9. Acessibilidade visual

- **Contraste mínimo:** <!-- LLM-ONLY: ex.: WCAG AA (4,5:1 em texto normal) -->
- **Foco visível:**
- **Área mínima de toque/clique:** <!-- LLM-ONLY: ex.: 44×44 -->
- **Movimento reduzido:**

## 10. Divergências conhecidas

<!-- LLM-ONLY: onde o código ainda não segue este documento (ex.: 2 telas usam spinner, mas o
  padrão é skeleton). "Nenhuma." se não houver. -->

- <onde> — <o que diverge> — <o que vale: este documento>

## 11. Próximo passo

- **ADR a registrar antes do primeiro spec com tela:**
  - `/agent-spec-adr "<título da decisão>"` — <o que ela decide, em poucas palavras>

<!-- LLM-ONLY: ADR aqui = decisão transversal recomendada por este documento que ainda não
  existe em docs/adr/; sem nenhuma → "**ADR a registrar antes do primeiro spec com tela:**
  nenhuma". -->
